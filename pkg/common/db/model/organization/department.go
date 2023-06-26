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

func (d *Department) GetParent(ctx context.Context, parentID string) ([]table.Department, error) {
	var ms []table.Department
	return ms, utils.Wrap(d.db.Where("parent_id = ?", parentID).Order("`order` ASC, `create_time` ASC").Find(&ms).Error, "")
}

func (d *Department) Update(ctx context.Context, department *table.Department) error {
	department.ChangeTime = time.Now()
	return utils.Wrap(d.db.Updates(department).Error, "")
}

func (d *Department) Create(ctx context.Context, departments ...*table.Department) error {
	return errs.Wrap(d.db.WithContext(ctx).Create(departments).Error)
}

func (d *Department) FindOne(ctx context.Context, departmentID string) (*table.Department, error) {
	var m table.Department
	return &m, utils.Wrap(d.db.Where("department_id = ?", departmentID).First(&m).Error, "")
}
