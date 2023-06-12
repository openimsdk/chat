package table

import "context"

// ClientConfig 客户端相关配置项
type ClientConfig struct {
	Key   string `gorm:"column:key;primary_key;type:varchar(255)"`
	Value string `gorm:"column:value;not null;type:text"`
}

func (ClientConfig) TableName() string {
	return "client_config"
}

type ClientConfigInterface interface {
	NewTx(tx any) ClientConfigInterface
	Set(ctx context.Context, config map[string]*string) error
	Get(ctx context.Context) (map[string]string, error)
}
