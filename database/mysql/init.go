package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func Init(conf *Config) (*gorm.DB, error) {
	gconf := &gorm.Config{
		SkipDefaultTransaction: true, // 全局禁用事务,需要时在临时开启/关闭
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&loc=Local", conf.User, conf.Password, conf.Url, conf.Dbname, conf.Charset)

	orm, err := gorm.Open(mysql.Open(dsn), gconf)
	if err != nil {
		return nil, err
	}
	sqlDB, err := orm.DB()
	if err != nil {
		return nil, err
	}
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(conf.MaxIdle)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(conf.MaxConn)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour * 8)

	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}

	orm.Logger.LogMode(2)
	return orm, err
}

func Transaction(db *gorm.DB, do func(tx *gorm.DB) error) (err error) {
	// 开始事务
	tx := db.Begin()
	defer func() {
		if e := recover(); e != nil {
			tx.Rollback()
			err = fmt.Errorf("%v", e)
		}
	}()
	err = do(tx)
	if err != nil {
		// 遇到错误时回滚事务
		tx.Rollback()
		return err
	}
	// 否则，提交事务
	tx.Commit()
	return nil
}
