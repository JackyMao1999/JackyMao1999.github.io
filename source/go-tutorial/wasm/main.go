package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

var i *interp.Interpreter

type stringWriter struct {
	b *strings.Builder
}

func (w *stringWriter) Write(p []byte) (int, error) {
	return w.b.Write(p)
}

func main() {
	i = newInterpreter()

	js.Global().Set("_goWasmRun", js.FuncOf(runGoJS))
	js.Global().Set("_goWasmTest", js.FuncOf(runGoTestJS))
	js.Global().Set("_goWasmReset", js.FuncOf(resetGoJS))

	fmt.Println("Go WASM runtime ready")

	select {}
}

func newInterpreter() *interp.Interpreter {
	return interp.New(interp.Options{})
}

func runGoJS(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return map[string]any{"error": "no code provided"}
	}

	code := args[0].String()

	var stdoutBuf, stderrBuf strings.Builder
	stdoutW := &stringWriter{&stdoutBuf}
	stderrW := &stringWriter{&stderrBuf}

	ip := interp.New(interp.Options{
		Stdout: stdoutW,
		Stderr: stderrW,
	})
	ip.Use(stdlib.Symbols)

	wrapped := wrapCode(code)

	_, evalErr := ip.Eval(wrapped)

	if evalErr != nil {
		stdout := strings.TrimSpace(stdoutBuf.String())
		errStr := evalErr.Error()
		return map[string]any{
			"error":  errStr,
			"stdout": stdout,
		}
	}

	stdout := strings.TrimSpace(stdoutBuf.String())
	stderr := strings.TrimSpace(stderrBuf.String())

	if stderr != "" {
		return map[string]any{
			"error":  stderr,
			"stdout": stdout,
		}
	}

	return map[string]any{
		"stdout": stdout,
	}
}

func runGoTestJS(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return map[string]any{"error": "missing arguments: functionCode, functionName, testCasesJSON"}
	}

	funcCode := args[0].String()
	funcName := args[1].String()
	testCasesJSON := args[2].String()

	result, err := evalTests(funcCode, funcName, testCasesJSON)
	if err != nil {
		return map[string]any{"error": err.Error()}
	}

	results := make([]any, len(result))
	for idx, r := range result {
		results[idx] = map[string]any{
			"passed":   r.passed,
			"got":      r.got,
			"expected": r.expected,
			"args":     r.args,
			"error":    r.err,
		}
	}

	return map[string]any{
		"results": results,
	}
}

func resetGoJS(this js.Value, args []js.Value) any {
	i = newInterpreter()
	return "reset"
}

type testResult struct {
	passed   bool
	got      string
	expected string
	args     string
	err      string
}

func wrapCode(code string) string {
	code = strings.TrimSpace(code)
	if strings.HasPrefix(code, "package ") {
		return code
	}
	if strings.HasPrefix(code, "import ") {
		return "package main\n" + code
	}
	return "package main\n" + code
}

func evalTests(funcCode, funcName, testCasesJSON string) ([]testResult, error) {
	wrapped := wrapCode(funcCode)

	ip := interp.New(interp.Options{})
	ip.Use(stdlib.Symbols)

	_, evalErr := ip.Eval(wrapped)
	if evalErr != nil {
		return nil, fmt.Errorf("compile error: %v", evalErr)
	}

	cases, err := parseTestCases(testCasesJSON)
	if err != nil {
		return nil, fmt.Errorf("parse error: %v", err)
	}

	results := make([]testResult, 0, len(cases))
	errorsCount := 0
	for _, tc := range cases {
		results = append(results, runTestCase(ip, funcName, tc, &errorsCount))
	}

	return results, nil
}

func parseTestCases(jsonStr string) ([]testCase, error) {
	jsonVal := js.Global().Get("JSON").Call("parse", jsonStr)
	if !jsonVal.Truthy() {
		return nil, fmt.Errorf("invalid JSON")
	}

	length := jsonVal.Get("length").Int()
	cases := make([]testCase, length)

	for i := 0; i < length; i++ {
		item := jsonVal.Index(i)
		noexpect := item.Get("noexpect")
		tc := testCase{
			NoExpect: noexpect.Type() == js.TypeBoolean && noexpect.Bool(),
		}

		expectVal := item.Get("expect")
		if expectVal.Type() != js.TypeUndefined {
			switch expectVal.Type() {
			case js.TypeString:
				tc.Expected = expectVal.String()
			case js.TypeNumber:
				tc.Expected = expectVal.Call("toString").String()
			case js.TypeBoolean:
				if expectVal.Bool() {
					tc.Expected = "true"
				} else {
					tc.Expected = "false"
				}
			default:
				tc.Expected = expectVal.Call("toString").String()
			}
		} else {
			expectedField := item.Get("expected")
			if expectedField.Type() != js.TypeUndefined {
				tc.Expected = expectedField.Call("toString").String()
			}
		}

		argsArr := item.Get("args")
		if argsArr.Type() == js.TypeObject {
			argsLen := argsArr.Get("length").Int()
			tc.Args = make([]any, argsLen)
			for j := 0; j < argsLen; j++ {
				arg := argsArr.Index(j)
				switch arg.Type() {
				case js.TypeNumber:
					tc.Args[j] = arg.Float()
				case js.TypeBoolean:
					tc.Args[j] = arg.Bool()
				case js.TypeString:
					tc.Args[j] = arg.String()
				case js.TypeNull:
					tc.Args[j] = nil
				default:
					tc.Args[j] = arg.String()
				}
			}
		}

		cases[i] = tc
	}

	return cases, nil
}

type testCase struct {
	Args     []any
	Expected string
	NoExpect bool
}

func runTestCase(ip *interp.Interpreter, funcName string, tc testCase, errorsCount *int) testResult {
	argsStr := argsToString(tc.Args)
	callExpr := fmt.Sprintf("%s(%s)", funcName, argsStr)

	v, err := ip.Eval(callExpr)
	if err != nil {
		*errorsCount++
		return testResult{
			passed:   false,
			got:      "",
			expected: tc.Expected,
			args:     argsStr,
			err:      err.Error(),
		}
	}

	got := formatValue(v)

	passed := false
	if !tc.NoExpect {
		passed = normalizeStr(got) == normalizeStr(tc.Expected)
	} else {
		passed = true
	}

	return testResult{
		passed:   passed,
		got:      got,
		expected: tc.Expected,
		args:     argsStr,
		err:      "",
	}
}

func argsToString(args []any) string {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = toGoLiteral(a)
	}
	return strings.Join(parts, ", ")
}

func toGoLiteral(v any) string {
	switch val := v.(type) {
	case nil:
		return "nil"
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%v", val)
	case string:
		return fmt.Sprintf("%q", val)
	default:
		return fmt.Sprintf("%q", fmt.Sprint(v))
	}
}

func formatValue(v interface{}) string {
	if v == nil {
		return "nil"
	}
	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%d", val)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case float32, float64:
		f := fmt.Sprintf("%v", val)
		return f
	case bool:
		if val {
			return "true"
		}
		return "false"
	case []interface{}:
		parts := make([]string, len(val))
		for i, item := range val {
			parts[i] = formatValue(item)
		}
		return "[" + strings.Join(parts, " ") + "]"
	case map[interface{}]interface{}:
		parts := make([]string, 0)
		for k, v := range val {
			parts = append(parts, fmt.Sprintf("%v:%v", formatValue(k), formatValue(v)))
		}
		return "map[" + strings.Join(parts, " ") + "]"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func normalizeStr(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"'")
	return s
}
