package projectvalidator

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type ValidationConfiguration struct {
	ScenarioFile string `json:"scenario_file" yaml:"scenario_file"`
	TestsDirectory string `json:"tests_directory" yaml:"tests_directory"`
}

type Specification struct {
	Project string `json:"project,omitempty" yaml:"project"`
	Scope string `json:"scope" yaml:"scope"`
	Stories []Story `json:"stories,omitempty" yaml:"stories"`
	//Markdown is a slice of markdowns to render
	MarkDown []string `json:"markdown" yaml:"markdown"`
}

type Story struct {
	Name string `json:"name,omitempty" yaml:"name"`
	//Tags indicate strings that will join tests to this Story
	Tags []string `json:"tags,omitempty" yaml:"tags"`
	//Risk as a value of Low, High, or Medium
	Risk string `json:"risk,omitempty" yaml:"risk"`
	Summary TestSummary `json:"test_summary,omitempty" yaml:"test_summary"`
	Tests []TestResult `json:"tests,omitempty" yaml:"tests"`
	//Markdown is a slice of markdowns to render
	MarkDown []string `json:"markdown" yaml:"markdown"`
}


type TestResult struct {
	Name	string `json:"name,omitempty" yaml:"name"`
	Tags	[]string `json:"tags,omitempty" yaml:"tags"`
	Passed	bool	`json:"passed,omitempty" yaml:"passed"`
	//Date is the string representation (ISO) for date of test execution
	Date	string	`json:"date,omitempty" yaml:"date"`
}

//{"Time":"2020-02-11T14:03:49.53841775-05:00","Action":"output","Package":"github.com/metrumresearchgroup/babylontest","Test":"TestSpecifiedConfigByAbsPathLoaded","Output":"--- PASS: TestSpecifiedConfigByAbsPathLoaded (6.23s)\n"}
//GoTestResult is the struct emitted by go test in JSON form
type GoTestResult struct {
	Time string
	Action string
	Package string
	Test string
	Output string
}

type TestSummary struct {
	TestCount int `json:"test_count,omitempty" yaml:"test_count"`
	PassedTests int `json:"passed_tests,omitempty" yaml:"passed_tests"`
	FailedTests int `json:"failed_tests,omitempty" yaml:"failed_tests"`
}


//LoadValidationScenarioFromFile reads the
func NewSpecification(file io.Reader) (*Specification,error) {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return &Specification{}, err
	}

	v := Specification{}

	err = json.Unmarshal(bytes,&v)

	if err != nil {
		return &Specification{}, err
	}

	return &v, nil
}