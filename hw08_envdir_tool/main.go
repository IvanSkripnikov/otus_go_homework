package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 2 {
		fmt.Println("Enter more arguments")
	} else {
		envs, err := ReadDir(os.Args[1])
		if err != nil {
			log.Fatalln(err)
		}

		RunCmd(os.Args[2:], envs)
	}
}
