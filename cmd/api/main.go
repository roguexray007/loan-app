//go:build boot
// +build boot

package main

import (
	"context"

	"github.com/roguexray007/loan-app/internal/boot"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	(&boot.API{}).Init(ctx)
}
