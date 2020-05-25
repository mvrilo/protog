package protog

import (
	"bytes"
	"errors"
	"strconv"
)

// tokens
const (
	tSyntax     = "syntax"
	tPackage    = "package"
	tMessage    = "message"
	tService    = "service"
	tBlockStart = "{"
	tBlockEnd   = "}"
	tSpace      = " "
	tEol        = ";"
)

// errors
var (
	errSyntaxType  = errors.New("syntax expects a string")
	errPackageType = errors.New("package expects a string")
	errMessageType = errors.New("message expects a map[string]interface{}")
)

// Encoder struct
type Encoder struct {
	buf   *bytes.Buffer
	lines []string

	Indent  bool
	Compact bool
}

// New returns an encoder object with initialized buffer and indentation enabled
func New() *Encoder {
	return &Encoder{
		buf:    bytes.NewBuffer(nil),
		Indent: true,
	}
}

// Write is a proxy for buf.WriteString
func (e *Encoder) Write(value string) {
	e.buf.WriteString(value)
}

func (e *Encoder) writeNL() {
	if !e.Compact {
		e.Write("\n")
	}
}

// writeTab writes a tab (\t)"`
func (e *Encoder) writeTab() {
	if e.Indent || !e.Compact {
		e.Write("\t")
	}
}

// WriteV writes a value inside `"`
func (e *Encoder) WriteValue(value string) {
	e.Write(`"` + value + `"`)
}

// WriteAssignment assigns value to name, e.g. `name = "value";`
func (e *Encoder) WriteAssignment(name, value string) {
	e.Write(name + ` = "` + value + `"`)
	e.Write(tEol)
}

// WriteMessage writes a message block
func (e *Encoder) WriteMessage(value interface{}) error {
	messages, err := e.parseMessages(value)
	if err != nil {
		return err
	}

	var i int
	for msgName, msgFields := range messages {
		if i > 0 {
			e.writeNL()
			e.writeNL()
		}
		e.writeMessage(msgName, msgFields)
		i++
	}

	e.lines[2] = e.buf.String()
	e.buf.Reset()
	return nil
}

func (e *Encoder) writeMessage(name string, fields map[string]string) {
	e.Write(tMessage)

	e.Write(tSpace)
	e.Write(name)
	e.Write(tSpace)

	e.Write(tBlockStart)

	var i int = 1
	for fieldName, fieldType := range fields {
		e.writeNL()
		e.writeNL()
		e.writeTab()
		e.Write(fieldType + " " + fieldName + " = " + strconv.Itoa(i))
		e.Write(tEol)
	}
	e.writeNL()
	e.writeNL()
	e.Write(tBlockEnd)
}

// WriteSyntax writes a syntax
func (e *Encoder) WriteSyntax(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errSyntaxType
	}

	e.WriteAssignment(tSyntax, v)
	e.lines[0] = e.buf.String()
	e.buf.Reset()
	return nil
}

// WritePackage writes a package
func (e *Encoder) WritePackage(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return errPackageType
	}

	e.writePackage(v)
	e.lines[1] = e.buf.String()
	e.buf.Reset()
	return nil
}

func (e *Encoder) writePackage(value string) {
	e.Write(tPackage)
	e.Write(tSpace)
	e.Write(value)
	e.Write(tEol)
}

func (e *Encoder) parseMessages(value interface{}) (map[string]map[string]string, error) {
	msgs, ok := value.(map[string]interface{})
	if !ok {
		return nil, errMessageType
	}

	messages := make(map[string]map[string]string)

	for msgName, msgFieldsI := range msgs {
		messages[msgName] = map[string]string{}

		fields, ok := msgFieldsI.(map[string]string)
		if !ok {
			return nil, errMessageType
		}

		for fieldName, fieldValue := range fields {
			messages[msgName][fieldName] = fieldValue
		}
	}

	return messages, nil
}

func (e *Encoder) Bytes() []byte {
	e.buf.Reset()

	for i, line := range e.lines {
		if i > 0 {
			e.writeNL()
			e.writeNL()
		}
		e.buf.WriteString(line)
	}

	return e.buf.Bytes()
}

func (e *Encoder) Encode(data map[string]interface{}) ([]byte, error) {
	e.lines = make([]string, len(data))

	for key, value := range data {
		switch key {
		case tSyntax:
			err := e.WriteSyntax(value)
			if err != nil {
				return nil, err
			}
		case tPackage:
			err := e.WritePackage(value)
			if err != nil {
				return nil, err
			}
		case tMessage:
			err := e.WriteMessage(value)
			if err != nil {
				return nil, err
			}
		default:
		}
	}

	return e.Bytes(), nil
}

// Encode encodes a map to []byte
func Encode(v map[string]interface{}) ([]byte, error) {
	return New().Encode(v)
}
