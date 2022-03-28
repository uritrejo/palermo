# Palermo




TODO: link to openAPI, add a little diagram

## Build & Run

``` shell
git clone https://github.com/uritrejo/palermo.git
cd palermo
./scripts/build.sh  # will run go build, go test, create binaries
./bin/palermo
>> time="2022-03-27T15:47:12-04:00" level=info msg="Palermo server is listening on localhost:4422"
```

todo: distinguish usage of binary to usage of api
### Usage
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

## Commands
add api routes

example usages