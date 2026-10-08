// napcatgen 根据 NapCat OpenAPI 生成 Go API 绑定。
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/zjutjh/onebot-sdk/internal/gen"
)

func main() {
	spec := flag.String("spec", "internal/openapi/4.18.33/openapi.json", "OpenAPI spec 路径")
	out := flag.String("out", "api", "生成输出目录")
	flag.Parse()

	if err := gen.GenerateFromFile(*spec, *out); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
