package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	f, err := excelize.OpenFile("superstor.cn.xlsx")
	if err != nil {
		panic(err)
	}

	rows, err := f.GetRows("abc")
	if err != nil {
		panic(err)
	}

	fmt.Println(rows)
}
