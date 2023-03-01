package unit_test

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
	"thoughtworks.com/lephora/support-kit/accecptance-test/common"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_report"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_suite"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type AssertCaseSuite struct{}

var _ = Suite(&AssertCaseSuite{})

func (s *AssertCaseSuite) SetUpTest(c *C) {
	_ = os.Setenv(common.PhaseEnv, string(common.Asserting))
	test_report.BuildReport()
	report := test_report.GetReport()
	report.AppendSuite("unit-test-suite")
	_ = os.Setenv("REPORT_STYLE", "testrunner-report.json")
}

func (s *AssertCaseSuite) TestUnmarshalCaseFile(c *C) {
	// given
	file, _ := os.ReadFile("./data/hello-world-test.yaml")
	tc := &test_suite.Case{}

	// when
	err := yaml.Unmarshal(file, tc)

	// then
	c.Check(err, IsNil)
	c.Check(tc.Stages, HasLen, 3)
	c.Check(tc.Stages[0].Request.Body, Equals, "{\"sequence\": \"hello world\"}")
}

func (s *AssertCaseSuite) TestCaseExecute(c *C) {
	// given
	file, _ := os.ReadFile("./data/hello-world-test.yaml")
	tc := &test_suite.Case{}

	// when
	_ = yaml.Unmarshal(file, tc)
	err := tc.Execute()

	// then
	c.Check(err, IsNil)
}

func (s *AssertCaseSuite) TestCaseExecuteWithFilter(c *C) {
	// given
	file, _ := os.ReadFile("./data/hello-world-test.yaml")
	tc := &test_suite.Case{}

	// when
	_ = yaml.Unmarshal(file, tc)
	tc.AddFilter(func(stage *test_suite.Stage) bool {
		return stage.Type == test_suite.API
	})
	err := tc.Execute()

	// then
	c.Check(err, IsNil)
}

func (s *AssertCaseSuite) TestCaseExecuteWithReportGen(c *C) {
	// given
	file, _ := os.ReadFile("./data/hello-world-test.yaml")
	tc := &test_suite.Case{}

	// when
	_ = yaml.Unmarshal(file, tc)
	tc.AddFilter(func(stage *test_suite.Stage) bool {
		return stage.Type == test_suite.API
	})
	err := tc.Execute()

	// then
	c.Check(err, IsNil)

	// when
	_, err = test_report.GetReport().Gen("gen")

	// then
	c.Check(err, IsNil)
}

func (s *AssertCaseSuite) TestCaseExecuteWithCaseVariable(c *C) {
	// given
	file, _ := os.ReadFile("./data/testsuite_case_level_vars/variables-test.yaml")
	tc := &test_suite.Case{}

	// when
	_ = yaml.Unmarshal(file, tc)
	err := tc.Execute()

	// then
	c.Check(err, IsNil)

	// when
	_, err = test_report.GetReport().Gen("gen/testsuite_case_level_vars")

	// then
	c.Check(err, IsNil)

	// when
	file, err = os.ReadFile(fmt.Sprintf("gen/testsuite_case_level_vars/%s.%s", common.ReportFile, common.JsonReportFileExt))

	// then
	c.Check(err, IsNil)

	// when
	report := make(map[string]any, 0)
	err = yaml.Unmarshal(file, &report)

	// then
	c.Check(err, IsNil)
	c.Check(report["pass"], Equals, true)
}

func (s *AssertCaseSuite) TestCaseExecuteWithSuiteVariable(c *C) {
	// given
	suite, err := test_suite.BuildTestSuite("./data/testsuite_suite_level_vars")

	// then
	c.Check(err, IsNil)

	// when
	err = suite.Execute()

	// then
	c.Check(err, IsNil)

	// when
	_, err = test_report.GetReport().Gen("./gen/testsuite_suite_level_vars")

	// then
	c.Check(err, IsNil)

	// when
	file, err := os.ReadFile(fmt.Sprintf("gen/testsuite_suite_level_vars/%s.%s", common.ReportFile, common.JsonReportFileExt))

	// then
	c.Check(err, IsNil)

	// when
	report := make(map[string]any, 0)
	err = yaml.Unmarshal(file, &report)

	// then
	c.Check(err, IsNil)
	c.Check(report["pass"], Equals, true)
}

func (s *AssertCaseSuite) TestCaseExecuteWithBothSuiteCaseVariable(c *C) {
	// given
	suite, err := test_suite.BuildTestSuite("./data/testsuite_both_suite_case_level_vars")

	// then
	c.Check(err, IsNil)

	// when
	err = suite.Execute()

	// then
	c.Check(err, IsNil)

	// when
	_, err = test_report.GetReport().Gen("./gen/testsuite_both_suite_case_level_vars")

	// then
	c.Check(err, IsNil)

	// when
	file, err := os.ReadFile(fmt.Sprintf("gen/testsuite_both_suite_case_level_vars/%s.%s", common.ReportFile, common.JsonReportFileExt))

	// then
	c.Check(err, IsNil)

	// when
	report := make(map[string]any, 0)
	err = yaml.Unmarshal(file, &report)

	// then
	c.Check(err, IsNil)
	c.Check(report["pass"], Equals, true)
}

func (s *AssertCaseSuite) TestCaseExecuteWithErrorRequest(c *C) {
	// given
	file, _ := os.ReadFile("./data/error-message-test.yaml")
	tc := &test_suite.Case{}

	// when
	_ = yaml.Unmarshal(file, tc)
	err := tc.Execute()

	// then
	c.Check(err, IsNil)

	// when
	report, err := test_report.GetReport().Gen("gen/testsuite_case_level_vars")

	// then
	c.Check(err, IsNil)

	// when
	file, err = os.ReadFile(report)

	// then
	c.Check(err, IsNil)
}

func (s *AssertCaseSuite) TestCaseExecuteWithBaseUrl(c *C) {
	// given
	suite, err := test_suite.BuildTestSuite("./data/testsuite_auto_fill_base_url")

	// then
	c.Check(err, IsNil)

	// when
	err = suite.Execute()

	// then
	c.Check(err, IsNil)

	// when
	_, err = test_report.GetReport().Gen("./gen/testsuite_auto_fill_base_url")

	// then
	c.Check(err, IsNil)

	// when
	file, err := os.ReadFile(fmt.Sprintf("gen/testsuite_auto_fill_base_url/%s.%s", common.ReportFile, common.JsonReportFileExt))

	// then
	c.Check(err, IsNil)

	// when
	report := make(map[string]any, 0)
	err = yaml.Unmarshal(file, &report)

	// then
	c.Check(err, IsNil)
	c.Check(report["pass"], Equals, true)
}
