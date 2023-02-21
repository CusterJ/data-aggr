package database

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"storage/modules/domain"
	"storage/modules/utils"
	"strings"
	"time"

	"github.com/CusterJ/data-aggr/proto/pb"
)

const esConstURL string = "http://localhost:9200/stats/"

type Elastic struct {
	URL string
}

func NewElastic() *Elastic {
	envURL := os.Getenv("ES_URL")

	if envURL == "" {
		return &Elastic{
			URL: esConstURL,
		}
	}

	return &Elastic{
		URL: envURL,
	}
}

func intervalToString(i int) string {
	var interval string

	switch i {
	case 0:
		interval = "hour"

	case 1:
		interval = "day"

	case 2:
		interval = "week"

	case 3:
		interval = "month"

	case 4:
		interval = "year"
	default:
		interval = "year"
	}

	return interval
}

func (e *Elastic) CheckIndex() error {
	log.Println(e.URL)
	res, err := http.Get(e.URL)
	utils.Check(err)

	defer res.Body.Close()

	if res.StatusCode != 200 {
		err := e.CreateIndex()
		if err != nil {
			// return fmt.Errorf("create index error, %s", err)
			log.Fatal("can't create index: ", err)
		}
	}
	return nil
}

func (e *Elastic) CreateIndex() error {

	query := (`{
		"mappings": {
		  "properties": {
			"data": {
			  "type": "text"
			},
			"id": {
			  "type": "text"
			},
			"signal": {
			  "type": "keyword"
			},
			"time": {
			  "type": "date",
			  "format": "epoch_second"
			}
		  }
		},
		"settings": {
		  "index": {
			"routing": {
			  "allocation": {
				"include": {
				  "_tier_preference": "data_content"
				}
			  }
			}
		  }
		}
	  }`)

	req, err := http.NewRequest(http.MethodPut, e.URL, strings.NewReader(query))
	utils.Check(err)

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	utils.Check(err)

	defer res.Body.Close()

	log.Println("Create Index status: ", res.Status)

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error creating new index, %s", res.Status)
	}

	return nil
}

func (e *Elastic) QueryStats(from, to, interval int) (domain.Aggrs, error) {
	utils.TimeTrack(time.Now(), "QueryStats")

	url := e.URL + "_search"
	size := 0

	query := fmt.Sprintf(`{
		"size": %d,
		"sort": [
		  {
			"signal": "asc"
		  }
		],
		"query": {
		  "bool": {
			"filter": [
			  {
				"range": {
				  "time": {
					"gte": %d000,
					"lte": %d000
				  }
				}
			  }
			]
		  }
		},
		"aggs": {
		  "histogram": {
			"date_histogram": {
			  "field": "time",
			  "calendar_interval": "%s",
			  "format": "yyyy-MM-dd"
			},
			"aggs": {
			  "signal-count": {
				"terms": {
				  "field": "signal"
				}
			  }
			}
		  }
		}
	  }`, size, from, to, intervalToString(interval))

	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(query))
	utils.Check(err)

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	utils.Check(err)
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	utils.Check(err)
	log.Printf("func QueryStats res: %+v", res.Status)

	esres := &domain.Aggrs{}

	err = json.Unmarshal(result, &esres)
	utils.Check(err)

	log.Printf("Total results: %+v, esres.Hits: %d\n", esres.Hits.Total.Value, len(esres.Aggregations.Histogram.HistoBuckets))

	return *esres, nil
}

func (e *Elastic) BulkWrite(data []*pb.FooData) error {
	err := e.CheckIndex()
	if err != nil {
		return err
	}

	url := e.URL + "_bulk"
	var payload string

	if len(data) > 0 {

		for _, v := range data {
			id := fmt.Sprintf(`{ "index": { "_id": "%s" }}`, v.Id)
			b, err := json.Marshal(v)
			utils.Check(err)
			payload += id + "\n" + string(b) + "\n"
		}

		// fmt.Printf("BulkWrite payload:\n%+v\n", payload)
		res, err := http.Post(url, "application/json", strings.NewReader(payload))
		utils.Check(err)
		defer res.Body.Close()
		log.Println("bulk save data status: ", res.Status)
	}
	return nil
}
