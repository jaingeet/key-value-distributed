import subprocess
import time
import os

def restartServer(serverIndex):
    subprocess.Popen('./server ' + str(serverIndex) + ' -r ' , shell=True)
    print("started server " + str(serverIndex))

def processExists(serverIndex):
    output = subprocess.check_output("ps aux | grep './server " + str(serverIndex) + "'", shell=True).decode("utf-8")
    if (len(output.split('\n')) < 4):
        return False
    return True

if __name__ == "__main__":
    while 1:
        print("checking server status")
        for i in range(3):
            if not processExists(i):
                restartServer(i)

        time.sleep(10)