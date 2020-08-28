package redis

import (
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
)

var (
	ScanTagName        = "db"
	DefaultLocation, _ = time.LoadLocation("Asia/Shanghai")
)

var DecodeHook = func(from reflect.Type, to reflect.Type, v interface{}) (interface{}, error) {
	if from.Kind() == reflect.String && to.Kind() == reflect.Struct {
		if to.String() == "time.Time" || to.String() == "*time.Time" {
			ss := v.(string)
			t := time.Time{}
			err := t.UnmarshalBinary([]byte(ss))
			if err != nil {
				var layout = "2006-01-02 15:04:05"
				if strings.Index(ss, "T") != -1 {
					layout = time.RFC3339
				}
				t, err = time.ParseInLocation(layout, ss, DefaultLocation)
				if err != nil {
					return nil, err
				}
			}
			return t, nil
		}
	}

	return v, nil
}

func scanStruct(src map[string]string, dest interface{}) error {
	config := &mapstructure.DecoderConfig{
		Result:           dest,
		TagName:          ScanTagName,
		WeaklyTypedInput: true,
		DecodeHook:       DecodeHook,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(src)
}

func structToMapInterface(m interface{}) map[string]interface{} {
	v := reflect.Indirect(reflect.ValueOf(m))
	t := v.Type()
	rs := make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		n := t.Field(i).Tag.Get(ScanTagName)
		switch s := f.Interface().(type) {
		case time.Time:
			rs[n] = s.Format("2006-01-02 15:04:05")
		default:
			rs[n] = f.Interface()
		}
	}

	return rs
}
