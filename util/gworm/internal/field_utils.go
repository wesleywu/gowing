package internal

import (
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
	"reflect"
)

const (
	OrmTagForStruct = "bson"
	// regularFieldNameRegPattern is the regular expression pattern for a string
	// which is a regular field name of table.
	regularFieldNameRegPattern = `^[\w\.\-]+$`
)

func GetFields(fieldNamesOrMapStruct ...interface{}) []string {
	length := len(fieldNamesOrMapStruct)
	if length == 0 {
		return nil
	}
	switch {
	// String slice.
	case length >= 2:
		return gconv.Strings(fieldNamesOrMapStruct)
	// It needs type asserting.
	case length == 1:
		structOrMap := fieldNamesOrMapStruct[0]
		switch r := structOrMap.(type) {
		case string:
			return []string{r}
		case []string:
			return r
		default:
			return getFieldsFromStructOrMap(structOrMap)
		}
	default:
		return nil
	}
}

func getFieldsFromStructOrMap(structOrMap interface{}) (fields []string) {
	fields = []string{}
	if IsStruct(structOrMap) {
		structFields, _ := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         structOrMap,
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		})
		for _, structField := range structFields {
			if tag := structField.Tag(OrmTagForStruct); tag != "" && gregex.IsMatchString(regularFieldNameRegPattern, tag) {
				fields = append(fields, tag)
			} else {
				fields = append(fields, structField.Name())
			}
		}
	} else {
		fields = gutil.Keys(structOrMap)
	}
	return
}

// IsStruct checks whether `value` is type of struct.
func IsStruct(value interface{}) bool {
	var reflectType = reflect.TypeOf(value)
	if reflectType == nil {
		return false
	}
	var reflectKind = reflectType.Kind()
	for reflectKind == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectKind = reflectType.Kind()
	}
	switch reflectKind {
	case reflect.Struct:
		return true
	}
	return false
}
