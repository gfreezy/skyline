package skyline

import (
	"regexp"
	"strconv"
	"time"
)

type lineStats struct {
	Matched bool
	Time    time.Time
	Number  float64
}

func parseLine(line []byte, matchStr *regexp.Regexp, t time.Time) lineStats {
	matches := matchStr.FindSubmatch(line)
	var (
		matched = false
		number  float64
	)
	switch len(matches) {
	case 1:
		matched = true
		number = 0
	case 2:
		matched = true
		n, err := strconv.ParseFloat(string(matches[1]), 64)
		if err != nil {
			n = 0
		}
		number = float64(n)
	}

	return lineStats{
		Matched: matched,
		Time:    t,
		Number:  number,
	}
}
