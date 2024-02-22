package core

import (
	"context"
	"github.com/agxmaster/atm"
	tools "github.com/agxmaster/atm/util"
	"github.com/agxmaster/atrg/util"
	"reflect"
	"regexp"
	"strings"
)

type MpModel struct {
	Model                    reflect.Type
	ColumnMapping            map[string]string          //db列采用别名作为查询条件， 展示的别名通过model的column的注解处理 这里不做单独处理
	Select                   []string                   //查询接口要展示的列
	Hidden                   []string                   //查询需要隐藏的列
	AddColumns               []map[string]ColumnAddFunc //用于查询接口 这里可以为数据集增加列、或者对列的数据做变形
	MaxPageSize              int                        //每页展示数据上限
	ListWithCustomScope      ListWithCustomScope        //针对list接口自定查询条件
	CreateParamsHandler      SaveParamsHandler
	CreateBatchParamsHandler SaveBatchParamsHandler
	UpdateParamsHandler      SaveParamsHandler
	Aggregations             []Aggregation
	CanGroupColumn           []string
	addSelect                []string
	MatchRouters             []CustomRoute
}

func (c *MpModel) NeedFormat() bool {
	return !(c.Hidden == nil || len(c.Hidden) == 0 || c.AddColumns == nil || len(c.AddColumns) == 0)
}

func (c *MpModel) FormatModelLine(ctx context.Context, res interface{}) (interface{}, error) {
	if !c.NeedFormat() {
		return res, nil
	}

	resMapRef, err := util.InterfaceToAny[atm.RowsMap](res)
	if err != nil {
		return res, err
	}

	if err != nil {
		return res, err
	}

	return c.FormatLine(ctx, *resMapRef)
}

func (c *MpModel) FormatLine(ctx context.Context, res atm.RowsMap) (atm.RowsMap, error) {
	if !c.NeedFormat() {
		return res, nil
	}

	c.FormatHandler(ctx, res)

	return res, nil
}

func (c *MpModel) FormatModelList(ctx context.Context, res any) (any, error) {

	if res == nil {
		res = make([]atm.RowsMap, 0)
	}

	if !c.NeedFormat() {
		return res, nil
	}

	resSliceMapRef, err := util.InterfaceToAny[[]atm.RowsMap](res)
	if err != nil {
		return res, err
	}

	resSliceMap, err := c.FormatRowsMapList(ctx, *resSliceMapRef)

	return resSliceMap, err
}

func (c *MpModel) FormatRowsMapList(ctx context.Context, res []atm.RowsMap) ([]atm.RowsMap, error) {

	if res == nil {
		res = make([]atm.RowsMap, 0)
	}

	if !c.NeedFormat() {
		return res, nil
	}
	for _, resMap := range res {
		c.FormatHandler(ctx, resMap)
	}

	return res, nil
}

func (c *MpModel) FormatHandler(ctx context.Context, resMap map[string]interface{}) {
	c.FormatSelect(ctx, resMap).FormatAddColumns(ctx, resMap).FormatHidden(ctx, resMap)
}

func (c *MpModel) FormatAddColumns(ctx context.Context, resMap map[string]interface{}) *MpModel {
	for _, fm := range c.AddColumns {
		for k, fn := range fm {
			resMap[k] = fn(ctx, resMap, k)
		}
	}
	return c
}

func (c *MpModel) FormatHidden(ctx context.Context, resMap map[string]interface{}) *MpModel {
	for _, hidColumn := range c.Hidden {
		if _, ok := resMap[hidColumn]; ok {
			delete(resMap, hidColumn)
		}
	}
	return c
}

func (c *MpModel) FormatSelect(ctx context.Context, resMap map[string]interface{}) *MpModel {

	for k, _ := range resMap {
		if !tools.Contains(k, c.Select) && !tools.Contains(k, c.addSelect) {
			delete(resMap, k)
		}
	}
	return c
}

func (c *MpModel) QueryChange(ctx context.Context, params map[string]interface{}) map[string]interface{} {
	if c.ColumnMapping != nil {
		for k, v := range c.ColumnMapping {
			if pv, ok := params[v]; ok {
				params[k] = pv
				delete(params, v)
			}
		}
	}
	return params
}

func (c *MpModel) AddSelect(ctx context.Context, selects []string) {
	c.addSelect = append(c.addSelect, selects...)
}

func (c *MpModel) IsCustomRoute(routePath string) bool {
	return strings.HasPrefix(routePath, Mp.CustomRouterPrefix)
}

func (c *MpModel) GetMpModel(method Method, routePath string, modelName string) *MpModel {
	if c.IsCustomRoute(routePath) {
		return c.GetCustomMpModel(method, routePath, modelName)
	} else {
		modelConf, _ := (*GetConf())[modelName]
		return &modelConf
	}
}

func (c *MpModel) GetCustomMpModel(method Method, routePath string, modelName string) *MpModel {
	if customModelMapCache != nil {
		if ModelMap, ok := customModelMapCache[modelName]; ok {
			for k, models := range ModelMap {
				matched, err := regexp.MatchString(k, GetRouteKey(method, routePath))
				if err != nil || !matched {
					return nil
				}
				return models
			}
		}
	}
	return nil
}
