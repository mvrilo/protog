# protog

Protobuf file generator for the command line.

## Installation

```
go get github.com/mvrilo/protog/cmd/protog
```

## Usage

```
$ protog -h
protog is a protobuf file generator for the command line

Usage:
  protog <name> [-dhfomnsp] [-n option_name:proto_name] [-m MessageName[field:type,field:type,...]] [-s ServiceName[MethodName:In:Out]]

Examples:
protog Greet.v1 -m HelloRequest[data:string]

Flags:
  -d, --dryrun                prints the generated proto to stdout
  -f, --force                 overwrite the file if it already exists
  -h, --help                  help for protog
  -m, --message stringArray   add a message and its fields
  -n, --option strings        add an option
  -o, --output string         output dir for the generated proto (default ".")
  -p, --package string        package name
  -s, --service stringArray   add a service and its methods
  -v, --version               version for protog
```

## Example Usage

Given the input:

```
./protog Greet.v1 \
    -n go_package:greet \
    -m HelloRequest[data:string] \
    -m HelloResponse[id:int64,data:string] \
    -s HelloService[SendHello:HelloRequest:HelloResponse,CheckHello]
```

You should get the file `greet.v1.proto` with the content:

```
syntax = "proto3";

package Greet.v1;

option "go_package" = "greet";

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

## Author

Murilo Santana <<mvrilo@gmail.com>>

## License

MIT
