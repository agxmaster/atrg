package core

type Config struct {
	ModelMap           map[string]MpModel
	CustomModelMap     map[string][]MpModel
	Iatr               Iatr
	RoutePrefix        string
	CustomRouterPrefix string
	SuccessCode        ResponseCode
	ErrorCode          ResponseCode
	Rules              Rules
}

type Option struct {
	Fn func(o *Config)
}

func (o *Config) Apply(opts []Option) {
	for _, op := range opts {
		op.Fn(o)
	}
}
