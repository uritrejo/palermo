# Palermo

RESTful API for managing messages. It stores and provides details about these messages, specifically whether a message is a palindrome

## Build & Run

``` shell
git clone https://github.com/uritrejo/palermo.git
cd palermo
./scripts/build.sh  # will run go build, go test, create binaries
./bin/palermo
>> time="2022-03-27T15:47:12-04:00" level=info msg="Palermo server is listening on localhost:4422"
```

### Options
Run `./bin/palermo -h` to see the flags available:
```shell
Usage of ./bin/palermo:
  -dbtype string
        -dbtype=<type>: types are 'basic' (local memory) and 'mongodb (default "basic")
  -loglevel string
        -loglevel=<level>: levels are info, debug, trace (default "debug")
  -mongodbport int
        -mongodbport=<port>: port where mongo db is listening (default 27017)
  -port int
        -port=<port>: port on which to listen and serve (default 4422)
```

## Architecture

![](docs/palermo-architecture-diagram.png)

## Paths
Refer to [api/APIdoc.html](api/APIdoc.html) for details on the API paths.
Summary & Examples with curl:
- /v1/createMsg POST
    - `curl -X POST localhost:4422/createMsg -H "Content-Type: application/json" -d '{"id":"1", "content":"kayak"}'`
- /v1/retrieveMsg/{id} GET
    - `curl localhost:4422/retrieveMsg/1`
- /v1/retrieveAllMsgs GET
    - `curl localhost:4422/retrieveAllMsgs`
- /v1/updateMsg/{id} POST
    - `curl -X POST localhost:4422/updateMsg/1 -H "Content-Type: application/json" -d '{"id":"1", "content":"canoe"}'`
- /v1/deleteMsg/{id} GET
    - `curl localhost:4422/deleteMsg/1`
    