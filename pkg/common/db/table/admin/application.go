package admin

import (
	"context"
	"github.com/openimsdk/tools/db/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Application struct {
	ID         primitive.ObjectID `bson:"_id"`
	Platform   string             `bson:"platform"`
	Hot        bool               `bson:"hot"`
	Version    string             `bson:"version"`
	Url        string             `bson:"url"`
	Text       string             `bson:"text"`
	Force      bool               `bson:"force"`
	Latest     bool               `bson:"latest"`
	CreateTime time.Time          `bson:"create_time"`
}

type ApplicationInterface interface {
	LatestVersion(ctx context.Context, platform string) (*Application, error)
	AddVersion(ctx context.Context, val *Application) error
	UpdateVersion(ctx context.Context, id primitive.ObjectID, update map[string]any) error
	DeleteVersion(ctx context.Context, id []primitive.ObjectID) error
	PageVersion(ctx context.Context, platforms []string, page pagination.Pagination) (int64, []*Application, error)
	FindPlatform(ctx context.Context, id []primitive.ObjectID) ([]string, error)
}
