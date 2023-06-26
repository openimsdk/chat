package organization

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	"gorm.io/gorm"
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
