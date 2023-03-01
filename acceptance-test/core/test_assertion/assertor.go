package test_assertion

import (
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type GenericAssertor struct {
}

func (ga *GenericAssertor) Assert(actual interface{}, checker Checker, expect interface{}) (bool, error) {
	return checker.Check(actual, expect)
}

type JsonAssertor struct {
	object any
}

func (ja *JsonAssertor) Object(object any) *JsonAssertor {
	ja.object = object
	return ja
}

func (ja *JsonAssertor) Assert(actual interface{}, checker Checker, expect interface{}) (bool, error) {
	if ja.object == nil {
		return false, errors.New("when you assert a json value, the object should not be null")
	}
	result := gjson.Get(ja.object.(string), actual.(string))
	return checker.Check(result.String(), expect)
}
