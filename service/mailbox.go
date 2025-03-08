package service

import (
	"context"
	"fmt"
	"strings"

	"mailbox-api/dto"
	"mailbox-api/model"
	"mailbox-api/repository"
)

type MailboxService interface {
	GetMailboxes(ctx context.Context, filter model.MailboxFilter) (*model.MailboxResponse, error)
	GetMailboxByIdentifier(ctx context.Context, identifier string) (*model.Mailbox, error)
	CalculateOrgMetrics(ctx context.Context) error
	GetMailboxesInSubOrg(ctx context.Context, managerIdentifier string, filter model.MailboxFilter) (*model.MailboxResponse, error)
	IsMailboxInSubOrg(ctx context.Context, managerIdentifier string, mailboxIdentifier string) (bool, error)
	ImportMailboxesFromCSV(ctx context.Context, csvData string) error
	ImportDepartmentsFromCSV(ctx context.Context, csvData string) error
}

type mailboxService struct {
	mailboxRepo    repository.MailboxRepository
	departmentRepo repository.DepartmentRepository
}

func NewMailboxService(mailboxRepo repository.MailboxRepository, departmentRepo repository.DepartmentRepository) MailboxService {
	return &mailboxService{
		mailboxRepo:    mailboxRepo,
		departmentRepo: departmentRepo,
	}
}

func (s *mailboxService) GetMailboxes(ctx context.Context, filter model.MailboxFilter) (*model.MailboxResponse, error) {
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}

	mailboxes, totalCount, err := s.mailboxRepo.GetMailboxes(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get mailboxes: %w", err)
	}

	totalPages := (totalCount + filter.PageSize - 1) / filter.PageSize
	pagination := &model.Pagination{
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}

	var result interface{} = mailboxes
	if len(filter.Fields) > 0 {
		result = dto.FilterMailboxFields(mailboxes, filter.Fields)
	}

	return dto.NewMailboxResponse(result, pagination), nil
}

func (s *mailboxService) GetMailboxByIdentifier(ctx context.Context, identifier string) (*model.Mailbox, error) {
	mailbox, err := s.mailboxRepo.GetMailboxByIdentifier(ctx, identifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get mailbox: %w", err)
	}

	return mailbox, nil
}

func (s *mailboxService) CalculateOrgMetrics(ctx context.Context) error {
	if err := s.mailboxRepo.CalculateOrgMetrics(ctx); err != nil {
		return fmt.Errorf("failed to calculate org metrics: %w", err)
	}
	return nil
}

func (s *mailboxService) GetMailboxesInSubOrg(ctx context.Context, managerIdentifier string, filter model.MailboxFilter) (*model.MailboxResponse, error) {
	manager, err := s.GetMailboxByIdentifier(ctx, managerIdentifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get manager: %w", err)
	}

	if manager == nil {
		return nil, fmt.Errorf("manager not found: %s", managerIdentifier)
	}

	allMailboxes, err := s.mailboxRepo.GetAllMailboxes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all mailboxes: %w", err)
	}

	mailboxMap := make(map[string]*model.Mailbox)
	for i := range allMailboxes {
		mailboxMap[allMailboxes[i].Identifier] = &allMailboxes[i]
	}

	subOrgMailboxes := []model.Mailbox{}
	findSubOrg(manager.Identifier, mailboxMap, &subOrgMailboxes)

	filteredMailboxes := []model.Mailbox{}
	for _, mailbox := range subOrgMailboxes {
		if filter.SearchTerm != "" {
			searchTerm := strings.ToLower(filter.SearchTerm)
			if !strings.Contains(strings.ToLower(mailbox.UserFullName), searchTerm) &&
				!strings.Contains(strings.ToLower(mailbox.JobTitle), searchTerm) &&
				!strings.Contains(strings.ToLower(mailbox.Department), searchTerm) {
				continue
			}
		}

		if filter.Department != 0 && mailbox.DepartmentID != filter.Department {
			continue
		}

		if filter.OrgDepthExact != nil && mailbox.OrgDepth != *filter.OrgDepthExact {
			continue
		}
		if filter.OrgDepthGt != nil && mailbox.OrgDepth <= *filter.OrgDepthGt {
			continue
		}
		if filter.OrgDepthLt != nil && mailbox.OrgDepth >= *filter.OrgDepthLt {
			continue
		}

		if filter.SubOrgSizeMin != nil && mailbox.SubOrgSize < *filter.SubOrgSizeMin {
			continue
		}
		if filter.SubOrgSizeMax != nil && mailbox.SubOrgSize > *filter.SubOrgSizeMax {
			continue
		}

		filteredMailboxes = append(filteredMailboxes, mailbox)
	}

	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}
	if filter.Page <= 0 {
		filter.Page = 1
	}

	totalCount := len(filteredMailboxes)
	startIndex := (filter.Page - 1) * filter.PageSize
	endIndex := startIndex + filter.PageSize
	if endIndex > totalCount {
		endIndex = totalCount
	}

	var pagedMailboxes []model.Mailbox
	if startIndex < totalCount {
		pagedMailboxes = filteredMailboxes[startIndex:endIndex]
	} else {
		pagedMailboxes = []model.Mailbox{}
	}

	totalPages := (totalCount + filter.PageSize - 1) / filter.PageSize
	pagination := &model.Pagination{
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalItems: totalCount,
		TotalPages: totalPages,
	}

	var result interface{} = pagedMailboxes
	if len(filter.Fields) > 0 {
		result = dto.FilterMailboxFields(pagedMailboxes, filter.Fields)
	}

	return dto.NewMailboxResponse(result, pagination), nil
}

func (s *mailboxService) IsMailboxInSubOrg(ctx context.Context, managerIdentifier string, mailboxIdentifier string) (bool, error) {
	if managerIdentifier == mailboxIdentifier {
		return true, nil // Manager can see their own mailbox
	}

	allMailboxes, err := s.mailboxRepo.GetAllMailboxes(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get all mailboxes: %w", err)
	}

	mailboxMap := make(map[string]*model.Mailbox)
	for i := range allMailboxes {
		mailboxMap[allMailboxes[i].Identifier] = &allMailboxes[i]
	}

	subOrgMailboxIdentifiers := make(map[string]bool)
	var subOrgMailboxes []model.Mailbox
	findSubOrg(managerIdentifier, mailboxMap, &subOrgMailboxes)

	for _, mailbox := range subOrgMailboxes {
		subOrgMailboxIdentifiers[mailbox.Identifier] = true
	}

	return subOrgMailboxIdentifiers[mailboxIdentifier], nil
}

func findSubOrg(managerID string, mailboxMap map[string]*model.Mailbox, result *[]model.Mailbox) {
	for id, mailbox := range mailboxMap {
		if mailbox.ManagerIdentifier == managerID {
			*result = append(*result, *mailbox)
			findSubOrg(id, mailboxMap, result)
		}
	}
}

func (s *mailboxService) ImportMailboxesFromCSV(ctx context.Context, csvData string) error {
	lines := strings.Split(csvData, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("invalid CSV data: no data rows found")
	}

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) < 5 {
			return fmt.Errorf("invalid CSV line at row %d: expected 5 fields, got %d", i+1, len(fields))
		}

		var departmentID int
		_, err := fmt.Sscanf(fields[3], "%d", &departmentID)
		if err != nil {
			return fmt.Errorf("invalid department ID at row %d: %w", i+1, err)
		}

		mailbox := model.Mailbox{
			Identifier:        fields[0],
			UserFullName:      fields[1],
			JobTitle:          fields[2],
			DepartmentID:      departmentID,
			ManagerIdentifier: fields[4],
			OrgDepth:          0,
			SubOrgSize:        0,
		}

		if mailbox.ManagerIdentifier == "null" {
			mailbox.ManagerIdentifier = ""
		}

		if err := s.mailboxRepo.CreateMailbox(ctx, mailbox); err != nil {
			return fmt.Errorf("failed to create mailbox at row %d: %w", i+1, err)
		}
	}

	if err := s.CalculateOrgMetrics(ctx); err != nil {
		return fmt.Errorf("failed to calculate org metrics: %w", err)
	}

	return nil
}

func (s *mailboxService) ImportDepartmentsFromCSV(ctx context.Context, csvData string) error {
	lines := strings.Split(csvData, "\n")
	if len(lines) < 2 {
		return fmt.Errorf("invalid CSV data: no data rows found")
	}

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		fields := strings.Split(line, ",")
		if len(fields) < 2 {
			return fmt.Errorf("invalid CSV line at row %d: expected 2 fields, got %d", i+1, len(fields))
		}

		var departmentID int
		_, err := fmt.Sscanf(fields[0], "%d", &departmentID)
		if err != nil {
			return fmt.Errorf("invalid department ID at row %d: %w", i+1, err)
		}

		department := model.Department{
			ID:   departmentID,
			Name: fields[1],
		}

		if err := s.departmentRepo.CreateDepartment(ctx, department); err != nil {
			return fmt.Errorf("failed to create department at row %d: %w", i+1, err)
		}
	}

	return nil
}
