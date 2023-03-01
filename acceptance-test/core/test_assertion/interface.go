package test_assertion

type (
	Assertor interface {
		Assert(actual any, checker Checker, expect any) (bool, error)
	}
	Checker interface {
		Check(actual, expect any) (result bool, error error)
	}
)

var (
	_ Checker  = (*isEqual)(nil)
	_ Checker  = (*isMatchRegex)(nil)
	_ Assertor = (*GenericAssertor)(nil)
	_ Assertor = (*JsonAssertor)(nil)
)