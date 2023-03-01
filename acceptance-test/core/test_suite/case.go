package test_suite

import (
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"thoughtworks.com/lephora/support-kit/accecptance-test/common"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_report"
	"thoughtworks.com/lephora/support-kit/accecptance-test/util"
	"time"
)

type Case struct {
	name   string
	Stages []*Stage          `yaml:"stages,omitempty"`
	Vars   map[string]string `yaml:"vars,omitempty"`

	option struct {
		relativePath string
		filters      []func(stage *Stage) bool
		variables    map[string]any
	}
}

func (c *Case) Name() string {
	return c.name
}

func (c *Case) SetVar(key string, value any) {
	c.setVar(key, value)
}

func (c *Case) setVar(key string, value any) {
	if c.option.variables == nil {
		c.option.variables = make(map[string]any, 0)
	}
	c.option.variables[key] = value
}

func (c *Case) Var() map[string]any {
	return c.option.variables
}

func (c *Case) AddFilter(filter func(stage *Stage) bool) *Case {
	if c.option.filters == nil {
		c.option.filters = make([]func(stage *Stage) bool, 0)
	}
	c.option.filters = append(c.option.filters, filter)
	return c
}

func (c *Case) Execute() error {
	switch common.CurrentPhase() {
	case common.Recording:

		if c.Stages != nil {
			if err := c.execute(); err != nil {
				return err
			}
		}
		return c.record()

	case common.Asserting:

		c.report()

		if c.Stages != nil {

			start := time.Now()

			if err := c.execute(); err != nil {
				return err
			}

			duration := time.Since(start)

			test_report.GetReport().SetCaseRunTime(duration)
		}
		return nil
	default:
		return nil
	}
}

func (c *Case) execute() error {
	for _, stage := range c.filter() {

		c.deliverCaseVariables(stage)

		stage.SetRelativePath(c.option.relativePath)

		if err := stage.Execute(); err != nil {
			return err
		}

		c.collectStageVariables(stage)

	}
	return nil
}

func (c *Case) report() {
	report := test_report.GetReport()
	report.AppendCase(c.name)
}

func (c *Case) filter() []*Stage {

	if c.option.filters == nil {
		return c.Stages
	}

	if c.Stages == nil {
		return []*Stage{}
	}

	var result []*Stage
	for _, stage := range c.Stages {
		for _, filter := range c.option.filters {
			if filter(stage) {
				result = append(result, stage)
				break
			}
		}
	}
	return result
}

func (c *Case) collectStageVariables(stage *Stage) {
	if stage.Type != API || stage.Actual == nil {
		return
	}

	c.setVar(stage.Name, map[string]any{"actual": map[string]any{"body": stage.Actual.Body, "status": stage.Actual.Status, "headers": stage.Actual.Headers}})
}

func (c *Case) deliverCaseVariables(s *Stage) {
	if c.option.variables != nil {
		for key, value := range c.option.variables {
			s.SetVar(key, value)
		}
	}

	if c.Vars != nil {
		for key, value := range c.Vars {
			s.SetVar(key, value)
		}
	}
}

func (c *Case) record() error {
	caseData, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(c.option.relativePath, caseData, 0766)
}

func (c *Case) SetName(name string) {
	c.name = filepath.Base(util.GetFileNameWithoutExtension(name))
	c.option.relativePath = name
}
