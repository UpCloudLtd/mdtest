package main

import (
	"os"

	"github.com/UpCloudLtd/mdtest/cmd"
)

func main() {
	code := cmd.Execute()
	os.Exit(code)
}
