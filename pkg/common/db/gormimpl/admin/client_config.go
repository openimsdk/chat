package admin

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	"gorm.io/gorm"
)

func NewClientConfig(db *gorm.DB) admin.ClientConfigInterface {
	return &ClientConfig{db: db}
}

type ClientConfig struct {
	db *gorm.DB
}

func (o *ClientConfig) NewTx(tx any) admin.ClientConfigInterface {
	return &ClientConfig{db: tx.(*gorm.DB)}
}

func (o *ClientConfig) Set(ctx context.Context, config map[string]*string) error {
	err := o.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for key, value := range config {
			if value == nil {
				if err := tx.Where("`key` = ?", key).Delete(&admin.ClientConfig{}).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Where("`key` = ?", key).Take(&admin.ClientConfig{}).Error; err == nil {
					if err := tx.Where("`key` = ?", key).Model(&admin.ClientConfig{}).Update("value", *value).Error; err != nil {
						return err
					}
				} else if err == gorm.ErrRecordNotFound {
					if err := tx.Create(&admin.ClientConfig{Key: key, Value: *value}).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}
		}
		return nil
	})
	return errs.Wrap(err)
}

func (o *ClientConfig) Get(ctx context.Context) (map[string]string, error) {
	var cs []*admin.ClientConfig
	if err := o.db.WithContext(ctx).Find(&cs).Error; err != nil {
		return nil, err
	}
	cm := make(map[string]string)
	for _, config := range cs {
		cm[config.Key] = config.Value
	}
	return cm, nil
}
