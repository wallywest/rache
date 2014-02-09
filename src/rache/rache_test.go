package rache

import (
  "fmt"
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  "testing"
)

var rs RouteSet

func BenchmarkRacheDenormalize15(b *testing.B) {
  benchmarkRacheDenormalize(b)
}

func BenchmarkRacheDenormalize100(b *testing.B) {
  benchmarkRacheDenormalize(b)
}
func BenchmarkRacheDenormalize300(b *testing.B) {
  benchmarkRacheDenormalize(b)
}
func BenchmarkRacheDenormalize1000(b *testing.B) {
  benchmarkRacheDenormalize(b)
}
func BenchmarkRacheDenormalize3000(b *testing.B) {
  benchmarkRacheDenormalize(b)
}
func BenchmarkRacheDenormalize10000(b *testing.B) {
  benchmarkRacheDenormalize(b)
}


func benchmarkRacheDenormalize(b *testing.B){
  setRouteSet()
  b.ResetTimer()

  for n:=0; n < b.N; n++ {
    rs.Denormalize()
  }
}

func setRouteSet(){
  const(
    RACHE_DB = "rache_test"
    RACHE_COLLECTION = "rache_test"
    MONGO_URL = "localhost:27017"
  )

  session,e := mgo.Dial(MONGO_URL)
  if e != nil {
    fmt.Println("Mongo session error")
    panic(e)
  }

  session.SetMode(mgo.Monotonic,true)
  collection := session.DB(RACHE_DB).C(RACHE_COLLECTION)


  q := bson.M{"app_id":8245,"vlabel.vlabel":"18181818181"}
  err := collection.Find(q).One(&rs)
  if err != nil {
    fmt.Println(err)
  }
  //return rs
}
