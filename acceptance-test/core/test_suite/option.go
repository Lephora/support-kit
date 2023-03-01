package test_suite

type Option struct {
	RootDir  string
	Suite    string
	CaseName string
}

var args Option

func GetOption() Option {
	return args
}

func SetOption(option Option) {
	args = option
}
func (o *Option) GetDir() string {
	if len(o.Suite) == 0 {
		return o.RootDir
	} else {
		return o.RootDir + "/" + o.Suite
	}
}

func (o *Option) HasCase() bool {
	return len(o.CaseName) != 0
}

func (o *Option) GetFullPath() string {
	if len(o.CaseName) == 0 {
		return o.GetDir()
	} else {
		return o.GetDir() + "/" + o.CaseName
	}
}
