package utils

import "fmt"

func Mylog(a ... interface{}) {
	fmt.Println(a...)
}
func MylogF(format string, a ... interface{}) {
	fmt.Printf(format, a...)
}


