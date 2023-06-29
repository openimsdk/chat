package organization

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	"gorm.io/gorm"
	"time"
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
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("department_id in (?)", departmentIdList).Order("`order` ASC, `create_time` ASC").Find(&ms).Error, "")
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

func (o *Department) GetMaxOrder(ctx context.Context, parentID string) (int32, error) {
	// SELECT IFNULL(MAX(`order`), 0) FROM department_members WHERE `user_id` = "22286361621"
	var order int32
	return order, utils.Wrap(o.db.WithContext(ctx).Model(&table.Department{}).Select("IFNULL(MAX(`order`), 0)").Where("parent_id = ?", parentID).Scan(&order).Error, "")
}

func (o *Department) UpdateOrderIncrement(ctx context.Context, parentID string, startOrder int32) error {
	return utils.Wrap(o.db.WithContext(ctx).Model(&table.Department{}).Where("parent_id = ? AND `order` >= ?", parentID, startOrder).Update("`order`", gorm.Expr("`order` + ?", 1)).Error, "")
}
func (o *Department) UpdateParentIDOrder(ctx context.Context, departmentID, parentID string, order int32) error {
	return utils.Wrap(o.db.WithContext(ctx).Model(&table.Department{}).Where("department_id = ?", departmentID).Updates(map[string]interface{}{
		"parent_id": parentID,
		"`order`":   order,
	}).Error, "")
}

func (o *Department) GetByName(ctx context.Context, name, parentID string) (*table.Department, error) {
	var m table.Department
	return &m, utils.Wrap(o.db.WithContext(ctx).Where("name = ? AND parent_id = ?", name, parentID).First(&m).Error, "")
}
