/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package main

import (
	"os"

	"github.com/eiladin/k8s-dotenv/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(version, os.Args[1:])
}
