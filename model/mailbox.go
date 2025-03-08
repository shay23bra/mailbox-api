package model

type Mailbox struct {
	Identifier        string `json:"mailbox_identifier" db:"mailbox_identifier"`
	UserFullName      string `json:"user_full_name" db:"user_full_name"`
	JobTitle          string `json:"job_title" db:"job_title"`
	DepartmentID      int    `json:"department_id" db:"department_id"`
	Department        string `json:"department" db:"department_name"`
	ManagerIdentifier string `json:"manager_mailbox_identifier" db:"manager_mailbox_identifier"`
	OrgDepth          int    `json:"org_depth" db:"org_depth"`
	SubOrgSize        int    `json:"sub_org_size" db:"sub_org_size"`
}

type MailboxFilter struct {
	SearchTerm     string   `form:"search"`
	Department     int      `form:"department"`
	OrgDepthExact  *int     `form:"org_depth_exact"`
	OrgDepthGt     *int     `form:"org_depth_gt"`
	OrgDepthLt     *int     `form:"org_depth_lt"`
	SubOrgSizeMin  *int     `form:"sub_org_size_min"`
	SubOrgSizeMax  *int     `form:"sub_org_size_max"`
	SortBy         []string `form:"sort_by"`
	SortDirections []string `form:"sort_dir"`
	Fields         []string `form:"fields"`
	Page           int      `form:"page"`
	PageSize       int      `form:"page_size"`
}

type MailboxResponse struct {
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
