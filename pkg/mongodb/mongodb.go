package mongodb

import (
	"errors"
	"reflect"
	"unsafe"

	"go.mongodb.org/mongo-driver/bson"
)

type Mgo struct{}

var (
	errModelTypeNotPointer = errors.New("Mongodb results argument must be a pointer to a slice ")
	errInvalidFilter       = errors.New("Invalid filter ")
	errInvalidSelection    = errors.New("Invalid selection ")
)

//var (
//	ErrNoRecord = errors.New("mongo: no documents in result")
//)

func (mgo *Mgo) validResult(model interface{}) error {
	t := reflect.ValueOf(model)
	if t.Kind() != reflect.Ptr {
		return errModelTypeNotPointer
	}
	return nil
}

func (mgo *Mgo) validWhere(where interface{}) (bson.M, bool) {
	var sql bson.M

	if where == nil {
		return sql, true
	}

	switch whereAssert := where.(type) {
	case bson.M:
		pp := whereAssert
		sql = *(*bson.M)(unsafe.Pointer(&pp))
	case map[string]interface{}:
		pp := whereAssert
		sql = *(*bson.M)(unsafe.Pointer(&pp))
	default:
		return nil, false
	}
	return sql, true
}

func (mgo *Mgo) valid(filter, selection, result interface{}) (bson.M, bson.M, error) {
	var (
		sql, where bson.M
		ok         bool
	)
	if err := mgo.validResult(result); err != nil {
		return nil, nil, err
	}

	if sql, ok = mgo.validWhere(filter); !ok {
		return nil, nil, errInvalidFilter
	}

	if where, ok = mgo.validWhere(selection); !ok {
		return nil, nil, errInvalidSelection
	}

	return sql, where, nil
}

//nolint
func (mgo *Mgo) validWithOutSelection(filter, result interface{}) (bson.M, bson.M, error) {
	var (
		sql, where bson.M //todo where is always nil
		ok         bool
	)
	if err := mgo.validResult(result); err != nil {
		return nil, nil, err
	}

	if sql, ok = mgo.validWhere(filter); !ok {
		return nil, nil, errInvalidFilter
	}

	return sql, where, nil
}

func (mgo *Mgo) FindOne(table string, filter, selection, result interface{}) error {

	sql, where, err := mgo.valid(filter, selection, result)
	if err != nil {
		return err
	}

	m := NewMongoInfo(sql, where, result)
	if table != "" {
		m.SetTable(table)
	}

	if err := m.Get(); err != nil {
		return err
	}

	return nil
}

func (mgo *Mgo) FindMany(table string, filter, selection, result interface{}) error {

	sql, where, err := mgo.valid(filter, selection, result)
	if err != nil {
		return err
	}

	m := NewMongoInfo(sql, where, result)
	if table != "" {
		m.SetTable(table)
	}

	if err := m.Get(); err != nil {
		return err
	}

	return nil

}

func (mgo *Mgo) FindManyWithExclude(table string, sort, exclude bson.M, filter, result interface{}) error {

	sql, where, err := mgo.validWithOutSelection(filter, result)
	if err != nil {
		return err
	}

	m := NewMongoInfo(sql, where, result)
	if table != "" {
		m.SetTable(table)
	}
	m = m.SetSelection(exclude).SetSort(sort)

	if err := m.Get(); err != nil {
		return err
	}

	return nil

}

func (mgo *Mgo) UpdateOne(table string, filter, data interface{}) error {
	var (
		sql bson.M
		ok  bool
	)

	if sql, ok = mgo.validWhere(filter); !ok {
		return errInvalidFilter
	}

	m := NewMongoInfo(sql, nil, data)
	if table != "" {
		m.SetTable(table)
	}

	if _, err := m.UpdateOne(); err != nil {
		return err
	}

	return nil
}

func (mgo *Mgo) UpdateMany(table string, filter, data interface{}) error {
	var (
		sql bson.M
		ok  bool
	)

	if sql, ok = mgo.validWhere(filter); !ok {
		return errInvalidFilter
	}

	m := NewMongoInfo(sql, nil, data)
	if table != "" {
		m.SetTable(table)
	}

	if _, err := m.UpdateMany(); err != nil {
		return err
	}

	return nil
}

func (mgo *Mgo) DeleteOne(table string, filter interface{}) error {
	var (
		sql bson.M
		ok  bool
	)
	if sql, ok = mgo.validWhere(filter); !ok {
		return errInvalidFilter
	}
	m := NewMongoInfo(sql, nil, nil)
	if table != "" {
		m.SetTable(table)
	}

	if _, err := m.DeleteOne(); err != nil {
		return err
	}

	return nil
}

func (mgo *Mgo) DeleteMany(table string, filter interface{}) error {

	var (
		sql bson.M
		ok  bool
	)
	if sql, ok = mgo.validWhere(filter); !ok {
		return errInvalidFilter
	}
	m := NewMongoInfo(sql, nil, nil)
	if table != "" {
		m.SetTable(table)
	}

	if _, err := m.DeleteAll(); err != nil {
		return err
	}

	return nil

}

// todo
//func (mgo *Mgo) FindOneAndUpdate(table string, filter interface{}) error {
//	var (
//		sql bson.M
//		ok bool
//	)
//
//	if sql, ok = mgo.validWhere(filter); ok {
//		return errInvalidFilter
//	}
//	m := NewMongoInfo(sql, nil, nil)
//	if table != "" {
//		m.SetTable(table)
//	}
//
//	return nil
//}

func (mgo *Mgo) Upsert(table string, filter, data interface{}) error {
	var (
		sql bson.M
		ok  bool
	)

	if sql, ok = mgo.validWhere(filter); !ok {
		return errInvalidFilter
	}

	m := NewMongoInfo(sql, nil, data)
	if table != "" {
		m.SetTable(table)
	}

	if _, err := m.Upsert(); err != nil {
		return err
	}

	return nil
}

func (mgo *Mgo) InsertOne(table string, data interface{}) error {

	m := NewMongoInfo(nil, nil, data)
	if table != "" {
		m.SetTable(table)
	}

	if _, err := m.InsertOne(); err != nil {
		return err
	}

	return nil
}

func (mgo *Mgo) InsertMany(table string, data []interface{}) error {

	m := NewMongoInfos(nil, nil, data)
	if table != "" {
		m.SetTable(table)
	}

	if _, err := m.InsertMany(); err != nil {
		return err
	}

	return nil

}

func (mgo *Mgo) Search(table string, result interface{}, q string, fields, excludes []string,
	filters map[string]interface{},
	page, limit int64, sort map[string]interface{}) (int64, error) {

	mg := NewMongoInfo(nil, nil, nil)
	if table != "" {
		mg.SetTable(table)
	}

	return mg.Find(result, q, fields, excludes, filters, page, limit, sort)

}

func (mgo *Mgo) Aggregate(table string, pipelines interface{}, result interface{}) error {
	mg := NewMongoInfo(nil, nil, result)
	if table != "" {
		mg.SetTable(table)
	}
	if pipelines != nil {
		mg.SetPipeline(pipelines)
	}
	return mg.Aggregate()
}

// 单纯的生成查询条件 为了配合 Aggregate 的match 操作
func (mgo *Mgo) AddFilter(exactMatch bool, q string, fields []string, filters map[string]interface{}) (bson.M, error) {

	var (
		andCondition []bson.M
		condition    = make(bson.M)
	)

	if len(fields) == 0 {
		err := errors.New("search fields empty")
		return condition, err
	}

	if q != "" {
		bor := make([]bson.M, len(fields))
		for i := 0; i < len(fields); i++ {
			if exactMatch {
				bor[i] = bson.M{fields[i]: q}
			} else {
				rq := bson.M{"$regex": q, "$options": "i"}
				bor[i] = bson.M{fields[i]: rq}
			}
		}
		condition["$or"] = bor

		andCondition = append(andCondition, condition)
		//db.Condition = condition
	}

	if len(filters) > 0 {
		andCondition = append(andCondition, filters)
		condition = bson.M{
			"$and": andCondition,
		}
	}

	return condition, nil
}
