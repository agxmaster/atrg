package core

type AggregationType string

const AggregationTypeCount AggregationType = "COUNT"
const AggregationTypeSum AggregationType = "SUM"
const AggregationTypeAvg AggregationType = "AVG"
const AggregationTypeMin AggregationType = "MIN"
const AggregationTypeMax AggregationType = "MAX"

type Aggregation struct {
	AggregationType AggregationType
	Column          string
	Alias           string
}
