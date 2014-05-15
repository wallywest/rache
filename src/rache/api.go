package rache

import (
  //"net/http"
  //"github.com/bmizerany/pat"
  "github.com/go-martini/martini"
)

//const addr = ":9000"

var cache *Cache

func StartApi(c *Cache) {
  cache = c
  //setupHandlers()
  //panic(http.ListenAndServe(":9000", nil))
  m := martini.Classic()
  setupHandlers(m)
  m.Run()
}

func setupHandlers(m *martini.ClassicMartini) {
  //m.Get("/routeset",RouteSetIndexHandler)
  //m.Get("/routeset/:app_id",RouteSetGetHandler)
  m.Get("/routeset/:app_id/:route",RouteSetGetHandler)
  m.Post("/routeset/:app_id/:route",RouteSetPostHandler)
  m.Delete("/routeset/:app_id/:route",RouteSetDeleteHandler)
}
