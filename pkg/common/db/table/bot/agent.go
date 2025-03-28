package bot

import (
	"context"
	"time"

	"github.com/openimsdk/tools/db/pagination"
)

type Agent struct {
	UserID     string    `bson:"user_id"`
	NickName   string    `bson:"nick_name"`
	FaceURL    string    `bson:"face_url"`
	Key        string    `bson:"key"`
	Url        string    `bson:"url"`
	Identity   string    `bson:"identity"`
	Model      string    `bson:"model"`
	Prompts    string    `bson:"prompts"`
	CreateTime time.Time `bson:"create_time"`
}

func (Agent) TableName() string {
	return "agent"
}

type AgentInterface interface {
	Create(ctx context.Context, elems ...*Agent) error
	Take(ctx context.Context, userID string) (*Agent, error)
	Find(ctx context.Context, userIDs []string) ([]*Agent, error)
	Update(ctx context.Context, userID string, data map[string]any) error
	Delete(ctx context.Context, userIDs []string) error
	Page(ctx context.Context, userIDs []string, pagination pagination.Pagination) (int64, []*Agent, error)
}
