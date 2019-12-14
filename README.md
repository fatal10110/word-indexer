# Word Indexer
Web app that index words from different sources, optimized for large inputs, it splits the data into batches and runs the index on each batch async

![go-docker](https://hackernoon.com/hn-images/1*JfSp7LWmVE1nj15IrxWSWQ.png)

### Supported sources
- Text - Http body
- URL - Web page content
- File - File content

# HowTo

1. Run `./run.sh` from project source
2. Http body example

```
curl -d 'lorem ipsum' -X POST http://localhost:8080/index?source=text
```
```
curl -X GET http://localhost:8080/index/lorem
```
3. URL example

```
curl -d 'http://www.lipsum.com' -X POST http://localhost:8080/index?source=url
```
```
curl -X GET http://localhost:8080/index/lorem
```
4. File example

```
curl -d '/tmp/ipsum.txt' -X POST http://localhost:8080/index?source=file
```
```
curl -X GET http://localhost:8080/index/lorem
```

# Components

**src/main.go**

App entrypoint, creats the store and broker and also runs the web server

**Broker**

There is a [simple interface](https://github.com/fatal10110/word-indexer/blob/master/src/broker.go#L12) that should be implemented to support other brokers, currently it supports only broker based on Go channels(in memory)
it also may be implemented on Redis(e.g using sets) or standard pub-sub.

**Store**

There is also a [store inteface](https://github.com/fatal10110/word-indexer/blob/master/src/store.go#L8) that should be implemented by statistic holder
currently it uses in-memory Go Map + RWMutex.

**Batching Job**

Async job that takes the byte slice input, and splits it into batches, it is needed to optimize the index time and use all cores.

**Indexer Job**

Async job, takes the batch slice and counts the number of word entries in it.

**Report Statistic Job**

Async job that takes the index result and adds it to the store

**API**

Supports two requests 
1. `POST /index` - Accepts also required qury param `source`, could not pass it as body since body may be to large to parse it.
2. `GET /index/:word` - Gets the stat value for specific word
