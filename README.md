MTA
===
Scalable, multi region, check system.

# Design
MTA has 3 distinct parts:
- scheduler
- worker
- consumer

All communication between this parts are made with kafka.
Worker and Consumer can be scale easily.

## Scheduler
 Scheduler send *CheckRequest* in kafka topic.
 
### backend
- db : get checks from database (PG)

## Worker
Worker get *CheckRequest* from kafka, run the
 check and send *CheckResponse* in kafka topic.
 
### Available Check

#### Ping
Just send echo request, and store RTT

#### Http
Do an http request and store different value :
- DNSLookup
- TCPConnection
- TLSHandshake
- ServerProcessing
- NameLookup
- Connect
- Pretransfer
- StartTransfer
- Total
 
## Consumer
Finally consumer get *CheckResponse* and do something ..
 
### backend
- db : Insert result in database (PG)
- stdout : Display result
- influx : Store resulte in influxdb

## What MTA means

Honestly, I really don't remeber , i'm so tired when I start this project 