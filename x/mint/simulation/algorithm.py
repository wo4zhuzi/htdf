import matplotlib.pyplot as plt

def plot(x=[1, 2, 3, 4],y=[1, 2, 3, 4]):
    plt.plot(x,y)
    plt.ylabel('block reward')
    plt.xlabel('block index')
    plt.show()

import os
from random import randint

BlkTime          = 5
AvgDaysPerMonth  = 30
DayinSecond      = 24 * 3600
AvgBlksPerMonth  = AvgDaysPerMonth * DayinSecond / BlkTime
MonthProvisions  = 75000.0
htdf2satoshi     = 100000000.0
AvgBlkReward     = MonthProvisions / AvgBlksPerMonth
AvgBlkRewardAsSatoshi = htdf2satoshi * AvgBlkReward

RATIO = 0.5

MAX_AMPLITUDE = 1.2
MIN_AMPLITUDE = 0.01
MAX_FREQUENCY = 
MIN_FREQUENCY = 

def estimatedAccumulatedSupply(index):
    return int(index * AvgBlkRewardAsSatoshi)

def simulateSimple(lastblkindex):
    actual = 0
    blkrewards = []
    for blkindex in range(1,lastblkindex):
        expected = estimatedAccumulatedSupply(blkindex)
        estimated = expected - actual if expected > actual else 0
        real = randint(0,int(estimated*RATIO)+estimated)
        actual += real
        blkrewards.append(real/htdf2satoshi)
    print("err=%f"%((estimatedAccumulatedSupply(lastblkindex)-actual)/htdf2satoshi))
    print("min(reward)=%f"%min(blkrewards))
    print("max(reward)=%f"%max(blkrewards))
    plot([i for i in range(1,lastblkindex)],blkrewards)
    return abs(actual-estimatedAccumulatedSupply(lastblkindex))

# simulateSimple(1001)
import math
def simulateSine():
    print(math.sin(math.pi))
    return

simulateSine()