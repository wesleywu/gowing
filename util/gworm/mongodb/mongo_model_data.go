package mongodb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gutil"
	"github.com/wesleywu/gowing/util/gworm/mongodb/internal/reflection"
)

// Data sets the operation data for the model.
// The parameter `data` can be type of string/map/gmap/slice/struct/*struct, etc.
// Note that, it uses shallow value copying for `data` if `data` is type of map/slice
// to avoid changing it inside function.
// Eg:
// Data("uid=10000")
// Data("uid", 10000)
// Data("uid=? AND name=?", 10000, "john")
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"}).
func (m *Model) Data(ctx context.Context, data ...interface{}) *Model {
	var (
		err error
	)
	if len(data) > 1 {
		if s := gconv.String(data[0]); gstr.Contains(s, "?") {
			m.data = s
			m.extraArgs = data[1:]
		} else {
			mp := make(map[string]interface{})
			for i := 0; i < len(data); i += 2 {
				mp[gconv.String(data[i])] = data[i+1]
			}
			m.data = mp
		}
		return m
	}
	onlyData := data[0]
	switch value := onlyData.(type) {
	case Result:
		m.data = value.List()

	case Record:
		m.data = value.Map()

	case List:
		list := make(List, len(value))
		for k, v := range value {
			list[k] = gutil.MapCopy(v)
		}
		m.data = list

	case Map:
		m.data = gutil.MapCopy(value)

	default:
		reflectInfo := reflection.OriginValueAndKind(value)
		switch reflectInfo.OriginKind {
		case reflect.Slice, reflect.Array:
			if reflectInfo.OriginValue.Len() > 0 {
				// If the `data` parameter is a DO struct,
				// it then adds `OmitNilData` option for this condition,
				// which will filter all nil parameters in `data`.
				if isDoStruct(reflectInfo.OriginValue.Index(0).Interface()) {
					m.OmitNilData()
					m.option |= optionOmitNilDataInternal
				}
			}
			list := make(List, reflectInfo.OriginValue.Len())
			for i := 0; i < reflectInfo.OriginValue.Len(); i++ {
				list[i], err = ConvertDataForRecord(ctx, reflectInfo.OriginValue.Index(i).Interface())
				if err != nil {
					panic(err)
				}
			}
			m.data = list

		case reflect.Struct:
			// If the `data` parameter is a DO struct,
			// it then adds `OmitNilData` option for this condition,
			// which will filter all nil parameters in `data`.
			if isDoStruct(value) {
				m.OmitNilData()
			}
			if v, ok := onlyData.(iInterfaces); ok {
				var (
					array = v.Interfaces()
					list  = make(List, len(array))
				)
				for i := 0; i < len(array); i++ {
					list[i], err = ConvertDataForRecord(ctx, array[i])
					if err != nil {
						panic(err)
					}
				}
				m.data = list
			} else {
				m.data, err = ConvertDataForRecord(ctx, onlyData)
				if err != nil {
					panic(err)
				}
			}

		case reflect.Map:
			m.data, err = ConvertDataForRecord(ctx, onlyData)
			if err != nil {
				panic(err)
			}

		default:
			m.data = onlyData
		}
	}
	return m
}

// isDoStruct checks and returns whether given type is a DO struct.
func isDoStruct(object interface{}) bool {
	// It checks by struct name like "XxxForDao", to be compatible with old version.
	// TODO remove this compatible codes in future.
	reflectType := reflect.TypeOf(object)
	if gstr.HasSuffix(reflectType.String(), modelForDaoSuffix) {
		return true
	}
	// It checks by struct meta for DO struct in version.
	if ormTag := gmeta.Get(object, OrmTagForStruct); !ormTag.IsEmpty() {
		match, _ := gregex.MatchString(
			fmt.Sprintf(`%s\s*:\s*([^,]+)`, OrmTagForDo),
			ormTag.String(),
		)
		if len(match) > 1 {
			return gconv.Bool(match[1])
		}
	}
	return false
}
