package skyline

import (
	"fmt"
	"strings"

	"strconv"

	"os/exec"

	"github.com/Knetic/govaluate"
)

type FilterWarnings struct {
	Warnings []*Warning
}

func (filterWarnings *FilterWarnings) addWarning(warning *Warning) {
	filterWarnings.Warnings = append(filterWarnings.Warnings, warning)
}

func (filterWarnings *FilterWarnings) Size() int {
	return len(filterWarnings.Warnings)
}

type Warning struct {
	conf         *WarningConf
	Formula      string
	FilterTotal  int
	FilterHigh   int
	AlertName    string
	AlertCommand string
	expression   *govaluate.EvaluableExpression
	IsWarning    bool
	MatchedCount int
}

func NewWarning(conf *WarningConf) *Warning {
	expr, err := govaluate.NewEvaluableExpression(conf.Formula)
	if err != nil {
		panic(err.Error())
	}

	high, total := parseWarningFilter(conf.WarningFilter)

	return &Warning{
		conf:         conf,
		Formula:      conf.Formula,
		FilterHigh:   high,
		FilterTotal:  total,
		expression:   expr,
		AlertName:    conf.AlertName,
		AlertCommand: conf.AlertCommand,
		IsWarning:    false,
	}
}

func parseWarningFilter(filterStr string) (int, int) {
	segments := strings.Split(filterStr, "/")
	if len(segments) != 2 {
		return 0, 0
	}

	high, err := strconv.ParseInt(segments[0], 10, 32)
	if err != nil {
		return 0, 0
	}
	total, err := strconv.ParseInt(segments[1], 10, 32)
	if err != nil {
		return 0, 0
	}

	return int(high), int(total)
}

func (warning *Warning) NeedTrigger(params map[string]interface{}) bool {
	matched, err := warning.expression.Evaluate(params)
	if err != nil {
		return false
	}
	return matched.(bool)
}

func (warning *Warning) Warn(msg string) {
	if warning.AlertCommand == "" {
		return
	}
	go func(msg string) {
		args := strings.Fields(warning.AlertCommand)
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = strings.NewReader(msg)
		err := cmd.Run()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(msg)
}

func (warning *Warning) itemName() string {

	for _, token := range warning.expression.Tokens() {
		if token.Kind != govaluate.VARIABLE {
			continue
		}
		itemName := removeStatsSuffix(token.Value.(string))
		return itemName
	}
	return ""
}

func removeStatsSuffix(s string) string {
	for _, suffix := range [...]string{"_cnt", "_avg", "_cps"} {
		trimed := strings.TrimSuffix(s, suffix)
		if trimed != s {
			return trimed
		}
	}
	return s
}

type WarningCenter struct {
	warnings map[string]*FilterWarnings
}

func NewWarningCenter(warningConfs []WarningConf) *WarningCenter {
	warnings := make(map[string]*FilterWarnings, 8)
	for _, conf := range warningConfs {
		warning := NewWarning(&conf)
		itemName := warning.itemName()
		if warnings[itemName] == nil {
			warnings[itemName] = &FilterWarnings{
				Warnings: make([]*Warning, 0),
			}
		}
		warnings[itemName].addWarning(warning)
	}
	return &WarningCenter{
		warnings: warnings,
	}
}

func (center *WarningCenter) FindfilterWarnings(name string) (*FilterWarnings, bool) {
	filterWarnings, ok := center.warnings[name]
	if !ok {
		return nil, false
	}
	return filterWarnings, true
}
