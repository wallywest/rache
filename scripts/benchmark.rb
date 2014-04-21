require 'http'
require 'timecop'
require 'pry'
require 'oj'
require 'json'
require 'benchmark'

bar = ["bar","18181818181"]
#Timecop.freeze(Date.today + 1) do
response = HTTP.get("http://localhost:9000/routeset", :params => {:vlabel => bar.last, :app_id => 8245, :time => Time.now.to_i})
p response

##Benchmark.bm(20) do |x|
  ##x.report("with caching 100 hits") do
    ##100.times do
      ##response = http.get "http://localhost:9000/routeset", :params => {:vlabel => bar.last, :app_id => 8245, :time => time.now.to_i}
    ##end
  ##end
  ##x.report("with caching 1000 hits") do
    ##1000.times do
      ##response = HTTP.get "http://localhost:9000/routeset",:params => {:vlabel => bar.last, :app_id => 8245, :time => Time.now.to_i}
    ##end
  ##end
##end
