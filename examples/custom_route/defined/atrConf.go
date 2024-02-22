package defined

import (
	"github.com/agxmaster/atrg/core"
	"github.com/agxmaster/atrg/examples/custom_route/model"
	"reflect"
)

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
