import time

start = time.time()  # current time in seconds

x = 0

while x < 10000000:
    x += 1

end = time.time()

print(end - start)
