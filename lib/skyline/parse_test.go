package skyline

import (
	"regexp"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	lineStr := []byte("wwwbaidu 1.23")
	matchStr := regexp.MustCompile("(\\d+.?\\d*)")
	line := parseLine(lineStr, matchStr, time.Now())
	if !(line.Matched && line.Number == 1.23) {
		t.Fail()
	}

	line = parseLine([]byte("wwwbaidu 0.3"), matchStr, time.Now())
	if !(line.Matched && line.Number == 0.3) {
		t.Fail()
	}

	line = parseLine([]byte("wwwbaidu 3"), matchStr, time.Now())
	if !(line.Matched && line.Number == 3) {
		t.Fail()
	}

	line = parseLine([]byte("wwwbaidu"), matchStr, time.Now())
	if line.Matched {
		t.Fail()
	}
}
