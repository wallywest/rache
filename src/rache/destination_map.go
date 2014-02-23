package rache

import (
)

type DestinationProperty struct {
  Id string `json:"id"`
  Type string `json:"type"`
  Properties map[string]string `json:"properties"`
}

var ids = []string{
  "0","1","2","3","4","5","6","7","8","9","10","11",
}

var types = []string{
  "Prompt","VlabelMap","Destination",
}


type DestinationMap struct {
  Map []DestinationProperty
}

func NewDestinationMap() DestinationMap{
  var m []DestinationProperty
  for _,t := range types {
    for _,v := range ids {
      d := DestinationProperty{
        Id: v,
        Type: t,
        Properties: RandomDestinationProperty(),
      }
      m = append(m,d)
    }
  }

  return DestinationMap{
    Map: m,
  }
}

func(d *DestinationMap) findDestination(id string, t string) *DestinationProperty{
  for _,k := range d.Map {
    if (k.Id == id && k.Type == t) {
      return &k
    }
  }
  return nil
}
