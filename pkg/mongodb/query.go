// Author liyan
// Date 2020-03-27 5:26 下午
// Mail liyana@hualala.com
// org base cloud platform

package mongodb

import (
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Query struct {
	db                *MongoInfo
	conditions        []bson.M
	selectColumn      bson.M
	selectedOneColumn bool
	paginated         bool
	sort              bson.M
	querySet          *mongo.Cursor
	collection        *mongo.Collection
	needReload        bool
	skip, limit       int64 //todo:  field `count` is unused
}

func (q *Query) Select(s bson.M) *Query {
	q.selectColumn = s
	q.selectedOneColumn = true
	q.needReload = true
	return q
}

func (q *Query) Or(s bson.M) *Query {
	q.conditions = append(q.conditions, s)
	q.needReload = true
	return q
}

func (q *Query) Sort(s bson.M) *Query {
	q.sort = s
	q.needReload = true
	return q
}

func (q *Query) Paginate(limit, page int64) *Query {
	q.paginated = true
	q.limit = limit
	q.skip = (page - 1) * limit
	q.needReload = true
	return q
}

func (q *Query) Limit(limit int64) *Query {
	q.paginated = true
	q.limit = limit
	q.skip = 0
	q.needReload = true
	return q
}

func (q *Query) Find() error {
	var err error
	DBCallWrap(func(database *mongo.Database) {

		err = q.prepare(database)
		if err != nil {
			return
		}

		err = q.Result()
		if err != nil {
			return
		}
	})

	return err
}

func (q *Query) Count() error {

	var err error

	DBCallWrap(func(database *mongo.Database) {

		err = q.prepare(database)
		if err != nil {
			return
		}

		col := q.db.Collection(database)
		query := q.parseQuery()
		q.db.Count, err = col.CountDocuments(ctx, query)
	})

	return err
}

//func (q *Query) Q(database *mongo.Database) *mongo.Collection {
//	q.prepare(database)
//	return q.querySet
//}

func (q *Query) parseQuery() bson.M {
	if len(q.conditions) < 2 {
		return q.conditions[0]
	}
	return bson.M{"$or": q.conditions}
}

func (q *Query) loadQuerySet(database *mongo.Database, query bson.M) error {

	col := q.db.Collection(database)
	q.collection = col

	var (
		cursor *mongo.Cursor
		err    error
	)
	findOptions := &options.FindOptions{}

	if q.selectedOneColumn {
		findOptions.SetProjection(q.selectColumn)
	}

	if q.paginated {
		findOptions.SetSkip(q.skip)
		findOptions.SetLimit(q.limit)
	}

	if q.sort != nil {
		findOptions.SetSort(q.sort)
	}

	cursor, err = col.Find(ctx, query, findOptions)
	if err != nil {
		return err
	}

	q.querySet = cursor
	return nil
}

func (q *Query) prepare(database *mongo.Database) error {

	if q.querySet == nil || q.needReload {
		query := q.parseQuery()
		if err := q.loadQuerySet(database, query); err != nil {
			return err
		}
	}
	return nil
}

func (q *Query) Result() error {
	if isSlice(q.db.Doc) {
		return q.querySet.All(ctx, q.db.Doc)
	}

	return q.querySet.All(ctx, q.db.Doc)
}

func isSlice(model interface{}) bool {
	s := reflect.ValueOf(model)
	if s.Kind() == reflect.Ptr {
		s = s.Elem()
	}
	return s.Kind() == reflect.Slice
}
