### cmds
    [min-gas-prices]
    hsd start --minimum-gas-prices=0.1htdf

    [fee & gasprice]
    hscli tx send [fromaddr] [toaddr] [amount] --gas-prices=0.005htdf
    hscli tx send [fromaddr] [toaddr] [amount] --gas=1000 --gas-prices=0.005htdf
    hscli tx send [fromaddr] [toaddr] [amount] --fees=1htdf
### references
#### [client cmd: fee & gas](https://cosmos.network/docs/gaia/gaiacli.html#fees-gas)
#### [gaiad.toml: minimum-gas-prices](https://cosmos.network/docs/gaia/join-mainnet.html#set-minimum-gas-prices)