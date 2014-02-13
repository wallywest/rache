local minute = tonumber(ARGV[1])
local set = ARGV[2]
local daySet = ARGV[3]
local ranges = {}
local out = {}
local ranges = redis.call("smembers",set)
local dayranges = {}

for k,v in pairs(ranges) do
  local i = tonumber(v)
  out[k] = i
end
table.sort(out)

local IsEven = function(num)
  return num % 2 == 0
end

local lookup = -1
local count = 0

for i,v in pairs(out) do
  if v >= minute then
    if IsEven(i) then
      lookup = out[i-1]
    else
      lookup = out[i]
    end
    break
  end
  count = i
end
if lookup == -1 then 
  lookup = out[count]
end
print(lookup)

dayranges = redis.call("zrangebyscore",daySet,lookup,lookup)

return dayranges
