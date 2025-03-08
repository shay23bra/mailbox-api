package repository

import (
	"context"
	"errors"
	"fmt"

	"mailbox-api/db"
	"mailbox-api/model"

	"github.com/jackc/pgx/v4"
)

type DepartmentRepository interface {
	GetDepartments(ctx context.Context) ([]model.Department, error)
	GetDepartmentByID(ctx context.Context, id int) (*model.Department, error)
	CreateDepartment(ctx context.Context, department model.Department) error
}

type departmentRepository struct {
	db *db.DB
}

func NewDepartmentRepository(db *db.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) GetDepartments(ctx context.Context) ([]model.Department, error) {
	query := `
	SELECT 
		department_id, 
		department_name
	FROM 
		departments
	ORDER BY 
		department_name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query departments: %w", err)
	}
	defer rows.Close()

	departments := []model.Department{}
	for rows.Next() {
		var department model.Department
		err := rows.Scan(
			&department.ID,
			&department.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan department: %w", err)
		}
		departments = append(departments, department)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over departments: %w", err)
	}

	return departments, nil
}

func (r *departmentRepository) GetDepartmentByID(ctx context.Context, id int) (*model.Department, error) {
	query := `
	SELECT 
		department_id, 
		department_name
	FROM 
		departments
	WHERE 
		department_id = $1`

	var department model.Department
	err := r.db.QueryRow(ctx, query, id).Scan(
		&department.ID,
		&department.Name,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}

	return &department, nil
}

func (r *departmentRepository) CreateDepartment(ctx context.Context, department model.Department) error {
	query := `
	INSERT INTO departments (
		department_id, 
		department_name
	) VALUES ($1, $2)`

	_, err := r.db.Exec(ctx, query,
		department.ID,
		department.Name,
	)

	if err != nil {
		return fmt.Errorf("failed to create department: %w", err)
	}

	return nil
}
