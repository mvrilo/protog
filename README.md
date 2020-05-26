# protog

Protobuf file generator for the command line.

## Installation

```
go get github.com/mvrilo/protog/cmd/protog
```

## Usage

```
protog <name> [-dhofomsp] [-m MessageName[field:type,field:type,...]] [-s ServiceName[MethodName:In:Out]]
```

## Example Usage

Given the input:

```
./protog Greet.v1 \
    -m HelloRequest[data:string] \
    -m HelloResponse[id:int64,data:string] \
    -s HelloService[SendHello:HelloRequest:HelloResponse,CheckHello]
```

You should get the file `greet.v1.proto` with the content:

```
syntax = "proto3";

package Greet.v1;

import "google/protobuf/empty.proto";

message HelloRequest {
	string data = 1;
}

message HelloResponse {
	string data = 1;
	int64 id = 2;
}

service HelloService {
	rpc SendHello (HelloRequest) returns (HelloResponse) {};
	rpc CheckHello (google.protobuf.Empty) returns (google.protobuf.Empty) {};
}
```

## License

MIT

## Author

Murilo Santana <<mvrilo@gmail.com>>
