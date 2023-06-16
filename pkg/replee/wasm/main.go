package main

import (
	"context"
	"encoding/json"
	"github.com/dop251/goja"
	"github.com/perdasilva/replee/pkg/replee/repl"
	"log"
	"strconv"
	"strings"
	"syscall/js"
)

const (
	colorDefault = `\x1b[32m`
	colorError   = `\x1b[31m`
	colorSuccess = `\x1b[90m`

	modeStartOfInput = "startOfInput"
	modeMidInput     = "midInput"
	modeReset        = "reset"
)

type Response struct {
	Mode   string      `json:"mode"`
	IsErr  bool        `json:"error"`
	Output interface{} `json:"output"`
	Color  string      `json:"color"`
	Indent int         `json:"indent"`
}

type Request struct {
	Mode   string `json:"mode"`
	Input  string `json:"input"`
	Indent int    `json:"indent"`
}

func noOp(mode string, indent int) *Response {
	return &Response{
		IsErr:  false,
		Output: "",
		Color:  colorDefault,
		Indent: indent,
		Mode:   mode,
	}
}

func errOp(msg string) *Response {
	return &Response{
		IsErr:  true,
		Output: msg,
		Color:  colorError,
		Mode:   modeStartOfInput,
	}
}

func stringify(response *Response) string {
	bytes, err := json.Marshal(response)
	if err != nil {
		return err.Error()
	}
	return string(bytes)
}

func validateNotNull(name string, value js.Value) *Response {
	if value.IsNull() || value.IsUndefined() {
		return errOp("replee error: '" + name + "' is null or undefined")
	}
	return nil
}

func isEmpty(value js.Value) bool {
	return strings.Trim(strings.TrimSpace(value.String()), "\t\n\r") == ""
}

func parseRequest(input js.Value) (Request, *Response) {
	request := Request{
		Mode:   "",
		Input:  "",
		Indent: 0,
	}

	if err := validateNotNull("input", input); err != nil {
		return request, err
	}

	mode := input.Get("mode")
	if err := validateNotNull("mode", mode); err != nil {
		return request, err
	}
	if mode.String() != modeStartOfInput || mode.String() != modeMidInput {
		return request, errOp("replee error: mode must be one of '" + modeStartOfInput + "' or '" + modeMidInput + "'")
	}
	request.Mode = mode.String()

	indent := input.Get("indent")
	if err := validateNotNull("indent", indent); err != nil {
		return request, err
	}
	i, err := strconv.Atoi(indent.String())
	if err != nil {
		return request, errOp("replee error: 'indent' must be an int")
	}
	request.Indent = i

	stdin := input.Get("input")
	if stdin.IsNull() || stdin.IsUndefined() {
		return request, noOp(request.Mode, request.Indent)
	}
	if isEmpty(stdin) {
		return request, noOp(request.Mode, request.Indent)
	}
	request.Input = stdin.String()

	return request, nil
}

func main() {
	ctx := context.Background()
	vm := goja.New()
	vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))

	if err := repl.BootstrapRepleeVM(ctx, vm); err != nil {
		log.Fatal(err)
	}

	js.Global().Set("replee", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		input := args[0].String()
		request := &Request{}
		if err := json.Unmarshal([]byte(input), request); err != nil {
			return "error"
		}

		response := map[string]interface{}{
			"isErr":  false,
			"output": "",
			"mode":   modeReset,
			"indent": 0,
			"color":  colorSuccess,
		}

		value, err := vm.RunString(request.Input)
		if err != nil {
			response["isErr"] = true
			response["output"] = err.Error()
		} else {
			if goja.IsNull(value) || goja.IsUndefined(value) {
				response["output"] = ""
			} else {
				response["output"] = value.String()
			}
		}

		return js.ValueOf(response)
	}))

	<-(make(chan struct{}))
}
