package test_runner

import (
	"github.com/pkg/errors"
	"io/fs"
	"path/filepath"
	"strings"
	"thoughtworks.com/lephora/support-kit/accecptance-test/common"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_report"
	"thoughtworks.com/lephora/support-kit/accecptance-test/test_suite"
	"thoughtworks.com/lephora/support-kit/accecptance-test/util"
)

type Testrunner struct {
}

func (tr *Testrunner) Run(option test_suite.Option) error {
	test_suite.SetOption(option)
	if common.CurrentPhase() == common.Asserting {
		test_report.BuildReport()
	}
	_, ok := util.FileExistWithExtensionName(option.GetFullPath(), common.SupportYamlExt...)
	if option.HasCase() && !ok {
		return errors.New("[Test_Runner] failed to find testcases " + option.CaseName + " in " + option.GetDir())
	}

	_, ok = util.FileExistWithExtensionName(filepath.Join(option.GetDir(), common.DescriptionFileName), common.SupportYamlExt...)
	if !ok {
		if len(option.Suite) != 0 {
			return errors.New("[Test_Runner] failed to recognize testsuite in " + option.GetDir() + ", no description.yaml be found")
		}
	}

	return filepath.WalkDir(option.GetDir(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.New("[Test_Runner] failed to find testsuite in " + option.GetDir())
		}
		if !d.IsDir() {
			return nil
		}
		if !test_suite.IsTestSuite(path) {
			return nil
		}
		suite, err := test_suite.BuildTestSuite(path)
		if err != nil {
			return err
		}
		return suite.Execute()
	})
}

func (tr *Testrunner) RunByRoot(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.New("[Test_Runner] failed to find testsuite in " + root)
		}
		if !d.IsDir() {
			return nil
		}
		if !test_suite.IsTestSuite(path) {
			return nil
		}
		suite, err := test_suite.BuildTestSuite(path)
		if err != nil {
			return err
		}
		return suite.Execute()
	})
}

func (tr *Testrunner) RunBySuite(path string) error {
	_, ok := util.FileExistWithExtensionName(filepath.Join(path, common.DescriptionFileName), common.SupportYamlExt...)
	if !ok {
		return errors.New("[Test_Runner] failed to recognize testsuite in " + path + ", no description.yaml be found")
	}

	if !test_suite.IsTestSuite(path) {
		return nil
	}
	suite, err := test_suite.BuildTestSuite(path)
	if err != nil {
		return err
	}
	return suite.Execute()
}

func (tr *Testrunner) RunByCase(name string) error {
	_, ok := util.FileExistWithExtensionName(name, common.SupportYamlExt...)
	if !ok {
		return errors.New("[Test_Runner] failed to find testcases " + name)
	}

	index := strings.LastIndex(name, "/")
	path := name[:index]
	caseName := name[index+1:]

	if !test_suite.IsTestSuite(path) {
		return nil
	}
	builtSuite, err := test_suite.BuildTestSuite(path)
	if err != nil {
		return err
	}
	builtSuite.AddFilter(func(caseObj *test_suite.Case) bool {
		return caseName == caseObj.Name()
	})
	return builtSuite.Execute()
}
