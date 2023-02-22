package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	fg "reader/modules/file_generator"
	"reader/modules/utils"
	"strconv"
	"strings"
	"time"

	"github.com/CusterJ/data-aggr/proto/pb"
)

func (s *Server) GetMainPageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	To generate new file go to
		/generate - generate new file with fake data 
	ex ===
	http://localhost:8002/generate?length=10
		/stats - get stats
	ex ===
	http://localhost:8002/stats?interval=year&start=1595575638&end=1637685638

	intervals are: hour, day, week, month, year
	`)
}

func (s *Server) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	defer utils.TimeTrack(time.Now(), "GetStats handler")
	interval := r.URL.Query().Get("interval")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	fd, err := strconv.Atoi(start)
	utils.Check(err)

	td, err := strconv.Atoi(end)
	utils.Check(err)

	iv, ok := pb.Interval_value["INTERVAL_"+strings.ToUpper(interval)]
	if !ok {
		iv = int32(pb.Interval_INTERVAL_YEAR)
	}

	log.Printf("interval get %s = %d, start: %s, end: %s\n", interval, iv, start, end)

	res, err := s.sc.GetStats(r.Context(), &pb.GetStatsRequest{
		FromDate: int32(fd),
		ToDate:   int32(td),
		Interval: pb.Interval(iv),
	})

	utils.Check(err)

	w.Header().Add("content-type", "application/json")

	b, err := json.Marshal(res.Aggrs)
	utils.Check(err)
	fmt.Fprintf(w, "%s", b)
}

func (s *Server) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "This handler will generate a file of a given length, then read it and save the data to the database. Max length is 50K\n=======")

	length := r.URL.Query().Get("length")

	if length == "" {
		fmt.Fprintln(w, "To generate file add length parameter: /generate?length=100")
		return
	}

	lg, err := strconv.Atoi(length)
	if err != nil {
		fmt.Fprintln(w, "Generate and read file error! Add a number for parameter length")
		return
	}

	if lg <= 0 || lg > 50000 {
		fmt.Fprintln(w, "Length must be greater than zero and not greater than 50.000")
		return
	}

	// log.Printf("GenerateHandler get = %d\n", lg)

	err = fg.GenerateNewFile(lg)
	utils.Check(err)

	err = s.SaveFileStream("data.json")
	if err != nil {
		fmt.Fprintf(w, "Save data error: %s", err)
		return
	}

	fmt.Fprintf(w, "Generate and read file with lengh: %d -- done", lg)
}
