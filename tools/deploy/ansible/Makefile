# junying-todo, 2020-02-10
CONFIG_DIR=/home/ubuntu/config/nosentry/htdf/observation

all: init copy run

init:
	@ansible observers -m shell -a "if ! [ -d ' /root/.hsd/config' ];then mkdir -p /root/.hsd/config; fi"

copy: cp.config cp.hsd cp.hscli cp.Makefile
cp.config:
	@ansible observers -m copy -a "src=/home/ubuntu/config/nosentry/htdf/observation/config.toml dest=/root/.hsd/config force=yes"
	@ansible observers -m copy -a "src=/home/ubuntu/config/nosentry/htdf/observation/genesis.json dest=/root/.hsd/config force=yes"
cp.hsd:	
	@ansible observers -m copy -a "src=/home/ubuntu/go/bin/hsd dest=/usr/local/bin force=yes mode=755"
cp.hscli:	
	@ansible observers -m copy -a "src=/home/ubuntu/go/bin/hscli dest=/usr/local/bin force=yes mode=755"
cp.Makefile:
	@ansible observers -m copy -a "src=/home/ubuntu/config/nosentry/htdf/observation/Makefile dest=/root force=yes"	

run:
	@ansible observers -m shell -a "make start-daemon start-rest"

