package admin

import (
	"context"
	admindb "github.com/openimsdk/chat/pkg/common/db/table/admin"
	"github.com/openimsdk/chat/pkg/common/mctx"
	"github.com/openimsdk/chat/pkg/protocol/admin"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func IsNotFound(err error) bool {
	switch errs.Unwrap(err) {
	case redis.Nil, mongo.ErrNoDocuments:
		return true
	default:
		return false
	}
}

func (o *adminServer) db2pbApplication(val *admindb.Application) *admin.ApplicationVersion {
	return &admin.ApplicationVersion{
		Id:         val.ID.Hex(),
		Platform:   val.Platform,
		Version:    val.Version,
		Url:        val.Url,
		Text:       val.Text,
		Force:      val.Force,
		Latest:     val.Latest,
		Hot:        val.Hot,
		CreateTime: val.CreateTime.UnixMilli(),
	}
}

func (o *adminServer) LatestApplicationVersion(ctx context.Context, req *admin.LatestApplicationVersionReq) (*admin.LatestApplicationVersionResp, error) {
	res, err := o.Database.LatestVersion(ctx, req.Platform)
	if err == nil {
		return &admin.LatestApplicationVersionResp{Version: o.db2pbApplication(res)}, nil
	} else if IsNotFound(err) {
		return &admin.LatestApplicationVersionResp{}, nil
	} else {
		return nil, err
	}
}

func (o *adminServer) AddApplicationVersion(ctx context.Context, req *admin.AddApplicationVersionReq) (*admin.AddApplicationVersionResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	val := &admindb.Application{
		ID:         primitive.NewObjectID(),
		Platform:   req.Platform,
		Version:    req.Version,
		Url:        req.Url,
		Text:       req.Text,
		Force:      req.Force,
		Latest:     req.Latest,
		Hot:        req.Hot,
		CreateTime: time.Now(),
	}
	if err := o.Database.AddVersion(ctx, val); err != nil {
		return nil, err
	}
	return &admin.AddApplicationVersionResp{}, nil
}

func (o *adminServer) UpdateApplicationVersion(ctx context.Context, req *admin.UpdateApplicationVersionReq) (*admin.UpdateApplicationVersionResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	oid, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, errs.ErrArgs.WrapMsg("invalid id " + err.Error())
	}
	update := make(map[string]any)
	putUpdate(update, "platform", req.Platform)
	putUpdate(update, "version", req.Version)
	putUpdate(update, "url", req.Url)
	putUpdate(update, "text", req.Text)
	putUpdate(update, "force", req.Force)
	putUpdate(update, "latest", req.Latest)
	putUpdate(update, "hot", req.Hot)
	if err := o.Database.UpdateVersion(ctx, oid, update); err != nil {
		return nil, err
	}
	return &admin.UpdateApplicationVersionResp{}, nil
}

func (o *adminServer) DeleteApplicationVersion(ctx context.Context, req *admin.DeleteApplicationVersionReq) (*admin.DeleteApplicationVersionResp, error) {
	if _, err := mctx.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	ids := make([]primitive.ObjectID, 0, len(req.Id))
	for _, id := range req.Id {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errs.ErrArgs.WrapMsg("invalid id " + err.Error())
		}
		ids = append(ids, oid)
	}
	if err := o.Database.DeleteVersion(ctx, ids); err != nil {
		return nil, err
	}
	return &admin.DeleteApplicationVersionResp{}, nil
}

func (o *adminServer) PageApplicationVersion(ctx context.Context, req *admin.PageApplicationVersionReq) (*admin.PageApplicationVersionResp, error) {
	total, res, err := o.Database.PageVersion(ctx, req.Platform, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &admin.PageApplicationVersionResp{
		Total:    total,
		Versions: datautil.Slice(res, o.db2pbApplication),
	}, nil
}

func putUpdate[T any](update map[string]any, name string, val interface{ GetValuePtr() *T }) {
	ptrVal := val.GetValuePtr()
	if ptrVal == nil {
		return
	}
	update[name] = *ptrVal
}
