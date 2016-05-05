package skyline

type Cycle struct {
	Number int64
	Count  int64
	Period int64
	total  float32
	Id     int64
}

func NewCycle(id int64, period int64) *Cycle {
	return &Cycle{
		Number: 0,
		Count:  0,
		Period: period,
		total:  0,
		Id:     id,
	}
}

func (self *Cycle) AddLine(line *lineStats) {
	self.Count += 1
	self.total += line.Number
}

func (self *Cycle) Averate() float32 {
	return self.total / float32(self.Count)
}

func (self *Cycle) Rate() int64 {
	return self.Count / self.Period
}
