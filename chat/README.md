# chat 

Chat server and client.

## Download, Build and Run 

Requires that [$GOPATH is set](https://golang.org/doc/code.html#GOPATH).

#### Download source and dependencies

```sh
$ go get -u github.com/tormoder/chat/...
```

#### Build and run server

```sh
$ cd $GOPATH/src/github.com/tormoder/chat/cmd/chatserver
$ go build
$ ./chatserver
```

#### Build and run client 

```sh
$ cd $GOPATH/src/github.com/tormoder/chat/cmd/chatclient
$ go build
$ ./chatclient
```

## Usage

#### Server

```
Usage of ./chatserver:
  -port port
        The chat server port (default 10000)
  -v    show verbose debugging output
```

#### Client

```
Usage of ./chatclient:
  -saddr string
        The chat server address in the format of host:port (default "127.0.0.1:10000")
```

## Dependencies

* Serialization: [Protocol Buffers](http://github.com/golang/protobuf/)
* Communication: [gRPC](http://www.grpc.io/) 
