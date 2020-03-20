[![CircleCI](https://circleci.com/gh/orientwalt/htdf/tree/master.svg?style=shield)](https://circleci.com/gh/orientwalt/htdf/tree/master)
[![](https://godoc.org/github.com/orientwalt/htdf?status.svg)](http://godoc.org/github.com/orientwalt/htdf) [![Go Report Card](https://goreportcard.com/badge/github.com/orientwalt/htdf)](https://goreportcard.com/report/github.com/orientwalt/htdf)
[![Travis](https://travis-ci.org/orientwalt/htdf.svg?branch=master)](https://travis-ci.org/orientwalt/htdf)
[![version](https://img.shields.io/github/tag/orientwalt/htdf.svg)](https://github.com/orientwalt/htdf/releases/latest)
[![](https://tokei.rs/b1/github/orientwalt/htdf?category=lines)](https://github.com/orientwalt/htdf)
[![License](https://img.shields.io/badge/License-Apache%202.0-green.svg)](https://opensource.org/licenses/Apache-2.0)

# HTDF
## Introduction
   HTDF is a cosmos-sdk application that provides fundamental crytocurrency functions including account management, transaction processing, and smart contract. It still uses BPOS of tendermint as its consensus algorithm. This project is now UNDER ACTIVE DEVELOPMENT.
   
   **Note**: Requires Go 12.9+
## Features
  * [x] **account**: ethereum-style
  * [x] **transaction**: 
    * [x]  cold-wallet functions(create, sign, broadcast)
    * [x]  fee & reward system
  * [x] **rest**: auth/query rest removal - tx/sign, encode, broadcast
  * [ ] **block**: non-empty block
  * [x] **docker**: standalone
  * [x] **docker-compose**: multi-node
  * [ ] **emergency system**
    * [ ] monitoring system
    * [ ] alert system
    * [x] urgent response system
      * [ ] hard fork: export-based function disabled
    * [ ] validator abnormality detection
  * [x] **security**
    * [x] sentry node architecture
    * [x] dynamic system issuer
  * [x] **validators & delegators**
  * [x] **guardian**
  * [x] **upgrade**
## Executables
```
hsd
hscli
```
## [Quick Start](https://github.com/orientwalt/htdf/blob/master/docs/build%20%26%20run.md)
Only one command is enough to set up a standalone blockchain on your local machine.
```
make new
```
