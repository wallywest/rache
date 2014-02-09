package rache

import (
  "strconv"
)

type Segment struct{
  Days string
  StartTime string
  EndTime string
}

type Allocation struct {
  Percentage string
  Destinations []Destination
}

type Destination struct {
  Id string
  Type string
}

func(s Segment) formattedString() (formatted string) {
  formatted = "|" + s.bitDays() + "|" + s.StartTime + "|" + s.EndTime
  return
}

func(s Segment) bitDays() (d string) {
  i, _ := strconv.ParseInt(s.Days,2,64)
  i = i << 1
  d = strconv.FormatInt(i,10)
  return
}
