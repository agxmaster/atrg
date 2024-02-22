
## 配置用法参考


参考examples/standard/defined/atrConf.go

```go

package defined

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/agxmaster/atm/clause"
	"github.com/agxmaster/atrg/core"
	"github.com/agxmaster/atrg/examples/simple/model"
	"github.com/agxmaster/atrg/util"
	"gorm.io/gorm"
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
			ColumnMapping: map[string]string{
				"gender": "sex",
			},
			Select: []string{"name"},
			//Select: []string{"name", "gender", "class", "created_at", "updated_at"},
			Hidden:              []string{"updated_at"},
			AddColumns:          StudentAddColumns,
			ListWithCustomScope: StudentListWithCustomScope,
			CreateParamsHandler: StudentCreateParamsHandler,
			UpdateParamsHandler: StudentUpdateParamsHandler,
			Aggregations:        []core.Aggregation{{core.AggregationTypeCount, "age", "count_name"}},
			CanGroupColumn:      []string{"name", "class"},
		},
	}
}

func StudentListWithCustomScope(ctx context.Context, customJson []byte) (clause.Scope, error) {
	var studentCustomStruct StudentCustomScope
	err := json.Unmarshal(customJson, &studentCustomStruct)
	if err != nil {
		return nil, err
	}

	if studentCustomStruct.MinAge >= studentCustomStruct.MaxAge {
		return nil, errors.New("min age can't >= max age")
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where("age >= ?", studentCustomStruct.MinAge).Where("age <= ?", studentCustomStruct.MaxAge)
	}, nil
}

func StudentUpdateParamsHandler(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	modelParams, err := util.InterfaceToAny[model.Student](params)
	if err != nil {
		return nil, err
	}
	if modelParams.Age < 5 {
		return nil, errors.New("age 5 not allowed")
	} else if modelParams.Age < 6 {
		modelParams.Age = 8
	}
	newParams, err := util.InterfaceToAny[map[string]interface{}](interface{}(*modelParams))
	return *newParams, err
}

func StudentCreateParamsHandler(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	mapParams, err := util.InterfaceToAny[map[string]interface{}](params)
	if err != nil {
		return nil, err
	}

	if age, ok := (*mapParams)["age"]; ok {
		if ageInt, ok := age.(int); ok {
			if ageInt < 5 {
				return nil, errors.New("age 5 not allowed")
			} else if ageInt < 6 {
				(*mapParams)["age"] = 8
			}
		}
	}
	return *mapParams, nil
}

var StudentAddColumns = []map[string]core.ColumnAddFunc{
	{
		"className": func(ctx context.Context, line map[string]interface{}, addKey string) interface{} {
			if class, ok := line["class"]; ok {
				if classId, ok := class.(string); ok {
					if classId == "3" {
						return "3班"
					}
				}
			}
			return "其他班级"
		},
		"name": func(ctx context.Context, line map[string]interface{}, addKey string) interface{} {
			if name, ok := line["name"]; ok {
				if nameCode, ok := name.(string); ok {
					if nameCode == "2" {
						return "小2"
					}
				}
			}
			return line["name"]
		},
	},
}


```

初始化时传入
```go
atrg.SetUp(dal.DB, g, atrg.WithModelConfig(func() map[string]core.MpModel {
    return defined.ConfigModelMap()
}))

```


#### 自定义参数解析
配置: ListWithCustomScope <br >

定义解析到的结构体StudentCustomScope
```go



type StudentCustomScope struct {
    MinAge int `json:"minAge"`
    MaxAge int `json:"maxAge"`
}

func StudentListWithCustomScope(ctx context.Context, customJson []byte) (clause.Scope, error) {
				var studentCustomStruct StudentCustomScope
				err := json.Unmarshal(customJson, &studentCustomStruct)
				if err != nil {
					return nil, err
				}

				if studentCustomStruct.MinAge >= studentCustomStruct.MaxAge {
					return nil, errors.New("min age can't >= max age")
				}

				return func(db *gorm.DB) *gorm.DB {
					return db.Where("age >= ?", studentCustomStruct.MinAge).Where("age <= ?", studentCustomStruct.MaxAge)
				}, nil
}
```

请求示例 `_custom={"minAge":8,"maxAge":11}`
```go
curl --location --globoff 'http://localhost:8888/list/student/1?_order=id%2Cdesc&_order=name%2Cdesc&_custom={%22minAge%22%3A8%2C%22maxAge%22%3A11}'

```

普通查询的key=value这种结构只支持等于操作 复杂的and/or/>=/<=等操作没有支持 ，此时就需要用到_custom参数

#### 指定查询字段

```go
    Select:              []string{"name", "gender", "class", "created_at", "updated_at"},

```
#### 隐藏的字段

隐藏的字段，可以指定了查询字段后再进行隐藏，意义时这个字段可以在其他计算中使用
```go
			Hidden:              []string{"updated_at"},
```

#### 展示字段扩充

`line map[string]interface{}` 参数是拿到的结果集，可以用来增加key、修改key的值，比如登陆的用户想要展示用户名其实可以通过在middleware中，把username写入到context中，然后在这里读出写入结果集，也可以进行其他表的查询把结果集拼接起来，也可以通过定义map对配置id映射到name进行展示

```go
var StudentAddColumns = []map[string]core.ColumnAddFunc{
{
    "className": func(ctx context.Context, line map[string]interface{}, addKey string) interface{} {
        if class, ok := line["class"]; ok {
            if classId, ok := class.(string); ok {
                if classId == "3" {
                    return "3班"
                }
                }
            }
        return "其他班级"
    },
    "name": func(ctx context.Context, line map[string]interface{}, addKey string) interface{} {
        if name, ok := line["name"]; ok {
            if nameCode, ok := name.(string); ok {
                if nameCode == "2" {
                    return "小2"
                }
            }
        }
        return line["name"]
        },
    },
}
```
#### 创建前数据完善
配置到: CreateParamsHandler, 在创建前对数据改造 可以加入一个默认的列的值，计算列的值、等操作
```go
func StudentCreateParamsHandler(ctx context.Context, params interface{}) (map[string]interface{}, error) {
	mapParams, err := util.InterfaceToAny[map[string]interface{}](params)
	if err != nil {
		return nil, err
	}

	if age, ok := (*mapParams)["age"]; ok {
		if ageInt, ok := age.(int); ok {
			if ageInt < 5 {
				return nil, errors.New("age 5 not allowed")
			} else if ageInt < 6 {
				(*mapParams)["age"] = 8
			}
		}
	}
	return *mapParams, nil
}

```

####  修改前数据完善
配置到 UpdateParamsHandler 使用方法同创建前数据完善
#### 批量创建前数据完善
配置到CreateBatchParamsHandler  使用方法同创建前数据完善

