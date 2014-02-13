package rache

import (
  "encoding/json"
  "time"
  "net/http"
)

func RouteSetHandler(w http.ResponseWriter, r *http.Request) {
  defer TimeTrack(time.Now(), "Api Request")

  w.Header().Set("Content-Type", "application/json")
  Logger.Infof("%s %s",r.Method, r.URL)
  params := r.URL.Query()
  vlabel := params.Get("vlabel")
  app_id := params.Get("app_id")
  ts := params.Get("time")
  values := cache.Get(vlabel,app_id,ts)
  enc := json.NewEncoder(w)
  enc.Encode(values)
}

func CacheHandler(w http.ResponseWriter, r *http.Request) {
  /*w.Header().Set("Content-Type", "application/json")*/
  ////res := //value
  //enc := json.NewEncoder(w)
  /*//enc.Encode(res)*/
}
