package parse

import (
	"context"
	"errors"
	"fmt"
	"github.com/agxmaster/atm/clause"
	tools "github.com/agxmaster/atm/util"
	"github.com/agxmaster/atrg/core"
	"github.com/agxmaster/atrg/util"
	"strconv"

	"strings"
)

const PageSizeKey = "_size"
const OrderKey = "_order"
const GroupKey = "_group"
const CustomKey = "_custom"

type ClauseParse struct {
	clauses   clause.Clauses
	ModelType *core.MpModel
	err       error
	ctx       context.Context
}

func (p *ClauseParse) Parse(params *map[string]interface{}, page int) *ClauseParse {

	p.parsePage(params, page).
		parseOrders(params).
		ParsSelect().
		ParseCustomScope(params)

	if len(p.clauses) > 0 {
		p.clauses = append(p.clauses, clause.ColumnMap(*params))
	}
	return p
}

func (p *ClauseParse) parsePage(params *map[string]interface{}, page int) *ClauseParse {

	var pageSize int

	if p.ModelType.MaxPageSize > 0 {
		pageSize = p.ModelType.MaxPageSize
	} else {
		pageSize = core.DefaultMaxPageSize
	}

	if size, ok := (*params)[PageSizeKey]; ok {
		switch sizeType := size.(type) {
		case []interface{}:
			pageSize = p.getSize(sizeType[0], pageSize)
		case []string:
			pageSize = p.getSize(sizeType[0], pageSize)
		case interface{}:
			pageSize = p.getSize(sizeType, pageSize)
		}
		delete(*params, PageSizeKey)
	}
	pageClause := clause.Page{NeedCount: true, PageNum: page, PageSize: pageSize}
	groupClause := p.parseGroup(params)

	if groupClause == nil || len(groupClause) == 0 {
		p.clauses = append(p.clauses, pageClause)
	} else {
		p.clauses = append(p.clauses, clause.PageWithGroup{
			Page:   pageClause,
			Groups: groupClause,
		})
	}
	return p

}

func (p *ClauseParse) getSize(sizeAny interface{}, curSize int) int {
	sizeInt, err := strconv.Atoi(fmt.Sprintf("%v", sizeAny))
	if err != nil {
		p.err = errors.New("_page must be a int")
	}
	if sizeInt > 0 && sizeInt < curSize {
		curSize = sizeInt
	}
	return curSize
}
func (p *ClauseParse) parseOrders(params *map[string]interface{}) *ClauseParse {

	if orders, ok := (*params)[OrderKey]; ok {
		clausOrders := clause.Orders{}
		switch ordersData := orders.(type) {
		case []interface{}:
			for _, orderString := range ordersData {
				order := p.parseOrder(orderString)
				if order != nil {
					clausOrders = append(clausOrders, *order)
				}
			}
			break
		case []string:
			for _, orderString := range ordersData {
				order := p.parseOrder(orderString)
				if order != nil {
					clausOrders = append(clausOrders, *order)
				}
			}
			break
		case interface{}:
			order := p.parseOrder(orders)
			if order != nil {
				clausOrders = append(clausOrders, *order)
			}
			break
		}
		delete(*params, OrderKey)
		if clausOrders != nil && len(clausOrders) != 0 {
			p.clauses = append(p.clauses, clausOrders)
		}
	}
	return p
}

func (p *ClauseParse) parseOrder(orderString interface{}) *clause.Order {

	orderCase := strings.Split(fmt.Sprintf("%v", orderString), ",")
	if len(orderCase) == 0 {
		return nil
	}
	claus := clause.Order{Field: orderCase[0]}
	if len(orderCase) > 1 && orderCase[1] == "desc" {
		claus.Desc = true
	} else {
		claus.Desc = false
	}

	return &claus
}

func (p *ClauseParse) checkEnableGroup() bool {
	return !(p.ModelType.CanGroupColumn == nil || len(p.ModelType.CanGroupColumn) == 0)
}

func (p *ClauseParse) parseGroup(params *map[string]interface{}) clause.Groups {

	clausGroup := clause.Groups{}
	if groups, ok := (*params)[GroupKey]; ok {
		if !p.checkEnableGroup() {
			p.err = errors.New(fmt.Sprintf("not support group"))
			return nil
		}

		switch groupData := groups.(type) {
		case []string:
			if groupData == nil || len(groupData) == 0 {
				return nil
			}
			clausGroup = p.getGroups(groupData[0])
		case []interface{}:
			if groupData == nil || len(groupData) == 0 {
				return nil
			}
			clausGroup = p.getGroups(groupData[0])
		case interface{}:
			if groupData == "" {
				return nil
			}
			clausGroup = p.getGroups(groupData)
		}
		delete(*params, GroupKey)
	}
	return clausGroup
}

func (p *ClauseParse) getGroups(groupAny interface{}) (clausGroup clause.Groups) {
	groupsColumns := strings.Split(fmt.Sprintf("%v", groupAny), ",")

	if !p.selectMustGroup(groupsColumns) {
		return
	}

	for _, groupString := range groupsColumns {
		if groupString != "" && p.checkSupportGroupColumn(groupString) {
			clausGroup = append(clausGroup, groupString)
		}
	}
	return
}

func (p *ClauseParse) selectMustGroup(groupsColumns []string) bool {
	for _, selColumn := range p.ModelType.Select {
		if !tools.Contains(selColumn, groupsColumns) {
			p.err = errors.New(fmt.Sprintf("select column must group %s", selColumn))
			return false
		}
	}
	return true
}

func (p *ClauseParse) checkSupportGroupColumn(groupString string) bool {
	if tools.Contains(fmt.Sprintf("%v", groupString), p.ModelType.CanGroupColumn) {
		return true
	}
	p.err = errors.New(fmt.Sprintf("not support group column %s", groupString))
	return false
}

func (p *ClauseParse) ParsSelect() *ClauseParse {

	var clauseSelect clause.Select

	if p.ModelType.Select != nil && len(p.ModelType.Select) > 0 {
		clauseSelect = p.ModelType.Select
	}

	selectAggregate := p.parseAggregations()
	if selectAggregate != nil && len(selectAggregate) != 0 {
		clauseSelect = append(clauseSelect, selectAggregate...)
	}

	if clauseSelect != nil && len(clauseSelect) > 0 {
		p.clauses = append(p.clauses, clauseSelect)
	}

	return p
}

func (p *ClauseParse) parseAggregations() clause.Select {

	var selectAggregate clause.Select

	if p.ModelType.Aggregations == nil && len(p.ModelType.Aggregations) == 0 {
		return selectAggregate
	}

	var selectColumn []string
	for _, aggregate := range p.ModelType.Aggregations {
		asName := fmt.Sprintf("%v_%v", aggregate.AggregationType, aggregate.Column)

		if aggregate.Alias != "" {
			asName = aggregate.Alias
		}
		selectAggregate = append(selectAggregate, fmt.Sprintf("%s(%s) as %s", aggregate.AggregationType, aggregate.Column, asName))
		selectColumn = append(selectColumn, asName)
	}
	p.ModelType.AddSelect(p.ctx, selectColumn)
	return selectAggregate
}

func (p *ClauseParse) ParseCustomScope(params *map[string]interface{}) *ClauseParse {
	if custom, ok := (*params)[CustomKey]; ok {
		var (
			customString interface{}
			scope        clause.Scope
			err          error
		)
		switch customData := custom.(type) {
		case []interface{}:
			customString = customData[0]
		case []string:
			customString = customData[0]
		case interface{}:
			customString = customData
		}

		if customString == "" {
			return p
		}
		if p.ModelType.ListWithCustomScope != nil {
			scope, err = p.ModelType.ListWithCustomScope(p.ctx, util.StringToBytes(fmt.Sprintf("%v", customString)))
			if err != nil {
				p.err = err
				return p
			}
			if scope != nil {
				p.clauses = append(p.clauses, scope)
			}
		}
		delete(*params, CustomKey)
	}
	return p
}

func (p *ClauseParse) GetClause() (clause.Clauses, error) {
	return p.clauses, p.err
}
