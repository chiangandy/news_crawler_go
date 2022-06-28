# times_crawler_go project

Repository is contains Taiwan Hot-Tops news media data crawler. 
使用GO Colly 抓取台灣頭條新聞媒體即時頭條新聞


## Function description
Using golang with gocolly to crawling data from internet and save to monogDB. The media is include TVBS, Mirrormedia, Chinatimes, STEN, Ettoday
	
## version control
	V 1.5 - release in 2021/3/1
	V 1.8 - release in 2021/7/30


## System module using list:
	github.com/PuerkitoBio/goquery v1.8.0 
	github.com/antchfx/htmlquery v1.2.4
	github.com/antchfx/xmlquery v1.3.10 
	github.com/gobwas/glob v0.2.3 
	github.com/gocolly/colly v1.2.0 
	github.com/kennygrant/sanitize v1.2.4 
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca 
	github.com/sirupsen/logrus v1.8.1 
	github.com/temoto/robotstxt v1.1.2 
	golang.org/x/net v0.0.0-20220425223048-2871e0cb64e4 
	google.golang.org/appengine v1.6.7 
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 


## Directory layout and location

* build - the final execution file 
* crawler - data crawlering program
* exec - execute entry main program
* go.mod - module manage files
* build_program.sh - build go program to be executtion file


## Using Database

- Using Mongodb V 4.0.3


## Configuration and Build

- it need to setup envvironment varibable enable program can connect to mongoDB

```
	export MONGODB_HOST="localhost"
	export MONGODB_PORT="27017"
	export MONGODB_USERNAME="admin"
	export MONGODB_PASSWD="password"
	export MONGODB_DB="go_crawler_data"
	
	./build_program.sh
```

## License
License is opensource. Please refer to Apache License 2.0（Apache 2.0）


## Project status
Project is done by first version, but it will keep developing.

