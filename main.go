package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/gfreezy/skyline/lib/skyline"
	"github.com/hpcloud/tail"
	getopt "github.com/kesselborn/go-getopt"
)

func main() {
	optionDefinition := getopt.Options{
		"Monitor logs",
		getopt.Definitions{
			{"config|c", "config file", getopt.Required, ""},
		},
	}

	options, _, _, e := optionDefinition.ParseCommandLine()

	help, wantsHelp := options["help"]

	if e != nil || wantsHelp {
		exitCode := 0

		switch {
		case wantsHelp && help.String == "usage":
			fmt.Print(optionDefinition.Usage())
		case wantsHelp && help.String == "help":
			fmt.Print(optionDefinition.Help())
		default:
			fmt.Printf("**** Error: %s\n\n%s", e.Error(), optionDefinition.Help())
			exitCode = e.ErrorCode
		}
		os.Exit(exitCode)
	}

	conf, err := skyline.LoadConfig(options["config"].String)
	if err != nil {
		fmt.Printf("**** Error: %s\n\n%s", err.Error(), optionDefinition.Help())
		os.Exit(-1)
	}

	warningCenter := skyline.NewWarningCenter(conf.Warnings)
	var wg sync.WaitGroup

	for _, monitorConf := range conf.Monitors {
		wg.Add(1)
		go func(c skyline.MonitorConf, center *skyline.WarningCenter) {
			monitor(c, center)
			wg.Done()
		}(monitorConf, warningCenter)
	}
	wg.Wait()
}

func monitor(monitorConf skyline.MonitorConf, warningCenter *skyline.WarningCenter) {
	var filename = monitorConf.LogFilePath
	t, err := tail.TailFile(filename, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		panic(err)
	}

	filters := make([]*skyline.Filter, 0, len(monitorConf.FilterItems))
	for _, filterConf := range monitorConf.FilterItems {
		filterName := fmt.Sprintf("%s_%s", monitorConf.LogNamePrefix, filterConf.ItemNamePrefix)
		filter := skyline.NewFilter(filterName, &filterConf)
		filters = append(filters, filter)
		filterWarnings, ok := warningCenter.FindfilterWarnings(filter.Name)
		if !ok {
			continue
		}
		fmt.Println(filterWarnings.Warnings[0].Formula)
		go warn(filterWarnings, filter.CycleStatsChannel)
	}

	for line := range t.Lines {
		for _, f := range filters {
			if f == nil {
				continue
			}
			f.AddLine([]byte(line.Text), line.Time, true)
		}
	}
}