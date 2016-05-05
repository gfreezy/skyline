package skyline

import (
	"regexp"
	"time"
)

type Filter struct {
	matchedStr         *regexp.Regexp
	cycle              int64
	threshold          float32
	namePrefix         string
	currentCycleNumber int64
	cycleStats         *Cycle
	CycleStatsChannel  chan Cycle
	cycleId            int64
}

func NewFilter(conf *FilterItemConf) *Filter {
	r := regexp.MustCompile(conf.MatchStr)
	ch := make(chan Cycle)
	return &Filter{
		matchedStr:         r,
		cycle:              conf.Cycle,
		threshold:          conf.Threshold,
		namePrefix:         conf.ItemNamePrefix,
		currentCycleNumber: 0,
		cycleStats:         NewCycle(0, conf.Cycle),
		CycleStatsChannel:  ch,
		cycleId:            0,
	}
}

func (self *Filter) AddLine(line []byte, t time.Time) {
	cycleNumber := t.Unix() / self.cycle
	if cycleNumber != self.currentCycleNumber {
		// fmt.Println("id", self.cycleStats.Id, "count", self.cycleStats.Count, "avg", self.cycleStats.AverateTime(), "rate", self.cycleStats.Rate())
		self.CycleStatsChannel <- *self.cycleStats
		self.cycleStats = NewCycle(self.cycleId, self.cycle)
		self.currentCycleNumber = cycleNumber
		self.cycleId += 1
	}
	lineStats := parseLine(line, self.matchedStr, t)
	if lineStats.Matched {
		self.cycleStats.AddLine(&lineStats)
	}
}
