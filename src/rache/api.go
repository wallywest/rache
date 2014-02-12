package rache

import (
  //"encoding/json"
  //"fmt"
  "net/http"
  "net/url"
)

const addr = ":9000"

var cache *Cache
func StartApi(c *Cache) {
  cache = c
  http.HandleFunc("/routeset", routeset_handler)
  http.HandleFunc("/cache", cache_handler)

  err := http.ListenAndServe(addr, nil)
  if err != nil {
    Logger.Error("ListenAndServe:", err)
  }
}

func fetch_routeset(v url.Values) {
  vlabel := v.Get("vlabel")
  app_id := v.Get("app_id")
  ts := v.Get("time")
  cache.Get(vlabel,app_id,ts)
  //Logger.Infof("%v",values)
}

func routeset_handler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  Logger.Infof("%s %s",r.Method, r.URL)
  params := r.URL.Query()

  fetch_routeset(params)

  //res := //value
  //enc := json.NewEncoder(w)
  //enc.Encode(res)
}

func cache_handler(w http.ResponseWriter, r *http.Request) {
  /*w.Header().Set("Content-Type", "application/json")*/
  ////res := //value
  //enc := json.NewEncoder(w)
  /*//enc.Encode(res)*/
}
