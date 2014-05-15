package rache

import (
  "encoding/json"
  "net/http"
  "fmt"
  "github.com/go-martini/martini"
)

func RouteSetHandler(w http.ResponseWriter, r *http.Request) {
  params := r.URL.Query()
  vlabel := params.Get("vlabel")
  app_id := params.Get("app_id")
  ts := params.Get("time")
  values := cache.Get(vlabel,app_id,ts)
  v := map[string][]DestinationRouteJson{
    "destinations":values,
  }
  enc := json.NewEncoder(w)
  enc.Encode(v)
}

func RouteSetIndexHandler(res http.ResponseWriter, req *http.Request){
}

func RouteSetGetHandler(res http.ResponseWriter, req *http.Request, params martini.Params){
  fmt.Println(params)
  //params := req.URL.Query()
  /*vlabel := params.Get("vlabel")*/
  //app_id := params.Get("app_id")
  //ts := params.Get("time")
  
  /*values := cache.Get(vlabel,app_id,ts)*/
  //v := map[string][]DestinationRouteJson{
    //"destinations":values,
  //}
  //enc := json.NewEncoder(w)
  /*enc.Encode(v)*/
}

func RouteSetPostHandler(res http.ResponseWriter, req *http.Request, params martini.Params){
  //Update routeset for a particular route
  //Update local caches

  fmt.Println(params)
}

func RouteSetDeleteHandler(res http.ResponseWriter, req *http.Request, params martini.Params){
  //Delete routeset for a particular route
  //Expire from local caches

  fmt.Println(params)
}


func CacheHandler(w http.ResponseWriter, r *http.Request) {
  /*w.Header().Set("Content-Type", "application/json")*/
  ////res := //value
  //enc := json.NewEncoder(w)
  /*//enc.Encode(res)*/
}
