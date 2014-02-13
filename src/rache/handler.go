package rache

import (
  //"encoding/json"
  //"fmt"
  "net/http"
  "net/url"
)

func RouteSetHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  Logger.Infof("%s %s",r.Method, r.URL)
  params := r.URL.Query()

  fetch_routeset(params)

  //res := //value
  //enc := json.NewEncoder(w)
  //enc.Encode(res)
}

func CacheHandler(w http.ResponseWriter, r *http.Request) {
  /*w.Header().Set("Content-Type", "application/json")*/
  ////res := //value
  //enc := json.NewEncoder(w)
  /*//enc.Encode(res)*/
}

func fetch_routeset(v url.Values) {
  vlabel := v.Get("vlabel")
  app_id := v.Get("app_id")
  ts := v.Get("time")
  cache.Get(vlabel,app_id,ts)
}
