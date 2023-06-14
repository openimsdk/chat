package database

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/db/model/chat"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	"gorm.io/gorm"
	"time"
)

type ChatDatabaseInterface interface {
	IsNotFound(err error) bool
	GetUser(ctx context.Context, userID string) (account *table.Account, err error)
	UpdateUseInfo(ctx context.Context, userID string, attribute map[string]any, fn func() error) (err error)
	FindAttribute(ctx context.Context, userIDs []string) ([]*table.Attribute, error)
	FindAttributeByAccount(ctx context.Context, accounts []string) ([]*table.Attribute, error)
	TakeAttributeByPhone(ctx context.Context, areaCode string, phoneNumber string) (*table.Attribute, error)
	TakeAttributeByAccount(ctx context.Context, account string) (*table.Attribute, error)
	TakeAttributeByUserID(ctx context.Context, userID string) (*table.Attribute, error)
	Search(ctx context.Context, keyword string, genders []int32, pageNumber int32, showNumber int32) (uint32, []*table.Attribute, error)
	CountVerifyCodeRange(ctx context.Context, account string, start time.Time, end time.Time) (uint32, error)
	AddVerifyCode(ctx context.Context, verifyCode *table.VerifyCode, fn func() error) error
	UpdateVerifyCodeIncrCount(ctx context.Context, id uint) error
	TakeLastVerifyCode(ctx context.Context, account string) (*table.VerifyCode, error)
	DelVerifyCode(ctx context.Context, id uint) error
	RegisterUser(ctx context.Context, register *table.Register, account *table.Account, attribute *table.Attribute, fn func() error) error
	GetAccount(ctx context.Context, userID string) (*table.Account, error)
	GetAttribute(ctx context.Context, userID string) (*table.Attribute, error)
	GetAttributeByAccount(ctx context.Context, account string) (*table.Attribute, error)
	GetAttributeByPhone(ctx context.Context, areaCode string, phoneNumber string) (*table.Attribute, error)
	LoginRecord(ctx context.Context, record *table.UserLoginRecord, verifyCodeID *uint) error
	UpdatePassword(ctx context.Context, userID string, password string) error
	UpdatePasswordAndDeleteVerifyCode(ctx context.Context, userID string, password string, code uint) error
}

func NewChatDatabase(db *gorm.DB) ChatDatabaseInterface {
	return &ChatDatabase{
		tx:              tx.NewGorm(db),
		register:        chat.NewRegister(db),
		account:         chat.NewAccount(db),
		attribute:       chat.NewAttribute(db),
		userLoginRecord: chat.NewUserLoginRecord(db),
		verifyCode:      chat.NewVerifyCode(db),
	}
}

type ChatDatabase struct {
	tx              tx.Tx
	register        table.RegisterInterface
	account         table.AccountInterface
	attribute       table.AttributeInterface
	userLoginRecord table.UserLoginRecordInterface
	verifyCode      table.VerifyCodeInterface
}

func (o *ChatDatabase) IsNotFound(err error) bool {
	return errs.Unwrap(err) == gorm.ErrRecordNotFound
}

func (o *ChatDatabase) GetUser(ctx context.Context, userID string) (account *table.Account, err error) {
	return o.account.Take(ctx, userID)
}

func (o *ChatDatabase) UpdateUseInfo(ctx context.Context, userID string, attribute map[string]any, fn func() error) (err error) {
	return o.tx.Transaction(func(tx any) error {
		if len(attribute) > 0 {
			if err := o.attribute.NewTx(tx).Update(ctx, userID, attribute); err != nil {
				return err
			}
		}
		if fn != nil {
			return fn()
		}
		return nil
	})
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

func (o *ChatDatabase) TakeAttributeByAccount(ctx context.Context, account string) (*table.Attribute, error) {
	return o.attribute.TakeAccount(ctx, account)
}

func (o *ChatDatabase) TakeAttributeByUserID(ctx context.Context, userID string) (*table.Attribute, error) {
	return o.attribute.Take(ctx, userID)
}

func (o *ChatDatabase) Search(ctx context.Context, keyword string, genders []int32, pageNumber int32, showNumber int32) (uint32, []*table.Attribute, error) {
	return o.attribute.Search(ctx, keyword, genders, pageNumber, showNumber)
}

func (o *ChatDatabase) CountVerifyCodeRange(ctx context.Context, account string, start time.Time, end time.Time) (uint32, error) {
	return o.verifyCode.RangeNum(ctx, account, start, end)
}

func (o *ChatDatabase) AddVerifyCode(ctx context.Context, verifyCode *table.VerifyCode, fn func() error) error {
	return o.tx.Transaction(func(tx any) error {
		if err := o.verifyCode.NewTx(tx).Add(ctx, []*table.VerifyCode{verifyCode}); err != nil {
			return err
		}
		if fn != nil {
			return fn()
		}
		return nil
	})
}

func (o *ChatDatabase) UpdateVerifyCodeIncrCount(ctx context.Context, id uint) error {
	return o.verifyCode.Incr(ctx, id)
}

func (o *ChatDatabase) TakeLastVerifyCode(ctx context.Context, account string) (*table.VerifyCode, error) {
	return o.verifyCode.TakeLast(ctx, account)
}

func (o *ChatDatabase) DelVerifyCode(ctx context.Context, id uint) error {
	return o.verifyCode.Delete(ctx, id)
}

func (o *ChatDatabase) RegisterUser(ctx context.Context, register *table.Register, account *table.Account, attribute *table.Attribute, fn func() error) error {
	return o.tx.Transaction(func(tx any) error {
		if err := o.register.NewTx(tx).Create(ctx, register); err != nil {
			return err
		}
		if err := o.account.NewTx(tx).Create(ctx, account); err != nil {
			return err
		}
		if err := o.attribute.NewTx(tx).Create(ctx, attribute); err != nil {
			return err
		}
		if fn != nil {
			return fn()
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

func (o *ChatDatabase) LoginRecord(ctx context.Context, record *table.UserLoginRecord, verifyCodeID *uint) error {
	return o.tx.Transaction(func(tx any) error {
		if err := o.userLoginRecord.NewTx(tx).Create(ctx, record); err != nil {
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

func (o *ChatDatabase) UpdatePasswordAndDeleteVerifyCode(ctx context.Context, userID string, password string, code uint) error {
	return o.tx.Transaction(func(tx any) error {
		if err := o.account.NewTx(tx).UpdatePassword(ctx, userID, password); err != nil {
			return err
		}
		if err := o.verifyCode.NewTx(tx).Delete(ctx, code); err != nil {
			return err
		}
		return nil
	})
}
