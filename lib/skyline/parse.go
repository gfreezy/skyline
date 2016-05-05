package skyline

import (
	"regexp"
	"strconv"
	"time"
)

type lineStats struct {
	Matched bool
	Time    time.Time
	Number  float32
}

func parseLine(line []byte, matchStr *regexp.Regexp, t time.Time) lineStats {
	matches := matchStr.FindSubmatch(line)
	var (
		matched = false
		number  float32
	)
	switch len(matches) {
	case 1:
		matched = true
		number = 0
	case 2:
		matched = true
		n, err := strconv.ParseFloat(string(matches[1]), 32)
		if err != nil {
			n = 0
		}
		number = float32(n)
	}

	return lineStats{
		Matched: matched,
		Time:    t,
		Number:  number,
	}
}
