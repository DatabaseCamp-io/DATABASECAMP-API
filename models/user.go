package models

type ChangePointMode string

var Mode = struct {
	Add    ChangePointMode
	Reduce ChangePointMode
}{
	"+",
	"-",
}
