package dto

import "mailbox-api/model"

func NewMailboxResponse(data interface{}, pagination *model.Pagination) *model.MailboxResponse {
	return &model.MailboxResponse{
		Data:       data,
		Pagination: pagination,
	}
}

func FilterMailboxFields(mailboxes []model.Mailbox, fields []string) []map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(mailboxes))

	for i, mailbox := range mailboxes {
		m := make(map[string]interface{})

		for _, field := range fields {
			switch field {
			case "mailbox_identifier":
				m["mailbox_identifier"] = mailbox.Identifier
			case "user_full_name":
				m["user_full_name"] = mailbox.UserFullName
			case "job_title":
				m["job_title"] = mailbox.JobTitle
			case "department_id":
				m["department_id"] = mailbox.DepartmentID
			case "department":
				m["department"] = mailbox.Department
			case "org_depth":
				m["org_depth"] = mailbox.OrgDepth
			case "sub_org_size":
				m["sub_org_size"] = mailbox.SubOrgSize
			case "manager_mailbox_identifier":
				m["manager_mailbox_identifier"] = mailbox.ManagerIdentifier
			}
		}

		result[i] = m
	}

	return result
}
