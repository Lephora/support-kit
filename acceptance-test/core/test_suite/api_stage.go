package test_suite

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"thoughtworks.com/lephora/support-kit/accecptance-test/common"
	"thoughtworks.com/lephora/support-kit/accecptance-test/spec_validate"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_assertion"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_report"
)

type Request struct {
	Url      string            `yaml:"url,omitempty" json:"url,omitempty"`
	Method   string            `yaml:"method,omitempty" json:"method,omitempty"`
	Headers  map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Params   map[string]string `yaml:"params,omitempty" json:"params,omitempty"`
	Body     string            `yaml:"body,omitempty" json:"body,omitempty"`
	File     string            `yaml:"file,omitempty" json:"file,omitempty"`
	Protocol string            `yaml:"protocol,omitempty" json:"protocol,omitempty"`
}

type Actual struct {
	Status       *int              `yaml:"status,omitempty" json:"status,omitempty"`
	Headers      map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Body         *string           `yaml:"body,omitempty" json:"body,omitempty"`
	ErrorMessage string            `yaml:"error_message,omitempty" json:"error_message,omitempty"`
}

type Assertion struct {
	Status   int               `yaml:"status,omitempty" json:"status,omitempty"`
	Headers  map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Body     string            `yaml:"body,omitempty" json:"body,omitempty"`
	JsonPath JsonPathValidator `yaml:"json_path,omitempty" json:"json_path,omitempty"`
	Actual   *Actual           `yaml:"actual,omitempty" json:"actual,omitempty"`
}

type JsonPathValidator []struct {
	Path  string `json:"path,omitempty" yaml:"path,omitempty"`
	Regex string `json:"regex,omitempty" yaml:"regex,omitempty"`
}

func (a *Actual) Success(status int, headers map[string]string, body *string) *Actual {
	a.Status = &status
	a.Headers = headers
	a.Body = body
	return a
}
func (a *Actual) Fail(errorMessage string) *Actual {
	a.ErrorMessage = errorMessage
	return a
}
func (s *Stage) executeApi() error {

	s.Actual = &Actual{}

	renderBody := ""
	if s.getRenderRequest().Method == "GET" && s.getRenderRequest().Body != "" {
		fmt.Println("[WARNING]: GET Method does not need body")
	} else {
		renderBody = s.getRenderRequest().Body
	}

	request, err := http.NewRequest(s.getRenderRequest().Method, s.getRenderRequest().Url, strings.NewReader(renderBody))

	if err != nil {
		s.report(false, err.Error())
		return nil
	}
	request.Header = recoverHeaders(s.getRenderRequest().Headers)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		s.Actual.Fail(err.Error())
		if common.CurrentPhase() == common.Asserting {
			s.report(false, err.Error())
		}
		return nil
	}

	buffer := bytes.NewBuffer([]byte{})
	_, _ = resp.Body.Read(buffer.Bytes())
	newReader := bytes.NewReader(buffer.Bytes())
	readCloser := io.NopCloser(newReader)
	spec_validate.CollectFlow(request, &http.Response{
		StatusCode: resp.StatusCode,
		Header:     resp.Header,
		Body:       readCloser,
	})

	statusCode := resp.StatusCode
	body := resp.Body

	defer body.Close()

	s.Actual.Success(resp.StatusCode, transferHeaders(resp.Header), transferBody(resp.Body))

	if common.CurrentPhase() == common.Asserting && s.Assertion != nil {
		if s.Assertion.Status != 0 {
			assertor := test_assertion.GenericAssertor{}
			assert, err := assertor.Assert(statusCode, test_assertion.Equals, s.Assertion.Status)
			if err != nil || !assert {
				s.report(assert, err.Error())
			} else {
				s.Assertion.Actual = s.Actual
				s.report(assert, map[string]any{"request": s.getRenderRequest(), "assertion": s.Assertion})
			}
		}

		if s.Assertion.Body != "" {
			assertor := test_assertion.GenericAssertor{}
			assert, err := assertor.Assert(*s.Actual.Body, test_assertion.Equals, s.Assertion.Body)
			if err != nil || !assert {
				s.report(assert, err.Error())
			} else {
				s.Assertion.Actual = s.Actual
				s.report(assert, map[string]any{"request": s.getRenderRequest(), "assertion": s.Assertion})
			}
		}

		if s.Assertion.JsonPath != nil {
			for _, jsonPath := range s.Assertion.JsonPath {
				assertor := test_assertion.JsonAssertor{}
				assertor.Object(*s.Actual.Body)
				assert, err := assertor.Assert(jsonPath.Path, test_assertion.Match, jsonPath.Regex)
				if err != nil {
					s.report(assert, err.Error())
				} else {
					s.Assertion.Actual = s.Actual
					s.report(assert, map[string]any{"request": s.getRenderRequest(), "assertion": s.Assertion})
				}
			}
		}
	}
	return nil
}

func transferBody(body io.ReadCloser) *string {
	if body == nil {
		return nil
	}
	content, err := io.ReadAll(body)
	if err != nil {
		msg := fmt.Sprintf("transfer body error: %s", err.Error())
		return &msg
	}
	result := string(content)
	return &result
}

func transferHeaders(header http.Header) map[string]string {
	if header == nil {
		return nil
	}
	result := make(map[string]string, 0)
	for key, values := range header {
		result[key] = strings.Join(values, ", ")
	}
	return result
}

func recoverHeaders(header map[string]string) http.Header {
	if header == nil {
		return nil
	}
	result := make(map[string][]string, 0)
	for key, values := range header {
		result[key] = strings.Split(values, ", ")
	}
	return result
}

func (s *Stage) report(pass bool, detail any) {
	test_report.GetReport().AppendStage(s.Name, pass, detail)
}
