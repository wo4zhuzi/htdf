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

MAX_AMPLITUDE = AvgBlkReward
MIN_AMPLITUDE = 0.001
MAX_CYCLE = 3000
MIN_CYCLE = 100

def estimatedAccumulatedSupply(index):
    return int(index * AvgBlkRewardAsSatoshi)

def simulateRandom(lastblkindex):
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

# simulateRandom(1001)

def sine(amp=1.0,cycle=1000,step=500):
    radian = 2 * math.pi * step /cycle
    return amp * math.sin(radian) + AvgBlkReward

def accumulate(amp=1.0,cycle=1000):
    blkrewards = []
    for i in range(cycle):
        blkrewards.append(sine(amp,cycle,i))
    accumulated = sum(blkrewards)
    gap = accumulated - AvgBlkReward*cycle
    return blkrewards,accumulated,gap

import math
def simulateSine(lastblkindex=30000):
    scales = 100
    step = 0
    accumulated = 0
    rewards = []
    while step < lastblkindex:
        amp = MIN_AMPLITUDE + (MAX_AMPLITUDE - MIN_AMPLITUDE) * randint(0,scales)/scales
        cycle = randint(MIN_CYCLE,MAX_CYCLE)
        cyclerewards,_,gap = accumulate(amp,cycle)
        accumulated += gap;rewards += cyclerewards
        step += cycle
    print("err:%f"%accumulated)
    plot([i for i in range(0,step)],rewards)
    return

simulateSine()

def testSineExp():
    y = []
    scycle = 100
    bcyle = 1000
    lastidx = 2000
    for i in range(lastidx):
        radian = 2 * math.pi * i /scycle
        amp = math.exp(-(i%bcyle)/300.0)*math.cos(radian)
        y.append(amp)
    plot([i for i in range(0,lastidx)],y)

# testSineExp()

def calcCycle(idx):
    scycle = 1000
    # bcyle = 10000
    # fdecay = 3000.0
    radian = 2 * math.pi * idx /scycle
    # return int(MAX_CYCLE * (math.exp(-(idx%bcyle)/fdecay)*math.cos(-radian) + 1.0))/2
    return int(MAX_CYCLE * (math.cos(-radian) + 1.1))/2


def simulateSineExp(lastblkindex=30000):
    scales = 100
    step = 0
    accumulated = 0
    rewards = []
    cycles = []
    while step < lastblkindex:
        amp = MIN_AMPLITUDE + (MAX_AMPLITUDE - MIN_AMPLITUDE) * randint(0,scales)/scales
        cycle =calcCycle(step)
        cyclerewards,_,gap = accumulate(amp,cycle)
        accumulated += gap;rewards += cyclerewards
        step += cycle
        cycles.append(cycle)
    print("err:%f"%accumulated)
    plot([i for i in range(0,step)],rewards)
    print(cycles)

# simulateSineExp()