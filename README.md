# atrg

### 简介
在gin框架中注册增删改查的路由，只需要增加一行代码即可拥有所有表的增删改查操作，列表查询支持字段名的别名、根据字段分组、分页、排序、自定义解析规则等便捷操作。

如examples/simple中main.go的代码
```go
func main() {
    g := gin.Default()
    dal.MysqlSetup()
    atrg.SetUp(dal.DB, g)
    g.Run(":8888")
}

```

### 创建测试用表

准备工作创建表atm.student 账号密码请修改example/*/dal/mysql.go 中的配置

```
var dsn = "root:12345678@tcp(localhost:3306)/atm?charset=utf8&parseTime=True&loc=Local"
```
```
create table student (
    id bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
    name varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' ,
    province varchar(128) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' ,
    gender int NOT NULL DEFAULT '0' ,
    age int NOT NULL DEFAULT '0',
    class int NOT NULL DEFAULT '0',
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
```

### 运行

#### 运行 examples/main.go main 方法


#### 创建接口

POST /atr/create/:model

创建5个student
```
  curl -X POST 127.0.0.1:8888/atr/create/student -H 'Content-Type: application/json' -d "{\"name\":\"Beatty\",\"gender\":1,\"age\":12,\"class\":\"1\",\"province\":\"hebei\"}"
  curl -X POST 127.0.0.1:8888/atr/create/student -H 'Content-Type: application/json' -d "{\"name\":\"Bell\",\"gender\":1,\"age\":8,\"class\":\"2\",\"province\":\"hunan\"}"
  curl -X POST 127.0.0.1:8888/atr/create/student -H 'Content-Type: application/json' -d "{\"name\":\"Abbey\",\"gender\":2,\"age\":8,\"class\":\"1\",\"province\":\"beijin\"}"
  curl -X POST 127.0.0.1:8888/atr/create/student -H 'Content-Type: application/json' -d "{\"name\":\"Adele\",\"gender\":2,\"age\":11,\"class\":\"2\",\"province\":\"beijin\"}"
  curl -X POST 127.0.0.1:8888/atr/create/student -H 'Content-Type: application/json' -d "{\"name\":\"Beckman\",\"gender\":1,\"age\":11,\"class\":\"2\",\"province\":\"henan\"}"
```


成功响应
```json
{"code":0,"data":null,"message":""}
```

失败响应
```json
{"code":-1,"data":null,"message":""}
```
##### list接口

Get /atr/list/:model/:page

```
curl http://localhost:8888/atr/list/student/2?class=1&_order=age,desc&_order=gender&_size=1
```

```json
{
  "code":0,
  "data":{
    "total":2,
    "data":[
      {
        "age":8,
        "class":1,
        "created_at":"2024-02-05T15:06:57+08:00",
        "gender":2,
        "id":3,
        "name":"Abbey",
        "province":"beijin",
        "updated_at":"2024-02-05T15:06:57+08:00"
      }
    ]
  },
  "message":""
}
```

##### info接口

GET /atr/info/:model/:id

```
curl http://localhost:8888/atr/info/student/1
```

```json
{
    "code":0,
    "data":{
        "age":12,
        "class":1,
        "created_at":"2024-02-05T15:06:57+08:00",
        "gender":1,
        "id":1,
        "name":"Beatty",
        "province":"hebei",
        "updated_at":"2024-02-05T15:06:57+08:00"
    },
    "message":""
}
```
POST /atr/update/:model

##### update接口
```
curl -X POST 127.0.0.1:8888/atr/update/student/5 -H 'Content-Type: application/json' -d "{\"name\":\"Sandra\",\"gender\":1,\"age\":13,\"class\":\"2\"}"
```

```
{"code":0,"data":null,"message":""}
```

##### delete 接口
POST /atr/delete/:model/:id

```
curl -X POST 127.0.0.1:8888/atr/delete/student/5
```

```
{"code":0,"data":null,"message":""}
```

#### 批量创建接口

POST /atr/batch/create/:model

```
 curl -X POST 127.0.0.1:8888/atr/batch/create/student -H 'Content-Type: application/json' -d "[{\"name\":\"Scarlet\",\"gender\":2,\"age\":8,\"class\":\"3\",\"province\":\"beijin\"},{\"name\":\"Tania\",\"gender\":2,\"age\":9,\"class\":\"2\",\"province\":\"hebei\"}]"
```

```
{"code":0,"data":null,"message":""}
```

### 路由分组
```go
RouterPrefix       = "/atr"
CustomRouterPrefix = "/custom"
```
默认路由有/atr的前缀 这个是可以设置的通过`WithRouterPrefix` 需要注意如果使用了自定义路由，自定义路由的前缀和默认路由前缀不能一样，因为绑定了不同的middleware分组。


### 更多功能
[绑定模型](./doc/bind_model.md) 绑定模型可以通过模型的注解完成参数校验，和对gorm注解的支持 <br>
[代码结构](./doc/struct.md) <br>
[自定义路由](./doc/custom.md#自定义路由) <br>
[自定义控制器](./doc/custom.md#自定义handler) <br>
[自定义参数解析](./doc/mconf.md#自定义参数解析)  list接口支持<br>
[创建前数据完善](./doc/mconf.md#创建前数据完善)  create/update/batchcreate 接口在插入数据前可以自定对参数的补充变形，插入和修改前的数据完善<br>
[批量创建前数据完善](./doc/mconf.md#批量创建前数据完善)  create/update/batchcreate 接口在插入数据前可以自定对参数的补充变形，插入和修改前的数据完善<br>
[修改前数据完善](./doc/mconf.md#修改前数据完善)  create/update/batchcreate 接口在插入数据前可以自定对参数的补充变形，插入和修改前的数据完善<br>
[展示字段配置](./doc/mconf.md#指定查询字段)  Select/Hidden/AddColumns 对展示数据进行配置、隐藏、扩充<br>
[控制访问](./doc/rules.md)  <br>

list支持分组和聚合函数 参考 examples/custom_route