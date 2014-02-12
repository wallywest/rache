package rache

import (
  "github.com/garyburd/redigo/redis"
  "fmt"
  seelog "github.com/cihub/seelog"

)

var cacheLogConfig= `
<seelog>
<outputs>
<file path="./log/cache.log"/>
</outputs>
</seelog>
`

type Cache struct {
  Conn redis.Conn
  Logger seelog.LoggerInterface
}

func NewCache() (cache *Cache){
  config := []byte(cacheLogConfig)
  l ,_ := seelog.LoggerFromConfigAsBytes(config)
  c, err := redis.Dial("tcp", ":6379")
  if err != nil {
    panic("NO REDIS")
  }
  cache = &Cache{Conn:c, Logger:l}
  return
}

func(c *Cache) Info() {
  _, err := c.Conn.Do("INFO")
  if err != nil {
    fmt.Println(err)
  }
}

func(c *Cache) Close() {
  c.Conn.Close()
}

func(c *Cache) Get(v string, appid string, ts string) {
  bday,minute := UnixTimeToSegment(ts)
  rkey := c.rangeKey(v,appid,bday)
  dkey := c.dayKey(v,appid,bday)
  c.findDay(minute,rkey,dkey)
}

func(c *Cache) rangeKey(v string, appid string, bday string) (key string){
  key = "index:" + v + ":" + appid + ":" + bday + ":ranges"
  return
}
func(c *Cache) dayKey(v string, appid string, bday string) (key string){
  key = "index:" + v + ":" + appid + ":" + bday
  return
}

func(c *Cache) findDay(m int, rindex string, dindex string) {
  var getScript = redis.NewScript(1,luaRanges)
  result,err := getScript.Do(c.Conn,0,m,rindex,dindex)
  c.Logger.Infof("%v %s %s",m,rindex,dindex)
  if err != nil {
    fmt.Println("Redis Error: %s", err)
  }
  if result == nil {
    c.Logger.Infof("Cache miss for %s",dindex)
  }else {
    r := result.([]interface{})
    for _,value := range r {
      c.Logger.Infof("Lookup found %s",string(value.([]byte)))
    }
  }
}
