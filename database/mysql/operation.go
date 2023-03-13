package mysql

import (
	"github.com/tieing/lemon/tools/convertor"
	"gorm.io/gorm"
	"strings"
)

type BaseQuery interface {
	Insert(db *gorm.DB, data interface{}) (int64, error)
	Update(db *gorm.DB, id int, model interface{}, values interface{}) error
	Updates(db *gorm.DB, model interface{}, values interface{}, filters ...interface{}) error
	FindOne(db *gorm.DB, model interface{}, selects string, filters ...interface{}) (err error)
	FindGroupList(db *gorm.DB, pageNum int, limit int, selects string, group string, order string, list interface{}, filters ...interface{}) (err error)
	FindList(db *gorm.DB, list interface{}, selects string, sort string, filters ...interface{}) (err error)
	Delete(db *gorm.DB, model interface{}, filters ...interface{}) (err error)
	FindByPage(db *gorm.DB, pageNum, limit int, selects string, list interface{}, sort string, filters ...interface{}) (count int64, err error)
}

type BaseQueryImpl struct{}

func (imp *BaseQueryImpl) Insert(db *gorm.DB, data interface{}) (id int64, err error) {
	tx := db.Create(data)
	return tx.RowsAffected, tx.Error
}

func (imp *BaseQueryImpl) Update(db *gorm.DB, id int, model interface{}, values interface{}) error {
	err := db.Model(model).Where("id in (?)", id).Updates(values).Error
	return err
}

func (imp *BaseQueryImpl) Updates(db *gorm.DB, model interface{}, values interface{}, filters ...interface{}) error {
	queryArr, values2 := whereArr(filters...)
	err := db.Model(model).Where(strings.Join(queryArr, " AND "), values2...).Updates(values).Error
	return err
}

func (imp *BaseQueryImpl) Delete(db *gorm.DB, model interface{}, filters ...interface{}) (err error) {
	queryArr, values := whereArr(filters...)
	err = db.Where(strings.Join(queryArr, " AND "), values...).Delete(model).Error
	return
}

func (imp *BaseQueryImpl) FindOne(db *gorm.DB, data any, selects string, filters ...interface{}) (err error) {
	queryArr, values := whereArr(filters...)
	err = db.Model(data).Select(selects).Where(strings.Join(queryArr, " AND "), values...).Order("id desc").First(data).Error
	return
}

func (imp *BaseQueryImpl) FindList(db *gorm.DB, list interface{}, selects string, sort string, filters ...interface{}) (err error) {
	queryArr, values := whereArr(filters...)
	query := db.Model(list)
	err = query.Select(selects).Where(strings.Join(queryArr, " AND "), values...).Order(sort).Find(list).Error
	return
}
func (imp *BaseQueryImpl) FindGroupList(db *gorm.DB, pageNum int, limit int, selects string, group string, order string, list interface{}, filters ...interface{}) (err error) {
	offset := (pageNum - 1) * limit
	queryArr, values := whereArr(filters...)
	query := db.Model(list)
	err = query.Select(selects).Where(strings.Join(queryArr, " AND "), values...).Limit(limit).Offset(offset).Order(order).Group(group).Find(list).Error
	return
}

func (imp *BaseQueryImpl) FindByPage(db *gorm.DB, pageNum, limit int, selects string, list interface{}, sort string, filters ...interface{}) (count int64, err error) {
	offset := (pageNum - 1) * limit
	queryArr, values := whereArr(filters...)
	query := db.Model(list)
	query.Where(strings.Join(queryArr, " AND "), values...).Count(&count)
	err = query.Select(selects).Order(sort).Limit(limit).Offset(offset).Find(list).Error
	return
}

func whereArr(filters ...interface{}) ([]string, []interface{}) {
	var queryArr []string
	var values []interface{}
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			queryArr = append(queryArr, convertor.ToString(filters[k]))
			values = append(values, filters[k+1])
		}
	}
	return queryArr, values
}
