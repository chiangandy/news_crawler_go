#!/bin/bash

go build -o build/ettoday exec/ettoday/spider.go 

go build -o build/setn exec/setn/spider.go 

go build -o build/mirrormedia exec/mirrormedia/spider.go 

go build -o build/chinatimes exec/chinatimes/spider.go 

go build -o build/tvbs exec/tvbs/spider.go 

go build -o build/push_info exec/push_info/main.go 
