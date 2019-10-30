# MainNet
### action
    make hstart
    make hstest
    make hstop
    sudo make clean
### hstest
    [alias]
    cp hs* /usr/local/bin
    alias hscli="hscli --home node4/.hscli"
    alias hsd="hsd --home node4/.hsd"

    [config/start]
    hscli config chain-id testchain
    hscli config trust-node true
    hscli config node http://192.168.10.2:26657
    # nohup hsd start & > /dev/null
    # clear

    [transactions]
    hscli accounts new 12345678
    hscli accounts list
    hscli query account $(hscli accounts list| sed -n '1p')
    hscli tx send $(hscli accounts list| sed -n '1p') $(hscli accounts list| sed -n '2p') 20000000stake --gas-prices=0.00001stake
    hscli tx send $(hscli accounts list| sed -n '1p') $(hscli accounts list| sed -n '2p') 20000000satoshi --gas-prices=0.00001satoshi

    [validators]
    - show yours
    hsd tendermint show-validator
    - check status
    hscli query staking validators
    hscli  query staking validator [cosmosvaloper]
    - confirm running
    hscli query tendermint-validator-set
    hscli query tendermint-validator-set | grep [cosmosvalcons/cosmosvalconspub]
    hscli query tendermint-validator-set | grep "$(hsd tendermint show-validator)"
    - start yours(tip: 1,000,000 for voting power 1, 10,000,000 for 10, 100,000,000 for 100)
    hscli tx staking create-validator $(hscli accounts list| sed -n '2p') \
                                       --pubkey=$(hsd tendermint show-validator)\
                                       --amount=10000000stake \
                                       --moniker=client \
                                       --commission-rate=0.10 \
                                       --commission-max-rate=0.20 \
                                       --commission-max-change-rate=0.01 \
                                       --min-self-delegation=1 \
                                       --gas-prices=0.00001stake
    hscli tx staking edit-validator $(hscli accounts list| sed -n '2p') --gas-prices=0.00001stake
    or
    hscli tx staking edit-validator $(hscli accounts list| sed -n '2p')\
                --moniker=client \
                --chain-id=testchain \
                --website="https://cosmos.network" \
                --identity=23870f5bb12ba2c4967c46db \
                --details="To infinity and beyond!" \
                --gas-prices=0.00001stake \
                --commission-rate=0.10
    - unjail
    hscli tx slashing unjail $(hscli accounts list| sed -n '2p') --gas-prices=0.00001stake
    - log
    hscli query slashing signing-info [cosmosvalconspub]
    - check

    [delegators]
    hscli tx staking delegate $(hscli accounts list| sed -n '2p') \
                               $(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '2p'|awk '{print $3F}' | cut -d'"' -f 2)\
                               10000000stake --gas-adjustment=1.5 --gas-prices=0.00001stake
    hscli tx staking redelegate $(hscli accounts list| sed -n '2p') \
                                 $(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '2p'|awk '{print $3F}' | cut -d'"' -f 2) \
                                 $(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '3p'|awk '{print $3F}' | cut -d'"' -f 2) \
                                 10000000stake --gas-adjustment=1.5 --gas-prices=0.00001stake
    hscli query distr rewards $(hscli accounts list| sed -n '2p')
    hscli query account $(hscli accounts list| sed -n '2p')
    hscli tx staking unbond $(hscli accounts list| sed -n '1p') \
                             $(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '1p'|awk '{print $3F}' | cut -d'"' -f 2) \
                             1000000satoshi \
                             --gas-adjustment 1.5 --gas-prices=0.00001satoshi
    hscli tx distr withdraw-rewards $(hscli accounts list| sed -n '2p') \
                                     $(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '2p'|awk '{print $3F}' | cut -d'"' -f 2) \ --gas-adjustment 1.5 --gas-prices=0.00001stake