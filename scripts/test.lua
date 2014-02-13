local range = {}
local set = tonumber(arg[1])

range[1] = 0
range[2] = 59
range[3] = 60
range[4] = 1400
local lookup = -1

function IsEven(num)
  return num % 2 == 0
end

local count = 0
for i,v in pairs(range) do
  if v >= set then
    if IsEven(i) then
      lookup = range[i-1]
    else
      lookup = range[i]
    end
    break
  end
  count = i
  print("running through loop")
end
if lookup == -1 then 
  lookup = range[count]
end
print(lookup)
