package core

import (
	"context"
	"fmt"
	"github.com/agxmaster/atm/clause"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
)

type ResponseCode int64

const (
	CodeSuccess  ResponseCode = 0
	CodeAtrError ResponseCode = -1
)

const DefaultMaxPageSize = 50

var (
	Db                  *gorm.DB
	Mp                  = &Config{}
	customModelMapCache map[string]map[string]*MpModel
)

type ColumnAddFunc func(ctx context.Context, line map[string]interface{}, addKey string) interface{}
type SaveParamsHandler func(ctx context.Context, params interface{}) (map[string]interface{}, error)
type SaveBatchParamsHandler func(ctx context.Context, params interface{}) ([]map[string]interface{}, error)
type ListWithCustomScope func(ctx context.Context, customJson []byte) (clause.Scope, error)

type SetModel func() map[string]MpModel

type SetCustomModel func() map[string][]MpModel

type CustomRoute struct {
	RouteType RouteType
	RoutePath string
	Handler   gin.HandlerFunc
	Method    Method
}

func GetRouteKey(method Method, routePath string) string {
	return fmt.Sprintf("%s_%s", method, routePath)
}

func GetRouteRegKey(method Method, routePath string) string {
	return fmt.Sprintf("^%s_%s$", method, regexp.MustCompile(`(:[^/.]*)`).ReplaceAllString(routePath, "([^/.]*)"))
}

func SetCustomModelCache() {
	customModelMapCache = make(map[string]map[string]*MpModel)

	for modelName, models := range Mp.CustomModelMap {
		for _, model := range models {
			if _, ok := customModelMapCache[modelName]; !ok {
				customModelMapCache[modelName] = make(map[string]*MpModel)
			}
			for _, route := range model.MatchRouters {
				customModelMapCache[modelName][GetRouteRegKey(route.Method, fmt.Sprintf("%s%s", Mp.CustomRouterPrefix, route.RoutePath))] = &model
			}
		}
	}
}

func GetConf() *map[string]MpModel {
	return &Mp.ModelMap
}

func GetCustomConf() *map[string][]MpModel {
	return &Mp.CustomModelMap
}
