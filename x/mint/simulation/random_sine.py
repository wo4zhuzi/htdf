import matplotlib.pyplot as plt

def plot(x=[1, 2, 3, 4],y=[1, 2, 3, 4]):
    plt.plot(x,y)
    plt.ylabel('block reward')
    plt.xlabel('block index')
    plt.show()

import os
from random import randint
import math

BlkTime          = 5
AvgDaysPerMonth  = 30
DayinSecond      = 24 * 3600
AvgBlksPerMonth  = AvgDaysPerMonth * DayinSecond / BlkTime
MonthProvisions  = 75000.0
htdf2satoshi     = float(10**8)
AvgBlkReward     = MonthProvisions / AvgBlksPerMonth
AvgBlkRewardAsSatoshi = htdf2satoshi * AvgBlkReward

def calcReward(amp=1.0,cycle=1000,step=500):
    radian = 2 * math.pi * step /cycle
    return amp * math.sin(radian) + AvgBlkReward

MAX_AMPLITUDE = AvgBlkReward
MIN_AMPLITUDE = 0.001
MAX_CYCLE     = 3000
MIN_CYCLE     = 100

def randomAmplitude(scales=100):
    return MIN_AMPLITUDE + MAX_AMPLITUDE * randint(0,100) /100.0

def randomCycle():
    return randint(MIN_CYCLE,MAX_CYCLE-MIN_CYCLE)

def testRandomSine(lastblkheight):
    totalSupply = 6*10**7
    curAmplitude, curCycle, curLastIndex = 0,0,0
    rewards = []
    for curBlkHeight in range(1,lastblkheight+MAX_CYCLE/2):
        #
        if totalSupply > 9.6*10**7: rewards.append(0);break
        #
        if curBlkHeight >= (curCycle + curLastIndex):
            curAmplitude = randomAmplitude()
            curCycle = randomCycle()
            curLastIndex = curBlkHeight
        BlockReward = calcReward(curAmplitude, curCycle, curBlkHeight-curLastIndex)
        rewards.append(BlockReward)
        totalSupply += BlockReward
    plot(range(1,curBlkHeight+1),rewards)
    return

testRandomSine(10000)