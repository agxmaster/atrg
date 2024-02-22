
## 绑定模型

### 自定义参数解析

参考example/bind_model/defined/atr_conf.go ， 并创建examples/bind_model/model/student.go gorm的模型文件。

```go

func ConfigModelMap() map[string]core.MpModel {
	return map[string]core.MpModel{
		"student": {Model: reflect.TypeOf(model.Student{}),},
	}
}

```

main中绑定模型配置
```go
func main() {
	z := server.Default()
	dal.MysqlSetup()
	atrg.SetUp(dal.DB, z, defined.ConfigModelMap, nil, nil)

	z.Spin()
}


```

* 接收参数、查询数据、返回数据 用到的都是定义到model中的结构。
* 绑定模型后就可以用用到gorm注解参考 https://gorm.io/zh_CN/docs/models.html 可以在插入数据时做一些校验、默认值、列名别名等特性。
* 可以用到 gin注解做参数校验，可以扩展注解校验参考 https://gin-gonic.com/zh-cn/docs/examples/binding-and-validation/

