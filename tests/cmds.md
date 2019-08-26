### account
    hscli accounts new 12345678
### query
    hscli query account $(hscli accounts list| sed -n '2p')
### bank
    hscli tx send $(hscli accounts list| sed -n '1p') $(hscli accounts list| sed -n '2p') 20000000stake --gas-prices=0.00001stake
    hscli tx create $(hscli accounts list| sed -n '1p') $(hscli accounts list| sed -n '2p') 1000stake --gas-prices=0.00001stake --encode=false
    hscli tx sign [rawcode] --encode=false
### staking
    [validating]
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

    [delegating]
    hscli tx staking delegate $(hscli accounts list| sed -n '2p') <validatorAddress> 10000000stake --gas-adjustment=1.5 --gas-prices=0.00001stake