//go:build js && wasm

package main

import (
	"context"
	"fmt"
	"syscall/js"

	"github.com/prnvbn/protoc-gen-sbexml/internal/protosource"
)

func main() {
	js.Global().Set("generateSBE", promiseFunc(generateSBE))
	select {}
}

func generateSBE(args []js.Value) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("generateSBE expects 1 argument, got %d", len(args))
	}
	return protosource.GenerateSBE(context.Background(), args[0].String())
}

func promiseFunc(fn func([]js.Value) (string, error)) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		executor := js.FuncOf(func(_ js.Value, promiseArgs []js.Value) any {
			resolve := promiseArgs[0]
			reject := promiseArgs[1]

			go func() {
				result, err := fn(args)
				if err != nil {
					reject.Invoke(js.Global().Get("Error").New(err.Error()))
					return
				}
				resolve.Invoke(result)
			}()

			return nil
		})
		promise := js.Global().Get("Promise").New(executor)
		executor.Release()
		return promise
	})
}
