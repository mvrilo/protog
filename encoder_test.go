package protog

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		in      map[string]interface{}
		expects []byte
		err     error
	}{
		{
			in:  map[string]interface{}{"syntax": 3},
			err: errSyntaxType,
		},
		{
			in:  map[string]interface{}{"message": "Hello"},
			err: errMessageType,
		},
		{
			in:  map[string]interface{}{"service": "Hello"},
			err: errServiceType,
		},
		{
			in:      map[string]interface{}{"syntax": "proto3"},
			expects: []byte(`syntax = "proto3";`),
		},
		{
			in: map[string]interface{}{
				"syntax":  "proto3",
				"package": "Hello",
			},
			expects: []byte("syntax = \"proto3\";\n\npackage Hello;"),
		},
		{
			in: map[string]interface{}{
				"syntax":  "proto3",
				"package": "Hello",
				"option": [][]string{
					[]string{"go_package", "proto"},
				},
			},
			expects: []byte("syntax = \"proto3\";\n\npackage Hello;\n\noption \"go_package\" = \"proto\";"),
		},
		{
			in: map[string]interface{}{
				"syntax":  "proto3",
				"package": "Hello",
				"message": map[string]interface{}{
					"Hello": map[string]string{
						"id": "int64",
					},
				},
			},
			expects: []byte("syntax = \"proto3\";\n\npackage Hello;\n\nmessage Hello {\n\tint64 id = 1;\n}"),
		},
		{
			in: map[string]interface{}{
				"syntax":  "proto3",
				"package": "Hello",
				"message": map[string]interface{}{
					"HelloRequest": map[string]string{
						"data": "string",
					},
				},
				"service": map[string]interface{}{
					"HelloService": map[string]interface{}{
						"SayHello": map[string]string{
							"in":  "HelloRequest",
							"out": "",
						},
					},
				},
			},
			expects: []byte("syntax = \"proto3\";\n\npackage Hello;\n\nimport \"google/protobuf/empty.proto\";\n\nmessage HelloRequest {\n\tstring data = 1;\n}\n\nservice HelloService {\n\trpc SayHello (HelloRequest) returns (google.protobuf.Empty) {};\n}"),
		},
		{
			in: map[string]interface{}{
				"syntax":  "proto3",
				"package": "Hello",
				"message": map[string]interface{}{
					"HelloRequest": map[string]string{
						"data": "string",
					},
					"HelloResponse": map[string]string{
						"id": "int64",
					},
				},
				"service": map[string]interface{}{
					"HelloService": map[string]interface{}{
						"SayHello": map[string]string{
							"in":  "HelloRequest",
							"out": "HelloResponse",
						},
					},
				},
			},
			expects: []byte("syntax = \"proto3\";\n\npackage Hello;\n\nmessage HelloRequest {\n\tstring data = 1;\n}\n\nmessage HelloResponse {\n\tint64 id = 1;\n}\n\nservice HelloService {\n\trpc SayHello (HelloRequest) returns (HelloResponse) {};\n}"),
		},
	}

	for _, tt := range tests {
		out, err := Encode(tt.in)
		if tt.err != nil && err != tt.err {
			t.Errorf("expects error %s, got %s", tt.expects, err)
		}

		if !bytes.Equal(out, tt.expects) {
			t.Errorf("expects %+v, got %+v\n\n", string(tt.expects), string(out))
		}
	}
}
