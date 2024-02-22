package core

type RuleType int

const (
	RuleTypeAllow RuleType = iota
	RuleTypeDeny
)

type Rules []Rule
type Rule struct {
	Table      []string
	RouteTypes []RouteType
	RuleType   RuleType
}

func (rs Rules) Allow(routeType RouteType, table string) bool {
	for _, rule := range rs {
		if rule.Match(routeType, table) {
			return rule.BeAllow()
		}
	}
	return true
}

func (r *Rule) Match(routeType RouteType, table string) bool {
	return r.TableContains(table) && r.RouteTypeContains(routeType)
}
func (r *Rule) BeAllow() bool {
	return r.RuleType == RuleTypeAllow
}

func (r *Rule) RouteTypeContains(routeType RouteType) bool {
	for _, v := range r.RouteTypes {
		if v == "*" || v == routeType {
			return true
		}
	}
	return false
}

func (r *Rule) TableContains(table string) bool {
	for _, v := range r.Table {
		if v == "*" || v == table {
			return true
		}
	}
	return false
}
