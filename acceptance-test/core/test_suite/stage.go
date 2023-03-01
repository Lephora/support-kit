package test_suite

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"html/template"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"thoughtworks.com/lephora/support-kit/accecptance-test/global_config"
	"thoughtworks.com/lephora/support-kit/accecptance-test/util"
)

type Stage struct {
	Type      StageType  `yaml:"type,omitempty"`
	Name      string     `yaml:"name,omitempty"`
	Request   Request    `yaml:"request,omitempty"`
	Actual    *Actual    `yaml:"actual,omitempty"`
	Assertion *Assertion `yaml:"assert,omitempty"`
	option    struct {
		relativePath string
		renderStage  *Stage
		variables    map[string]any
	}
}

func (s *Stage) SetRelativePath(path string) {
	s.option.relativePath = path
}

func (s *Stage) RelativePath() string {
	return s.option.relativePath
}

func (s *Stage) SetVar(key string, value any) {
	if s.option.variables == nil {
		s.option.variables = make(map[string]any, 0)
	}
	s.option.variables[key] = value
}

func (s *Stage) Var() map[string]any {
	return s.option.variables
}

func (s *Stage) GetRequest() Request {
	if s.option.renderStage != nil {
		return s.getRenderRequest()
	}
	return s.Request
}

func (s *Stage) getRenderRequest() Request {
	return s.option.renderStage.Request
}

func (s *Stage) Execute() error {
	s.render()

	if s.Request.File != "" {
		err := s.readRequestBody(s.calculateBodyFilePath())
		if err != nil {
			return err
		}
	}

	if s.Type == API {
		s.autoFillUrl()
		s.autoAddProtocol()
		return s.executeApi()
	}
	return nil
}

func (s *Stage) render() {
	marshal, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}

	temp, err := template.New(s.Name).Parse(string(marshal))
	if err != nil {
		panic(err)
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := temp.Execute(buffer, s.option.variables); err != nil {
		panic(err)
	}

	if err := json.Unmarshal([]byte(escapeHTMLChar(buffer.String())), &s.option.renderStage); err != nil {
		panic(err)
	}
}

func escapeHTMLChar(o string) string {
	return strings.Replace(o, "&#34;", "\\\"", -1)
}

func (s *Stage) autoFillUrl() {
	u, err := url.Parse(s.option.renderStage.Request.Url)
	if err != nil || u.Scheme == "" || u.Host == "" {
		host := global_config.GetConfiguration().Host()
		s.option.renderStage.Request.Url = host + s.option.renderStage.Request.Url
	}
}

func (s *Stage) autoAddProtocol() {
	if s.option.renderStage.Request.Protocol != "" {
		protocol := s.option.renderStage.Request.Protocol
		s.option.renderStage.Request.Url = strings.Replace(s.option.renderStage.Request.Url, "https", protocol, 1)
		s.option.renderStage.Request.Url = strings.Replace(s.option.renderStage.Request.Url, "http", protocol, 1)
	}
}

func (s *Stage) readRequestBody(realDir string) error {
	ok := util.FileOrDirExist(realDir)
	if !ok {
		return errors.New("[Test_Runner] body file not exist: " + realDir)
	}

	readFile, err := os.ReadFile(realDir)
	if err != nil {
		return errors.New("[Test_Runner] cannot read body file")
	}

	s.option.renderStage.Request.Body = string(readFile)

	return nil
}

func (s *Stage) calculateBodyFilePath() string {

	path := s.RelativePath()

	dir := filepath.Dir(path)

	return filepath.Join(dir, s.Request.File)

}
