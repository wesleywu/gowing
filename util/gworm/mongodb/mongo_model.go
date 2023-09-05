package mongodb

import (
	"context"
	"reflect"

	"github.com/WesleyWu/gowing/util/gworm/mongodb/internal/reflection"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Model struct {
	tablesInit     string
	collection     *mongo.Collection
	data           interface{}   // Data for operation, which can be type of map/[]map/struct/*struct/string, etc.
	extraArgs      []interface{} // Extra custom arguments for sql, which are prepended to the arguments before sql committed to underlying driver.
	filter         bson.D
	fieldsIncluded *gset.StrSet
	fieldsExcluded *gset.StrSet
	groupBy        string        // Used for "group by" statement.
	orderBy        string        // Used for "order by" statement.
	having         []interface{} // Used for "having..." statement.
	start          int           // Used for "select ... start, limit ..." statement.
	limit          int           // Used for "select ... start, limit ..." statement.
	option         int           // Option for extra operation features.
	offset         int           // Offset statement for some databases grammar.
	safe           bool
}

const (
	linkTypeMaster           = 1
	linkTypeSlave            = 2
	defaultFields            = "*"
	whereHolderOperatorWhere = 1
	whereHolderOperatorAnd   = 2
	whereHolderOperatorOr    = 3
	whereHolderTypeDefault   = "Default"
	whereHolderTypeNoArgs    = "NoArgs"
	whereHolderTypeIn        = "In"
)

func NewModel(ctx context.Context, collectionName string) *Model {
	return &Model{
		collection: GetCollection(ctx, collectionName),
		filter:     bson.D{},
	}
}

// Fields appends `fields` to the operation fields of the model
//
// Also see FieldsEx.
//
// reference: https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/project/#include-a-field
func (m *Model) Fields(fields ...string) *Model {
	if len(fields) == 0 {
		return m
	}
	if m.fieldsIncluded == nil {
		m.fieldsIncluded = gset.NewStrSet()
	}
	m.fieldsIncluded.Add(fields...)
	return m
}

// FieldsEx appends `fields` to the excluded operation fields of the model
//
// Also see Fields.
//
// reference: https://www.mongodb.com/docs/drivers/go/current/fundamentals/crud/read-operations/project/#exclude-a-field
func (m *Model) FieldsEx(fields ...string) *Model {
	if len(fields) == 0 {
		return m
	}
	if m.fieldsExcluded == nil {
		m.fieldsExcluded = gset.NewStrSet()
	}
	m.fieldsExcluded.Add(fields...)
	return m
}

func (m *Model) WhereEq(key string, value interface{}) *Model {
	m.filter = append(m.filter, bson.E{
		Key:   key,
		Value: value,
	})
	return m
}

func (m *Model) WhereNE(key string, value interface{}) *Model {
	m.filter = append(m.filter, bson.E{
		Key:   key,
		Value: bson.M{"$ne": value},
	})
	return m
}

func (m *Model) WhereGT(key string, value interface{}) *Model {
	m.filter = append(m.filter, bson.E{
		Key:   key,
		Value: bson.M{"$gt": value},
	})
	return m
}

func (m *Model) WhereGTE(key string, value interface{}) *Model {
	m.filter = append(m.filter, bson.E{
		Key:   key,
		Value: bson.M{"$gte": value},
	})
	return m
}

func (m *Model) WhereLT(key string, value interface{}) *Model {
	m.filter = append(m.filter, bson.E{
		Key:   key,
		Value: bson.M{"$lt": value},
	})
	return m
}

func (m *Model) WhereLTE(key string, value interface{}) *Model {
	m.filter = append(m.filter, bson.E{
		Key:   key,
		Value: bson.M{"$lte": value},
	})
	return m
}

func (m *Model) WhereIn(key string, value ...interface{}) *Model {
	m.filter = append(m.filter, bson.E{
		Key:   key,
		Value: bson.M{"$in": value},
	})
	return m
}

func (m *Model) WhereNotIn(key string, value ...interface{}) *Model {
	m.filter = append(m.filter, bson.E{
		Key:   key,
		Value: bson.M{"$nin": value},
	})
	return m
}

func (m *Model) WhereBetween(key string, min, max interface{}) *Model {
	m.filter = append(m.filter, bson.E{
		Key:   key,
		Value: bson.M{"$gte": min, "$lte": max},
	})
	return m
}

func (m *Model) WhereNotBetween(key string, min, max interface{}) *Model {
	m.filter = append(m.filter, bson.E{
		Key: key,
		Value: bson.D{
			{"$or",
				bson.A{
					bson.D{{key, bson.D{{"$gt", max}}}},
					bson.D{{key, bson.D{{"$lt", min}}}},
				},
			},
		},
	})
	return m
}

func (m *Model) WhereLike(key string, like string) *Model {
	m.filter = append(m.filter, bson.E{
		Key:   key,
		Value: bson.M{"$regex": like, "$options": "im"},
	})
	return m
}

func (m *Model) WhereNotLike(key string, like string) *Model {
	m.filter = append(m.filter, bson.E{
		Key: key,
		Value: bson.M{
			"$not": bson.M{"$regex": like, "$options": "im"},
		},
	})
	return m
}

func (m *Model) WhereNull(key ...string) *Model {
	for _, oneKey := range key {
		m.filter = append(m.filter, bson.E{
			Key:   oneKey,
			Value: nil,
		})
	}
	return m
}

func (m *Model) WhereNotNull(key ...string) *Model {
	for _, oneKey := range key {
		m.filter = append(m.filter, bson.E{
			Key:   oneKey,
			Value: bson.M{"$ne": nil},
		})
	}
	return m
}

// WherePri does the same logic as Model.Where except that if the parameter `where`
// is a single condition like int/string/float/slice, it treats the condition as the primary
// key value. That is, if primary key is "id" and given `where` parameter as "123", the
// WherePri function treats the condition as "id=123", but Model.Where treats the condition
// as string "123".
func (m *Model) WherePri(args []string) *Model {
	lenArgs := len(args)
	switch lenArgs {
	case 0:
		return m
	case 1:
		m.filter = append(m.filter, bson.E{Key: "_id", Value: args[0]})
	default:
		m.filter = append(m.filter, bson.E{Key: "_id", Value: bson.E{Key: "$in", Value: args}})
	}
	return m
}

func (m *Model) Count(ctx context.Context) (int64, error) {
	return m.collection.CountDocuments(ctx, m.filter, nil)
}

func (m *Model) Scan(ctx context.Context, pointer interface{}) error {
	reflectInfo := reflection.OriginTypeAndKind(pointer)
	if reflectInfo.InputKind != reflect.Ptr {
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`the parameter "pointer" for function Scan should type of pointer`,
		)
	}
	switch reflectInfo.OriginKind {
	case reflect.Slice, reflect.Array:
		return m.doStructs(ctx, pointer)

	case reflect.Struct, reflect.Invalid:
		return m.doStruct(ctx, pointer)

	default:
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`element of parameter "pointer" for function Scan should type of struct/*struct/[]struct/[]*struct`,
		)
	}
}

func (m *Model) doStruct(ctx context.Context, pointer interface{}) error {
	cursor, err := m.collection.Find(ctx, m.filter)
	if err != nil {
		return err
	}
	if cursor.TryNext(ctx) {
		err = cursor.Decode(pointer)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (m *Model) doStructs(ctx context.Context, pointer interface{}) error {
	cursor, err := m.collection.Find(ctx, m.filter)
	if err != nil {
		return err
	}
	err = cursor.All(ctx, pointer)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.collection.InsertOne(ctx, document, opts...)
}

func (m *Model) Save(ctx context.Context, document interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	update := bson.D{
		{
			"$set", document,
		},
	}
	return m.collection.UpdateOne(ctx, m.filter, update, opts...)
}

func (m *Model) Delete(ctx context.Context,
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return m.collection.DeleteMany(ctx, m.filter, opts...)
}
