package main

import(
  "runtime/pprof"
  //"strconv"
  "os"
  "flag"
  "rache"
  "fmt"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
)
const(
  RACHE_DB = "rache_test"
  RACHE_COLLECTION = "rache_test"
  MONGO_URL = "localhost:27017"
)


var testConfig = `
<seelog>
<outputs>
<file path="./log/main.log"/>
</outputs>
</seelog>
`

var prof = flag.String("prof", "", "write cpu profile to file")
var limit = flag.String("limit", "1", "write cpu profile to file")
//var fillcache = flagString("cache","false","prefill cache with values")

func main(){
  flag.Parse()

  config := []byte(testConfig)

  rache.NewLogger(config)
  defer rache.FlushLog()

  if *prof != "" {
    fmt.Println("Profiling CPU")
    f, err := os.Create(*prof)
    if err != nil {
      panic(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
  }

  session := setupDB()
  defer session.Close()

  /*rc := rache.NewRedisClient()*/
  /*defer rc.Close()*/

  //collection := session.DB(RACHE_DB).C(RACHE_COLLECTION)
  /*route := findRouteSet(collection)*/
  //done := make(chan bool)
  /*cache_entry := make(chan rache.Entry)*/


  /*count := 0*/
  //max,_ := strconv.Atoi(*limit)
  //for i:= 0; i < max; i++ {
    //go route.Denormalize(done,cache_entry)
  /*}*/

  cache := rache.NewCache()
  defer cache.Close()
  rache.StartApi(cache)

  /*for {*/
    //select {
    //case entry := <- cache_entry:
      //rc.FillCache(entry)
      //go func(e rache.Entry) {
       //logger.Info(entry)
      //}(entry)
    //case <-done:
      //if count == max-1 {
        //fmt.Println("quitting")
        //os.Exit(1)
      //}else{
        //count = count + 1
      //}
    //}
  /*}*/
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

func findRouteSet(c *mgo.Collection) (rs rache.RouteSet){
  q := bson.M{"app_id":8245,"vlabel.vlabel":"18181818181"}
  err := c.Find(q).One(&rs)
  if err != nil {
    fmt.Println(err)
  }
  return rs
}
