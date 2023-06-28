package chat

import "context"

type IPForbiddenInterface interface {
	Restriction(ctx context.Context, ip string, isLogin bool) (bool, error)
}
