## 访问控制

```go
atrz.WithRules(core.Rules{
			{Table: []string{"student"}, RouteTypes: []core.RouteType{core.RouteTypeList}, RuleType: core.RuleTypeAllow},
			{Table: []string{"*"}, RouteTypes: []core.RouteType{"*"}, RuleType: core.RuleTypeDeny},
}),
```

#### 结构
```go

type Rules []Rule
type Rule struct {
	Table      []string
	RouteTypes []RouteType
	RuleType   RuleType
}

```

#### 规则类型
```go
const (
	RuleTypeAllow RuleType = iota //允许
	RuleTypeDeny //禁止
)
```
RouteType 是支持的默认路由类型，在访问控制的时候会用到
```go
	RouteTypeInfo       RouteType = "INFO"
	RouteTypeList       RouteType = "LIST"
	RouteTypeCreate     RouteType = "CREATE"
	RouteTypeUpdate     RouteType = "UPDATE"
	RouteTypeDelete     RouteType = "DELETE"
	RouteTypeCrateBatch RouteType = "CREATE_BATCH"
```

访问控制 通过atrz.WithRules配置，可以配置多条，从上到下依次匹配 <br>
`Table:[]string{"*"}` 应用到所有表<br>
`RouteTypes: []core.RouteType{"*"}` 应用到所有路由<br>
在不设置的时候会默认支持所有表所有路由，如果不想开启所有表建议在最后一条配置 `{Table: []string{"*"}, RouteTypes: []core.RouteType{"*"}, RuleType: core.RuleTypeDeny}`