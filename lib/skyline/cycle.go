package skyline

import (
	"fmt"
)

type Cycle struct {
	Name   string
	Number int64
	Count  int64
	Period int64
	total  float64
	Id     int64
	Params map[string]interface{}
}

func NewCycle(id int64, name string, period int64) *Cycle {
	return &Cycle{
		Number: 0,
		Count:  0,
		Period: period,
		total:  0,
		Id:     id,
		Name:   name,
		Params: make(map[string]interface{}, 3),
	}
}

func (self *Cycle) AddLine(line *lineStats) {
	self.Count += 1
	self.total += line.Number
	cnt_key := fmt.Sprintf("%s_cnt", self.Name)
	avg_key := fmt.Sprintf("%s_avg", self.Name)
	cps_key := fmt.Sprintf("%s_cps", self.Name)
	self.Params[cnt_key] = self.Count
	self.Params[cps_key] = self.Rate()
	self.Params[avg_key] = self.Averate()
}

func (self *Cycle) Averate() float64 {
	return self.total / float64(self.Count)
}

func (self *Cycle) Rate() int64 {
	return self.Count / self.Period
}
