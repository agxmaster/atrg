package defined

import (
	"github.com/agxmaster/atrg/core"
	"github.com/agxmaster/atrg/examples/bind_model/model"
	"reflect"
)

type StudentCustomScope struct {
	MinAge int `json:"minAge"`
	MaxAge int `json:"maxAge"`
}

func ConfigModelMap() map[string]core.MpModel {
	return map[string]core.MpModel{
		"student": {
			Model: reflect.TypeOf(model.Student{}),
		},
	}
}
