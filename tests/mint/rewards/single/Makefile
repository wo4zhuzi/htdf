
acc.rewards:
	@acc=$(shell hscli accounts list| sed -n '1p');\
	 reward=$$(hscli query distr rewards $${acc}  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$reward;

second.rewards:
	@acc=$(shell hscli accounts list| sed -n '2p');\
	 reward=$$(hscli query distr rewards $${acc}  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$reward;

community.pool:
	@reward=$$(hscli query distr community-pool  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$reward;

commision:
	@validator_address=$$(findstr validator_address ~/.hsd/config/genesis.json  | fromstr ': "'|endstr '"');\
	 reward=$$(hscli query distr commission $$validator_address  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$reward;


rewards:
	@read -p "accaddr: " accaddr; \
	 rewards=$$(hscli query distr rewards $$accaddr  | sed -n '4p' | fromstr ': "'|endstr '"');\
	 echo $$rewards

print.rewards:
	@r_old=0;cp_old=0;cm_old=0;\
	 for i in 1 2 3 4 5; \
	 do \
	 r=$$($(MAKE) -s acc.rewards);\
	 cp=$$($(MAKE) -s community.pool);\
	 cm=$$($(MAKE) -s commision);\
	 d1=$$(python -c "print $$r-$$r_old");\
	 d2=$$(python -c "print $$cp-$$cp_old");\
	 d3=$$(python -c "print $$cm-$$cm_old");\
	 r_old=$$r;cp_old=$$cp;cm_old=$$cm;\
	 echo $$r1 $$r2 $$cp $$cm;\
	 mr=$$(python -c "print $$d1+$$d2+$$d3");\
	 echo delta--- $$d1 $$d2 $$d3 $$mr;\
	 echo ratio--- $$(python -c "print '%.2f'%($$d1/$$mr)") $$(python -c "print '%.2f'%($$d2/$$mr)") $$(python -c "print '%.2f'%($$d3/$$mr)");\
	 sleep 5;\
	 done\
