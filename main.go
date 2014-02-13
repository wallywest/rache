package main

import(
  "runtime/pprof"
  "os"
  "flag"
  "rache"
  "fmt"
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

  /*count := 0*/
  //max,_ := strconv.Atoi(*limit)
  //for i:= 0; i < max; i++ {
    //go route.Denormalize(done,cache_entry)
  /*}*/

  cache := rache.NewCache()
  defer cache.Close()
  rache.StartApi(cache)
}
