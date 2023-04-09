package mongodb

import (
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"regexp"
	"time"
)

// iString is the type assert api for String.
type iString interface {
	String() string
}

// iIterator is the type assert api for Iterator.
type iIterator interface {
	Iterator(f func(key, value interface{}) bool)
}

// iInterfaces is the type assert api for Interfaces.
type iInterfaces interface {
	Interfaces() []interface{}
}

// iTableName is the interface for retrieving table name fro struct.
type iTableName interface {
	TableName() string
}

const (
	OrmTagForStruct    = "orm"
	OrmTagForTable     = "table"
	OrmTagForWith      = "with"
	OrmTagForWithWhere = "where"
	OrmTagForWithOrder = "order"
	OrmTagForDo        = "do"
)

var (
	// quoteWordReg is the regular expression object for a word check.
	quoteWordReg = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

	// Priority tags for struct converting for orm field mapping.
	structTagPriority = append([]string{OrmTagForStruct}, gconv.StructTagPriority...)
)

// DataToMapDeep converts `value` to map type recursively(if attribute struct is embedded).
// The parameter `value` should be type of *map/map/*struct/struct.
// It supports embedded struct definition for struct.
func DataToMapDeep(value interface{}) map[string]interface{} {
	m := gconv.Map(value, structTagPriority...)
	for k, v := range m {
		switch v.(type) {
		case time.Time, *time.Time, gtime.Time, *gtime.Time, gjson.Json, *gjson.Json:
			m[k] = v

		default:
			// Use string conversion in default.
			if s, ok := v.(iString); ok {
				m[k] = s.String()
			} else {
				m[k] = v
			}
		}
	}
	return m
}
