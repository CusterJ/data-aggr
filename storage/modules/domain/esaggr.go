package domain

type Aggrs struct {
	Took         int          `json:"took"`
	TimedOut     bool         `json:"timed_out"`
	Hits         Hits         `json:"hits"`
	Aggregations Aggregations `json:"aggregations"`
}

type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}
type Hits struct {
	Total Total `json:"total"`
}

type Aggregations struct {
	Histogram Histogram `json:"histogram"`
}

type Histogram struct {
	HistoBuckets []HistoBuckets `json:"buckets"`
	Interval     string         `json:"interval"`
}

type HistoBuckets struct {
	KeyAsString string      `json:"key_as_string"`
	Key         int64       `json:"key"`
	DocCount    int         `json:"doc_count"`
	SignalCount SignalCount `json:"signal-count"`
}

type SignalCount struct {
	DocCountErrorUpperBound int       `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int       `json:"sum_other_doc_count"`
	Buckets                 []Buckets `json:"buckets"`
}

type Buckets struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
}
