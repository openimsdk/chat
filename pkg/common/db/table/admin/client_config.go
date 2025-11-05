package admin

import "context"

// ClientConfig config
type ClientConfig struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

func (ClientConfig) TableName() string {
	return "client_config"
}

type ClientConfigInterface interface {
	Set(ctx context.Context, config map[string]string) error
	Get(ctx context.Context) (map[string]string, error)
	Del(ctx context.Context, keys []string) error
}
