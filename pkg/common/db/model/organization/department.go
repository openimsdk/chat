package organization

import (
	"context"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	"gorm.io/gorm"
)

func NewDepartment(db *gorm.DB) *Department {
	return &Department{
		db: db,
	}
}

type Department struct {
	db *gorm.DB
}

func (d *Department) GetParent(ctx context.Context, parentID string) ([]*table.Department, error) {
	var ms []*table.Department
	return ms, utils.Wrap(d.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("`order` ASC, `create_time` ASC").Find(ms).Error, "")
}

func (d *Department) Update(ctx context.Context, department *table.Department) error {
	department.ChangeTime = time.Now()
	return utils.Wrap(d.db.WithContext(ctx).Updates(department).Error, "")
}

func (d *Department) Create(ctx context.Context, departments ...*table.Department) error {
	return errs.Wrap(d.db.WithContext(ctx).Create(departments).Error)
}

func (d *Department) FindOne(ctx context.Context, departmentID string) (*table.Department, error) {
	var m table.Department
	return &m, utils.Wrap(d.db.WithContext(ctx).Where("department_id = ?", departmentID).First(&m).Error, "")
}

func (o *Department) GetList(ctx context.Context, departmentIdList []string) ([]*table.Department, error) {
	if len(departmentIdList) == 0 {
		return []*table.Department{}, nil
	}
	var ms []*table.Department
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("department_id in (?)", departmentIdList).Order("`order` ASC, `create_time` ASC").Find(ms).Error, "")
}

func (o *Department) UpdateParentID(ctx context.Context, oldParentID, newParentID string) error {
	return utils.Wrap(o.db.WithContext(ctx).Model(&table.Department{}).Where("parent_id = ?", oldParentID).Update("parent_id", newParentID).Error, "")
}

func (o *Department) Delete(ctx context.Context, departmentIDList []string) error {
	if len(departmentIDList) == 0 {
		return nil
	}
	return utils.Wrap(o.db.WithContext(ctx).Where("department_id in (?)", departmentIDList).Delete(&table.Department{}).Error, "")
}

func (o *Department) GetDepartment(ctx context.Context, departmentId string) (*table.Department, error) {
	var m table.Department
	return &m, utils.Wrap(o.db.WithContext(ctx).Where("department_id = ?", departmentId).First(&m).Error, "")
}

func (o *Department) GetByName(ctx context.Context, name, parentID string) (*table.Department, error) {
	var m table.Department
	return &m, utils.Wrap(o.db.Where("name = ? AND parent_id = ?", name, parentID).First(&m).Error, "")
}
