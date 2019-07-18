package main

type Course struct {
	Name    string `json:"KCM"`
	Teacher string `json:"SKJS"`
	Data    []*Detail
}

func NewCourse(name string, teacher string) *Course {
	return &Course{Name: name, Teacher: teacher}
}

type Courses []*Course

type Detail struct {
	Week     []int  `json:"ZCMC"`
	IndexDay int    `json:"SKXQ"`
	Begin    int    `json:"KSJC"`
	End      int    `json:"JSJC"`
	Place    string `json:"JASMC"`
}

func NewDetail(week []int, indexDay int, begin int, end int, place string) *Detail {
	return &Detail{Week: week, IndexDay: indexDay, Begin: begin, End: end, Place: place}
}

func (c *Course) itemParser(s map[string]interface{}) {
	c.Data = append(c.Data, parseDetail(s))
}

func (cs *Courses) CheckRepeat(name string, teacher string) (bool, int) {
	for i, v := range *cs {
		if v.Name == name && v.Teacher == teacher {
			return true, i
		}
	}
	return false, -1
}
