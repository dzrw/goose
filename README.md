goose
=====

## Installation

```bash
$ go get "github.com/politician/goose"
$ go install
```

## Usage

```bash
$ goose
```

Starts up the various services, and waits for connections.

### Create a watch

```bash
$ curl -X POST --data "{\"tag\":\"opaque-id-1234\",\"dataSourceName\":\"data-access-service\",\"matchExpr\":{\"method\":\"GET\",\"path\":\"/foo\"},\"echo\":{\"status\":200,\"headers\":{},\"body\":\"hello, world\"}}" http://127.0.0.1:8080/watches
{"id": 1}
```

### Query a watched endpoint

```bash
$ curl http://127.0.0.1:8081/foo
hello, world
```

## Redis Integration

If you're running Redis, goose will store any incoming requests as hashes and `LPUSH` the hash keys onto the list specified as the `tag` field in a watch.

### Example

In the following example, we'll start up goose as a background job, set up a watch such that requests will be saved in a redis list named *exampletag*, curl the watched endpoint, then verify that redis contains the request using telnet.

```bash
$ goose -redis localhost:6379 &
$ curl -X POST --data "{\"tag\":\"exampletag\",\"dataSourceName\":\"data-access-service\",\"matchExpr\":{\"method\":\"GET\",\"path\":\"/foo\"},\"echo\":{\"status\":200,\"headers\":{},\"body\":\"hello, world\"}}" http://127.0.0.1:8080/watches
{"id": 1}
$ curl -X GET --data "stuff & things" http://127.0.0.1:8081/foo
hello, world
$ telnet localhost 6379
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
exists exampletag
:1
lrange exampletag 0 1
*1
$39
goose:requests:Uizt47II+ERfNZgucG+RVQ==
hgetall goose:requests:Uizt47II+ERfNZgucG+RVQ==
*6
$6
method
$3
GET
$3
url
$4
/foo
$4
body
$20
c3R1ZmYgJiB0aGluZ3M=
quit
+OK
Connection closed by foreign host.
$
```

