package test_report

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"strconv"
	"thoughtworks.com/lephora/support-kit/accecptance-test/global_config"
	"thoughtworks.com/lephora/support-kit/accecptance-test/util"
	"time"
)

var report Report

func GetReport() *Report {
	return &report
}

func BuildReport() *Report {
	report = Report{}
	return &report
}

type stageReport struct {
	Pass   bool   `json:"pass"`
	Name   string `json:"name"`
	Detail any    `json:"detail,omitempty"`
}

type caseReport struct {
	Pass    bool           `json:"pass"`
	Name    string         `json:"name"`
	Stages  []*stageReport `json:"stages,omitempty"`
	RunTime time.Duration  `json:"run-time,omitempty"`
}

type suiteReport struct {
	Pass  bool          `json:"pass"`
	Name  string        `json:"name"`
	Cases []*caseReport `json:"cases,omitempty"`
}

type Report struct {
	Pass      bool           `json:"pass"`
	Suites    []*suiteReport `json:"suites,omitempty"`
	Statistic Statistic      `json:"statistic,omitempty"`
}

type Statistic struct {
	Pass      bool          `json:"pass"`
	Suite     statistic     `json:"suites"`
	Case      statistic     `json:"cases"`
	TotalTime time.Duration `json:"total-time,omitempty"`
}

type statistic struct {
	Total   int `json:"total"`
	Success int `json:"success"`
	Failed  int `json:"failed"`
}

func (r *Report) AppendSuite(name string) {
	if r.Suites == nil {
		r.Suites = make([]*suiteReport, 0)
	}
	r.Suites = append(r.Suites, &suiteReport{Name: name})
}

func (r *Report) getLastSuite() *suiteReport {
	return r.Suites[len(r.Suites)-1]
}

func (r *Report) AppendCase(name string) {
	suite := r.getLastSuite()
	if r.getLastSuite().Cases == nil {
		suite.Cases = make([]*caseReport, 0)
	}
	suite.Cases = append(suite.Cases, &caseReport{Name: name})
}
func (r *Report) SetCaseRunTime(duration time.Duration) {
	lastCase := r.getLastCase()
	lastCase.RunTime = duration
}
func (r *Report) getLastCase() *caseReport {
	suite := r.getLastSuite()
	return suite.Cases[len(suite.Cases)-1]
}

func (r *Report) AppendStage(name string, pass bool, detail any) {
	lastCase := r.getLastCase()
	if lastCase.Stages == nil {
		lastCase.Stages = make([]*stageReport, 0)
	}
	lastCase.Stages = append(lastCase.Stages, &stageReport{Pass: pass, Name: name, Detail: detail})
}

func (cr *caseReport) integrate() {
	cr.Pass = true
	for _, stage := range cr.Stages {
		if !stage.Pass {
			fmt.Printf("Fail: %s: %s", stage.Name, stage.Detail)
			cr.Pass = false
			break
		}
	}
}

func (sr *suiteReport) integrate() {
	sr.Pass = true
	for _, testcase := range sr.Cases {
		testcase.integrate()
		if !testcase.Pass {
			fmt.Printf(" in case %s", testcase.Name)
			sr.Pass = false
		}
	}
}

func (r *Report) integrate() {
	r.Pass = true
	for _, suite := range r.Suites {
		suite.integrate()
		if !suite.Pass {
			fmt.Printf(" in suite %s\n", suite.Name)
			r.Pass = false
		}
	}
}

func (r *Report) statistic() {
	r.Statistic.Pass = true
	for _, suite := range r.Suites {
		suitePass := true
		for _, testcase := range suite.Cases {
			if !testcase.Pass {
				suitePass = false
				r.Statistic.Case.Failed += 1
			} else {
				r.Statistic.Case.Success += 1
			}
			r.Statistic.Case.Total += 1
		}
		if !suitePass {
			r.Statistic.Pass = false
			r.Statistic.Suite.Failed += 1
		} else {
			r.Statistic.Suite.Success += 1
		}
		r.Statistic.Suite.Total += 1
	}
}

func (r *Report) Gen(path string) (string, error) {
	if !util.FileOrDirExist(path) {
		if err := os.MkdirAll(path, 0766); err != nil {
			return "", errors.Wrap(err, "failed to generate report")
		}
	}

	r.integrate()
	r.statistic()

	if global_config.GetConfiguration().ReportType() == global_config.XmlReport {
		xmlReport := r.XML()
		marshal, err := xml.Marshal(xmlReport)
		if err != nil {
			return "", errors.Wrap(err, "failed to marshal report")
		}
		reportPath := filepath.Join(path, fmt.Sprintf(global_config.GetConfiguration().ReportName()))
		if err := os.WriteFile(reportPath, marshal, 0766); err != nil {
			return "", errors.Wrap(err, "failed to write report")
		}
		return reportPath, nil
	}

	marshal, err := json.Marshal(r)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal report")
	}
	reportPath := filepath.Join(path, fmt.Sprintf(global_config.GetConfiguration().ReportName()))
	if err := os.WriteFile(reportPath, marshal, 0766); err != nil {
		return "", errors.Wrap(err, "failed to write report")
	}

	return reportPath, nil
}

func (r *Report) SetTotalRunTime(duration time.Duration) {
	r.Statistic.TotalTime = duration
}

func (r *Report) Print() {
	data := []string{strconv.Itoa(r.Statistic.Suite.Total), strconv.Itoa(r.Statistic.Case.Total), strconv.Itoa(r.Statistic.Case.Success),
		strconv.Itoa(r.Statistic.Case.Failed), strconv.FormatBool(r.Statistic.Pass), r.Statistic.TotalTime.String()}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"suites", "cases", "success cases", "failed cases", "pass", "run time"})

	table.Append(data)
	table.Render()
}

func (r *Report) XML() *xmlReport {
	xr := &xmlReport{Suites: make([]*XmlSuite, 0)}
	for _, jsonSuite := range r.Suites {
		xmlSuite := &XmlSuite{Name: jsonSuite.Name, Testcase: make([]*XmlCase, 0)}
		var time time.Duration
		for _, jsonCase := range jsonSuite.Cases {
			xmlSuite.Total++
			xmlCase := &XmlCase{
				Name: jsonCase.Name,
				Time: jsonCase.RunTime.String(),
			}
			var text string
			for _, jsonStage := range jsonCase.Stages {
				marshal, _ := json.Marshal(jsonStage.Detail)
				text = fmt.Sprintf("%s\n==========\nstage name: %s \npass: %v\ndetail: %s\n", text, jsonStage.Name, jsonStage.Pass, string(marshal))
			}
			if jsonCase.Pass {
				xmlCase.Log = &Message{Text: text}
			} else {
				xmlSuite.Failures++
				xmlCase.Failure = &Message{Text: text}
			}
			xmlSuite.Testcase = append(xmlSuite.Testcase, xmlCase)
			time += jsonCase.RunTime
		}
		xmlSuite.Time = time.String()
		xr.Suites = append(xr.Suites, xmlSuite)
	}
	return xr
}
