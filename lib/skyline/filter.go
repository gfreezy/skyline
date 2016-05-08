package skyline

import (
	"fmt"
	"regexp"
	"time"
)

type Filter struct {
	Name               string
	matchedStr         *regexp.Regexp
	cyclePeriod        int64
	threshold          float64
	namePrefix         string
	currentCycleNumber int64
	cycleStats         *Cycle
	CycleStatsChannel  chan Cycle
	cycleID            int64
	cycleRealCount     int64
}

func NewFilter(name string, conf *FilterItemConf) *Filter {
	r := regexp.MustCompile(conf.MatchStr)
	ch := make(chan Cycle)
	return &Filter{
		Name:               name,
		matchedStr:         r,
		cyclePeriod:        conf.Cycle,
		threshold:          conf.Threshold,
		namePrefix:         conf.ItemNamePrefix,
		currentCycleNumber: 0,
		cycleStats:         NewCycle(0, name, conf.Cycle),
		CycleStatsChannel:  ch,
		cycleID:            0,
	}
}

func (self *Filter) AddLine(line []byte, t time.Time) {
	cycleNumber := t.Unix() / self.cyclePeriod
	if cycleNumber != self.currentCycleNumber {
		if Debug {
			fmt.Printf(
				"[%s] name: %s, cycle: %d, count: %d/%d, rate: %d/%d, avg: %.3f\n",
				t.Format(time.RFC3339),
				self.Name,
				self.cycleStats.Id,
				self.cycleStats.Count,
				self.cycleRealCount,
				self.cycleStats.Rate(),
				self.circleRealRate(),
				self.cycleStats.Averate(),
			)
		}
		self.CycleStatsChannel <- *self.cycleStats
		self.cycleStats = NewCycle(self.cycleID, self.Name, self.cyclePeriod)
		self.currentCycleNumber = cycleNumber
		self.cycleID++
		self.cycleRealCount = 0
	}
	lineStats := parseLine(line, self.matchedStr, t)

	aboveThreshold := false
	if self.threshold == 0 || lineStats.Number >= self.threshold {
		aboveThreshold = true
	}
	if lineStats.Matched && aboveThreshold {
		self.cycleStats.AddLine(&lineStats)
	}
	self.cycleRealCount++
}

func (self *Filter) circleRealRate() int64 {
	return self.cycleRealCount / self.cyclePeriod
}
