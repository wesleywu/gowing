package mongodb

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
)

// ConvertDataForRecord is a very important function, which does converting for any data that
// will be inserted into table/collection as a record.
//
// The parameter `value` should be type of *map/map/*struct/struct.
// It supports embedded struct definition for struct.
func ConvertDataForRecord(ctx context.Context, value interface{}) (map[string]interface{}, error) {
	var (
		err  error
		data = DataToMapDeep(value)
	)
	for k, v := range data {
		data[k], err = ConvertDataForRecordValue(ctx, v)
		if err != nil {
			return nil, gerror.Wrapf(err, `ConvertDataForRecordValue failed for value: %#v`, v)
		}
	}
	return data, nil
}

func ConvertDataForRecordValue(ctx context.Context, value interface{}) (interface{}, error) {
	var (
		err            error
		convertedValue = value
	)
	// If `value` implements interface `driver.Valuer`, it then uses the interface for value converting.
	if valuer, ok := value.(driver.Valuer); ok {
		if convertedValue, err = valuer.Value(); err != nil {
			if err != nil {
				return nil, err
			}
		}
		return convertedValue, nil
	}
	// Default value converting.
	var (
		rvValue = reflect.ValueOf(value)
		rvKind  = rvValue.Kind()
	)
	for rvKind == reflect.Ptr {
		rvValue = rvValue.Elem()
		rvKind = rvValue.Kind()
	}
	switch rvKind {
	case reflect.Slice, reflect.Array, reflect.Map:
		// It should ignore the bytes type.
		if _, ok := value.([]byte); !ok {
			// Convert the value to JSON.
			convertedValue, err = json.Marshal(value)
			if err != nil {
				return nil, err
			}
		}

	case reflect.Struct:
		switch r := value.(type) {
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		case time.Time:
			if r.IsZero() {
				convertedValue = nil
			}

		case gtime.Time:
			if r.IsZero() {
				convertedValue = nil
			} else {
				convertedValue = r.Time
			}

		case *gtime.Time:
			if r.IsZero() {
				convertedValue = nil
			} else {
				convertedValue = r.Time
			}

		case *time.Time:
			// Nothing to do.

		default:
			// Use string conversion in default.
			if s, ok := value.(iString); ok {
				convertedValue = s.String()
			} else {
				// Convert the value to JSON.
				convertedValue, err = json.Marshal(value)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return convertedValue, nil
}
