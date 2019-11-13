module github.com/orientwalt/htdf

go 1.12

require (
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/aristanetworks/goarista v0.0.0-20191023202215-f096da5361bb
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/bgentry/speakeasy v0.1.0
	github.com/btcsuite/btcd v0.20.0-beta
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/deckarep/golang-set v1.7.1
	github.com/emicklei/proto v1.8.0
	github.com/ethereum/go-ethereum v1.8.27
	github.com/fortytw2/leaktest v1.3.0 // indirect
	github.com/go-kit/kit v0.8.0
	github.com/gogo/protobuf v1.2.1
	github.com/golang/protobuf v1.3.2
	github.com/hashicorp/golang-lru v0.5.3 // indirect
	github.com/mattn/go-isatty v0.0.10
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/pelletier/go-toml v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.1.0
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190826022208-cac0b30c2563 // indirect
	github.com/rjeczalik/notify v0.9.2
	github.com/rs/cors v1.7.0 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	github.com/tendermint/btcd v0.1.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/iavl v0.12.4
	github.com/tendermint/tendermint v0.32.1
	github.com/tendermint/tmlibs v0.9.0
	golang.org/x/crypto v0.0.0-20190404164418-38d8ce5564a5
	gopkg.in/fatih/set.v0 v0.2.1 // indirect
	gopkg.in/karalabe/cookiejar.v2 v2.0.0-20150724131613-8dcd6a7f4951 // indirect
)

replace (
	github.com/tendermint/iavl v0.12.4 => github.com/orientwalt/iavl v0.12.4
	github.com/tendermint/tendermint v0.31.5 => github.com/orientwalt/tendermint v90.90.90+incompatible
	golang.org/x/crypto => github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
)
