# MainNet
### action
    make hstart
    make hstest
    make hstop
    sudo make clean
### hstest
    [alias]
    cp hs* /usr/local/bin
    alias ahscli="hscli --home node4/.hscli"
    alias ahsd="hsd --home node4/.hsd"

    [config/start]
    ahscli config chain-id testchain
    ahscli config trust-node true
    ahscli config node http://192.168.10.2:26657
    # nohup hsd start & > /dev/null
    # clear

    [transactions]
    ahscli accounts new 12345678
    ahscli accounts list
    ahscli query account $(ahscli accounts list| sed -n '1p')
    ahscli tx send $(ahscli accounts list| sed -n '1p') $(ahscli accounts list| sed -n '2p') 20000000stake --gas-prices=0.00001stake
    ahscli tx send $(ahscli accounts list| sed -n '1p') $(ahscli accounts list| sed -n '2p') 20000000satoshi --gas-prices=0.00001satoshi

    [validators]
    - show yours
    ahsd tendermint show-validator
    - check status
    ahscli query staking validators
    ahscli  query staking validator [cosmosvaloper]
    - confirm running
    ahscli query tendermint-validator-set
    ahscli query tendermint-validator-set | grep [cosmosvalcons/cosmosvalconspub]
    ahscli query tendermint-validator-set | grep "$(hsd tendermint show-validator)"
    - start yours(tip: 1,000,000 for voting power 1, 10,000,000 for 10, 100,000,000 for 100)
    ahscli tx staking create-validator $(ahscli accounts list| sed -n '2p') \
                                       --pubkey=$(hsd tendermint show-validator)\
                                       --amount=10000000stake \
                                       --moniker=client \
                                       --commission-rate=0.10 \
                                       --commission-max-rate=0.20 \
                                       --commission-max-change-rate=0.01 \
                                       --min-self-delegation=1 \
                                       --gas-prices=0.00001stake
    ahscli tx staking edit-validator $(ahscli accounts list| sed -n '2p') --gas-prices=0.00001stake
    or
    ahscli tx staking edit-validator $(ahscli accounts list| sed -n '2p')\
                --moniker=client \
                --chain-id=testchain \
                --website="https://cosmos.network" \
                --identity=23870f5bb12ba2c4967c46db \
                --details="To infinity and beyond!" \
                --gas-prices=0.00001stake \
                --commission-rate=0.10
    - unjail
    ahscli tx slashing unjail $(ahscli accounts list| sed -n '2p') --gas-prices=0.00001stake
    - log
    ahscli query slashing signing-info [cosmosvalconspub]
    - check

    [delegators]
    ahscli tx staking delegate $(ahscli accounts list| sed -n '2p') \
                               $(grep -nr validator_address  node0/.hsd/config/genesis.json |sed -n '2p'|awk '{print $3F}' | cut -d'"' -f 2)\
                               10000000stake --gas-adjustment=1.5 --gas-prices=0.00001stake
    ahscli tx staking redelegate $(ahscli accounts list| sed -n '2p') \
                                 $(grep -nr validator_address  node0/.hsd/config/genesis.json |sed -n '2p'|awk '{print $3F}' | cut -d'"' -f 2) \
                                 $(grep -nr validator_address  node0/.hsd/config/genesis.json |sed -n '3p'|awk '{print $3F}' | cut -d'"' -f 2) \
                                 10000000stake --gas-adjustment=1.5 --gas-prices=0.00001stake
    ahscli query distr rewards $(ahscli accounts list| sed -n '2p')
    ahscli query account $(ahscli accounts list| sed -n '2p')
    ahscli tx staking unbond $(ahscli accounts list| sed -n '2p') \
                             $(grep -nr validator_address  node0/.hsd/config/genesis.json |sed -n '2p'|awk '{print $3F}' | cut -d'"' -f 2) \
                             1000000stake
                             --gas-adjustment 1.5 --gas-prices=0.00001stake
    ahscli tx distr withdraw-rewards $(ahscli accounts list| sed -n '2p') \
                                     $(grep -nr validator_address  node0/.hsd/config/genesis.json |sed -n '2p'|awk '{print $3F}' | cut -d'"' -f 2) \ --gas-adjustment 1.5 --gas-prices=0.00001stake