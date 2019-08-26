CHAIN_ID = testchain

check:
        @ps -ef | grep hsd
        @ps -ef | grep hscli

clean:
        @rm -rf /root/.hsd /root/.hsd/app.log

stop-daemon:
        @pkill hsd

stop-rest:
        @pkill hscli

start-daemon:
        @nohup hsd start >> ${HOME}/.hsd/app.log  2>&1  &

start-rest:
        @nohup hscli rest-server --chain-id=${CHAIN_ID} --trust-node=true --laddr=tcp://0.0.0.0:1317 >> ${HOME}/.hsd/restServer.log  2>&1  &

stopall: stop-daemon stop-rest
restart-daemon: stop-daemon start-daemon
restart-rest: stop-rest start-rest
