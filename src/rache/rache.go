package rache

import (
  seelog "github.com/cihub/seelog"
  "fmt"
  "strconv"
)


var logger seelog.LoggerInterface

type RouteSet struct{
  Vlabel map[string]string
  AppId int `app_id`
  RouteSet map[string]interface{} `route_set`
}

func NewLogger(c []byte) {
  l ,_ := seelog.LoggerFromConfigAsBytes(c)
  logger = l
}

func FlushLog() {
  logger.Flush()
}

func (r RouteSet) Denormalize(done chan bool) {
  vlabel_s := r.Vlabel["vlabel"]+"|"+strconv.Itoa(r.AppId)
  tree := r.RouteSet["tree"]
  m := tree.(map[string]interface{})
  for s, a:= range m {
    seg := r.findSegment(s)
    for _,id := range a.([]interface{}) {
      allocation := r.findAllocation(id.(string))
      r.buildEntries(vlabel_s,seg,allocation)
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

func(r RouteSet) buildEntries(vlabel string, seg *Segment, alloc *Allocation) {
  entry := vlabel + seg.formattedString()
  for i,d := range alloc.Destinations {
    alloc_string := "|" + alloc.Percentage + "|" + d.Id + "|" + strconv.Itoa(i+1)
    fmt.Println(entry + alloc_string)
  }
}
