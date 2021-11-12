package models

type content struct {
	id       int
	name     string
	activity []*int
}

type group struct {
	id      int
	name    string
	content map[int]*content
}

type Overview struct {
	group              map[int]*group
	activityContentMap map[int]int
	activityCount      map[int]int
}
