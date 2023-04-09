package mongodb

import (
	"context"
	"fmt"
	"github.com/WesleyWu/gowing/util/gworm/mongodb/internal/reflection"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
)

type Model struct {
	db             Core
	tablesInit     string
	collection     *mongo.Collection
	data           interface{}   // Data for operation, which can be type of map/[]map/struct/*struct/string, etc.
	extraArgs      []interface{} // Extra custom arguments for sql, which are prepended to the arguments before sql committed to underlying driver.
	whereBuilder   *WhereBuilder // Condition builder for where operation.
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

func NewModel(collection *mongo.Collection) *Model {
	return &Model{
		collection: collection,
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

func (m *Model) Where(where interface{}, args ...interface{}) *Model {
	// todo
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

func (m *Model) Count(ctx context.Context, where ...interface{}) (int64, error) {
	return m.collection.CountDocuments(ctx, where)
}

func (m *Model) InsertOne(ctx context.Context, document interface{},
	opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return m.collection.InsertOne(ctx, document, opts...)
}

func (m *Model) Update(ctx context.Context, dataAndWhere ...interface{}) (*mongo.UpdateResult, error) {
	if len(dataAndWhere) > 0 {
		if len(dataAndWhere) > 2 {
			return m.Data(ctx, dataAndWhere[0]).Where(dataAndWhere[1], dataAndWhere[2:]...).Update(ctx)
		} else if len(dataAndWhere) == 2 {
			return m.Data(ctx, dataAndWhere[0]).Where(dataAndWhere[1]).Update(ctx)
		} else {
			return m.Data(ctx, dataAndWhere[0]).Update(ctx)
		}
	}
	if m.data == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "updating table with empty data")
	}
	var (
		updateData                                    = m.data
		reflectInfo                                   = reflection.OriginTypeAndKind(updateData)
		fieldNameUpdate                               = m.getSoftFieldNameUpdated()
		conditionWhere, conditionExtra, conditionArgs = m.formatCondition(ctx, false, false)
		conditionStr                                  = conditionWhere + conditionExtra
		err                                           error
	)
	if m.unscoped {
		fieldNameUpdate = ""
	}

	switch reflectInfo.OriginKind {
	case reflect.Map, reflect.Struct:
		var dataMap map[string]interface{}
		dataMap, err = ConvertDataForRecord(ctx, m.data)
		if err != nil {
			return nil, err
		}
		// Automatically update the record updating time.
		if fieldNameUpdate != "" {
			dataMap[fieldNameUpdate] = gtime.Now().String()
		}
		updateData = dataMap

	default:
		updates := gconv.String(m.data)
		// Automatically update the record updating time.
		if fieldNameUpdate != "" {
			if fieldNameUpdate != "" && !gstr.Contains(updates, fieldNameUpdate) {
				updates += fmt.Sprintf(`,%s='%s'`, fieldNameUpdate, gtime.Now().String())
			}
		}
		updateData = updates
	}
	newData, err := m.filterDataForInsertOrUpdate(updateData)
	if err != nil {
		return nil, err
	}

	if !gstr.ContainsI(conditionStr, " WHERE ") {
		intlog.Printf(
			ctx,
			`sql condition string "%s" has no WHERE for UPDATE operation, fieldNameUpdate: %s`,
			conditionStr, fieldNameUpdate,
		)
		return nil, gerror.NewCode(
			gcode.CodeMissingParameter,
			"there should be WHERE condition statement for UPDATE operation",
		)
	}
	return m.collection.UpdateMany(ctx, m.filter, req, opts)
}
