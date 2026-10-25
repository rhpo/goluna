 # I made it to compare between the speed of Luna and python knowing that they're both interpreted

import time

start = time.time() 

x = 0

while x < 10000000:
    x += 1

end = time.time()

print(end - start)
