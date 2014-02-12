package rache

import (
  "github.com/dchest/uniuri"
  seelog "github.com/cihub/seelog"
  /*"io"*/
  //"encoding/hex"
  /*"crypto/sha1"*/
  "strconv"
  "github.com/garyburd/redigo/redis"
  "fmt"
)
//var Redis redis.Conn

var redisConfig= `
<seelog>
<outputs>
<file path="./log/redis.log"/>
</outputs>
</seelog>
`

const(
  UUIDLen = 20
)
type RedisClient struct {
  Conn *redis.Conn
  Logger seelog.LoggerInterface
}

func NewRedisClient() *RedisClient{
  config := []byte(redisConfig)
  l ,_ := seelog.LoggerFromConfigAsBytes(config)

  conn, err := redis.Dial("tcp", ":6379")
  if err != nil {
    panic("NO REDIS")
  }
  client := new(RedisClient)
  client.Conn = &conn
  client.Logger = l

  return client
}

func(r *RedisClient) Info() {
  reply, err := (*r.Conn).Do("INFO")
  if err != nil {
    fmt.Println(err)
  }
  strs,err := redis.String(reply,err)
  fmt.Print(strs)
}

func(r *RedisClient) Close() {
  (*r.Conn).Close()
}

func(r *RedisClient) FillCache(entry Entry) {
  md := r.GenerateKey(entry)
  r.GenerateIndexes(md,entry)
  r.AddValue(md,entry)
}

func(r *RedisClient) GenerateKey(entry Entry) (k string){
  k = uniuri.NewLen(UUIDLen)
  /*key := sha1.New()*/
  //io.WriteString(key,key_str)
  //md := key.Sum(nil)
  //k = hex.EncodeToString(md)
  /*r.CheckDuplicateKey(k)*/
  return
}

func(r *RedisClient) CheckDuplicateKey (md string) {
  result,err := (*r.Conn).Do("EXISTS",md)
  if err != nil {
    r.Logger.Info("Redis Error: ", err)
  }else{
    boo := result.(int64)
    if boo == 1 {
      r.Logger.Infof("Redis Key %s exists",md)
    }
  }
}

func(r *RedisClient) GenerateIndexes(md string, entry Entry) {
  r.AddIndex(md,entry)
  r.AddAppIdIndex(md,entry)
  r.AddDayIndex(md,entry)
}

func(r *RedisClient) AddIndex(key string, entry Entry){
  _,err := (*r.Conn).Do("SADD",entry.VlabelIndex,key)
  if err != nil {
    fmt.Println(err)
    r.Logger.Error("Redis Error: ", err)
  }
  return
}

func(r *RedisClient) AddAppIdIndex(key string, entry Entry){
  _,err := (*r.Conn).Do("SADD",entry.AppIdIndex,key)
  if err != nil {
    r.Logger.Error("Redis Error: ", err)
  }
  return
}


func(r *RedisClient) AddDayIndex (key string, entry Entry){
  score,_ := strconv.Atoi(entry.StartTime)
  score2,_ := strconv.Atoi(entry.EndTime)
  _,err := (*r.Conn).Do("ZADD",entry.DayBinaryIndex,score,key)
  _,err2 := (*r.Conn).Do("SADD",entry.DayBinaryIndex+":ranges",score,score2)
  if err != nil {
    r.Logger.Debugf("Redis Error: %s", err)
  }
  if err2 != nil {
    r.Logger.Debugf("Redis Error: %s", err2)
  }
  return
}

func(r *RedisClient) AddValue(key string, entry Entry) {
  for _,e := range entry.Value{
    _,err := (*r.Conn).Do("RPUSH",key,e)
    if err != nil {
      r.Logger.Error("Redis Error: ", err)
    }
  }
}
