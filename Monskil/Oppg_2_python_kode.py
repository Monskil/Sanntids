from threading import Thread
import threading

i = 0
lock = threading.Lock()

def someThreadFunction1():
    global i
    lock.acquire()
    for j in range (0, 1000000):
        i=i+1
        #j = i
        # j += 1
        # i= j $ apt-get install git
    lock.release()

def someThreadFunction2():
    global i
    lock.acquire()
    for k in range (0, 1000002):
        i=i-1
    lock.release()
    
def main():
    
    Thread1 = Thread(target = someThreadFunction1, args = (),)
    Thread1.start()
    
    Thread2 = Thread(target = someThreadFunction2, args = (),)
    
    Thread2.start()

    Thread1.join()
    Thread2.join()
    print i

main()
