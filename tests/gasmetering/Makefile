# written by junying
# 2019.10.28
# index
# 1.htdfservice tx
# - non-contract tx
# - contract tx
# 2.non-contract tx

TMP_PATH = /tmp/result.json
BLK_TIME = 5
# htdfservice tx
## non-contract tx
### test cases
#### gas < DefaultTxGas(21000)
#### gasprice < 1satoshi
TX_AMOUNT = 1000000
UNIT = satoshi

test.normal:
	@acc1=$$(hscli accounts list|row 1);\
	 acc2=$$(hscli accounts list|row 2);\
	 hscli query account $$acc1 > ${TMP_PATH};\
	 orgn1=$$(findkey amount ${TMP_PATH});\
	 hscli query account $$acc2 > ${TMP_PATH};\
	 orgn2=$$(findkey amount ${TMP_PATH});\
	 hscli tx send $$acc1 $$acc2 ${TX_AMOUNT}${UNIT} > ${TMP_PATH};\
	 txhash=$$(findstr txhash ${TMP_PATH}|fromstr ": "|extractstr '"' '"');\
	 sleep ${BLK_TIME};\
	 hscli query tx $$txhash > ${TMP_PATH};\
	 success=$$(findkey success ${TMP_PATH});\
	 gas_wanted=$$(findkey gas_wanted ${TMP_PATH});\
	 gas_used=$$(findkey gas_used ${TMP_PATH});\
	 GasPrice=$$(findkey GasPrice ${TMP_PATH});\
	 Gas=$$(findkey Gas ${TMP_PATH});\
	 hscli query account $$acc1 > ${TMP_PATH};\
	 curr1=$$(findkey amount ${TMP_PATH});\
	 hscli query account $$acc2 > ${TMP_PATH};\
	 curr2=$$(findkey amount ${TMP_PATH});\
	 margin=$$(python -c "print $$orgn1-$$curr1");\
	 python -c "result='passed!' if $$margin == ($$GasPrice*$$gas_used+${TX_AMOUNT}) else 'failed!'; print result";\
	 echo "$$orgn1 - $$curr1 == $$margin == $$gas_used * $$GasPrice + ${TX_AMOUNT}";


## contract tx
### test cases
#### 1. gas > gasUsed
#### 2. Intrinsic Gas < gas < gasUsed
#### 3. gas < Intrinsic Gas
test.contract.creation:
	@acc1=$$(hscli accounts list|row 1);\
	 acc2=$$(hscli accounts list|row 2);\
	 hscli query account $$acc1 > ${TMP_PATH};\
	 orgn1=$$(findkey amount ${TMP_PATH});\
	 hscli query account $$acc2 > ${TMP_PATH};\
	 orgn2=$$(findkey amount ${TMP_PATH});\
	 echo $$(cat ../evm/coin/coin_sol_Coin.bin);\
	 $(MAKE) -sC ../evm create.contract.default;\
	 txhash=$$(findkey txhash ${TMP_PATH});\
	 success=$$(findkey success ${TMP_PATH});\
	 gas_wanted=$$(findkey gas_wanted ${TMP_PATH});\
	 gas_used=$$(findkey gas_used ${TMP_PATH});\
	 GasPrice=$$(findkey GasPrice ${TMP_PATH});\
	 Gas=$$(findkey Gas ${TMP_PATH});\
	 hscli query account $$acc1 > ${TMP_PATH};\
	 curr1=$$(findkey amount ${TMP_PATH});\
	 hscli query account $$acc2 > ${TMP_PATH};\
	 curr2=$$(findkey amount ${TMP_PATH});\
	 margin=$$(python -c "print $$orgn1-$$curr1");\
	 python -c "result='passed!' if $$margin == $$GasPrice*$$gas_used else 'failed!'; print result";\
	 echo "$$orgn1 - $$curr1 == $$margin == $$gas_used * $$GasPrice?";

# non-htdfservice tx
# staking unbond testing failed
# after-unbounding balance = before-balance
GAS_PRICE = 0.00001

test.staking:
	@acc=$$(hscli accounts list|row 1);\
	 hscli query account $$acc1 > ${TMP_PATH};\
	 orgn=$$(findkey amount ${TMP_PATH});\
	 hscli query staking validators > ${TMP_PATH};\
	 tokens=$$(findkey tokens ${TMP_PATH});\
	 val_addr=$$(grep -nr validator_address  ~/.hsd/config/genesis.json |sed -n '1p'|awk '{print $3F}' | cut -d'"' -f 2);\
	 rewards=$$(hscli query distr rewards $$acc);\
	 hscli tx staking unbond $$acc $$val_addr ${TX_AMOUNT}${UNIT} --gas-adjustment 1.5 --gas-prices=${GAS_PRICE}${UNIT} > ${TMP_PATH};\
	 txhash=$$(findstr txhash ${TMP_PATH}|fromstr ": "|extractstr '"' '"');\
	 sleep ${BLK_TIME};\
	 hscli query tx $$txhash > ${TMP_PATH};\
	 success=$$(findkey success ${TMP_PATH});\
	 gas_wanted=$$(findkey gas_wanted ${TMP_PATH});\
	 gas_used=$$(findkey gas_used ${TMP_PATH});\
	 hscli query account $$acc > ${TMP_PATH};\
	 curr=$$(findkey amount ${TMP_PATH});\
	 margin=$$(python -c "print $$orgn-$$curr");\
	 python -c "result='passed!' if $$margin == $$GasPrice*$$gas_used else 'failed!'; print result";\