package unit_test

import (
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_suite"

	. "gopkg.in/check.v1"
)

type StageSuite struct{}

var _ = Suite(&StageSuite{})

func (s *StageSuite) TestCaseExecuteWithRender(c *C) {
	// given
	stage := test_suite.Stage{
		Name: "test_stage",
		Request: test_suite.Request{
			Url:     "{{ .other_stage.response.url }}",
			Headers: map[string]string{"content-type": "{{ .Header }}"},
			Body:    "other_stage body is {{ .other_stage.response.body }}",
		},
	}

	stage.SetVar("other_stage", map[string]map[string]string{"response": {"url": "https://www.baidu.com", "body": `{"name":"unbelievable"}`}})
	stage.SetVar("Header", "application/xml")

	// when
	err := stage.Execute()

	// then
	c.Check(err, IsNil)
	c.Check(stage.Name, Equals, "test_stage")
	c.Check(stage.GetRequest().Url, Equals, "https://www.baidu.com")
	c.Check(stage.GetRequest().Body, Equals, `other_stage body is {"name":"unbelievable"}`)
	c.Check(stage.GetRequest().Headers["content-type"], Equals, "application/xml")

}
