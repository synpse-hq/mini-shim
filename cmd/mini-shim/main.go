package main

import (
	"context"
	"fmt"
	"os"

	"github.com/synpse-hq/mini-shim/pkg/proxy"
)

func main() {

	cfg, err := proxy.Load()
	if err != nil {
		fmt.Println("failed to load proxy config: ", err)
		os.Exit(1)
	}

	p := proxy.New(cfg)

	ctx := context.Background()

	err = p.Serve(ctx)
	if err != nil {
		fmt.Println("error: ", err)
	}
}
