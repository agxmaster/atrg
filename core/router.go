package core

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"regexp"
	"sync"
)

type Iatr interface {
	Info(c *gin.Context)
	List(c *gin.Context)
	Delete(c *gin.Context)
	Create(c *gin.Context)
	BatchCreate(c *gin.Context)
	Update(c *gin.Context)
}

type RouteType string
type Method string

const (
	PageKey  = "page"
	IdKey    = "id"
	ModelKey = "model"

	RouterPrefix       = "/atr"
	CustomRouterPrefix = "/custom"

	RouteTypeInfo       RouteType = "INFO"
	RouteTypeList       RouteType = "LIST"
	RouteTypeCreate     RouteType = "CREATE"
	RouteTypeUpdate     RouteType = "UPDATE"
	RouteTypeDelete     RouteType = "DELETE"
	RouteTypeCrateBatch RouteType = "CREATE_BATCH"

	MethodPost Method = "POST"
	MethodGet  Method = "GET"
)

var RouteBindInstance RouteBind

type RouteBind struct {
	RouteMap   map[RouteType]gin.HandlerFunc
	RegPathMap map[string]RouteType
	mu         sync.Mutex
}

func initRoute() {
	SetDefaultRoute()
	RouteBindInstance.SetRouteBind()
}

var defaultRoute []CustomRoute

func (r *RouteBind) SetRouteBind() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.RouteMap = make(map[RouteType]gin.HandlerFunc)
	r.RegPathMap = make(map[string]RouteType)
	for _, route := range defaultRoute {
		r.RouteMap[route.RouteType] = route.Handler
		r.RegPathMap[GetRouteRegKey(route.Method, fmt.Sprintf("%s%s", Mp.RoutePrefix, route.RoutePath))] = route.RouteType
	}
}

func (r *RouteBind) AddBind(routeType RouteType, HandlerFunc gin.HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.RouteMap[routeType] = HandlerFunc
}

func (r *RouteBind) GetDefaultHandler(routeType RouteType) gin.HandlerFunc {
	return r.RouteMap[routeType]
}
func AtrRouter(r *gin.Engine) {

	initRoute()
	dr := r.Group(Mp.RoutePrefix, func(ctx *gin.Context) {
		for reg, routeType := range RouteBindInstance.RegPathMap {
			if ok, _ := regexp.MatchString(reg, GetRouteKey(Method(ctx.Request.Method), ctx.Request.URL.Path)); ok {
				if !Mp.Rules.Allow(routeType, ctx.Param(ModelKey)) {
					Error(ctx, errors.New("model not support this route"))
					ctx.Abort()
				} else {
					return
				}
			}
		}
	})

	for _, route := range defaultRoute {
		if route.Method == MethodGet {
			dr.GET(route.RoutePath, route.Handler)
		}
		if route.Method == MethodPost {
			dr.POST(route.RoutePath, route.Handler)
		}
	}

	cr := r.Group(Mp.CustomRouterPrefix, func(ctx *gin.Context) {
		if models := (&MpModel{}).GetCustomMpModel(Method(ctx.Request.Method), ctx.Request.URL.Path, ctx.Param(ModelKey)); models == nil {
			Error(ctx, errors.New("model not support this route"))
			ctx.Abort()
		}
	})

	customConf := *GetCustomConf()
	for _, confs := range customConf {
		for _, confOne := range confs {
			if confOne.MatchRouters == nil || len(confOne.MatchRouters) == 0 {
				return
			}
			for _, routeConfig := range confOne.MatchRouters {

				handler := routeConfig.Handler
				if handler == nil {
					handler = RouteBindInstance.GetDefaultHandler(routeConfig.RouteType)
				}

				if routeConfig.Method == MethodGet {
					cr.GET(routeConfig.RoutePath, handler)
				}
				if routeConfig.Method == MethodPost {
					cr.POST(routeConfig.RoutePath, handler)
				}
			}

		}
	}
}

func SetDefaultRoute() {
	defaultRoute = []CustomRoute{
		{
			RouteType: RouteTypeInfo,
			RoutePath: "/info/:model/:id",
			Handler:   Mp.Iatr.Info,
			Method:    MethodGet,
		},
		{
			RouteType: RouteTypeList,
			RoutePath: "/list/:model/:page",
			Handler:   Mp.Iatr.List,
			Method:    MethodGet,
		},
		{
			RouteType: RouteTypeCreate,
			RoutePath: "/create/:model",
			Handler:   Mp.Iatr.Create,
			Method:    MethodPost,
		},
		{
			RouteType: RouteTypeDelete,
			RoutePath: "/delete/:model/:id",
			Handler:   Mp.Iatr.Delete,
			Method:    MethodPost,
		},
		{
			RouteType: RouteTypeCrateBatch,
			RoutePath: "/batch/create/:model",
			Handler:   Mp.Iatr.BatchCreate,
			Method:    MethodPost,
		},
		{
			RouteType: RouteTypeUpdate,
			RoutePath: "/update/:model",
			Handler:   Mp.Iatr.Update,
			Method:    MethodPost,
		},
	}
}
