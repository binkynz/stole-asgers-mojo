//go:build gen

package main

import (
	sqlc "github.com/sqlc-dev/sqlc/pkg/cli"
)

func main() {
	sqlc.Run([]string{"generate", "-f", ".sqlc.yaml"})
}
