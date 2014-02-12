package rache

import (
  "github.com/garyburd/redigo/redis"
  "fmt"
)
//var Redis redis.Conn

type Cache struct {
  Conn redis.Conn
}

func NewCache() (cache *Cache){
  c, err := redis.Dial("tcp", ":6379")
  if err != nil {
    panic("NO REDIS")
  }
  defer c.Close()

  cache = &Cache{Conn:c}
  return
}

func(c *Cache) Info() {
  n, err := c.Conn.Do("INFO")
  if err != nil {
    fmt.Println(err)
  }
  fmt.Println(n)
}
