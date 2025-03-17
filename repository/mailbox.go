package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"mailbox-api/db"
	"mailbox-api/model"

	"github.com/jackc/pgx/v4"
)

type MailboxRepository interface {
	GetMailboxes(ctx context.Context, filter model.MailboxFilter) ([]model.Mailbox, int, error)
	GetMailboxByIdentifier(ctx context.Context, identifier string) (*model.Mailbox, error)
	GetMailboxesByRole(ctx context.Context, role string) ([]model.Mailbox, error)
	GetAllMailboxes(ctx context.Context) ([]model.Mailbox, error)
	CreateMailbox(ctx context.Context, mailbox model.Mailbox) error
	UpdateOrgDepth(ctx context.Context, identifier string, depth int) error
	UpdateSubOrgSize(ctx context.Context, identifier string, size int) error
	CalculateOrgMetrics(ctx context.Context) error
}

type mailboxRepository struct {
	db *db.DB
}

func NewMailboxRepository(db *db.DB) MailboxRepository {
	return &mailboxRepository{db: db}
}

func (r *mailboxRepository) GetMailboxes(ctx context.Context, filter model.MailboxFilter) ([]model.Mailbox, int, error) {
	query := `
	SELECT 
		m.mailbox_identifier, 
		m.user_full_name, 
		m.job_title, 
		m.department_id, 
		d.department_name, 
		m.manager_mailbox_identifier, 
		m.org_depth, 
		m.sub_org_size
	FROM 
		mailboxes m
	JOIN 
		departments d ON m.department_id = d.department_id
	WHERE 1=1`

	countQuery := `
	SELECT 
		COUNT(*)
	FROM 
		mailboxes m
	JOIN 
		departments d ON m.department_id = d.department_id
	WHERE 1=1`

	params := []interface{}{}
	paramIndex := 1

	if filter.SearchTerm != "" {
		searchCondition := fmt.Sprintf(`
		AND (
			m.user_full_name ILIKE $%d 
			OR m.job_title ILIKE $%d 
			OR d.department_name ILIKE $%d
		)`, paramIndex, paramIndex, paramIndex)
		query += searchCondition
		countQuery += searchCondition
		params = append(params, "%"+filter.SearchTerm+"%")
		paramIndex++
	}

	if filter.Department != 0 {
		departmentCondition := fmt.Sprintf(`
		AND m.department_id = $%d`, paramIndex)
		query += departmentCondition
		countQuery += departmentCondition
		params = append(params, filter.Department)
		paramIndex++
	}

	if filter.OrgDepthExact != nil {
		orgDepthCondition := fmt.Sprintf(`
		AND m.org_depth = $%d`, paramIndex)
		query += orgDepthCondition
		countQuery += orgDepthCondition
		params = append(params, *filter.OrgDepthExact)
		paramIndex++
	}

	if filter.OrgDepthGt != nil {
		orgDepthGtCondition := fmt.Sprintf(`
		AND m.org_depth > $%d`, paramIndex)
		query += orgDepthGtCondition
		countQuery += orgDepthGtCondition
		params = append(params, *filter.OrgDepthGt)
		paramIndex++
	}

	if filter.OrgDepthLt != nil {
		orgDepthLtCondition := fmt.Sprintf(`
		AND m.org_depth < $%d`, paramIndex)
		query += orgDepthLtCondition
		countQuery += orgDepthLtCondition
		params = append(params, *filter.OrgDepthLt)
		paramIndex++
	}

	if filter.SubOrgSizeMin != nil {
		subOrgSizeMinCondition := fmt.Sprintf(`
		AND m.sub_org_size >= $%d`, paramIndex)
		query += subOrgSizeMinCondition
		countQuery += subOrgSizeMinCondition
		params = append(params, *filter.SubOrgSizeMin)
		paramIndex++
	}

	if filter.SubOrgSizeMax != nil {
		subOrgSizeMaxCondition := fmt.Sprintf(`
		AND m.sub_org_size <= $%d`, paramIndex)
		query += subOrgSizeMaxCondition
		countQuery += subOrgSizeMaxCondition
		params = append(params, *filter.SubOrgSizeMax)
		paramIndex++
	}

	if len(filter.SortBy) > 0 && len(filter.SortBy) == len(filter.SortDirections) {
		query += " ORDER BY "
		sorts := []string{}

		for i, field := range filter.SortBy {
			direction := "ASC"
			if strings.ToUpper(filter.SortDirections[i]) == "DESC" {
				direction = "DESC"
			}

			var dbField string
			switch field {
			case "mailbox_identifier":
				dbField = "m.mailbox_identifier"
			case "user_full_name":
				dbField = "m.user_full_name"
			case "job_title":
				dbField = "m.job_title"
			case "department_id":
				dbField = "m.department_id"
			case "department":
				dbField = "d.department_name"
			case "org_depth":
				dbField = "m.org_depth"
			case "sub_org_size":
				dbField = "m.sub_org_size"
			default:
				dbField = "m.user_full_name"
			}

			sorts = append(sorts, fmt.Sprintf("%s %s", dbField, direction))
		}

		query += strings.Join(sorts, ", ")
	} else {
		query += " ORDER BY m.user_full_name ASC"
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
		params = append(params, filter.PageSize, offset)
		paramIndex += 2
	}

	var totalCount int
	err := r.db.QueryRow(ctx, countQuery, params[:len(params)-2]...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count mailboxes: %w", err)
	}

	rows, err := r.db.Query(ctx, query, params...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query mailboxes: %w", err)
	}
	defer rows.Close()

	mailboxes := []model.Mailbox{}
	for rows.Next() {
		var mailbox model.Mailbox
		var managerId sql.NullString // Use sql.NullString to handle NULL values
		err := rows.Scan(
			&mailbox.Identifier,
			&mailbox.UserFullName,
			&mailbox.JobTitle,
			&mailbox.DepartmentID,
			&mailbox.Department,
			&managerId, // Scan into NullString
			&mailbox.OrgDepth,
			&mailbox.SubOrgSize,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan mailbox: %w", err)
		}

		// Convert NullString to string
		if managerId.Valid {
			mailbox.ManagerIdentifier = managerId.String
		} else {
			mailbox.ManagerIdentifier = "" // Empty string for NULL
		}
		mailboxes = append(mailboxes, mailbox)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over mailboxes: %w", err)
	}

	return mailboxes, totalCount, nil
}

func (r *mailboxRepository) GetMailboxByIdentifier(ctx context.Context, identifier string) (*model.Mailbox, error) {
	query := `
	SELECT 
		m.mailbox_identifier, 
		m.user_full_name, 
		m.job_title, 
		m.department_id, 
		d.department_name, 
		m.manager_mailbox_identifier, 
		m.org_depth, 
		m.sub_org_size
	FROM 
		mailboxes m
	JOIN 
		departments d ON m.department_id = d.department_id
	WHERE 
		m.mailbox_identifier = $1`

	var mailbox model.Mailbox
	err := r.db.QueryRow(ctx, query, identifier).Scan(
		&mailbox.Identifier,
		&mailbox.UserFullName,
		&mailbox.JobTitle,
		&mailbox.DepartmentID,
		&mailbox.Department,
		&mailbox.ManagerIdentifier,
		&mailbox.OrgDepth,
		&mailbox.SubOrgSize,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get mailbox: %w", err)
	}

	return &mailbox, nil
}

func (r *mailboxRepository) GetMailboxesByRole(ctx context.Context, role string) ([]model.Mailbox, error) {
	query := `
	SELECT 
		m.mailbox_identifier, 
		m.user_full_name, 
		m.job_title, 
		m.department_id, 
		d.department_name, 
		m.manager_mailbox_identifier, 
		m.org_depth, 
		m.sub_org_size
	FROM 
		mailboxes m
	JOIN 
		departments d ON m.department_id = d.department_id
	WHERE 
		m.job_title ILIKE $1`

	rows, err := r.db.Query(ctx, query, role)
	if err != nil {
		return nil, fmt.Errorf("failed to query mailboxes by role: %w", err)
	}
	defer rows.Close()

	mailboxes := []model.Mailbox{}
	for rows.Next() {
		var mailbox model.Mailbox
		var managerId sql.NullString // Use sql.NullString to handle NULL values
		err := rows.Scan(
			&mailbox.Identifier,
			&mailbox.UserFullName,
			&mailbox.JobTitle,
			&mailbox.DepartmentID,
			&mailbox.Department,
			&managerId, // Scan into NullString
			&mailbox.OrgDepth,
			&mailbox.SubOrgSize,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mailbox: %w", err)
		}

		// Convert NullString to string
		if managerId.Valid {
			mailbox.ManagerIdentifier = managerId.String
		} else {
			mailbox.ManagerIdentifier = "" // Empty string for NULL
		}
		mailboxes = append(mailboxes, mailbox)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over mailboxes: %w", err)
	}

	return mailboxes, nil
}

func (r *mailboxRepository) GetAllMailboxes(ctx context.Context) ([]model.Mailbox, error) {
	query := `
	SELECT 
		m.mailbox_identifier, 
		m.user_full_name, 
		m.job_title, 
		m.department_id, 
		d.department_name, 
		m.manager_mailbox_identifier, 
		m.org_depth, 
		m.sub_org_size
	FROM 
		mailboxes m
	JOIN 
		departments d ON m.department_id = d.department_id
	ORDER BY m.mailbox_identifier`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all mailboxes: %w", err)
	}
	defer rows.Close()

	mailboxes := []model.Mailbox{}
	for rows.Next() {
		var mailbox model.Mailbox
		var managerId sql.NullString // Use sql.NullString to handle NULL values
		err := rows.Scan(
			&mailbox.Identifier,
			&mailbox.UserFullName,
			&mailbox.JobTitle,
			&mailbox.DepartmentID,
			&mailbox.Department,
			&managerId, // Scan into NullString
			&mailbox.OrgDepth,
			&mailbox.SubOrgSize,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan mailbox: %w", err)
		}
		// Convert NullString to string
		if managerId.Valid {
			mailbox.ManagerIdentifier = managerId.String
		} else {
			mailbox.ManagerIdentifier = "" // Empty string for NULL
		}
		mailboxes = append(mailboxes, mailbox)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over mailboxes: %w", err)
	}

	return mailboxes, nil
}

func (r *mailboxRepository) CreateMailbox(ctx context.Context, mailbox model.Mailbox) error {
	query := `
	INSERT INTO mailboxes (
		mailbox_identifier, 
		user_full_name, 
		job_title, 
		department_id, 
		manager_mailbox_identifier, 
		org_depth, 
		sub_org_size
	) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.db.Exec(ctx, query,
		mailbox.Identifier,
		mailbox.UserFullName,
		mailbox.JobTitle,
		mailbox.DepartmentID,
		mailbox.ManagerIdentifier,
		mailbox.OrgDepth,
		mailbox.SubOrgSize,
	)

	if err != nil {
		return fmt.Errorf("failed to create mailbox: %w", err)
	}

	return nil
}

func (r *mailboxRepository) UpdateOrgDepth(ctx context.Context, identifier string, depth int) error {
	query := `
	UPDATE mailboxes 
	SET org_depth = $1 
	WHERE mailbox_identifier = $2`

	_, err := r.db.Exec(ctx, query, depth, identifier)
	if err != nil {
		return fmt.Errorf("failed to update org depth: %w", err)
	}

	return nil
}

func (r *mailboxRepository) UpdateSubOrgSize(ctx context.Context, identifier string, size int) error {
	query := `
	UPDATE mailboxes 
	SET sub_org_size = $1 
	WHERE mailbox_identifier = $2`

	_, err := r.db.Exec(ctx, query, size, identifier)
	if err != nil {
		return fmt.Errorf("failed to update sub-org size: %w", err)
	}

	return nil
}

func (r *mailboxRepository) CalculateOrgMetrics(ctx context.Context) error {
	mailboxes, err := r.GetAllMailboxes(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all mailboxes: %w", err)
	}

	mailboxMap := make(map[string]*model.Mailbox)
	for i := range mailboxes {
		mailboxMap[mailboxes[i].Identifier] = &mailboxes[i]
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	for _, mailbox := range mailboxes {
		depth := 0
		current := mailbox.Identifier
		visited := make(map[string]bool)

		for current != "" && !visited[current] {
			visited[current] = true
			parent := mailboxMap[current].ManagerIdentifier
			if parent == "" {
				break
			}
			depth++
			current = parent
		}

		_, err := tx.Exec(ctx, `UPDATE mailboxes SET org_depth = $1 WHERE mailbox_identifier = $2`, depth, mailbox.Identifier)
		if err != nil {
			return fmt.Errorf("failed to update org depth for %s: %w", mailbox.Identifier, err)
		}
	}

	for _, mailbox := range mailboxes {
		size := calculateSubOrgSize(mailbox.Identifier, mailboxMap)

		_, err := tx.Exec(ctx, `UPDATE mailboxes SET sub_org_size = $1 WHERE mailbox_identifier = $2`, size, mailbox.Identifier)
		if err != nil {
			return fmt.Errorf("failed to update sub-org size for %s: %w", mailbox.Identifier, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func calculateSubOrgSize(identifier string, mailboxMap map[string]*model.Mailbox) int {
	size := 0

	for _, mailbox := range mailboxMap {
		if mailbox.ManagerIdentifier == identifier {
			size++
			size += calculateSubOrgSize(mailbox.Identifier, mailboxMap)
		}
	}

	return size
}
