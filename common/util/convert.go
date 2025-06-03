package util

import (
	"encoding/json"
	"errors"
	"github.com/jinzhu/copier"
)

func Convert(from, to any) error {
	return copier.CopyWithOption(to, from, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
		Converters: []copier.TypeConverter{
			StringsToString(),
			StringToStrings(),
		},
	})
}

func StringsToString() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: []string{},
		DstType: copier.String,
		Fn: func(src interface{}) (dst interface{}, err error) {
			val, ok := src.([]string)
			if !ok {
				return nil, errors.New("src type is not []string")
			}
			marshaledData, err := json.Marshal(val)
			if err != nil {
				return nil, errors.New("failed to marshal []string")
			}
			return string(marshaledData), nil
		},
	}
}

func StringToStrings() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: copier.String,
		DstType: []string{},
		Fn: func(src interface{}) (dst interface{}, err error) {
			val, ok := src.(string)
			if !ok {
				return nil, errors.New("src type is not string")
			}
			var unmarshaledData []string
			err = json.Unmarshal([]byte(val), &unmarshaledData)
			if err != nil {
				return nil, errors.New("failed to unmarshal string")
			}
			return unmarshaledData, nil
		},
	}
}
