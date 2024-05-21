package main

import "fmt"

func main() {
	ok, err := say("Hello World")
	if err != nil {
		panic(err.Error())
	}

	switch ok {
	case true:
		fmt.Println("deu certo")
	default:
		fmt.Println("Deu errado")
	}
}

func say(what string) (bool, error) {
	if what == "" {
		return false, fmt.Errorf("Empty string")
	}
	fmt.Println(what)
	return true, nil
}
