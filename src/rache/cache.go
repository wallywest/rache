package rache

import (
  "github.com/garyburd/redigo/redis"
  "github.com/dchest/uniuri"
  "fmt"
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
    c.Logger.Error("Null result")
  } else {
    r := result.([]interface{})
    if len(r) == 0 {
      c.Logger.Infof("Filling cache for for %s",dindex)
      c.fillCache()
    } else{
      for _,value := range r {
        c.Logger.Infof("Lookup found %s",string(value.([]byte)))
      }
    }
  }
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
  collection := c.Session.DB(RACHE_DB).C(RACHE_COLLECTION)
  q := bson.M{"app_id":8245,"vlabel.vlabel":"18181818181"}
  err := collection.Find(q).One(&rs)
  if err != nil {
    fmt.Println(err)
  }
  return rs
}

func(c *Cache) Save(entry Entry) {
  md := c.GenerateKey(entry)
  c.GenerateIndexes(md,entry)
  c.AddValue(md,entry)
}

func(c *Cache) GenerateKey(entry Entry) (k string){
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

func(c *Cache) AddValue(key string, entry Entry) {
  for _,e := range entry.Value{
    _,err := (c.Conn).Do("RPUSH",key,e)
    if err != nil {
      c.Logger.Error("Redis Error: ", err)
    }
  }
}
