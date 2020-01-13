DELEGATOR_CNTS = $$(findstr delegator_address ~/.hsd/config/genesis.json | linecount)

initialize:
	@for i in $$(python -c "print ' '.join(str(item) for item in range(${VALIDATOR_CNTS}))"); \
	 do \
	 validator=$$(findstr validator_address ~/.hsd/config/genesis.json | sed -n '${i}p' | fromstr ': "'|endstr '"');\
	 if [ -d "$$validator" ]; then rm -rf $$validator; mkdir $$validator; fi;\
	 delegator=$$(findstr delegator_address ~/.hsd/config/genesis.json | sed -n '${i}p' | fromstr ': "'|endstr '"');\
	 echo $$delegator >> $$validator/delegators.list;\
	 done

VALIDATOR_CNTS = $$(ls|grep htdf|linecount)

validator.rewards:
	@rewards=$$(echo Vrewards:);\
	 for i in $$(python -c "print ' '.join(str(item) for item in range(${VALIDATOR_CNTS}))"); \
	 do \
	 validator=$$(ls|grep htdf|sed -n '${i}p');\
	 echo $$validator;\
	 reward=$$(hscli query distr commission $$validator  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 rewards=$$(echo $$rewards $$reward);\
	 done;\
	 echo $$rewards;

delegator.rewards:
	@rewards=$$(echo Drewards:);\
	 for i in $$(python -c "print ' '.join(str(item) for item in range(${VALIDATOR_CNTS}))"); \
	 do \
	 validator=$$(ls|grep htdf|sed -n '${i}p');\
	 echo ------------------------------------------------------;\
	 echo V$$i:$$validator;\
	 delcnts=$$(linecount $$validator/delegators.list);\
	 for j in $$(python -c "print ' '.join(str(item) for item in range($$delcnts))"); \
	 do \
	 delegator=$$(cat $$validator/delegators.list | sed -n '${j}p');\
	 echo D$$j:$$delegator;\
	 reward=$$(hscli query distr rewards $$delegator  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 rewards=$$(echo $$rewards $$reward);\
	 done;\
	 echo $$rewards;\
	 done;\
	 
community.pool:
	@reward=$$(hscli query distr community-pool  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$reward;

check.rewards:
	@read -p "accaddr: " accaddr; \
	 rewards=$$(hscli query distr rewards $$accaddr  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$rewards

height:
	@height=$$(hscli query block | sed -n "16p"| fromstr ': "'|endstr '"');\
	 echo $$height

total_txs:
	@total_txs=$$(hscli query block | sed -n "19p"| fromstr ': "'|endstr '"');\
	 echo $$total_txs

num_txs:
	@num_txs=$$(hscli query block | sed -n "18p"| fromstr ': "'|endstr '"');\
	 echo $$num_txs

print.rewards:
	@r_old=0;cp_old=0;cm_old=0;num_txs_old=0;total_txs_old=0;\
	 for i in $$(python -c "print ' '.join(str(item) for item in range(100))"); \
	 do \
	 r=$$($(MAKE) -s acc.rewards);\
	 cp=$$($(MAKE) -s community.pool);\
	 cm=$$($(MAKE) -s commision);\
	 d1=$$(python -c "print $$r-$$r_old");\
	 d2=$$(python -c "print $$cp-$$cp_old");\
	 d3=$$(python -c "print $$cm-$$cm_old");\
	 r_old=$$r;cp_old=$$cp;cm_old=$$cm;\
	 echo $$r $$cp $$cm;\
	 mr=$$(python -c "print $$d1+$$d2+$$d3");\
	 height=$$($(MAKE) -s height);\
	 total_txs=$$($(MAKE) -s total_txs);\
	 num_txs=$$($(MAKE) -s num_txs);\
	 echo delta:$$height:$$num_txs:$$total_txs--- $$d1 $$d2 $$d3 $$mr $$(python -c "print $$mr-$$num_txs_old*3000000");\
	 num_txs_old=$$num_txs;total_txs_old=$$total_txs;\
	 echo ratio:$$height:$$num_txs:$$total_txs--- $$(python -c "print '%.2f'%($$d1/$$mr)") $$(python -c "print '%.2f'%($$d2/$$mr)") $$(python -c "print '%.2f'%($$d3/$$mr)");\
	 sleep 4;\
	 done\
