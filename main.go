package main

import (
	"container/ring"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/hpcloud/tail"
)

type ParsedResult struct {
	Number  float32
	Matched bool
	Time    time.Time
}

type CycleStats struct {
	CycleNumber int64
	Count       int64
	Period      int64
	Rate        float32
	AverateTime float32
}

func main() {
	var filename = "test.log"
	t, err := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		panic(err)
	}

	var (
		cycle       int64   = 5
		cycleNumber int64   = 0
		count       int64   = 0
		sum         float32 = 0
		lastCycle   int64   = 0
	)

	var statsChan = make(chan *ParsedResult)
	var cycleStatsChan = make(chan *CycleStats)

	go parseWorker(t.Lines, statsChan)
	go warn(cycleStatsChan)

	for stats := range statsChan {
		timestamp := stats.Time.Unix()
		currentNumber := timestamp % cycle

		if currentNumber != lastCycle {
			cycleStatsChan <- &CycleStats{
				Count:       count,
				Period:      cycle,
				AverateTime: sum / float32(count),
				Rate:        float32(count) / float32(cycle),
				CycleNumber: cycleNumber,
			}

			fmt.Printf("time: %s cycle: %d count: %d rate: %f avg: %f\n",
				stats.Time, cycleNumber, count, float32(count)/float32(cycle), sum/float32(count))
			lastCycle = currentNumber
			cycleNumber += 1
			count = 0
			sum = 0
		}
		count += 1
		sum += stats.Number
	}
}

func parseWorker(lines <-chan *tail.Line, rets chan<- *ParsedResult) {
	for line := range lines {
		var ret = parse(line)
		if !ret.Matched {
			continue
		}

		rets <- ret
	}
}

func warn(statsChan <-chan *CycleStats) {
	const ringLen = 100
	var warning = false
	var ringStats = ring.New(ringLen)
	for cycle := range statsChan {
		ringStats.Value = cycle
		needWarn := checkWarning(ringStats, 5)
		if needWarn {
			warning = true
			fmt.Println("warning")
		}
		if !needWarn && warning {
			warning = false
			fmt.Println("warning cleared")
		}
		ringStats = ringStats.Next()
	}
}

func checkWarning(ringStats *ring.Ring, count int) bool {
	var (
		exceedTimes int = 0
	)
	for index := 0; index < count; index++ {
		cycleStats, ok := ringStats.Value.(*CycleStats)
		if !ok {
			continue
		}

		if cycleStats.AverateTime > 1.0 {
			exceedTimes += 1
		}

		if exceedTimes > 3 {
			return true
		}

		ringStats = ringStats.Prev()
	}
	return false
}

func parse(line *tail.Line) *ParsedResult {
	re := regexp.MustCompile("(\\d+\\.?\\d*)")
	results := re.FindSubmatch([]byte(line.Text))
	if len(results) == 0 {
		return &ParsedResult{0, false, line.Time}
	}
	n, err := strconv.ParseFloat(string(results[1]), 32)
	if err != nil {
		panic(err)
	}

	return &ParsedResult{Number: float32(n), Matched: true, Time: line.Time}
}
