// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package database

import (
	"context"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/db/dbutil"
	"github.com/OpenIMSDK/tools/pagination"
	"go.mongodb.org/mongo-driver/mongo"

	constant2 "github.com/OpenIMSDK/chat/pkg/common/constant"
	admin2 "github.com/OpenIMSDK/chat/pkg/common/db/model/admin"
	"github.com/OpenIMSDK/chat/pkg/common/db/model/chat"
	"github.com/OpenIMSDK/chat/pkg/common/db/table/admin"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"github.com/OpenIMSDK/tools/tx"
)

type ChatDatabaseInterface interface {
	IsNotFound(err error) bool
	GetUser(ctx context.Context, userID string) (account *table.Account, err error)
	UpdateUseInfo(ctx context.Context, userID string, attribute map[string]any) (err error)
	FindAttribute(ctx context.Context, userIDs []string) ([]*table.Attribute, error)
	FindAttributeByAccount(ctx context.Context, accounts []string) ([]*table.Attribute, error)
	TakeAttributeByPhone(ctx context.Context, areaCode string, phoneNumber string) (*table.Attribute, error)
	TakeAttributeByEmail(ctx context.Context, Email string) (*table.Attribute, error)
	TakeAttributeByAccount(ctx context.Context, account string) (*table.Attribute, error)
	TakeAttributeByUserID(ctx context.Context, userID string) (*table.Attribute, error)
	Search(ctx context.Context, normalUser int32, keyword string, gender int32, pagination pagination.Pagination) (int64, []*table.Attribute, error)
	SearchUser(ctx context.Context, keyword string, userIDs []string, genders []int32, pagination pagination.Pagination) (int64, []*table.Attribute, error)
	CountVerifyCodeRange(ctx context.Context, account string, start time.Time, end time.Time) (int64, error)
	AddVerifyCode(ctx context.Context, verifyCode *table.VerifyCode, fn func() error) error
	UpdateVerifyCodeIncrCount(ctx context.Context, id string) error
	TakeLastVerifyCode(ctx context.Context, account string) (*table.VerifyCode, error)
	DelVerifyCode(ctx context.Context, id string) error
	RegisterUser(ctx context.Context, register *table.Register, account *table.Account, attribute *table.Attribute) error
	GetAccount(ctx context.Context, userID string) (*table.Account, error)
	GetAttribute(ctx context.Context, userID string) (*table.Attribute, error)
	GetAttributeByAccount(ctx context.Context, account string) (*table.Attribute, error)
	GetAttributeByPhone(ctx context.Context, areaCode string, phoneNumber string) (*table.Attribute, error)
	GetAttributeByEmail(ctx context.Context, email string) (*table.Attribute, error)
	LoginRecord(ctx context.Context, record *table.UserLoginRecord, verifyCodeID *string) error
	UpdatePassword(ctx context.Context, userID string, password string) error
	UpdatePasswordAndDeleteVerifyCode(ctx context.Context, userID string, password string, codeID string) error
	NewUserCountTotal(ctx context.Context, before *time.Time) (int64, error)
	UserLoginCountTotal(ctx context.Context, before *time.Time) (int64, error)
	UserLoginCountRangeEverydayTotal(ctx context.Context, start *time.Time, end *time.Time) (map[string]int64, int64, error)
	UploadLogs(ctx context.Context, logs []*table.Log) error
	DeleteLogs(ctx context.Context, logID []string, userID string) error
	SearchLogs(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*table.Log, error)
	GetLogs(ctx context.Context, LogIDs []string, userID string) ([]*table.Log, error)
}

func NewChatDatabase(db *mongo.Database) (ChatDatabaseInterface, error) {
	register, err := chat.NewRegister(db)
	if err != nil {
		return nil, err
	}
	account, err := chat.NewAccount(db)
	if err != nil {
		return nil, err
	}
	attribute, err := chat.NewAttribute(db)
	if err != nil {
		return nil, err
	}
	userLoginRecord, err := chat.NewUserLoginRecord(db)
	if err != nil {
		return nil, err
	}
	verifyCode, err := chat.NewVerifyCode(db)
	if err != nil {
		return nil, err
	}
	forbiddenAccount, err := admin2.NewForbiddenAccount(db)
	if err != nil {
		return nil, err
	}
	log, err := chat.NewLogs(db)
	if err != nil {
		return nil, err
	}
	return &ChatDatabase{
		tx:               tx.NewMongo(db.Client()),
		register:         register,
		account:          account,
		attribute:        attribute,
		userLoginRecord:  userLoginRecord,
		verifyCode:       verifyCode,
		forbiddenAccount: forbiddenAccount,
		log:              log,
	}, nil
}

type ChatDatabase struct {
	tx               tx.CtxTx
	register         table.RegisterInterface
	account          table.AccountInterface
	attribute        table.AttributeInterface
	userLoginRecord  table.UserLoginRecordInterface
	verifyCode       table.VerifyCodeInterface
	forbiddenAccount admin.ForbiddenAccountInterface
	log              table.LogInterface
}

func (o *ChatDatabase) GetLogs(ctx context.Context, LogIDs []string, userID string) ([]*table.Log, error) {
	return o.log.Get(ctx, LogIDs, userID)
}

func (o *ChatDatabase) DeleteLogs(ctx context.Context, logID []string, userID string) error {
	return o.log.Delete(ctx, logID, userID)
}

func (o *ChatDatabase) SearchLogs(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*table.Log, error) {
	return o.log.Search(ctx, keyword, start, end, pagination)
}

func (o *ChatDatabase) UploadLogs(ctx context.Context, logs []*table.Log) error {
	return o.log.Create(ctx, logs)
}

func (o *ChatDatabase) IsNotFound(err error) bool {
	return dbutil.IsDBNotFound(err)
}

func (o *ChatDatabase) GetUser(ctx context.Context, userID string) (account *table.Account, err error) {
	return o.account.Take(ctx, userID)
}

func (o *ChatDatabase) UpdateUseInfo(ctx context.Context, userID string, attribute map[string]any) (err error) {
	return o.attribute.Update(ctx, userID, attribute)
}

func (o *ChatDatabase) FindAttribute(ctx context.Context, userIDs []string) ([]*table.Attribute, error) {
	return o.attribute.Find(ctx, userIDs)
}

func (o *ChatDatabase) FindAttributeByAccount(ctx context.Context, accounts []string) ([]*table.Attribute, error) {
	return o.attribute.FindAccount(ctx, accounts)
}

func (o *ChatDatabase) TakeAttributeByPhone(ctx context.Context, areaCode string, phoneNumber string) (*table.Attribute, error) {
	return o.attribute.TakePhone(ctx, areaCode, phoneNumber)
}

func (o *ChatDatabase) TakeAttributeByEmail(ctx context.Context, email string) (*table.Attribute, error) {
	return o.attribute.TakeEmail(ctx, email)
}

func (o *ChatDatabase) TakeAttributeByAccount(ctx context.Context, account string) (*table.Attribute, error) {
	return o.attribute.TakeAccount(ctx, account)
}

func (o *ChatDatabase) TakeAttributeByUserID(ctx context.Context, userID string) (*table.Attribute, error) {
	return o.attribute.Take(ctx, userID)
}

func (o *ChatDatabase) Search(ctx context.Context, normalUser int32, keyword string, genders int32, pagination pagination.Pagination) (total int64, attributes []*table.Attribute, err error) {
	var forbiddenIDs []string
	if int(normalUser) == constant2.NormalUser {
		forbiddenIDs, err = o.forbiddenAccount.FindAllIDs(ctx)
		if err != nil {
			return 0, nil, err
		}
	}
	total, totalUser, err := o.attribute.SearchNormalUser(ctx, keyword, forbiddenIDs, genders, pagination)
	if err != nil {
		return 0, nil, err
	}
	return total, totalUser, nil
}

func (o *ChatDatabase) SearchUser(ctx context.Context, keyword string, userIDs []string, genders []int32, pagination pagination.Pagination) (int64, []*table.Attribute, error) {
	return o.attribute.SearchUser(ctx, keyword, userIDs, genders, pagination)
}

func (o *ChatDatabase) CountVerifyCodeRange(ctx context.Context, account string, start time.Time, end time.Time) (int64, error) {
	return o.verifyCode.RangeNum(ctx, account, start, end)
}

func (o *ChatDatabase) AddVerifyCode(ctx context.Context, verifyCode *table.VerifyCode, fn func() error) error {
	return o.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := o.verifyCode.Add(ctx, []*table.VerifyCode{verifyCode}); err != nil {
			return err
		}
		if fn != nil {
			return fn()
		}
		return nil
	})
}

func (o *ChatDatabase) UpdateVerifyCodeIncrCount(ctx context.Context, id string) error {
	return o.verifyCode.Incr(ctx, id)
}

func (o *ChatDatabase) TakeLastVerifyCode(ctx context.Context, account string) (*table.VerifyCode, error) {
	return o.verifyCode.TakeLast(ctx, account)
}

func (o *ChatDatabase) DelVerifyCode(ctx context.Context, id string) error {
	return o.verifyCode.Delete(ctx, id)
}

func (o *ChatDatabase) RegisterUser(ctx context.Context, register *table.Register, account *table.Account, attribute *table.Attribute) error {
	return o.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := o.register.Create(ctx, register); err != nil {
			return err
		}
		if err := o.account.Create(ctx, account); err != nil {
			return err
		}
		if err := o.attribute.Create(ctx, attribute); err != nil {
			return err
		}
		return nil
	})
}

func (o *ChatDatabase) GetAccount(ctx context.Context, userID string) (*table.Account, error) {
	return o.account.Take(ctx, userID)
}

func (o *ChatDatabase) GetAttribute(ctx context.Context, userID string) (*table.Attribute, error) {
	return o.attribute.Take(ctx, userID)
}

func (o *ChatDatabase) GetAttributeByAccount(ctx context.Context, account string) (*table.Attribute, error) {
	return o.attribute.TakeAccount(ctx, account)
}

func (o *ChatDatabase) GetAttributeByPhone(ctx context.Context, areaCode string, phoneNumber string) (*table.Attribute, error) {
	return o.attribute.TakePhone(ctx, areaCode, phoneNumber)
}

func (o *ChatDatabase) GetAttributeByEmail(ctx context.Context, email string) (*table.Attribute, error) {
	return o.attribute.TakeEmail(ctx, email)
}

func (o *ChatDatabase) LoginRecord(ctx context.Context, record *table.UserLoginRecord, verifyCodeID *string) error {
	return o.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := o.userLoginRecord.Create(ctx, record); err != nil {
			return err
		}
		if verifyCodeID != nil {
			if err := o.verifyCode.Delete(ctx, *verifyCodeID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (o *ChatDatabase) UpdatePassword(ctx context.Context, userID string, password string) error {
	return o.account.UpdatePassword(ctx, userID, password)
}

func (o *ChatDatabase) UpdatePasswordAndDeleteVerifyCode(ctx context.Context, userID string, password string, codeID string) error {
	return o.tx.Transaction(ctx, func(ctx context.Context) error {
		if err := o.account.UpdatePassword(ctx, userID, password); err != nil {
			return err
		}
		if err := o.verifyCode.Delete(ctx, codeID); err != nil {
			return err
		}
		return nil
	})
}

func (o *ChatDatabase) NewUserCountTotal(ctx context.Context, before *time.Time) (int64, error) {
	return o.register.CountTotal(ctx, before)
}

func (o *ChatDatabase) UserLoginCountTotal(ctx context.Context, before *time.Time) (int64, error) {
	return o.userLoginRecord.CountTotal(ctx, before)
}

func (o *ChatDatabase) UserLoginCountRangeEverydayTotal(ctx context.Context, start *time.Time, end *time.Time) (map[string]int64, int64, error) {
	return o.userLoginRecord.CountRangeEverydayTotal(ctx, start, end)
}
