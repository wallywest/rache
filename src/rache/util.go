package rache

import(
  "time"
  "strconv"
)

var DAYS = map[string]string{
  "Sunday":"128",
  "Monday":"64",
  "Tuesday":"32",
  "Wednesday":"16",
  "Thursday":"8",
  "Friday":"4",
  "Saturday":"2",
}

func UnixTimeToSegment(ts string) (day string, m int){
  t ,_:= strconv.ParseInt(ts,10,64)
  tdate := time.Unix(t,0)
  day = DAYS[tdate.Weekday().String()]
  m = HourToMinutes(tdate.Hour()) + tdate.Minute()
  return
}

func HourToMinutes(hour int) (m int){
  m = 60 * hour
  return
}

