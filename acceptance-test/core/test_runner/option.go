package test_runner

import "github.com/pkg/errors"

type Option struct {
	Roots  []string
	Suites []string
	Cases  []string
}

type RunnerType string

var (
	RunOption = &Option{}
	ROOT      = RunnerType("root")
	SUITE     = RunnerType("suite")
	CASE      = RunnerType("case")
)

func (receiver *Option) Status() (RunnerType, error) {
	if len(receiver.Roots) != 0 && len(receiver.Suites) != 0 && len(receiver.Cases) != 0 {
		return "", errors.New("you must appoint a testrunner type")
	}
	if (len(receiver.Roots) != 0 && len(receiver.Suites) != 0) || (len(receiver.Roots) != 0 && len(receiver.Cases) != 0) || (len(receiver.Cases) != 0 && len(receiver.Suites) != 0) {
		return "", errors.New("only one type can be appointed")
	}
	switch {
	case len(receiver.Roots) != 0:
		return ROOT, nil
	case len(receiver.Suites) != 0:
		return SUITE, nil
	case len(receiver.Cases) != 0:
		return CASE, nil
	default:
		return "", errors.New("unknown status!")
	}
}
