package mongodb

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	docsNull      = errors.New("mongodb: docs is null")
	conditionNull = errors.New("mongodb: condition is null")
)

type MongoInfo struct {
	// 数据库表名
	Table string

	// 返回的个数
	Count int64

	// 查询条件
	Condition bson.M

	// 指定返回字段
	Selection bson.M

	// 聚合操作的pipeline
	Pipeline interface{}

	// 文档类型，新增，更新单个情况(当没有指定Table，会根据结构体反射出Table名)
	Doc interface{}

	// 多个文档类型,用于插入个数情况
	Docs []interface{}

	// 自定义更新setting
	setting bool

	Sort bson.M
}

func NewMongoInfo(condition, selection bson.M, doc interface{}) *MongoInfo {
	if selection == nil {
		selection = bson.M{
			"_id": false,
		}
	}

	if doc == nil {
		var re []interface{}
		doc = &re
	}

	return &MongoInfo{
		Selection: selection,
		Condition: condition,
		Doc:       doc,
	}
}
func NewMongoInfos(condition, selection bson.M, doc []interface{}) *MongoInfo {
	if selection == nil {
		selection = bson.M{
			"_id": false,
		}
	}

	//if doc == nil {
	//	var re []interface{}
	//	doc = &re
	//}

	return &MongoInfo{
		Selection: selection,
		Condition: condition,
		Docs:      doc,
	}
}

// https://golangtc.com/t/57330f47b09ecc050a0000f1
//var dbCoLimit = make(chan struct{}, 10)

func DBCallWrap(fun func(database *mongo.Database)) {
	// todo: 全局控制读写锁，防止数据并发操作导致数据原子性问题

	fun(MgoClient.Database(dbName))
}

func (db *MongoInfo) SetTable(table string) *MongoInfo {
	db.Table = table
	return db
}

func (db *MongoInfo) SetCondition(condition bson.M) *MongoInfo {
	db.Condition = condition
	return db
}

func (db *MongoInfo) SetSelection(selection bson.M) *MongoInfo {
	db.Selection = selection
	return db
}

func (db *MongoInfo) SetSort(sort bson.M) *MongoInfo {
	db.Sort = sort
	return db
}

//func (db *MongoInfo) SetPipeline(pipeline []bson.M) *MongoInfo {
//	db.Pipeline = pipeline
//	return db
//}

func (db *MongoInfo) SetPipeline(pipeline interface{} /*[]bson.M*/) *MongoInfo {
	db.Pipeline = pipeline
	return db
}

func (db *MongoInfo) SetDoc(doc interface{}) *MongoInfo {
	db.Doc = doc
	return db
}

func (db *MongoInfo) SetDocs(docs []interface{}) *MongoInfo {
	db.Docs = docs
	return db
}

func (db *MongoInfo) SetSetting() *MongoInfo {
	db.setting = true
	return db
}

func (db *MongoInfo) Collection(database *mongo.Database) *mongo.Collection {
	var (
		col *mongo.Collection
	)
	if db.Table != "" {
		col = database.Collection(db.Table)
	} else {
		col = database.Collection(colName(db.Doc))
	}
	return col
}

func (db *MongoInfo) Where() *Query {
	var (
		query = new(Query)
	)

	query.db = db
	query.selectColumn = db.Selection
	query.conditions = append(query.conditions, db.Condition)

	return query
}

func (db *MongoInfo) Get() error {
	var (
		err       error
		cursor    *mongo.Cursor
		sigResult *mongo.SingleResult
	)
	DBCallWrap(func(database *mongo.Database) {

		col := db.Collection(database)

		defer func() {
			if cursor != nil {
				_ = cursor.Close(ctx)
			}
		}()

		if isSlice(db.Doc) {
			ops := &options.FindOptions{Projection: db.Selection, Sort: db.Sort}
			//if db.Selection != nil {
			//	cursor, err = col.Find(ctx, db.Condition, ops)
			//	if err != nil {
			//		return
			//	}
			//} else {
			cursor, err = col.Find(ctx, db.Condition, ops)
			if err != nil {
				return
			}
			//}
			err = cursor.All(ctx, db.Doc)
			if err != nil {
				return
			}

		} else {
			if db.Selection != nil {
				sigResult = col.FindOne(ctx, db.Condition, &options.FindOneOptions{Projection: db.Selection})
			} else {
				sigResult = col.FindOne(ctx, db.Condition)
			}
			err = sigResult.Decode(db.Doc)
			if err != nil {
				return
			}
		}
	})

	return err
}

func (db *MongoInfo) UpdateOne() (*mongo.UpdateResult, error) {

	if db.Doc == nil {
		return nil, docsNull
	}

	var (
		err          error
		updateResult *mongo.UpdateResult
	)
	DBCallWrap(func(database *mongo.Database) {

		col := db.Collection(database)

		updateResult, err = col.UpdateOne(ctx, db.Condition, db.Doc)
		if err != nil {
			return
		}
	})

	return updateResult, err
}

func (db *MongoInfo) UpdateMany() (*mongo.UpdateResult, error) {

	if db.Doc == nil {
		return nil, docsNull
	}

	var (
		err          error
		updateResult *mongo.UpdateResult
	)
	DBCallWrap(func(database *mongo.Database) {

		col := db.Collection(database)

		updateResult, err = col.UpdateMany(ctx, db.Condition, db.Doc)
		if err != nil {
			return
		}
	})

	return updateResult, err
}

// todo
//func (db *MongoInfo) Fin() {
//	DBCallWrap(func(database *mongo.Database) {
//		col := db.Collection(database)
//		col.FindOneAndUpdate(ctx, )
//		 UpdateOne(ctx context.Context, filter interface{},
//		update interface{}, opts ...*options.FindOneAndUpdateOptions)
//		 UpdateOne(ctx context.Context, filter interface{},
//		 update interface{}, opts ...*options.UpdateOptions)
//
//	})
//}

func (db *MongoInfo) Upsert() (*mongo.UpdateResult, error) {

	if db.Doc == nil {
		return nil, docsNull
	}

	var (
		err          error
		updateResult *mongo.UpdateResult
	)
	DBCallWrap(func(database *mongo.Database) {

		col := db.Collection(database)

		updateResult, err = col.UpdateOne(ctx, db.Condition, bson.M{"$set": db.Doc}, options.Update().SetUpsert(true))
		if err != nil {
			return
		}
	})

	return updateResult, err
}

func (db *MongoInfo) InsertOne() (*mongo.InsertOneResult, error) {
	if db.Doc == nil {
		return nil, docsNull
	}

	var (
		err             error
		insertOneResult *mongo.InsertOneResult
	)
	DBCallWrap(func(database *mongo.Database) {

		col := db.Collection(database)
		insertOneResult, err = col.InsertOne(ctx, db.Doc)
		if err != nil {
			return
		}
	})

	return insertOneResult, err
}

func (db *MongoInfo) InsertMany() (*mongo.InsertManyResult, error) {
	if db.Docs == nil {
		return nil, docsNull
	}

	var (
		err              error
		insertManyResult *mongo.InsertManyResult
	)
	DBCallWrap(func(database *mongo.Database) {

		col := db.Collection(database)
		insertManyResult, err = col.InsertMany(ctx, db.Docs)
		if err != nil {
			return
		}
	})

	return insertManyResult, err
}

func (db *MongoInfo) Aggregate() error {

	var (
		err    error
		cursor *mongo.Cursor
	)
	DBCallWrap(func(database *mongo.Database) {

		col := db.Collection(database)

		defer func() {
			if cursor != nil {
				_ = cursor.Close(ctx)
			}
		}()

		cursor, err = col.Aggregate(ctx, db.Pipeline)
		if err != nil {
			return
		}

		if isSlice(db.Doc) {
			if err = cursor.All(ctx, db.Doc); err != nil {
				return
			}
		} else {
			for cursor.Next(ctx) {
				if err = cursor.Decode(db.Doc); err != nil {
					return
				}
			}
		}
	})

	return err
}

func (db *MongoInfo) DeleteOne() (*mongo.DeleteResult, error) {
	if db.Condition == nil {
		return nil, conditionNull
	}

	var (
		err             error
		deleteOneResult *mongo.DeleteResult
	)
	DBCallWrap(func(database *mongo.Database) {

		col := db.Collection(database)
		deleteOneResult, err = col.DeleteOne(ctx, db.Condition)
		if err != nil {
			return
		}
	})

	return deleteOneResult, err
}

func (db *MongoInfo) DeleteAll() (*mongo.DeleteResult, error) {
	if db.Condition == nil {
		return nil, conditionNull
	}

	var (
		err          error
		deleteResult *mongo.DeleteResult
	)
	DBCallWrap(func(database *mongo.Database) {

		col := db.Collection(database)
		deleteResult, err = col.DeleteMany(ctx, db.Condition)
		if err != nil {
			return
		}
	})

	return deleteResult, err
}

func (db *MongoInfo) Find(
	re interface{}, q string, fields, excludes []string, filters map[string]interface{}, page, limit int64, sort bson.M) (int64, error) {

	db.Doc = re

	if e := db.findCondition(q, fields, filters); e != nil {
		return 0, e
	}

	if e := db.findExclude(excludes); e != nil {
		return 0, e
	}

	qs := db.Where()
	if sort != nil {
		qs.Sort(sort)
	}

	if e := qs.Count(); e != nil {
		return 0, e
	}

	qs.Paginate(limit, page)

	if e := qs.Find(); e != nil {
		return 0, e
	}

	return db.Count, nil
}

//nolint
func (db *MongoInfo) findExclude(excludes []string) error {
	ex := bson.M{}
	for i, n := 0, len(excludes); i < n; i++ {
		ex[excludes[i]] = false
	}

	db.Selection = ex
	return nil
}

func (db *MongoInfo) findCondition(q string, fields []string, filters map[string]interface{}) error {
	var (
		condition    = make(bson.M)
		andCondition []bson.M
	)

	if len(fields) != 0 && q != "" {
		rq := bson.M{"$regex": q, "$options": "i"}
		bor := make([]bson.M, len(fields))
		for i := 0; i < len(fields); i++ {
			bor[i] = bson.M{fields[i]: rq}
		}
		condition["$or"] = bor
		db.Condition = condition
		andCondition = append(andCondition, condition)
		//db.Condition = condition
	}

	if len(filters) > 0 {
		andCondition = append(andCondition, filters)
		db.Condition = bson.M{
			"$and": andCondition,
		}
	}

	return nil

}

func (db *MongoInfo) Distinct(key string, re interface{}) error {
	var (
		err error
	)

	DBCallWrap(func(database *mongo.Database) {
		col := db.Collection(database)
		re, err = col.Distinct(ctx, key, db.Condition)
		if err != nil {
			return
		}
	})
	return err
}

type ResultCount struct {
	Total int `bson:"total"`
}
