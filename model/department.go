package model

type Department struct {
	ID   int    `json:"department_id" db:"department_id"`
	Name string `json:"department_name" db:"department_name"`
}
