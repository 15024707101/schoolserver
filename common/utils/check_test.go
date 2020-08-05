package utils

import (
	"fmt"
	"testing"
)

func TestCheckCard(t *testing.T) {
	// 阮中伟 41030519800209353X 男	37 河南省 洛阳市 涧西区
	var idcard = "110105710923582"
	fmt.Println(CheckCard(idcard))
}

func TestCheckString(t *testing.T) {
	var str = "321@#$%^*"
	fmt.Println("检验字符串强度 start:")
	s := CheckString(str)
	fmt.Printf("检验字符串强度 end :%s\n", s)
}
