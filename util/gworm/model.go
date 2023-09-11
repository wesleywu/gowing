/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gworm

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/wesleywu/gowing/util/gworm/internal"
	"github.com/wesleywu/gowing/util/gworm/mongodb"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModelType uint8

const (
	// GF_ORM GoFrame ORM
	GF_ORM ModelType = iota
	// MONGO MongoDB via native driver
	MONGO
)

type Model struct {
	Type       ModelType
	GfModel    *gdb.Model
	MongoModel *mongodb.Model
}

type Result struct {
	Type           ModelType
	SqlResult      sql.Result
	MongoResult    *mongodb.MongoResult
	LastInsertedId string
	RowsAffected   int64
}

func (m *Model) WithAll() *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WithAll()
		return m
	case MONGO:
		return m
	}
	return m
}

func (m *Model) Fields(fieldNamesOrMapStruct ...interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.Fields(fieldNamesOrMapStruct...)
		return m
	case MONGO:
		fields := internal.GetFields(fieldNamesOrMapStruct...)
		m.MongoModel = m.MongoModel.Fields(fields...)
		return m
	}
	return m
}

func (m *Model) FieldsEx(fieldNamesOrMapStruct ...interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.FieldsEx(fieldNamesOrMapStruct...)
		return m
	case MONGO:
		fields := internal.GetFields(fieldNamesOrMapStruct...)
		m.MongoModel = m.MongoModel.FieldsEx(fields...)
		return m
	}
	return m
}

func (m *Model) Where(where interface{}, args ...interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.Where(where, args...)
		return m
	case MONGO:
		m.MongoModel.WhereEq(gconv.String(where), args[0])
		return m
	}
	return m
}

func (m *Model) WhereNot(column string, value interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereNot(column, value)
		return m
	case MONGO:
		m.MongoModel.WhereNE(column, value)
		return m
	}
	return m
}

func (m *Model) WherePri(where interface{}, args ...interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WherePri(where, args...)
		return m
	case MONGO:
		newWhere := gconv.Strings(where)
		m.MongoModel = m.MongoModel.WherePri(newWhere)
		return m
	}
	return m
}

func (m *Model) WhereGT(column string, value interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereGT(column, value)
		return m
	case MONGO:
		m.MongoModel.WhereGT(column, value)
		return m
	}
	return m
}

func (m *Model) WhereGTE(column string, value interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereGTE(column, value)
		return m
	case MONGO:
		m.MongoModel.WhereGTE(column, value)
		return m
	}
	return m
}

func (m *Model) WhereLT(column string, value interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereLT(column, value)
		return m
	case MONGO:
		m.MongoModel.WhereLT(column, value)
		return m
	}
	return m
}

func (m *Model) WhereLTE(column string, value interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereLTE(column, value)
		return m
	case MONGO:
		m.MongoModel.WhereLTE(column, value)
		return m
	}
	return m
}

func (m *Model) WhereBetween(column string, min, max interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereBetween(column, min, max)
		return m
	case MONGO:
		m.MongoModel.WhereBetween(column, min, max)
		return m
	}
	return m
}

func (m *Model) WhereNotBetween(column string, min, max interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereNotBetween(column, min, max)
		return m
	case MONGO:
		m.MongoModel.WhereNotBetween(column, min, max)
		return m
	}
	return m
}

func (m *Model) WhereIn(column string, in []interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereIn(column, in)
		return m
	case MONGO:
		m.MongoModel.WhereIn(column, in)
		return m
	}
	return m
}

func (m *Model) WhereNotIn(column string, in []interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereNotIn(column, in)
		return m
	case MONGO:
		m.MongoModel.WhereNotIn(column, in)
		return m
	}
	return m
}

func (m *Model) WhereLike(column string, like string) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereLike(column, like)
		return m
	case MONGO:
		m.MongoModel.WhereLike(column, like)
		return m
	}
	return m
}

func (m *Model) WhereNotLike(column string, like string) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereNotLike(column, like)
		return m
	case MONGO:
		m.MongoModel.WhereNotLike(column, like)
		return m
	}
	return m
}

func (m *Model) WhereNull(column ...string) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereNull(column...)
		return m
	case MONGO:
		m.MongoModel.WhereNull(column...)
		return m
	}
	return m
}

func (m *Model) WhereNotNull(column ...string) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.WhereNotNull(column...)
		return m
	case MONGO:
		m.MongoModel.WhereNotNull(column...)
		return m
	}
	return m
}

//
//func (m *Model) Data(data ...interface{}) *Model {
//	switch m.Type {
//	case GF_ORM:
//		m.GfModel = m.GfModel.Data(data...)
//		return m
//	case MONGO:
//		return m
//	}
//	return m
//}

func (m *Model) Page(page, limit int) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.Page(page, limit)
		return m
	case MONGO:
		return m
	}
	return m
}

func (m *Model) Order(orderBy ...interface{}) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.Order(orderBy...)
		return m
	case MONGO:
		return m
	}
	return m
}

func (m *Model) Limit(limit ...int) *Model {
	switch m.Type {
	case GF_ORM:
		m.GfModel = m.GfModel.Limit(limit...)
		return m
	case MONGO:
		return m
	}
	return m
}

func (m *Model) Scan(ctx context.Context, pointer interface{}, where ...interface{}) error {
	switch m.Type {
	case GF_ORM:
		return m.GfModel.Scan(pointer, where)
	case MONGO:
		return m.MongoModel.Scan(ctx, pointer)
	}
	return errors.New("Not supported model type: " + string(m.Type))
}

func (m *Model) Count(ctx context.Context) (int64, error) {
	switch m.Type {
	case GF_ORM:
		count, err := m.GfModel.Count()
		return int64(count), err
	case MONGO:
		return m.MongoModel.Count(ctx)
	}
	return 0, errors.New("Not supported model type: " + string(m.Type))
}

func (m *Model) InsertOne(ctx context.Context, req interface{}) (*Result, error) {
	switch m.Type {
	case GF_ORM:
		result, err := m.GfModel.Insert(req)
		return &Result{
			Type:      m.Type,
			SqlResult: result,
		}, err
	case MONGO:
		result, err := m.MongoModel.InsertOne(ctx, req)
		if err != nil {
			return nil, err
		}
		return &Result{
			Type:           m.Type,
			LastInsertedId: gconv.String(result.InsertedID),
			RowsAffected:   1,
			MongoResult: &mongodb.MongoResult{
				InsertedID:   gconv.String(result.InsertedID),
				AffectedRows: 1,
			},
		}, err
	}
	return nil, errors.New("Not supported model type: " + string(m.Type))
}

func (m *Model) Update(ctx context.Context, data interface{}) (*Result, error) {
	switch m.Type {
	case GF_ORM:
		result, err := m.GfModel.Save(data)
		return &Result{
			Type:        m.Type,
			SqlResult:   result,
			MongoResult: nil,
		}, err
	case MONGO:
		result, err := m.MongoModel.Save(ctx, data)
		if err != nil {
			return nil, err
		}
		return &Result{
			Type:         m.Type,
			RowsAffected: result.ModifiedCount,
			MongoResult: &mongodb.MongoResult{
				AffectedRows: result.ModifiedCount,
			},
		}, nil
	}
	return nil, errors.New("Not supported model type: " + string(m.Type))
}

func (m *Model) Upsert(ctx context.Context, data interface{}) (*Result, error) {
	switch m.Type {
	case GF_ORM:
		result, err := m.GfModel.Save(data)
		return &Result{
			Type:        m.Type,
			SqlResult:   result,
			MongoResult: nil,
		}, err
	case MONGO:
		result, err := m.MongoModel.Save(ctx, data, options.Update().SetUpsert(true))
		if err != nil {
			return nil, err
		}
		return &Result{
			Type:         m.Type,
			RowsAffected: result.ModifiedCount,
			MongoResult: &mongodb.MongoResult{
				AffectedRows: result.ModifiedCount,
			},
		}, nil
	}
	return nil, errors.New("Not supported model type: " + string(m.Type))
}

func (m *Model) Delete(ctx context.Context, where ...interface{}) (*Result, error) {
	switch m.Type {
	case GF_ORM:
		result, err := m.GfModel.Delete(where...)
		return &Result{
			Type:        m.Type,
			SqlResult:   result,
			MongoResult: nil,
		}, err
	case MONGO:
		result, err := m.MongoModel.Delete(ctx)
		if err != nil {
			return nil, err
		}
		return &Result{
			Type:         m.Type,
			RowsAffected: result.DeletedCount,
			MongoResult: &mongodb.MongoResult{
				AffectedRows: result.DeletedCount,
			},
		}, err
	}
	return nil, errors.New("Not supported model type: " + string(m.Type))
}

//
//func (r *Result) InsertedId() (string, error) {
//	switch r.Type {
//	case GF_ORM:
//		id, err := r.SqlResult.LastInsertId()
//		return strconv.FormatInt(id, 10), err
//	case MONGO:
//		id := gconv.String(r.MongoResult.InsertedID)
//		return id, nil
//	}
//	return "", errors.New("Not supported model type: " + string(r.Type))
//}
//
//func (r *Result) LastInsertedIdInt64() (int64, error) {
//	switch r.Type {
//	case GF_ORM:
//		id, err := r.SqlResult.LastInsertId()
//		return id, err
//	case MONGO:
//		id := gconv.Int64(r.MongoResult.InsertedID)
//		return id, nil
//	}
//	return 0, errors.New("Not supported model type: " + string(r.Type))
//}
//
//func (r *Result) RowsAffected() (int64, error) {
//	switch r.Type {
//	case GF_ORM:
//		rows, err := r.SqlResult.RowsAffected()
//		return rows, err
//	case MONGO:
//		return 0, nil
//	}
//	return 0, errors.New("Not supported model type: " + string(r.Type))
//}
