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
	 hscli query account $$acc1 > ${TMP_PATH};\
	 orgn1=$$(findkey amount ${TMP_PATH});\
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
	 margin=$$(python -c "print $$orgn1-$$curr1");\
	 python -c "result='passed!' if $$margin == $$GasPrice*$$gas_used else 'failed!'; print result";\
	 echo "$$orgn1 - $$curr1 == $$margin == $$gas_used * $$GasPrice?";

# non-htdfservice tx
# staking unbond testing failed
# after-unbounding balance = before-balance
UNBOND_AMOUNT = 100000
GAS_PRICE = 1.0
GENESIS_JSON_PATH = ~/.hsd/config/genesis.json
test.staking.unbond:
	@acc=$$(hscli accounts list|row 1);\
	 hscli query account $$acc > ${TMP_PATH};\
	 orgn=$$(findkey amount ${TMP_PATH});\
	 hscli query staking validators > ${TMP_PATH};\
	 tokens=$$(findkey tokens ${TMP_PATH});\
	 echo balance:$$orgn;\
	 echo tokens:$$tokens;\
	 val_addr=$$(findkey validator_address ${GENESIS_JSON_PATH});\
	 hscli query distr rewards $$acc > ${TMP_PATH};\
	 rewards=$$(findstr amount ${TMP_PATH}|fromstr ": "|extractstr '"' '"');\
	 echo rewards:$$rewards;\
	 hscli tx staking unbond $$acc $$val_addr ${UNBOND_AMOUNT}${UNIT} --gas-adjustment 1.5 --gas-prices=${GAS_PRICE}${UNIT} > ${TMP_PATH};\
	 txhash=$$(findstr txhash ${TMP_PATH}|fromstr ": "|extractstr '"' '"');\
	 echo txhash:$$txhash;\
	 sleep ${BLK_TIME};\
	 hscli query tx $$txhash > ${TMP_PATH};\
	 success=$$(findkey success ${TMP_PATH});\
	 echo success:$$success;\
	 if [ $$success = "False" ]; then echo success:$$success please press ctrl+c to exit; sleep 1000;fi;\
	 gas_wanted=$$(findkey gas_wanted ${TMP_PATH});\
	 echo gas_wanted:$$gas_wanted;\
	 gas_used=$$(findkey gas_used ${TMP_PATH});\
	 echo gas_used:$$gas_used;\
	 hscli query account $$acc > ${TMP_PATH};\
	 curr=$$(findkey amount ${TMP_PATH});\
	 margin=$$(python -c "print $$curr-$$orgn");\
	 fee=$$(python -c "print int($$gas_used*${GAS_PRICE})+1");\
	 python -c "result='passed!' if $$margin == $$rewards+${UNBOND_AMOUNT}-$$fee else 'failed!'; print result";\
	 echo $$curr - $$orgn == $$margin == $$rewards+${UNBOND_AMOUNT}-$$fee?;

	