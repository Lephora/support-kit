package unit_test

import (
	"gopkg.in/yaml.v2"
	"os"
	"thoughtworks.com/lephora/support-kit/accecptance-test/common"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_suite"

	. "gopkg.in/check.v1"
)

type RecordCaseSuite struct{}

var _ = Suite(&RecordCaseSuite{})

func (s *RecordCaseSuite) SetUpTest(c *C) {
	_ = os.Setenv(common.PhaseEnv, string(common.Recording))
}

func (s *RecordCaseSuite) TestCaseExecuteWithCiteStageResult(c *C) {
	// given
	file, _ := os.ReadFile("./data/cite-stage-result-test.yaml")
	tc := &test_suite.Case{}
	// when
	err := yaml.Unmarshal(file, &tc)
	tc.SetName("data/cite-stage-result-test.yaml")
	// then
	c.Check(err, IsNil)

	// when
	err = tc.Execute()
	// then
	c.Check(err, IsNil)

}
