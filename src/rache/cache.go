package rache

import (
  "time"
  "github.com/garyburd/redigo/redis"
  "github.com/dchest/uniuri"
  "fmt"
  "encoding/json"
  seelog "github.com/cihub/seelog"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "strconv"
)

const(
  RACHE_DB = "rache_test"
  RACHE_COLLECTION = "rache_test"
  MONGO_URL = "localhost:27017"
  UUIDLen = 20
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
  Session *mgo.Session
}

func NewCache() (cache *Cache){
  config := []byte(cacheLogConfig)
  l ,_ := seelog.LoggerFromConfigAsBytes(config)
  session := setupDB()

  c, err := redis.Dial("tcp", ":6379")
  if err != nil {
    panic("NO REDIS")
  }
  cache = &Cache{Conn:c, Logger:l, Session:session}
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
  c.Session.Close()
}

func(c *Cache) Get(v string, appid string, ts string) []DestinationRouteJson{
  bday,minute := UnixTimeToSegment(ts)
  rkey := c.rangeKey(v,appid,bday)
  dkey := c.dayKey(v,appid,bday)
  return c.findDay(minute,rkey,dkey)
}

func(c *Cache) rangeKey(v string, appid string, bday string) (key string){
  key = "index:" + v + ":" + appid + ":" + bday + ":ranges"
  return
}
func(c *Cache) dayKey(v string, appid string, bday string) (key string){
  key = "index:" + v + ":" + appid + ":" + bday
  return
}

func(c *Cache) findDay(m int, rindex string, dindex string) (results []DestinationRouteJson){
  var getScript = redis.NewScript(1,luaRanges)
  result,err := getScript.Do(c.Conn,0,m,rindex,dindex)
  c.Logger.Infof("%v %s %s",m,rindex,dindex)
  if err != nil {
    fmt.Println("Redis Error: %s", err)
  }

  if result == nil {
    c.Logger.Error("Null result")
  } else {
    switch result.(type) {
    case []interface{}:
      r := result.([]interface{})
      if len(r) == 0 {
        c.Logger.Infof("Filling cache for for %s",dindex)
        c.fillCache()
      } else{
        c.findRoutes(r,results)
      }
    case redis.Error:
      c.fillCache()
    }
  }
  return
}

func(c *Cache) findRoutes(r []interface{}, d []DestinationRouteJson) {
  defer TimeTrack(time.Now(), "Route Lookup")
  for _,value := range r {
    s := string(value.([]byte))
    c.Logger.Infof("Lookup found %s",string(value.([]byte)))
    rd := c.findRoute(s)
    d = append(d,rd)
  }
}

func(c *Cache) findRoute(key string) DestinationRouteJson{
  var d DestinationRoute
  reply, _ := redis.Values(c.Conn.Do("HGETALL",key))
  redis.ScanStruct(reply,&d)
  var destination map[string]string
  json.Unmarshal(d.Destination,&destination)
  djson := DestinationRouteJson{
    Percentage: d.Percentage,
    Route_order: d.Route_order,
    Destination: destination,
  }

  //d.Destination = destination
  return djson
}

func(c *Cache) fillCache() {
  rs := c.findRouteSet()
  done := make(chan bool)
  entry_chan := make(chan Entry)
  go rs.Denormalize(done,entry_chan)

  for {
    select {
    case entry := <- entry_chan:
      c.Save(entry)
      c.Logger.Infof("Saving entry %s",entry)
    case <-done:
      c.Logger.Info("Finished filling cache")
      return
    }
  }
}

func setupDB() *mgo.Session{
  session,err := mgo.Dial(MONGO_URL)
  if err != nil {
    fmt.Println("Mongo session error")
    panic(err)
  }

  session.SetMode(mgo.Monotonic,true)
  return session
}

func(c *Cache) findRouteSet() (rs RouteSet){
  defer TimeTrack(time.Now(), "Query Mongo")
  collection := c.Session.DB(RACHE_DB).C(RACHE_COLLECTION)
  q := bson.M{"app_id":8245,"vlabel.vlabel":"18181818181"}
  err := collection.Find(q).One(&rs)
  if err != nil {
    fmt.Println(err)
  }
  return rs
}

func(c *Cache) Save(entry Entry) {
  defer TimeTrack(time.Now(), "Saving to Cache")
  for _,e := range entry.Value{
    md := c.GenerateKey()
    c.GenerateIndexes(md,entry)
    c.AddValue(md,e)
  }
}

func(c *Cache) GenerateKey() (k string){
  k = uniuri.NewLen(UUIDLen)
  /*key := sha1.New()*/
  //io.WriteString(key,key_str)
  //md := key.Sum(nil)
  //k = hex.EncodeToString(md)
  /*r.CheckDuplicateKey(k)*/
  return
}

func(c *Cache) CheckDuplicateKey (md string) {
  result,err := (c.Conn).Do("EXISTS",md)
  if err != nil {
    c.Logger.Info("Redis Error: ", err)
  }else{
    boo := result.(int64)
    if boo == 1 {
      c.Logger.Infof("Redis Key %s exists",md)
    }
  }
}

func(c *Cache) GenerateIndexes(md string, entry Entry) {
  c.AddIndex(md,entry)
  c.AddAppIdIndex(md,entry)
  c.AddDayIndex(md,entry)
}

func(c *Cache) AddIndex(key string, entry Entry){
  _,err := (c.Conn).Do("SADD",entry.VlabelIndex,key)
  if err != nil {
    fmt.Println(err)
    c.Logger.Error("Redis Error: ", err)
  }
  return
}

func(c *Cache) AddAppIdIndex(key string, entry Entry){
  _,err := (c.Conn).Do("SADD",entry.AppIdIndex,key)
  if err != nil {
    c.Logger.Error("Redis Error: ", err)
  }
  return
}


func(c *Cache) AddDayIndex (key string, entry Entry){
  score,_ := strconv.Atoi(entry.StartTime)
  score2,_ := strconv.Atoi(entry.EndTime)
  _,err := (c.Conn).Do("ZADD",entry.DayBinaryIndex,score,key)
  _,err2 := (c.Conn).Do("SADD",entry.DayBinaryIndex+":ranges",score,score2)
  if err != nil {
    c.Logger.Debugf("Redis Error: %s", err)
  }
  if err2 != nil {
    c.Logger.Debugf("Redis Error: %s", err2)
  }
  return
}

func(c *Cache) AddValue(key string, e DestinationRoute) {
  c.Logger.Infof("Hash is %s",e)
  _,err := (c.Conn).Do("HMSET",key,"percentage",e.Percentage,"route_order",e.Route_order,"destination",e.Destination)
  if err != nil {
    c.Logger.Error("Redis Error: ", err)
  }
}
