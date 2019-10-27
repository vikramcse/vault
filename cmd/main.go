package main

import (
	"fmt"

	"github.com/vikramcse/vault"
)

func main() {
	v := vault.NewFileVault("my-fake-key", "sec.txt")
	err := v.Set("twitter", "this is a secret key for twitter developer apis")
	if err != nil {
		panic(err)
	}

	plain, err := v.Get("twitter")
	if err != nil {
		panic(err)
	}

	fmt.Println("plain text is:", plain)
}
