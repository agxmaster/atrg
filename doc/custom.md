### 自定义路由

参考 examples/custom_route

自定义路由访问时需要增加固定默认的path前缀 用来匹配分组

```go
const CustomRouterPrefix = "/custom"
```

> 示例配置中 CanGroupColumn 表明此路由支持已哪些字段进行分组，请求中传入_group即可使用，Aggregations什么需要查询的聚合列 ，需要注意使用分组后根据sql的语法不在分组中的字段无法单独查询。 <br>
> CustomRoute.Handler 未指明的情况下会根据RouteType匹配默认的handler 当然这里也可以写自定义的了。

```go
func ConfigCustomModel() map[string][]core.MpModel {
	return map[string][]core.MpModel{
		"student": {
			{
				Model:          reflect.TypeOf(model.Student{}),
				Select:         []string{"class"},
				Aggregations:   []core.Aggregation{{core.AggregationTypeCount, "age", "count_name"}},
				CanGroupColumn: []string{"name", "class"},
				MatchRouters: []core.CustomRoute{
					{
						RouteType: core.RouteTypeList,
						RoutePath: "/group/list/:model/:page",
						Method:    core.MethodGet,
					},
				},
			},
		},
	}
}

```
初始化时传入ConfigCustomModel 方法

```go
atrg.SetUp(dal.DB, g, atrg.WithCustomRoute(func() map[string][]core.MpModel {
    return defined.ConfigCustomModel()
}))

```


#### 示例请求


```go
curl http://localhost:8888/custom/group/list/student/1\?_group\=class
```

```json
{
    "code":0,
    "data":{
        "total":3,
        "data":[
            {
                "class":"1",
                "count_name":2
            },
            {
                "class":"2",
                "count_name":3
            },
            {
                "class":"3",
                "count_name":1
            }
        ]
    },
    "message":""
}
```

```go
curl http://localhost:8888/custom/group/list/student/1\?_group\=name,class
```

```json
{
  "code":0,
  "data":{
    "total":6,
    "data":[
      {
        "class":"1",
        "count_name":3
      },
      {
        "class":"2",
        "count_name":1
      },
      {
        "class":"1",
        "count_name":1
      },
      {
        "class":"2",
        "count_name":1
      },
      {
        "class":"3",
        "count_name":1
      },
      {
        "class":"2",
        "count_name":1
      }
    ]
  },
  "message":""
}
```

### 自定义handler

参考 examples/custom_handler <br>
自定义handler 是通过go结构体的继承来覆写想要 自定义的默认类型的路由

```go

var CusAtr Custom

type Custom struct {
    *handler.Atr
}


func (a Custom) List(c *gin.Context) {
response.Success(c, "Custom")
}

```

示例中默认hadler List被覆写，其他的handler都不受影响<br>
初始化时传入Custom 结构 
```go
	atrg.SetUp(dal.DB, g, atrg.WithCustomHandler(defined.Custom{}))
```

#### 示例请求


```go
curl http://localhost:8888/list/student/1

```

```json
{"code":0,"data":"Custom","message":""}
```