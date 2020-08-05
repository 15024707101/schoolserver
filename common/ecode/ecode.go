package ecode

import "fmt"

var codeMap = map[int]struct{}{}

type ECode struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (c ECode) Error() string {
	return c.Msg
}

func (c ECode) Res() (int, string) {
	return c.Code, c.Msg
}

func New(i int, m string) ECode {
	if _, exist := codeMap[i]; exist {
		panic(fmt.Sprintf("ECode: %d already exist", i))
	}
	codeMap[i] = struct{}{}

	return ECode{i, m}
}
