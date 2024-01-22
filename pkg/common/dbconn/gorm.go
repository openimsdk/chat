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

package dbconn

import (
	"fmt"
	"time"

	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mw/specialerror"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMysqlGormDB() (*gorm.DB, error) {
	// Construct the DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		*config.Config.Mysql.Username, *config.Config.Mysql.Password, (*config.Config.Mysql.Address)[0], "mysql")

	// First attempt to open a new database connection
	db, err := gorm.Open(mysql.Open(dsn), nil)
	if err != nil {
		time.Sleep(time.Duration(30) * time.Second)
		db, err = gorm.Open(mysql.Open(dsn), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to open initial database connection with DSN %s: %w", dsn, err)
		}

	}

	sqlDB, err := db.DB()
	if err != nil {
		// Include specific function name and the DSN used
		return nil, fmt.Errorf("failed to get underlying sql.DB from GORM with DSN %s: %w", dsn, err)
	}
	defer sqlDB.Close()
	sql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s default charset utf8mb4 COLLATE utf8mb4_unicode_ci;", *config.Config.Mysql.Database)
	err = db.Exec(sql).Error
	if err != nil {
		return nil, fmt.Errorf("init db %w", err)
	}
	// Reconnect with the specific database
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		*config.Config.Mysql.Username, *config.Config.Mysql.Password, (*config.Config.Mysql.Address)[0], *config.Config.Mysql.Database)

	// Custom logger
	sqlLogger := log.NewSqlLogger(logger.LogLevel(*config.Config.Mysql.LogLevel), true, time.Duration(*config.Config.Mysql.SlowThreshold)*time.Millisecond)

	// Second attempt to open a new database connection
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: sqlLogger,
	})
	if err != nil {
		// Include specific function name and the DSN used
		return nil, fmt.Errorf("failed to open database connection with specific database using DSN %s: %w", dsn, err)
	}

	sqlDB, err = db.DB()
	if err != nil {
		// Include specific function name and the DSN used
		return nil, fmt.Errorf("failed to get underlying sql.DB from GORM with specific database using DSN %s: %w", dsn, err)
	}

	// Database connection configuration
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(*config.Config.Mysql.MaxLifeTime))
	sqlDB.SetMaxOpenConns(*config.Config.Mysql.MaxOpenConn)
	sqlDB.SetMaxIdleConns(*config.Config.Mysql.MaxIdleConn)

	return db, nil
}

func NewGormDB() (*gorm.DB, error) {
	specialerror.AddReplace(gorm.ErrRecordNotFound, errs.ErrRecordNotFound)
	specialerror.AddErrHandler(replaceDuplicateKey)
	return NewMysqlGormDB()
}

func replaceDuplicateKey(err error) errs.CodeError {
	if IsMysqlDuplicateKey(err) {
		return errs.ErrDuplicateKey
	}
	return nil
}

func IsMysqlDuplicateKey(err error) bool {
	if mysqlErr, ok := err.(*mysqlDriver.MySQLError); ok {
		return mysqlErr.Number == 1062
	}
	return false
}
