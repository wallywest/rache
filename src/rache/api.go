package rache

import (
  "net/http"
  "github.com/bmizerany/pat"
)

const addr = ":9000"

var cache *Cache

func StartApi(c *Cache) {
  cache = c

  setupHandlers()

  panic(http.ListenAndServe(":9000", nil))
}

func setupHandlers() {
  m := pat.New()
  m.Get("/routeset",http.HandlerFunc(RouteSetHandler))
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
    w.Header().Add("X-Frame-Options", "DENY")
    w.Header().Add("X-Content-Type-Options", "nosniff")
    w.Header().Add("X-XSS-Protection", "1; mode=block")
    // ok, now we're ready to serve the request.
    m.ServeHTTP(w, r)
  })
}
