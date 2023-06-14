package dbutil

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"gorm.io/gorm"
)

func IsGormNotFound(err error) bool {
	return errs.Unwrap(err) == gorm.ErrRecordNotFound
}
