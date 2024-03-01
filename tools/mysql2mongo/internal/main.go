package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/db/model/admin"
	"github.com/OpenIMSDK/chat/pkg/common/db/model/chat"
	"github.com/OpenIMSDK/chat/pkg/common/dbconn"
	"github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v3"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"reflect"
	"strconv"
)

const (
	versionTable = "chatver"
	versionKey   = "data_version"
	versionValue = 1
)

func SetMongoDataVersion(db *mongo.Database, curver string) error {
	filter := bson.M{"key": versionKey, "value": curver}
	update := bson.M{"$set": bson.M{"key": versionKey, "value": strconv.Itoa(versionValue)}}
	_, err := db.Collection(versionTable).UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
	return err
}

func InitConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &config.Config)
}

func GetMysql() (*gorm.DB, error) {
	conf := config.Config.Mysql
	mysqlDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", *conf.Username, *conf.Password, (*conf.Address)[0], *conf.Database)
	return gorm.Open(gormmysql.Open(mysqlDSN), &gorm.Config{Logger: logger.Discard})
}

func getColl(obj any) (_ *mongo.Collection, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("not found %+v", e)
		}
	}()
	stu := reflect.ValueOf(obj).Elem()
	typ := reflect.TypeOf(&mongo.Collection{}).String()
	for i := 0; i < stu.NumField(); i++ {
		field := stu.Field(i)
		if field.Type().String() == typ {
			return (*mongo.Collection)(field.UnsafePointer()), nil
		}
	}
	return nil, errors.New("not found model collection")
}

// NewTask A mysql table B mongodb model C mongodb table
func NewTask[A interface{ TableName() string }, B any, C any](gormDB *gorm.DB, mongoDB *mongo.Database, mongoDBInit func(db *mongo.Database) (B, error), convert func(v A) C) error {
	obj, err := mongoDBInit(mongoDB)
	if err != nil {
		return err
	}
	var zero A
	tableName := zero.TableName()
	coll, err := getColl(obj)
	if err != nil {
		return fmt.Errorf("get mongo collection %s failed, err: %w", tableName, err)
	}
	var count int
	defer func() {
		log.Printf("completed convert chat %s total %d\n", tableName, count)
	}()
	const batch = 100
	for page := 0; ; page++ {
		res := make([]A, 0, batch)
		if err := gormDB.Limit(batch).Offset(page * batch).Find(&res).Error; err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1146 {
				return nil // table not exist
			}
			return fmt.Errorf("find mysql table %s failed, err: %w", tableName, err)
		}
		if len(res) == 0 {
			return nil
		}
		temp := make([]any, len(res))
		for i := range res {
			temp[i] = convert(res[i])
		}
		if err := insertMany(coll, temp); err != nil {
			return fmt.Errorf("insert mongo table %s failed, err: %w", tableName, err)
		}
		count += len(res)
		if len(res) < batch {
			return nil
		}
		log.Printf("current convert chat %s completed %d\n", tableName, count)
	}
}

func insertMany(coll *mongo.Collection, objs []any) error {
	if _, err := coll.InsertMany(context.Background(), objs); err != nil {
		if !mongo.IsDuplicateKeyError(err) {
			return err
		}
	}
	for i := range objs {
		_, err := coll.InsertOne(context.Background(), objs[i])
		switch {
		case err == nil:
		case mongo.IsDuplicateKeyError(err):
		default:
			return err
		}
	}
	return nil
}

func Main(path string) error {
	if err := InitConfig(path); err != nil {
		return err
	}
	if config.Config.Mysql == nil {
		log.Println("mysql config is nil")
		return nil
	}
	mongoDB, err := dbconn.NewMongo()
	if err != nil {
		return err
	}
	var version struct {
		Key   string `bson:"key"`
		Value string `bson:"value"`
	}
	switch mongoDB.Collection(versionTable).FindOne(context.Background(), bson.M{"key": versionKey}).Decode(&version) {
	case nil:
		if ver, _ := strconv.Atoi(version.Value); ver >= versionValue {
			return nil
		}
	case mongo.ErrNoDocuments:
	default:
		return err
	}
	mysqlDB, err := GetMysql()
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1049 {
			if err := SetMongoDataVersion(mongoDB, version.Value); err != nil {
				return err
			}
			log.Println("set chat version config")
			return nil // database not exist
		}
		return err
	}

	var (
		cc convertChat
		ca convertAdmin
	)

	var tasks []func() error
	tasks = append(tasks,
		// chat
		func() error { return NewTask(mysqlDB, mongoDB, chat.NewAccount, cc.Account) },
		func() error { return NewTask(mysqlDB, mongoDB, chat.NewAttribute, cc.Attribute) },
		func() error { return NewTask(mysqlDB, mongoDB, chat.NewLogs, cc.Log) },
		func() error { return NewTask(mysqlDB, mongoDB, chat.NewRegister, cc.Register) },
		func() error { return NewTask(mysqlDB, mongoDB, chat.NewUserLoginRecord, cc.UserLoginRecord) },
		// admin
		func() error { return NewTask(mysqlDB, mongoDB, admin.NewAdmin, ca.Admin) },
		func() error { return NewTask(mysqlDB, mongoDB, admin.NewApplet, ca.Applet) },
		func() error { return NewTask(mysqlDB, mongoDB, admin.NewClientConfig, ca.ClientConfig) },
		func() error { return NewTask(mysqlDB, mongoDB, admin.NewForbiddenAccount, ca.ForbiddenAccount) },
		func() error { return NewTask(mysqlDB, mongoDB, admin.NewInvitationRegister, ca.InvitationRegister) },
		func() error { return NewTask(mysqlDB, mongoDB, admin.NewIPForbidden, ca.IPForbidden) },
		func() error { return NewTask(mysqlDB, mongoDB, admin.NewLimitUserLoginIP, ca.LimitUserLoginIP) },
		func() error { return NewTask(mysqlDB, mongoDB, admin.NewRegisterAddFriend, ca.RegisterAddFriend) },
		func() error { return NewTask(mysqlDB, mongoDB, admin.NewRegisterAddGroup, ca.RegisterAddGroup) },
	)

	for _, task := range tasks {
		if err := task(); err != nil {
			return err
		}
	}

	if err := SetMongoDataVersion(mongoDB, version.Value); err != nil {
		return err
	}

	return nil
}
