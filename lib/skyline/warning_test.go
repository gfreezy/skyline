package skyline

import "testing"

func TestRemoveStatsSuffix(t *testing.T) {
	if removeStatsSuffix("adflk_cnt") != "adflk" {
		t.Fail()
	}

	if removeStatsSuffix("a__cnt") != "a_" {
		t.Fail()
	}

	if removeStatsSuffix("adflk_avg") != "adflk" {
		t.Fail()
	}

	if removeStatsSuffix("adflk_cps") != "adflk" {
		t.Fail()
	}
}

func TestWarningItemName(t *testing.T) {
	warning := NewWarning(&WarningConf{
		AlertName:     "alert",
		Formula:       "nginx_a_cnt > 2",
		WarningName:   "warning-name",
		WarningFilter: "3/5",
	})
	if warning.itemName() != "nginx_a" {
		t.Fail()
	}

	warning2 := NewWarning(&WarningConf{
		AlertName:     "alert",
		Formula:       "nginx_a_avg > 2",
		WarningName:   "warning-name",
		WarningFilter: "3/5",
	})
	if warning2.itemName() != "nginx_a" {
		t.Fail()
	}
}

func TestEvaluate(t *testing.T) {
	warning := NewWarning(&WarningConf{
		AlertName:     "alert",
		Formula:       "nginx_a_avg > 2",
		WarningName:   "warning-name",
		WarningFilter: "3/5",
	})
	params := make(map[string]interface{}, 3)
	params["nginx_a_avg"] = 1
	if warning.NeedTrigger(params) {
		t.Fail()
	}

	params["nginx_a_avg"] = 3
	if !warning.NeedTrigger(params) {
		t.Fail()
	}

	warning = NewWarning(&WarningConf{
		AlertName:     "alert",
		Formula:       "nginx_a_avg > 1.1",
		WarningName:   "warning-name",
		WarningFilter: "3/5",
	})
	params["nginx_a_avg"] = 0.5
	if warning.NeedTrigger(params) {
		t.Fail()
	}

	params["nginx_a_avg"] = 3.9
	if !warning.NeedTrigger(params) {
		t.Fail()
	}
}

func TestFindFilterWarnings(t *testing.T) {
	confs := make([]WarningConf, 0)
	confs = append(confs, WarningConf{
		AlertName:     "alert",
		Formula:       "nginx_a_avg > 2",
		WarningName:   "warning-name",
		WarningFilter: "3/5",
	}, WarningConf{
		AlertName:     "alert2",
		Formula:       "nginx_a_cnt > 2",
		WarningName:   "warning-name",
		WarningFilter: "3/5",
	})
	center := NewWarningCenter(confs)
	filterWarnings, ok := center.FindfilterWarnings("nginx_a")
	if !ok || filterWarnings.Size() != 2 {
		t.Fail()
	}
}
