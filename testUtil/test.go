package testUtil

import (
	"fmt"
	"strings"
)

func Test(fileName string) {
	temp := readFile(fileName)
	// 去除最后的空行
	temp = temp[:len(temp)-1]
	var slice = make([][]string, len(temp))
	for i := 0; i < len(temp); i++ {
		slice[i] = make([]string, 4)
		a := strings.Split(temp[i], ",")
		copy(slice[i], a)
	}
	num1, num2 := 0, 0
	tempStr := "a"
	for _, single := range slice {
		a, b, c, d := single[0], single[1], single[2], single[3]
		tempStr += (a + b + c + d)
		if a == "False" && b == "True" {
			if c == "True" {
				num1++
			} else {
				num2++
			}
		}
	}
	fmt.Println("*********************\n", num1, num2)
}
