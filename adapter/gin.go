package adapter

type GinCtxCore interface {
	JSON(httpCode int, body interface{})
	BindJSON(body interface{}) error
	Bind(body interface{}) error
	Param(key string) string
}
