package main

import (
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

const (
	iterArrayLen = 50
)

func compile(src string, opts []cel.EnvOption) (cel.Program, error) {
	env, err := cel.NewEnv(opts...)
	if err != nil {
		return nil, err
	}
	ast, issues := env.Compile(src)
	if issues != nil {
		return nil, issues.Err()
	}
	program, err := env.Program(ast)
	if err != nil {
		return nil, err
	}
	return program, nil
}

func BenchmarkNoop(b *testing.B) {
	noopProg, err := compile("true", []cel.EnvOption{})
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := noopProg.Eval(map[string]interface{}{})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func iterate(n int) {

}

func genArray(len int, data string) []string {
	a := make([]string, len)
	for i := 0; i < len; i++ {
		a[i] = data
	}
	return a
}

func BenchmarkIterate(b *testing.B) {
	iterProg, err := compile("l.all(x, x == 'golang')", []cel.EnvOption{
		cel.Declarations(decls.NewVar("l", decls.NewListType(decls.String))),
	})
	if err != nil {
		b.Fatal(err)
	}
	arr := genArray(iterArrayLen, "golang")
	input := map[string]interface{}{
		"l": arr,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := iterProg.Eval(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAccess(b *testing.B) {
	accessProg, err := compile("a.b == 'some str'", []cel.EnvOption{
		cel.Declarations(decls.NewVar("a", decls.NewMapType(decls.String, decls.String))),
	})
	if err != nil {
		b.Fatal(err)
	}
	input := map[string]interface{}{
		"a": map[string]string{
			"b": "golang",
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := accessProg.Eval(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
