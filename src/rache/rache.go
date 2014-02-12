package rache

import (
  seelog "github.com/cihub/seelog"
  "strconv"
)


var Logger seelog.LoggerInterface

type RouteSet struct{
  Vlabel map[string]string
  AppId int `app_id`
  RouteSet map[string]interface{} `route_set`
  EntryChan chan Entry
}

func NewLogger(c []byte) (l seelog.LoggerInterface){
  l ,_ = seelog.LoggerFromConfigAsBytes(c)
  Logger = l
  return
}

func FlushLog() {
  Logger.Flush()
}

func (r RouteSet) Denormalize(done chan bool, cache_chan chan Entry) {
  r.EntryChan = cache_chan
  tree := r.RouteSet["tree"]
  m := tree.(map[string]interface{})
  for s, a:= range m {
    seg := r.findSegment(s)
    for _,id := range a.([]interface{}) {
      allocation := r.findAllocation(id.(string))
      r.buildEntries(seg,allocation)
    }
  }
  done <- true
}


func(r RouteSet) findAllocation(s string) *Allocation{
  var alloc Allocation
  allocations := r.RouteSet["allocations"].([]interface{})
  for _,t := range allocations {
    value := t.(map[string]interface{})
    if value["id"] == s {
      destinations := r.newDestinations(value["destinations"].([]interface{}))
      alloc = Allocation{
        Percentage:value["percentage"].(string),
        Destinations:destinations,
      }
    }
  }

  return &alloc
}

func(r RouteSet) newDestinations(ds []interface{}) (destinations []Destination){
  for _,destination := range ds {
    d_interface := destination.(map[string]interface{})
    d := &Destination{
      Id:d_interface["destination_id"].(string),
      Type:d_interface["type"].(string),
    }
    destinations = append(destinations,*d)
  }
  return
}

func(r RouteSet) findSegment(s string) *Segment{
  var seg Segment
  segments := r.RouteSet["segments"].([]interface{})
  for _,t := range segments {
    value := t.(map[string]interface{})
    if value["id"] == s {
      seg = Segment{
        Days:value["days"].(string), 
        StartTime:value["start_time"].(string),
        EndTime:value["end_time"].(string),
      }
    }
  }
  return &seg
}

func(r RouteSet) buildEntries(seg *Segment, alloc *Allocation) {
  var allocs []string
  key :=   r.Vlabel["vlabel"] + strconv.Itoa(r.AppId) + seg.bitDays() + seg.StartTime
  vlabelindex := "index:"+r.Vlabel["vlabel"]
  appidindex := "index:"+r.Vlabel["vlabel"] + ":" + strconv.Itoa(r.AppId)
  dayindex:= "index:"+r.Vlabel["vlabel"] + ":" + strconv.Itoa(r.AppId) + ":" + seg.bitDays()

  e := Entry{
    Key: key,
    VlabelIndex: vlabelindex,
    AppIdIndex: appidindex,
    DayBinaryIndex: dayindex,
    StartTime: seg.StartTime,
    EndTime: seg.EndTime,
    Value: allocs,
  }

  for i,d := range alloc.Destinations {
    alloc_string := alloc.Percentage + "|" + d.Id + "|" + d.Type + "|" + strconv.Itoa(i+1)
    e.Value = append(e.Value,alloc_string)
  }
  r.EntryChan <- e
}
