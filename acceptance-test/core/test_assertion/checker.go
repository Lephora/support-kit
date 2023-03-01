package test_assertion

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
)

var Equals isEqual
var Match isMatchRegex

type isEqual struct {
}

func (i isEqual) Check(actual, expect any) (result bool, error error) {
	if actual == expect {
		return true, nil
	}
	return false, errors.New(fmt.Sprintf("%T(%v) is not equals %T(%v)", actual, actual, expect, expect))
}

type isMatchRegex struct {
}

func (i isMatchRegex) Check(actual, expect interface{}) (result bool, err error) {
	compile, err := regexp.Compile(expect.(string))
	if err != nil {
		return false, err
	}
	isMatch := compile.MatchString(actual.(string))
	return isMatch, nil
}
