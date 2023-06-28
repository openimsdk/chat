package chat

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"gorm.io/gorm"
)

func NewIPForbidden(db *gorm.DB) *IPForbidden {
	return &IPForbidden{db: db}
}

type IPForbidden struct {
	db *gorm.DB
}

func (tb *IPForbidden) Restriction(ctx context.Context, ip string, isLogin bool) (bool, error) {
	var m admin.IPForbidden
	switch err := tb.db.WithContext(ctx).Model(&m).Where("ip = ?", ip).First(&m).Error; err {
	case nil:
		if isLogin {
			return m.LimitLogin == true, nil
		}
		return m.LimitRegister == true, nil
	case gorm.ErrRecordNotFound:
		return false, nil
	default:
		return false, utils.Wrap(err, "")
	}
}
