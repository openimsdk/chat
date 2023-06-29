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

func (o *DepartmentMember) FindByDepartmentID(ctx context.Context, departmentIDList []string) ([]*table.DepartmentMember, error) {
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
	return utils.Wrap(o.db.WithContext(ctx).Where("user_id = ? AND department_id = ?", m.UserID, m.DepartmentID).Updates(m).Error, "")
}

func (o *DepartmentMember) FindByUserID(ctx context.Context, userIDList []string) ([]*table.DepartmentMember, error) {
	if len(userIDList) == 0 {
		return []*table.DepartmentMember{}, nil
	}
	var ms []*table.DepartmentMember
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("user_id in (?)", userIDList).Find(ms).Error, "")
}

func (o *DepartmentMember) GetUserListInDepartment(ctx context.Context, departmentID string, userIDList []string) ([]*table.DepartmentMember, error) {
	if len(userIDList) == 0 {
		return []*table.DepartmentMember{}, nil
	}
	var ms []*table.DepartmentMember
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("department_id = ? AND user_id in (?)", departmentID, userIDList).Find(&ms).Error, "")
}

func (o *DepartmentMember) GetByDepartmentID(ctx context.Context, departmentID string) ([]*table.DepartmentMember, error) {
	var ms []*table.DepartmentMember
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("department_id = ?", departmentID).Find(ms).Error, "")
}

func (o *DepartmentMember) CreateList(ctx context.Context, ms []*table.DepartmentMember) error {
	now := time.Now()
	for i := 0; i < len(ms); i++ {
		ms[i].CreateTime = now
		ms[i].ChangeTime = now
	}
	return utils.Wrap(o.db.WithContext(ctx).Create(&ms).Error, "")
}

func (o *DepartmentMember) GetByKey(ctx context.Context, userID, departmentID string) (*table.DepartmentMember, error) {
	var ms *table.DepartmentMember
	return ms, utils.Wrap(o.db.WithContext(ctx).Where("user_id = ? and department_id = ?", userID, departmentID).First(ms).Error, "")
}
