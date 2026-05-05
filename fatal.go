package main

import (
	"fmt"
	"os"
)

func fatalErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func fatalMsg(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
