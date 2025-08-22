package basemodel

import "time"

type TDateTime struct {
	time.Time
}

type TDateOnly struct {
	time.Time
}

type TString string
