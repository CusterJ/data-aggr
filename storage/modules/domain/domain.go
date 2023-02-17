package domain

type FooData struct {
	ID     string `json:"id"`
	Time   int    `json:"time"`
	Signal string `json:"signal,omitempty"`
	Data   string `json:"data,omitempty"`
}

type StatsIntervalResponse struct {
	IntervalData []IntervalData `json:"content,omitempty"`
}

type IntervalData struct {
	Interval string `json:"interval,omitempty"`
	Start    string `json:"start,omitempty"`
	End      string `json:"end,omitempty"`
	FooData  `json:"foo_data,omitempty"`
}

type EsSearchResponse struct {
	Hits struct {
		Hits []struct {
			FooData `json:"_source,omitempty"`
		} `json:"hits,omitempty"`
		Total struct {
			Value int `json:"value,omitempty"`
		} `json:"total,omitempty"`
	} `json:"hits,omitempty"`
}
