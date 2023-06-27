package database

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/chat/pkg/common/db/model/organization"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	rpc "github.com/OpenIMSDK/chat/pkg/proto/organization"
	"gorm.io/gorm"
)

type OrganizationDatabaseInterface interface {
	//department
	GetDepartmentByID(ctx context.Context, departmentID string) (*table.Department, error)
	CreateDepartment(ctx context.Context, department ...*table.Department) error
	UpdateDepartment(ctx context.Context, department *table.Department) error
	GetParent(ctx context.Context, parentID string) ([]*table.Department, error)
	GetDepartment(ctx context.Context, departmentID string) (*table.Department, error)
	GetList(ctx context.Context, departmentIDList []string) ([]*table.Department, error)
	DeleteDepartment(ctx context.Context, departmentIDList []string) error
	UpdateParentID(ctx context.Context, oldParentID, newParentID string) error
	//departmentMember
	FindDepartmentMember(ctx context.Context, departmentIDList []string) ([]*table.DepartmentMember, error)
	GetDepartmentMemberByUserID(ctx context.Context, userID string) ([]*table.DepartmentMember, error)
	CreateDepartmentMember(ctx context.Context, DepartmentMember *table.DepartmentMember) error
	DeleteDepartmentIDList(ctx context.Context, departmentIDList []string) error
	DeleteDepartmentMemberByUserID(ctx context.Context, userID string) error
	DeleteDepartmentMemberByKey(ctx context.Context, userID string, departmentID string) error
	UpdateDepartmentMember(ctx context.Context, DepartmentMember *table.DepartmentMember) error
	FindDepartmentMemberByUserID(ctx context.Context, userIDList []string) ([]*table.DepartmentMember, error)
	GetUserListInDepartment(ctx context.Context, departmentID string, userIDList []string) ([]*table.DepartmentMember, error)
	GetDepartmentMemberByDepartmentID(ctx context.Context, departmentID string) ([]*table.DepartmentMember, error)
	//organizationUser
	CreateOrganizationUser(ctx context.Context, OrganizationUser *table.OrganizationUser) error
	UpdateOrganizationUser(ctx context.Context, OrganizationUser *table.OrganizationUser) error
	DeleteOrganizationUser(ctx context.Context, userID string) error
	GetOrganizationUser(ctx context.Context, userID string) (*table.OrganizationUser, error)
	SearchPage(ctx context.Context, positionList, userIDList []string, text string, sort []*rpc.GetSearchUserListSort, pageNumber uint32, showNumber uint32) (uint32, []*table.OrganizationUser, error)
	GetNoDepartmentUserIDList(ctx context.Context) ([]string, error)
	GetOrganizationUserList(ctx context.Context, userIDList []string) ([]*table.OrganizationUser, error)
	SearchOrganizationUser(ctx context.Context, positionList, userIDList []string, text string, sort []*rpc.GetSearchUserListSort) ([]*table.OrganizationUser, error)
	//organizaiton
	SetOrganization(ctx context.Context, Organization *table.Organization) error
	GetOrganization(ctx context.Context) (*table.Organization, error)
}

func NewOrganizationDatabase(db *gorm.DB) OrganizationDatabaseInterface {
	return &OrganizationDatabase{
		tx:               tx.NewGorm(db),
		Department:       organization.NewDepartment(db),
		DepartmentMember: organization.NewDepartmentMember(db),
		Organization:     organization.NewOrganization(db),
		OrganizationUser: organization.NewOrganizationUser(db),
	}
}

type OrganizationDatabase struct {
	tx               tx.Tx
	Department       table.DepartmentInterface
	DepartmentMember table.DepartmentMemberInterface
	OrganizationUser table.OrganizationUserInterface
	Organization     table.OrganizationInterface
}

func (o *OrganizationDatabase) SearchOrganizationUser(ctx context.Context, positionList, userIDList []string, text string, sort []*rpc.GetSearchUserListSort) ([]*table.OrganizationUser, error) {
	return o.OrganizationUser.Search(ctx, positionList, userIDList, text, sort)
}

func (o *OrganizationDatabase) GetDepartmentMemberByDepartmentID(ctx context.Context, departmentID string) ([]*table.DepartmentMember, error) {
	return o.DepartmentMember.GetByDepartmentID(ctx, departmentID)
}

func (o *OrganizationDatabase) GetUserListInDepartment(ctx context.Context, departmentID string, userIDList []string) ([]*table.DepartmentMember, error) {
	return o.DepartmentMember.GetUserListInDepartment(ctx, departmentID, userIDList)
}

func (o *OrganizationDatabase) GetOrganizationUserList(ctx context.Context, userIDList []string) ([]*table.OrganizationUser, error) {
	return o.OrganizationUser.GetList(ctx, userIDList)
}

func (o *OrganizationDatabase) GetNoDepartmentUserIDList(ctx context.Context) ([]string, error) {
	return o.OrganizationUser.GetNoDepartmentUserIDList(ctx)
}

func (o *OrganizationDatabase) GetOrganization(ctx context.Context) (*table.Organization, error) {
	return o.Organization.Get(ctx)
}

func (o *OrganizationDatabase) SetOrganization(ctx context.Context, Organization *table.Organization) error {
	return o.Organization.Set(ctx, Organization)
}

func (o *OrganizationDatabase) FindDepartmentMemberByUserID(ctx context.Context, userIDList []string) ([]*table.DepartmentMember, error) {
	return o.DepartmentMember.FindByUserID(ctx, userIDList)
}

func (o *OrganizationDatabase) SearchPage(ctx context.Context, positionList, userIDList []string, text string, sort []*rpc.GetSearchUserListSort, pageNumber uint32, showNumber uint32) (uint32, []*table.OrganizationUser, error) {
	return o.OrganizationUser.SearchPage(ctx, positionList, userIDList, text, sort, pageNumber, showNumber)
}

func (o *OrganizationDatabase) GetDepartmentMemberByUserID(ctx context.Context, userID string) ([]*table.DepartmentMember, error) {
	return o.DepartmentMember.Get(ctx, userID)
}

func (o *OrganizationDatabase) CreateDepartmentMember(ctx context.Context, DepartmentMember *table.DepartmentMember) error {
	return o.DepartmentMember.Create(ctx, DepartmentMember)
}

func (o *OrganizationDatabase) DeleteDepartmentMemberByUserID(ctx context.Context, userID string) error {
	return o.DepartmentMember.DeleteByUserID(ctx, userID)
}

func (o *OrganizationDatabase) DeleteDepartmentMemberByKey(ctx context.Context, userID string, departmentID string) error {
	return o.DepartmentMember.DeleteByKey(ctx, userID, departmentID)
}

func (o *OrganizationDatabase) UpdateDepartmentMember(ctx context.Context, DepartmentMember *table.DepartmentMember) error {
	return o.DepartmentMember.Update(ctx, DepartmentMember)
}

func (o *OrganizationDatabase) GetOrganizationUser(ctx context.Context, userID string) (*table.OrganizationUser, error) {
	return o.OrganizationUser.Get(ctx, userID)
}

func (o *OrganizationDatabase) DeleteOrganizationUser(ctx context.Context, userID string) error {
	return o.OrganizationUser.Delete(ctx, userID)
}

func (o *OrganizationDatabase) UpdateOrganizationUser(ctx context.Context, OrganizationUser *table.OrganizationUser) error {
	return o.OrganizationUser.Update(ctx, OrganizationUser)
}

func (o *OrganizationDatabase) CreateOrganizationUser(ctx context.Context, OrganizationUser *table.OrganizationUser) error {
	return o.OrganizationUser.Create(ctx, OrganizationUser)
}

func (o *OrganizationDatabase) DeleteDepartmentIDList(ctx context.Context, departmentIDList []string) error {
	return o.DepartmentMember.DeleteDepartmentIDList(ctx, departmentIDList)
}

func (o *OrganizationDatabase) DeleteDepartment(ctx context.Context, departmentIDList []string) error {
	return o.Department.Delete(ctx, departmentIDList)
}

func (o *OrganizationDatabase) UpdateParentID(ctx context.Context, oldParentID, newParentID string) error {
	return o.Department.UpdateParentID(ctx, oldParentID, newParentID)
}

func (o *OrganizationDatabase) GetList(ctx context.Context, departmentIDList []string) ([]*table.Department, error) {
	return o.Department.GetList(ctx, departmentIDList)
}

func (o *OrganizationDatabase) FindDepartmentMember(ctx context.Context, departmentIDList []string) ([]*table.DepartmentMember, error) {
	return o.DepartmentMember.FindByDepartmentID(ctx, departmentIDList)
}

func (o *OrganizationDatabase) GetParent(ctx context.Context, parentID string) ([]*table.Department, error) {
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

func (o *OrganizationDatabase) GetDepartment(ctx context.Context, departmentID string) (*table.Department, error) {
	return o.Department.GetDepartment(ctx, departmentID)
}
