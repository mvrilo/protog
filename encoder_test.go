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
			in: map[string]interface{}{
				"syntax": 3,
			},
			err: errSyntaxType,
		},

		{
			in: map[string]interface{}{
				"syntax": "proto3",
			},
			expects: []byte(`syntax = "proto3";`),
		},

		{
			in: map[string]interface{}{
				"syntax":  "proto3",
				"package": "Hello",
			},
			expects: []byte("syntax = \"proto3\";\npackage Hello;"),
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
			expects: []byte("syntax = \"proto3\";\npackage Hello;\nmessage Hello {\n\tint64 id = 1;\n}"),
		},
	}

	for _, tt := range tests {
		out, err := Encode(tt.in)
		if tt.err != nil && err != tt.err {
			t.Errorf("expects error %s, got %s", tt.expects, err)
		}

		if !bytes.Equal(out, tt.expects) {
			t.Errorf("expects %+v, got %+v\n", string(tt.expects), string(out))
		}
	}
}
