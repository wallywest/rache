require 'http'
require 'timecop'
require 'pry'
require 'oj'
require 'json'

bar = ["bar","18181818181"]
#Timecop.freeze(Date.today + 1) do
response = HTTP.get "http://localhost:9000/routeset", 
  :params => {:vlabel => bar.last, :app_id => 8245, :time => Time.now.to_i}
puts response
#end
