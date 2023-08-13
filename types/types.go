package types

type Metric struct {
	CPU    float64
	Memory float64
}

type NameSpaceMetric map[string]Metric

type MetricResponse struct {
	Metrics NameSpaceMetric `json:"metrics"`
}