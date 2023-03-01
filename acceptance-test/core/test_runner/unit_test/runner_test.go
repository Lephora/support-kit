package unit_test

import (
	"encoding/json"
	"github.com/joshdk/go-junit"
	"gopkg.in/yaml.v2"
	"os"
	"testing"
	"thoughtworks.com/lephora/support-kit/accecptance-test/common"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_report"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_runner"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_suite"
	"time"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type RunnerSuite struct{}

var _ = Suite(&RunnerSuite{})

func (s *RunnerSuite) SetUpTest(c *C) {
	_ = os.Setenv(common.PhaseEnv, string(common.Asserting))
	_ = os.Setenv("REPORT_STYLE", "testrunner-report.json")
}

//func (s *RunnerSuite) TestRunnerCorrectCheck(c *C) {
//	// given
//	testrunner := test_runner.Testrunner{}
//
//	// when
//	description, err := testrunner.CheckDescription("./data/testsuite_correct_import")
//
//	// then
//	c.Check(err, IsNil)
//	c.Check(description.Import, HasLen, 1)
//	c.Check(description.Import[0], Equals, "hello-test")
//}
//
//func (s *RunnerSuite) TestRunnerErrorCheck(c *C) {
//	// given
//	testrunner := test_runner.Testrunner{}
//
//	// when
//	_, err := testrunner.CheckDescription("./data/testsuite_nil_description")
//
//	// then
//	c.Check(err, NotNil)
//	c.Assert(err.Error(), Equals, "you must import your testcase into description")
//}

func (s *RunnerSuite) TestRunnerRun(c *C) {
	// given
	testrunner := test_runner.Testrunner{}

	// when
	err := testrunner.RunByRoot("./data/testsuite_correct_import")

	// then
	c.Check(err, IsNil)
}

func (s *RunnerSuite) TestRunnerExec(c *C) {
	// given
	testrunner := test_runner.Testrunner{}

	// when
	err := testrunner.RunByRoot("./data/testsuite_correct_baidu_request")

	// then
	c.Check(err, IsNil)

	// when
	_, err = test_report.GetReport().Gen("./gen")

	// then
	c.Check(err, IsNil)
}

func (s *RunnerSuite) TestRunnerReportWithCorrectRequest(c *C) {
	// given
	testrunner := test_runner.Testrunner{}
	test_report.BuildReport()

	// when
	err := testrunner.RunByRoot("./data/testsuite_correct_baidu_request")

	// then
	c.Check(err, IsNil)

	// when
	_, err = test_report.GetReport().Gen("./gen")

	// then
	file, _ := os.ReadFile("./gen/testrunner-report.json")
	report := test_report.Report{}
	_ = json.Unmarshal(file, &report)
	c.Check(report, NotNil)
	c.Check(report.Pass, Equals, true)
	c.Check(report.Suites, HasLen, 1)
	c.Check(report.Suites[0].Pass, Equals, true)
	c.Check(report.Suites[0].Cases, HasLen, 2)
	c.Check(report.Suites[0].Cases[0].RunTime, NotNil)
	c.Check(report.Suites[0].Cases[0].Pass, Equals, true)
	c.Check(report.Suites[0].Cases[0].Stages, HasLen, 2)

}

func (s *RunnerSuite) TestRunnerReportWithIncorrectRequest(c *C) {
	// given
	testrunner := test_runner.Testrunner{}
	test_report.BuildReport()

	// when
	err := testrunner.RunByRoot("./data/testsuite_incorrect_baidu_request")

	// then
	c.Check(err, IsNil)

	// when
	_, err = test_report.GetReport().Gen("./gen")

	// then
	file, _ := os.ReadFile("./gen/testrunner-report.json")
	report := test_report.Report{}
	_ = json.Unmarshal([]byte(file), &report)
	c.Check(report, NotNil)
	c.Check(report.Pass, Equals, false)
	c.Check(report.Suites, HasLen, 1)
	c.Check(report.Suites[0].Pass, Equals, false)
	c.Check(report.Suites[0].Cases, HasLen, 2)
	c.Check(report.Suites[0].Cases[0].RunTime, NotNil)
	c.Check(report.Suites[0].Cases[0].Pass, Equals, false)
	c.Check(report.Suites[0].Cases[0].Stages, HasLen, 2)

}

func (s *RunnerSuite) TestRunnerRewriteInputFileWithErrorMessage(c *C) {
	// given
	_ = os.Setenv(common.PhaseEnv, string(common.Recording))
	testrunner := test_runner.Testrunner{}

	// when
	err := testrunner.RunByRoot("./data/testsuite_incorrect_request")

	// then
	c.Check(err, IsNil)

	// then
	file, err := os.ReadFile("./data/testsuite_incorrect_request/cite-first-call-result.yaml")

	testcase := test_suite.Case{}
	err = yaml.Unmarshal(file, &testcase)
	c.Check(testcase, NotNil)
	c.Check(testcase.Stages[2].Actual.ErrorMessage, NotNil)
}

//todo: finish the unit-test case
//func (s *RunnerSuite) TestRunnerReportWithNilAssertion(c *C) {
//	// given
//	testrunner := test_runner.Testrunner{}
//
//	// when
//	err := testrunner.Run("./data/testsuite_nil_assertion")
//
//	// then
//	c.Check(err, IsNil)
//
//	// when
//	_, err = test_report.GetReport().Gen("./gen")
//
//	// then
//	file, _ := os.ReadFile("./gen/testrunner-report.json")
//	report := test_report.Report{}
//	_ = json.Unmarshal([]byte(file), &report)
//	c.Check(report, NotNil)
//	c.Check(report.Pass, Equals, false)
//	c.Check(report.Suites, HasLen, 1)
//	c.Check(report.Suites[0].Pass, Equals, false)
//	c.Check(report.Suites[0].Cases, HasLen, 2)
//	c.Check(report.Suites[0].Cases[0].RunTime, NotNil)
//	c.Check(report.Suites[0].Cases[0].Pass, Equals, false)
//	c.Check(report.Suites[0].Cases[0].Stages, HasLen, 2)
//
//}

func (s *RunnerSuite) TestRunnerWithMultiLevelVars(c *C) {
	// given
	testrunner := test_runner.Testrunner{}
	// when
	err := testrunner.RunByRoot("./data/testsuite_multi_level_vars/only_testcase")
	// then
	c.Check(err, IsNil)
	// when
	_, err = test_report.GetReport().Gen("./gen/multi_level_vars/only_testcase")
	// then
	file, _ := os.ReadFile("./gen/testrunner-report.json")
	report := make(map[string]any, 0)
	_ = json.Unmarshal(file, &report)
	c.Check(report, NotNil)

	// given
	err = testrunner.RunByRoot("./data/testsuite_multi_level_vars/only_description")
	// then
	c.Check(err, IsNil)
	// when
	_, err = test_report.GetReport().Gen("./gen/multi_level_vars/only_description")
	// then
	file, _ = os.ReadFile("./gen/testrunner-report.json")
	report = make(map[string]any, 0)
	_ = json.Unmarshal(file, &report)
	c.Check(report, NotNil)

	// given
	err = testrunner.RunByRoot("./data/testsuite_multi_level_vars/both_desc_testcase")
	// then
	c.Check(err, IsNil)
	// when
	_, err = test_report.GetReport().Gen("./gen/multi_level_vars/both_desc_testcase")
	// then
	file, _ = os.ReadFile("./gen/testrunner-report.json")
	report = make(map[string]any, 0)
	_ = json.Unmarshal(file, &report)
	c.Check(report, NotNil)
}

func (s *RunnerSuite) TestRunnerWithRecordMode(c *C) {
	// given
	_ = os.Setenv(common.PhaseEnv, string(common.Recording))
	testrunner := test_runner.Testrunner{}
	// when
	err := testrunner.RunByRoot("./data/testsuite_ready_for_record")
	// then
	c.Check(err, IsNil)
}
func (s *RunnerSuite) TestRunnerWithWrongSuitePath(c *C) {
	// given
	_ = os.Setenv(common.PhaseEnv, string(common.Recording))
	testrunner := test_runner.Testrunner{}
	// when
	err := testrunner.Run(test_suite.Option{RootDir: "./data/testsuite_correct_import", Suite: "/wrongtestsuitepath"})
	// then
	c.Check(err, NotNil)
}

func (s *RunnerSuite) TestRunnerWithAutoFillBaseUrl(c *C) {
	// given
	testrunner := test_runner.Testrunner{}

	// when
	err := testrunner.RunByRoot("./data/testsuite_auto_fill_base_url")

	// then
	c.Check(err, IsNil)
}
func (s *RunnerSuite) TestRunnerWithThreeArgs(c *C) {
	// given
	testrunner := test_runner.Testrunner{}

	// when
	err := testrunner.Run(test_suite.Option{RootDir: "./data", Suite: "testsuite_with_three_args", CaseName: "first_case"})

	// then
	c.Check(err, IsNil)
}
func (s *RunnerSuite) TestRunnerThrowErrorNilDescription(c *C) {
	// given
	testrunner := test_runner.Testrunner{}

	// when
	err := testrunner.Run(test_suite.Option{RootDir: "./data", Suite: "testsuite_nil_description"})

	// then
	c.Check(err, NotNil)
}
func (s *RunnerSuite) TestRunnerWithTwoArgs(c *C) {
	// given
	testrunner := test_runner.Testrunner{}

	// when
	err := testrunner.Run(test_suite.Option{RootDir: "./data", Suite: "testsuite_with_three_args"})

	// then
	c.Check(err, IsNil)
}

func (s *RunnerSuite) TestRunnerWithRequestBodyFromFile(c *C) {
	// given
	testrunner := test_runner.Testrunner{}

	// when
	err := testrunner.RunByRoot("./data/testsuite_get_request_body_from_file")

	// then
	c.Check(err, IsNil)
}

//func (s *RunnerSuite) TestRunnerWithAutoAddProtocol(c *C) {
//	// given
//	testrunner := test_runner.Testrunner{}
//
//	// when
//	err := testrunner.RunByRoot("./data/testsuite_auto_add_protocol")
//
//	// then
//	c.Check(err, IsNil)
//}

func (s *RunnerSuite) TestRunnerWithXMLReport(c *C) {
	// given
	// => runForTest()
	os.Setenv("REPORT_TYPE", "XML")
	os.Setenv("REPORT_NAME", "for-unit-test.xml")

	// when
	err := runForTest("./data/testsuite_with_xml_report")

	// then
	c.Check(err, IsNil)

	// when
	result, err := junit.IngestFile("./gen/for-unit-test.xml")
	c.Check(err, IsNil)
	c.Check(result, HasLen, 3)
	c.Check(result[0].Tests, HasLen, 2)

}

func runForTest(root string) error {
	_ = os.Setenv(common.PhaseEnv, string(common.Asserting))
	testrunner := test_runner.Testrunner{}

	if common.CurrentPhase() == common.Asserting {
		test_report.BuildReport()
	}
	start := time.Now()

	err := testrunner.RunByRoot(root)
	if err != nil {
		return err
	}

	if common.CurrentPhase() == common.Asserting {
		duration := time.Since(start)
		test_report.GetReport().SetTotalRunTime(duration)
		if _, err := test_report.GetReport().Gen("./gen"); err != nil {
			return err
		}
		test_report.GetReport().Print()
	}
	return nil
}
