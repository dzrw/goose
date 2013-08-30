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
