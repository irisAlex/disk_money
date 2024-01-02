// Author liyan
// Date 2020-03-27 5:25 下午
// Mail liyana@hualala.com
// org base cloud platform

package mongodb

import (
	"sync"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	mgo *Mgo

	once sync.Once
)

func NewMongodb() Modeler {
	once.Do(func() {
		mgo = new(Mgo)
	})
	return mgo
}

type Modeler interface {
	// 插入一条记录
	// table 表名, mongodb 的 collection
	// data 插入的数据
	InsertOne(table string, data interface{}) error

	// 插入多条数据
	// table 表名, mongodb 的 collection
	// data 插入的数据
	InsertMany(table string, data []interface{}) error

	// Update if exist, insert if not exist
	Upsert(table string, filter, data interface{}) error

	// 更新单条记录
	// table 表名, mongodb 的 collection
	// where 更新条件,mongodb 的 filter
	// data 插入的数据
	UpdateOne(table string, filter, data interface{}) error

	// 更新单条记录
	// table 表名, mongodb 的 collection
	// where 更新条件,mongodb 的 filter
	// data 插入的数据
	UpdateMany(table string, filter, data interface{}) error

	// 删除单条数据
	// table 表名, mongodb 的 collection
	// where 更新条件,mongodb 的 filter
	DeleteOne(table string, filter interface{}) error

	// 删除多条数据
	// table 表名, mongodb 的 collection
	// where 更新条件,mongodb 的 filter
	DeleteMany(table string, filter interface{}) error

	// 查询单条数据
	// table 表名, mongodb 的 collection
	// where 更新条件,mongodb 的 filter
	// result 查询结果
	FindOne(table string, filter, selection, result interface{}) error

	// 查询多条数据
	// table 表名, mongodb 的 collection
	// where 更新条件,mongodb 的 filter
	// result 查询结果
	FindMany(table string, filter, selection, result interface{}) error

	FindManyWithExclude(table string, sort, selection bson.M, filter, result interface{}) error

	// 根据过滤条件搜索, 支持分页、排序
	// table 表名
	// result 返回结果
	// q 搜索关键字, 根据 fields 提供的字段
	// excludes 排除
	// filters 过滤条件
	// page 分页数
	// limit 每页大小
	// sort 排序字段
	Search(table string, result interface{}, q string, fields, excludes []string,
		filters map[string]interface{},
		page, limit int64, sort map[string]interface{}) (int64, error)

	Aggregate(table string, pipelines interface{}, result interface{}) error
}
