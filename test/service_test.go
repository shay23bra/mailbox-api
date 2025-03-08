package test

import (
	"mailbox-api/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMailboxFilter tests the mailbox filter functionality
func TestMailboxFilter(t *testing.T) {
	// Test empty filter
	filter := model.MailboxFilter{}
	assert.Equal(t, 0, filter.Department)
	assert.Equal(t, "", filter.SearchTerm)
	assert.Nil(t, filter.OrgDepthExact)
	assert.Nil(t, filter.OrgDepthGt)
	assert.Nil(t, filter.OrgDepthLt)
	assert.Nil(t, filter.SubOrgSizeMin)
	assert.Nil(t, filter.SubOrgSizeMax)

	// Test filter with values
	orgDepthExact := 2
	orgDepthGt := 1
	orgDepthLt := 3
	subOrgSizeMin := 5
	subOrgSizeMax := 10

	filter = model.MailboxFilter{
		SearchTerm:     "test",
		Department:     1,
		OrgDepthExact:  &orgDepthExact,
		OrgDepthGt:     &orgDepthGt,
		OrgDepthLt:     &orgDepthLt,
		SubOrgSizeMin:  &subOrgSizeMin,
		SubOrgSizeMax:  &subOrgSizeMax,
		SortBy:         []string{"user_full_name"},
		SortDirections: []string{"asc"},
		Fields:         []string{"mailbox_identifier", "user_full_name"},
		Page:           1,
		PageSize:       10,
	}

	assert.Equal(t, "test", filter.SearchTerm)
	assert.Equal(t, 1, filter.Department)
	assert.Equal(t, 2, *filter.OrgDepthExact)
	assert.Equal(t, 1, *filter.OrgDepthGt)
	assert.Equal(t, 3, *filter.OrgDepthLt)
	assert.Equal(t, 5, *filter.SubOrgSizeMin)
	assert.Equal(t, 10, *filter.SubOrgSizeMax)
	assert.Equal(t, []string{"user_full_name"}, filter.SortBy)
	assert.Equal(t, []string{"asc"}, filter.SortDirections)
	assert.Equal(t, []string{"mailbox_identifier", "user_full_name"}, filter.Fields)
	assert.Equal(t, 1, filter.Page)
	assert.Equal(t, 10, filter.PageSize)
}

// TestMailboxModel tests the mailbox model
func TestMailboxModel(t *testing.T) {
	mailbox := model.Mailbox{
		Identifier:        "test@example.com",
		UserFullName:      "Test User",
		JobTitle:          "Test Title",
		DepartmentID:      1,
		Department:        "Test Department",
		ManagerIdentifier: "manager@example.com",
		OrgDepth:          2,
		SubOrgSize:        3,
	}

	assert.Equal(t, "test@example.com", mailbox.Identifier)
	assert.Equal(t, "Test User", mailbox.UserFullName)
	assert.Equal(t, "Test Title", mailbox.JobTitle)
	assert.Equal(t, 1, mailbox.DepartmentID)
	assert.Equal(t, "Test Department", mailbox.Department)
	assert.Equal(t, "manager@example.com", mailbox.ManagerIdentifier)
	assert.Equal(t, 2, mailbox.OrgDepth)
	assert.Equal(t, 3, mailbox.SubOrgSize)
}

// TestDepartmentModel tests the department model
func TestDepartmentModel(t *testing.T) {
	department := model.Department{
		ID:   1,
		Name: "Test Department",
	}

	assert.Equal(t, 1, department.ID)
	assert.Equal(t, "Test Department", department.Name)
}
