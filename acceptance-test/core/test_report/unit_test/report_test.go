package unit_test

import (
	"fmt"
	"github.com/joshdk/go-junit"
	. "gopkg.in/check.v1"
	"os"
	"testing"
	"thoughtworks.com/lephora/support-kit/accecptance-test/common"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_report"
)

func Test(t *testing.T) { TestingT(t) }

type CaseSuite struct{}

var _ = Suite(&CaseSuite{})

func (s *CaseSuite) SetUpTest(c *C) {
	_ = os.Setenv("REPORT_STYLE", "testrunner-report.json")
}

func (s *CaseSuite) TestReportGen(c *C) {
	// given
	report := test_report.Report{}

	// when
	gen, err := report.Gen("./gen")

	// then
	c.Check(err, IsNil)
	c.Check(gen, Equals, fmt.Sprintf("gen/%s.%s", common.ReportFile, common.JsonReportFileExt))
}

func (s *CaseSuite) TestReportAppendStageAndGen(c *C) {
	// given
	report := test_report.Report{}

	// when
	report.AppendSuite("first_suite")
	report.AppendCase("first_case")
	report.AppendStage("first_stage", true, "I'm pass")
	gen, err := report.Gen("./gen/success")

	// then
	c.Check(err, IsNil)
	c.Check(gen, Equals, fmt.Sprintf("gen/success/%s.%s", common.ReportFile, common.JsonReportFileExt))
	c.Check(report.Pass, Equals, true)

	// when
	report.AppendStage("second_stage", false, "I was failed")
	c.Check(report.Pass, Equals, true)
	gen, err = report.Gen("./gen/failed")

	// then
	c.Check(err, IsNil)
	c.Check(gen, Equals, fmt.Sprintf("gen/failed/%s.%s", common.ReportFile, common.JsonReportFileExt))
	c.Check(report.Pass, Equals, false)
}

func (s *CaseSuite) TestReportGenerateXML(c *C) {
	// given
	os.Setenv("REPORT_TYPE", "XML")
	os.Setenv("REPORT_NAME", "testrunner-report.xml")
	report := test_report.Report{}

	// when
	report.AppendSuite("first_suite")
	report.AppendCase("first_case")
	report.AppendStage("first_stage", true, "I'm pass")
	gen, err := report.Gen("./gen/success")

	// then
	c.Check(err, IsNil)
	c.Check(gen, Equals, fmt.Sprintf("gen/success/%s.%s", common.ReportFile, common.XmlReportFileExt))
	c.Check(report.Pass, Equals, true)

	// when
	suites, err := junit.IngestFile(gen)
	c.Check(err, IsNil)
	c.Check(suites, HasLen, 1)
	c.Check(suites[0].Tests, HasLen, 1)
}
