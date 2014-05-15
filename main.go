package main

import(
  "runtime/pprof"
  "os"
  "flag"
  "rache"
  "fmt"
)

var prof = flag.String("prof", "", "write cpu profile to file")
var limit = flag.String("limit", "1", "write cpu profile to file")
//var fillcache = flagString("cache","false","prefill cache with values")

func main(){
  flag.Parse()

  if *prof != "" {
    fmt.Println("Profiling CPU")
    f, err := os.Create(*prof)
    if err != nil {
      panic(err)
    }
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
  }
  rache.Run()
}
