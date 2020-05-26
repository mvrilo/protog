package protog

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

// tokens
const (
	tSyntax      = "syntax"
	tPackage     = "package"
	tMessage     = "message"
	tService     = "service"
	tImport      = "import"
	tOption      = "option"
	tRPC         = "rpc"
	tReturns     = "returns"
	tImportEmpty = "google/protobuf/empty.proto"
	tEmpty       = "google.protobuf.Empty"
	tBlockStart  = "{"
	tBlockEnd    = "}"
	tNewline     = "\n"
	tSpace       = " "
	tTab         = "\t"
	tEol         = ";"
)

// errors
var (
	errSyntaxType        = errors.New("syntax expects a string")
	errPackageType       = errors.New("package expects a string")
	errOptionType        = errors.New("package expects a []string")
	errMessageType       = errors.New("invalid message")
	errServiceType       = errors.New("invalid service")
	errServiceMethodType = errors.New("invalid service method")
)

// Encoder struct
type Encoder struct {
	buf         *bytes.Buffer
	lines       []string
	importEmpty bool

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
		e.Write(tNewline)
	}
}

func (e *Encoder) writeSpace() {
	e.Write(tSpace)
}

func (e *Encoder) writeTab() {
	if e.Indent || !e.Compact {
		e.Write(tTab)
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

func (e *Encoder) writePackage(value string) {
	e.Write(tPackage)
	e.Write(tSpace)
	e.Write(value)
	e.Write(tEol)
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

// WriteOption writes an option
func (e *Encoder) WriteOption(value interface{}) error {
	v, ok := value.([]string)
	if !ok {
		return errOptionType
	}

	if len(v) < 2 {
		return nil
	}

	e.Write(tOption)
	e.Write(tSpace)

	e.Write(`"` + v[0] + `" = "` + v[1] + `"`)
	e.Write(tEol)
	e.lines[2] = e.buf.String()
	e.buf.Reset()
	return nil
}

// WriteImport writes a new import
func (e *Encoder) WriteImport(value string) error {
	e.Write(tImport)
	e.Write(tSpace)
	e.WriteValue(value)
	e.Write(tEol)

	e.lines[3] = e.buf.String()
	e.buf.Reset()
	return nil
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

func (e *Encoder) writeMessage(name string, fields map[string]string) {
	e.Write(tMessage)

	e.Write(tSpace)
	e.Write(name)
	e.Write(tSpace)

	e.Write(tBlockStart)

	var i int
	for fieldName, fieldType := range fields {
		i++
		e.writeNL()
		e.writeTab()
		e.Write(fieldType)
		e.writeSpace()
		e.Write(fieldName)
		e.writeSpace()
		e.Write("=")
		e.writeSpace()
		e.Write(strconv.Itoa(i))
		e.Write(tEol)
	}
	e.writeNL()
	e.Write(tBlockEnd)
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

	e.lines[4] = e.buf.String()
	e.buf.Reset()
	return nil
}

func (e *Encoder) parseServices(value interface{}) (map[string]map[string]map[string]string, error) {
	svcs, ok := value.(map[string]interface{})
	if !ok {
		return nil, errServiceType
	}

	services := make(map[string]map[string]map[string]string)

	for svcName, svcMethods := range svcs {
		services[svcName] = make(map[string]map[string]string)
		methods, ok := svcMethods.(map[string]interface{})
		if !ok {
			return nil, errServiceMethodType
		}

		for methodName, methodData := range methods {
			data, ok := methodData.(map[string]string)
			if !ok {
				return nil, errServiceMethodType
			}

			in := data["in"]
			if strings.TrimSpace(in) == "" {
				e.importEmpty = true
				in = tEmpty
			}

			out := data["out"]
			if strings.TrimSpace(out) == "" {
				e.importEmpty = true
				out = tEmpty
			}

			services[svcName][methodName] = map[string]string{
				"in":  in,
				"out": out,
			}
		}
	}

	return services, nil
}

func (e *Encoder) writeService(name string, methods map[string]map[string]string) {
	e.Write(tService)

	e.Write(tSpace)
	e.Write(name)
	e.Write(tSpace)

	e.Write(tBlockStart)

	var i int
	for methodName, methodData := range methods {
		i++
		e.writeNL()
		e.writeTab()

		e.Write(tRPC)
		e.Write(tSpace)
		e.Write(methodName)
		e.writeSpace()
		e.Write("(")
		e.Write(methodData["in"])
		e.Write(")")
		e.writeSpace()
		e.Write(tReturns)
		e.writeSpace()
		e.Write("(")
		e.Write(methodData["out"])
		e.Write(")")
		e.writeSpace()
		e.Write(tBlockStart)
		e.Write(tBlockEnd)
		e.Write(tEol)
	}
	e.writeNL()
	e.Write(tBlockEnd)
}

// WriteService writes a message block
func (e *Encoder) WriteService(value interface{}) error {
	services, err := e.parseServices(value)
	if err != nil {
		return err
	}

	var i int
	for svcName, svcMethods := range services {
		if i > 0 {
			e.writeNL()
		}
		e.writeService(svcName, svcMethods)
		i++
	}

	e.lines[5] = e.buf.String()
	e.buf.Reset()
	return nil
}

func (e *Encoder) Bytes() []byte {
	e.buf.Reset()

	for i, line := range e.lines {
		if len(line) == 0 {
			continue
		}

		if i > 0 {
			e.writeNL()
			e.writeNL()
		}
		e.buf.WriteString(line)
	}

	return e.buf.Bytes()
}

func (e *Encoder) Encode(data map[string]interface{}) ([]byte, error) {
	e.lines = make([]string, 6)

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
		case tOption:
			err := e.WriteOption(value)
			if err != nil {
				return nil, err
			}
		case tMessage:
			err := e.WriteMessage(value)
			if err != nil {
				return nil, err
			}
		case tService:
			err := e.WriteService(value)
			if err != nil {
				return nil, err
			}

			if e.importEmpty {
				err := e.WriteImport(tImportEmpty)
				if err != nil {
					return nil, err
				}
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
