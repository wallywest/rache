package rache

import (
  "strconv"
)

type Entry struct {
  Key string
  VlabelIndex string
  AppIdIndex string
  DayBinaryIndex string
  StartTime string
  EndTime string
  Value []DestinationRoute
}

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

type DestinationRoute struct {
  Percentage string `redis:"percentage"`
  Route_order string `redis:"route_order"`
  Destination []byte `redis:"destination"`
}

type DestinationCollection struct {
  Destinations []DestinationRouteJson `json:"destinations"`
}

type DestinationRouteJson struct {
  Percentage string `json:"percentage"`
  Route_order string `json:"route_order"`
  Destination map[string]string `json:"destination"`
}
func(s Segment) formattedString() (formatted string) {
  formatted = s.StartTime + "|" + s.EndTime
  return
}

func(s Segment) bitDays() (d string) {
  i, _ := strconv.ParseInt(s.Days,2,64)
  i = i << 1
  d = strconv.FormatInt(i,10)
  return
}
