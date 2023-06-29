package organization

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	"github.com/OpenIMSDK/chat/pkg/proto/organization"
	"gorm.io/gorm"
	"strings"
	"time"
)

func NewOrganizationUser(db *gorm.DB) *OrganizationUser {
	return &OrganizationUser{
		db: db,
	}
}

type OrganizationUser struct {
	db *gorm.DB
}

func (tb *OrganizationUser) Create(ctx context.Context, m *table.OrganizationUser) error {
	m.CreateTime = time.Now()
	m.ChangeTime = time.Now()
	return utils.Wrap(tb.db.WithContext(ctx).Create(m).Error, "")
}

func (tb *OrganizationUser) Update(ctx context.Context, m *table.OrganizationUser) error {
	m.ChangeTime = time.Now()
	return utils.Wrap(tb.db.WithContext(ctx).Where("user_id = ?", m.UserID).Updates(&m).Error, "")
}

func (tb *OrganizationUser) Delete(ctx context.Context, userID string) error {
	return utils.Wrap(tb.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&table.OrganizationUser{}).Error, "")
}

func (tb *OrganizationUser) Get(ctx context.Context, userID string) (*table.OrganizationUser, error) {
	var m table.OrganizationUser
	return &m, utils.Wrap(tb.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error, "")
}

func (tb *OrganizationUser) SearchPage(ctx context.Context, positionList, userIDList []string, text string, sort []*organization.GetSearchUserListSort, pageNumber uint32, showNumber uint32) (uint32, []*table.OrganizationUser, error) {
	db := tb.db
	if len(positionList) > 0 {
		db = db.Where("position in (?)", positionList)
	}
	if len(userIDList) > 0 {
		db = db.Where("user_id in (?)", userIDList)
	}
	if text != "" {
		fields := []string{"user_id", "nickname", "english_name", "mobile", "telephone", "email"}
		//values := make([]interface{}, len(fields))
		for i := 0; i < len(fields); i++ {
			fields[i] = fmt.Sprintf("`%s` like '%%%s%%'", fields[i], text)
			//values[i] = text
		}
		db = db.Where(strings.Join(fields, " OR "))
		//db = db.Where(strings.Join(fields, " OR "), values...)
	} else {
		db = db.Where("1=1")
	}
	var count int64
	if err := db.Model(&table.OrganizationUser{}).Count(&count).Error; err != nil {
		return 0, nil, utils.Wrap(err, "")
	}
	if showNumber > 0 {
		db = db.Offset(int(pageNumber * showNumber)).Limit(int(showNumber))
	}
	if len(sort) > 0 { // DESC: 降序   ASC: 升序
		arr := make([]string, len(sort))
		for i, s := range sort {
			if s.Rule == "" {
				arr[i] = fmt.Sprintf("`%s`", s.Field)
			} else {
				arr[i] = fmt.Sprintf("`%s` %s", s.Field, s.Rule)
			}
		}
		db = db.Order(strings.Join(arr, ","))
	}
	db = db.Order("`order` ASC, `create_time` ASC")
	var ms []*table.OrganizationUser
	return uint32(count), ms, utils.Wrap(db.Find(ms).Error, "")
}

func (tb *OrganizationUser) GetNoDepartmentUserIDList(ctx context.Context) ([]string, error) {
	type Temp struct {
		UserID string
	}
	var ts []Temp
	err := tb.db.WithContext(ctx).Raw("SELECT ou.user_id,dm.user_id AS dm_user_id,dm.department_id AS dm_department_id FROM organization_users ou LEFT JOIN department_members dm ON dm.user_id=ou.user_id HAVING dm_user_id IS NULL OR dm_department_id=''").Scan(&ts).Error
	if err != nil {
		return nil, utils.Wrap(err, "")
	}

	userIDList := make([]string, len(ts))
	for i, t := range ts {
		userIDList[i] = t.UserID
	}
	return userIDList, nil
}

func (tb *OrganizationUser) GetList(ctx context.Context, userIDList []string) ([]*table.OrganizationUser, error) {
	if len(userIDList) == 0 {
		return []*table.OrganizationUser{}, nil
	}
	var ms []*table.OrganizationUser
	return ms, utils.Wrap(tb.db.Where("user_id in (?)", userIDList).Find(ms).Error, "")
}

func (tb *OrganizationUser) Search(ctx context.Context, positionList, userIDList []string, text string, sort []*organization.GetSearchUserListSort) ([]*table.OrganizationUser, error) {
	db := tb.db
	if len(positionList) > 0 {
		db = db.WithContext(ctx).Where("position in (?)", positionList)
	}
	if len(userIDList) > 0 {
		db = db.WithContext(ctx).Where("user_id in (?)", userIDList)
	}
	if text != "" {
		fields := []string{"user_id", "nickname", "english_name", "mobile", "telephone", "email"}
		//values := make([]interface{}, len(fields))
		for i := 0; i < len(fields); i++ {
			fields[i] = fmt.Sprintf("`%s` like '%%%s%%'", fields[i], text)
			//values[i] = text
		}
		db = db.WithContext(ctx).Where(strings.Join(fields, " OR "))
		//db = db.Where(strings.Join(fields, " OR "), values...)
	} else {
		db = db.WithContext(ctx).Where("1=1")
	}

	if len(sort) > 0 { // DESC: 降序   ASC: 升序
		arr := make([]string, len(sort))
		for i, s := range sort {
			if s.Rule == "" {
				arr[i] = fmt.Sprintf("`%s`", s.Field)
			} else {
				arr[i] = fmt.Sprintf("`%s` %s", s.Field, s.Rule)
			}
		}
		db = db.WithContext(ctx).Order(strings.Join(arr, ","))
	}
	db = db.WithContext(ctx).Order("`order` ASC, `create_time` ASC")
	var ms []*table.OrganizationUser
	return ms, utils.Wrap(db.Find(ms).Error, "")
}

func (tb *OrganizationUser) GetPage(ctx context.Context, pageNumber, showNumber int) (int64, []*table.OrganizationUser, error) {
	var total int64
	err := tb.db.Model(&table.OrganizationUser{}).Count(&total).Error
	if err != nil {
		return 0, nil, utils.Wrap(err, "")
	}

	var users []*table.OrganizationUser
	err = tb.db.WithContext(ctx).Model(&table.OrganizationUser{}).Order("create_time DESC").Offset(pageNumber * showNumber).Limit(showNumber).Find(users).Error
	return total, users, utils.Wrap(err, "")
}

func (tb *OrganizationUser) SearchV2(ctx context.Context, keyword string, orUserIDList []string, pageNumber, showNumber int) (int64, []*table.OrganizationUser, error) {
	db := tb.db.Model(&table.OrganizationUser{})
	if keyword != "" {
		vague := "%" + keyword + "%"
		db = db.WithContext(ctx).Where("user_id in (?) OR mobile = ? OR telephone = ? OR email = ? OR nickname like ? OR english_name like ?", append(orUserIDList, keyword), keyword, keyword, keyword, vague, vague)
	}
	var count int64
	if err := db.WithContext(ctx).Count(&count).Error; err != nil {
		return 0, nil, utils.Wrap(err, "")
	}
	db = db.WithContext(ctx).Order("`order` ASC, `create_time` ASC").Offset(int(pageNumber) * int(showNumber)).Limit(int(showNumber))
	var ms []*table.OrganizationUser
	return count, ms, utils.Wrap(db.WithContext(ctx).Find(&ms).Error, "")
}
