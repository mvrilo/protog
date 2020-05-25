# protog

## Description

Protobuf file generator for the command line.

## Installation

```
go get github.com/mvrilo/protog/cmd/protog
```

## Usage

```
protog <name> [-m Message[id:int64,fields:types,...]]
```

## Example Usage

Given the input:

```
./protog Greet.v1 \
    -m HelloRequest[data:string] \
    -m HelloResponse[id:int64,data:string]
```

You should get the file `greet.v1.proto` with the content:

```
syntax = "proto3";

package Greet.v1;

message HelloRequest {
	string data = 1;
}

message HelloResponse {
	int64 id = 1;
	string data = 2;
}
```

## Author

Murilo Santana <<mvrilo@gmail.com>>
