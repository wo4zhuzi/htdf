    delegator-addr: htdf1zf07fyt2an2ral8zve0u4y7lzqa6x4lqfeyl8m
    validator-addr: htdfvaloper1zf07fyt2an2ral8zve0u4y7lzqa6x4lqrquxss
    amount: 10000000stake
    samount: 100satoshi
### 设置chain-id
    hscli config chain-id [chain-id]
### 如果代表地址里面没钱的话，转账一笔
    hscli query accounts [delegator-addr]
    hscli tx send [fromaddr] [delegator-addr] [samount] --gas-prices=0.00001satoshi
### [stake-抵押，解绑](https://github.com/orientwalt/htdf/blob/master/x/staking/client/cli/tx.go)   
    [抵押]
    hscli tx staking delegate [delegator-addr] [validator-addr] [amount] --gas-adjustment=1.5 --gas-prices=0.00001satoshi
    
    [解绑]
    hscli tx staking unbond [delegator-addr] [validator-addr] [amount] --gas-adjustment 1.5 --gas-prices=0.00001satoshi
### [奖励-查询，回收](https://github.com/orientwalt/htdf/blob/master/x/distribution/client/cli/tx.go)
    [查询]
    hscli query distr rewards [delegator-addr]

    [回收]
    hscli tx distr withdraw-rewards [delegator-addr] [validator-addr] --gas-adjustment 1.5 --gas-prices=0.00001satoshi
