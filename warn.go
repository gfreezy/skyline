package main

import (
	"container/ring"
	"fmt"
	"time"

	"os"

	"github.com/gfreezy/skyline/lib/skyline"
)

func warn(filterWarnings *skyline.FilterWarnings, statsChan <-chan skyline.Cycle) {
	const ringLen = 100
	var ringStats = ring.New(ringLen)
	for cycle := range statsChan {
		ringStats.Value = cycle

		for _, warning := range filterWarnings.Warnings {
			needWarn := checkWarning(ringStats, warning)

			if needWarn {
				warning.IsWarning = true
				msg := fmt.Sprintf("[%s]WARN: %s[%s]\n", time.Now(), warning.AlertName, warning.Formula)
				warning.Warn(msg)
				fmt.Fprintln(os.Stderr, msg)
			} else if warning.IsWarning {
				warning.IsWarning = false
				msg := fmt.Sprintf("[%s]OK: %s[%s]\n", time.Now(), warning.AlertName, warning.Formula)
				warning.Warn(msg)
				fmt.Fprintln(os.Stderr, msg)
			}
		}

		ringStats = ringStats.Next()
	}
}

func checkWarning(ringStats *ring.Ring, warning *skyline.Warning) bool {
	var (
		exceedTimes int
	)
	for index := 0; index < warning.FilterTotal; index++ {
		cycleStats, ok := ringStats.Value.(skyline.Cycle)
		if !ok {
			continue
		}

		if warning.NeedTrigger(cycleStats.Params) {
			exceedTimes++
		}

		if exceedTimes >= warning.FilterHigh {
			return true
		}

		ringStats = ringStats.Prev()
	}
	return false
}
