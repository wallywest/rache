package main

import(
  "runtime/pprof"
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
  collection := session.DB(RACHE_DB).C(RACHE_COLLECTION)
  route := findRouteSet(collection)
  done := make(chan bool)

  count := 0
  limit := 100
  for i:= 0; i < limit; i++ {
    go route.Denormalize(done)
  }

  for {
    select {
    case <-done:
      if count == limit-1 {
        fmt.Println("quitting")
        os.Exit(1)
      }else{
        count = count + 1
      }
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

func findRouteSet(c *mgo.Collection) (rs rache.RouteSet){
  q := bson.M{"app_id":8245,"vlabel.vlabel":"18181818181"}
  err := c.Find(q).One(&rs)
  if err != nil {
    fmt.Println(err)
  }
  return rs
}
