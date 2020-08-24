package framework

import (
	"errors"
	"strings"
)

var WrongStateError = errors.New("can not take the operation in the current state")

// 自定义一个error类,并实现Error()接口
type ServicesError struct {
	errArr []error
}

func (se ServicesError) Error() string {
	var ret []string
	for _, err := range se.errArr {
		ret = append(ret, err.Error())
	}
	return strings.Join(ret, ";")
}
