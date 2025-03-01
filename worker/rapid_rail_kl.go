package worker

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mholt/archives"
)

type MrtWorker struct {
	httpCli *http.Client
}

var _ Worker = (*MrtWorker)(nil)

func init() {
	mrtWorker := MrtWorker{
		httpCli: NewWorkerHttpClient(),
	}

	Workers = append(Workers, &mrtWorker)
}

func (w *MrtWorker) Run() error {
	res, err := w.httpCli.Get("https://api.data.gov.my/gtfs-static/prasarana?category=rapid-rail-kl")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("worker/rapid_rail_kl: gtfs returned status code %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	fsys, err := archives.FileSystem(context.Background(), "", bytes.NewReader(b))
	if err != nil {
		return err
	}

	ff, err := fsys.Open("stop_times.txt")
	if err != nil {
		return err
	}
	defer ff.Close()

	r := csv.NewReader(ff)

	records, err := r.ReadAll()
	if err != nil {
		return fmt.Errorf("worker/rapid_rail_kl: could not read record: %s", err)
	}

	stopTimes := make([]StopTime, len(records[1:]))
	for i, record := range records[1:] {
		arrivalTime, err := time.Parse("15:04:05", record[3])
		if err != nil {
			return fmt.Errorf("worker/rapid_rail_kl: could not parse arrival time: %s", err)
		}

		departureTime, err := time.Parse("15:04:05", record[4])
		if err != nil {
			return fmt.Errorf("worker/rapid_rail_kl: could not parse departure time: %s", err)
		}

		directionID, err := strconv.Atoi(record[1])
		if err != nil {
			return fmt.Errorf("worker/rapid_rail_kl: could not parse direction id: %s", err)
		}

		stopSequence, err := strconv.Atoi(record[6])
		if err != nil {
			return fmt.Errorf("worker/rapid_rail_kl: could not parse direction id: %s", err)
		}

		stopTimes[i] = StopTime{
			RouteID:       record[0],
			DirectionID:   directionID,
			TripID:        record[2],
			ArrivalTime:   arrivalTime,
			DepartureTime: departureTime,
			StopID:        record[5],
			StopSequence:  stopSequence,
		}
	}

	stb, _ := json.Marshal(stopTimes)

	return os.WriteFile("./data/rapid_rail_kl/stop_times.json", stb, 0755)
}
