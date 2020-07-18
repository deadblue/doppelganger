package app

import "fmt"

func Run() (err error) {
	fmt.Println(Version())
	return nil
}
