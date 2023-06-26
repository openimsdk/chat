package organization

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	"gorm.io/gorm"
	"time"
)

func NewDepartmentMember(db *gorm.DB) *DepartmentMember {
	return &DepartmentMember{
		db: db,
	}
}

type DepartmentMember struct {
	db *gorm.DB
}

func (o *DepartmentMember) Find(ctx context.Context, departmentIDList []string) ([]*table.DepartmentMember, error) {
	if len(departmentIDList) == 0 {
		return []*table.DepartmentMember{}, nil
	}
	var ms []*table.DepartmentMember
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("department_id in (?)", departmentIDList).Find(ms).Error, "")
}

func (o *DepartmentMember) DeleteDepartmentIDList(ctx context.Context, departmentIDList []string) error {
	return utils.Wrap(o.db.WithContext(ctx).Where("department_id in (?)", departmentIDList).Delete(&table.DepartmentMember{}).Error, "")
}

func (o *DepartmentMember) DeleteByUserID(ctx context.Context, userID string) error {
	return utils.Wrap(o.db.WithContext(ctx).Where("user_id = ? ", userID).Delete(&table.DepartmentMember{}).Error, "")
}

func (o *DepartmentMember) Create(ctx context.Context, m *table.DepartmentMember) error {
	m.CreateTime = time.Now()
	m.ChangeTime = time.Now()
	return utils.Wrap(o.db.WithContext(ctx).Create(m).Error, "")
}

func (o *DepartmentMember) Get(ctx context.Context, userID string) ([]*table.DepartmentMember, error) {
	var ms []*table.DepartmentMember
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("user_id = ?", userID).Find(ms).Error, "")
}

func (o *DepartmentMember) DeleteByKey(ctx context.Context, userID, departmentID string) error {
	return utils.Wrap(o.db.WithContext(ctx).Where("user_id = ? AND department_id = ?", userID, departmentID).Delete(&table.DepartmentMember{}).Error, "")
}

func (o *DepartmentMember) Update(ctx context.Context, m *table.DepartmentMember) error {
	m.ChangeTime = time.Now()
	return utils.Wrap(o.db.Where("user_id = ? AND department_id = ?", m.UserID, m.DepartmentID).Updates(m).Error, "")
}
