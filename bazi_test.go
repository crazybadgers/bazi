package bazi

import (
	"testing"
	"time"
)

func TestBaziYear(t *testing.T) {
	// a := assert.New(t)
	cases := []struct {
		date       time.Time
		sex        int
		expectedYP string
	}{
		{time.Date(1979, 11, 18, 0, 0, 0, 0, time.UTC), 0, ""},
	}
	for _, c := range cases {
		b := GetBazi(c.date.Year(), int(c.date.Month()), c.date.Day(), c.date.Hour(), c.date.Minute(), 0, c.sex)
		b.Print()
	}
}
