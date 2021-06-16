package main

import "fmt"

func main() {
	strs := []string{"a"}

	for _, str := range strs {
		strs = append(strs, str)
	}
	fmt.Println(strs)
}
