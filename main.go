package main

import (
	"container/ring"
	"fmt"

	"github.com/gfreezy/skyline/lib/skyline"
	"github.com/hpcloud/tail"
)

func main() {
	var conf = skyline.LoadConfig("conf/config.json")
	monitor(conf.Monitors[0])
	// for _, monitorConf := range conf.Monitors {
	// 	go monitor(monitorConf)
	// }
}

func monitor(monitorConf skyline.MonitorConf) {
	var filename = monitorConf.LogFilePath
	t, err := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		panic(err)
	}

	filters := make([]*skyline.Filter, len(monitorConf.FilterItems))
	for _, filterConf := range monitorConf.FilterItems {
		filter := skyline.NewFilter(&filterConf)
		filters = append(filters, filter)
		go warn(filterConf.ItemNamePrefix, filter.CycleStatsChannel)
	}

	for line := range t.Lines {
		for _, f := range filters {
			if f == nil {
				continue
			}
			f.AddLine([]byte(line.Text), line.Time)
		}
	}
}

func warn(name string, statsChan <-chan skyline.Cycle) {
	const ringLen = 100
	var warning = false
	var ringStats = ring.New(ringLen)
	for cycle := range statsChan {
		ringStats.Value = cycle
		needWarn := checkWarning(ringStats, 5)
		if needWarn {
			warning = true
			fmt.Println(name, "warning")
		}
		if !needWarn && warning {
			warning = false
			fmt.Println(name, "warning cleared")
		}
		ringStats = ringStats.Next()
	}
}

func checkWarning(ringStats *ring.Ring, count int) bool {
	var (
		exceedTimes int
	)
	for index := 0; index < count; index++ {
		cycleStats, ok := ringStats.Value.(skyline.Cycle)
		if !ok {
			continue
		}

		if cycleStats.Averate() > 1.0 {
			exceedTimes++
		}

		if exceedTimes > 3 {
			return true
		}

		ringStats = ringStats.Prev()
	}
	return false
}
