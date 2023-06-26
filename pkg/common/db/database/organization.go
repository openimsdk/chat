package database

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/chat/pkg/common/db/model/organization"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	"gorm.io/gorm"
)

type OrganizationDatabaseInterface interface {
	GetDepartmentByID(ctx context.Context, departmentID string) (*table.Department, error)
	CreateDepartment(ctx context.Context, department ...*table.Department) error
	UpdateDepartment(ctx context.Context, department *table.Department) error
	GetParent(ctx context.Context, parentID string) ([]table.Department, error)
}

func NewOrganizationDatabase(db *gorm.DB) OrganizationDatabaseInterface {
	return &OrganizationDatabase{
		tx:         tx.NewGorm(db),
		Department: organization.NewDepartment(db),
	}
}

type OrganizationDatabase struct {
	tx         tx.Tx
	Department table.DepartmentInterface
}

func (o *OrganizationDatabase) GetParent(ctx context.Context, parentID string) ([]table.Department, error) {
	return o.Department.GetParent(ctx, parentID)
}

func (o *OrganizationDatabase) UpdateDepartment(ctx context.Context, department *table.Department) error {
	return o.Department.Update(ctx, department)
}

func (o *OrganizationDatabase) GetDepartmentByID(ctx context.Context, departmentID string) (*table.Department, error) {
	return o.Department.FindOne(ctx, departmentID)
}

func (o *OrganizationDatabase) CreateDepartment(ctx context.Context, department ...*table.Department) error {
	return o.Department.Create(ctx, department...)
}
