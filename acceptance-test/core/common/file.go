package common

import "os"

type TestrunnerPhase string

const (
	PhaseEnv                 = "TESTRUNNER_PHASE"
	DescriptionFileName      = "description"
	Recording                = TestrunnerPhase("record")
	Asserting                = TestrunnerPhase("assert")
	ReportFile               = "testrunner-report"
	JsonReportFileExt        = "json"
	XmlReportFileExt         = "xml"
	ENV_REPORT_TYPE          = "REPORT_TYPE"
	ENV_REPORT_NAME          = "REPORT_NAME"
	ENV_DEFAULT_HOST         = "DEFAULT_HOST"
	SUPPORT_REPORT_TYPE_XML  = "XML"
	SUPPORT_REPORT_TYPE_JSON = "JSON"
)

var (
	SupportYamlExt = []string{"yaml", "yml"}
)

func CurrentPhase() TestrunnerPhase {
	phase := os.Getenv(PhaseEnv)
	if phase == "" || TestrunnerPhase(phase) == Asserting {
		return Asserting
	}
	if TestrunnerPhase(phase) == Recording {
		return Recording
	}
	return ""
}
