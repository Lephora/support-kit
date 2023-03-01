package global_config

import (
	"fmt"
	"os"
	"thoughtworks.com/lephora/support-kit/accecptance-test/common"
	"time"
)

type Configuration interface {
	Type() ConfigurationType
	Host() string
	ReportName() string
	ReportType() ReportType
}

const env = ConfigurationType("env")
const XmlReport = ReportType("XML")
const JsonReport = ReportType("JSON")

var FixJsonReportName = fmt.Sprintf("%s.%s", common.ReportFile, common.JsonReportFileExt)
var DefaultJsonReportName = fmt.Sprintf("%s-%s.%s", common.ReportFile, time.Now().Format(time.RFC3339), common.JsonReportFileExt)
var DefaultXmlReportName = fmt.Sprintf("%s-%s.%s", common.ReportFile, time.Now().Format(time.RFC3339), common.XmlReportFileExt)

type ConfigurationType string
type ReportType string

var _ Configuration = (*EnvConfiguration)(nil)

type EnvConfiguration struct {
}

func (e EnvConfiguration) ReportType() ReportType {
	env := os.Getenv(common.ENV_REPORT_TYPE)
	if env == "" {
		return JsonReport
	}
	return ReportType(env)
}

func (e EnvConfiguration) Type() ConfigurationType {
	return env
}

func (e EnvConfiguration) Host() string {
	env := os.Getenv(common.ENV_DEFAULT_HOST)
	if env == "" {
		return "http://127.0.0.1:8884"
	}
	return env
}

func (e EnvConfiguration) ReportName() string {
	env := os.Getenv(common.ENV_REPORT_NAME)
	if env != "" {
		return env
	}
	reportType := os.Getenv(common.ENV_REPORT_TYPE)
	switch reportType {
	case common.SUPPORT_REPORT_TYPE_XML:
		return DefaultXmlReportName
	case common.SUPPORT_REPORT_TYPE_JSON:
		return DefaultJsonReportName
	default:
		return FixJsonReportName
	}
}

func GetConfiguration() Configuration {
	return EnvConfiguration{}
}
