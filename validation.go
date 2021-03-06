package goProjectValidator

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ValidationConfiguration struct {
	ScenarioFile string `json:"scenario_file" yaml:"scenario_file"`
	CommitsFile string `json:"commits_file" yaml:"commits_file"`
	TestsDirectory string `json:"tests_directory" yaml:"tests_directory"`
	OutputDirectory string `json:"output_directory" yaml:"output_directory"`
}

type Specification struct {
	Project string `json:"project,omitempty" yaml:"project"`
	Release string `json:"release,omitempty" yaml:"release"`
	Commits []CommitInfo `json:"commits,omitempty" yaml:"commits"`
	Scope string `json:"scope" yaml:"scope"`
	Stories []*Story `json:"stories,omitempty" yaml:"stories"`
	//Markdown is a slice of markdowns to render
	MarkDown []*Markdown `json:"markdown" yaml:"markdown"`

}

type CommitInfo struct {
	Repo string `json:"repo" yaml:"repo"`
	Commit string `json:"commit" yaml:"commit"`
}

type Story struct {
	Name string `json:"name,omitempty" yaml:"name"`
	//Tags indicate strings that will join tests to this Story
	Tags []string `json:"tags,omitempty" yaml:"tags"`
	//Risk as a value of Low, High, or Medium
	Risk string `json:"risk,omitempty" yaml:"risk"`
	Summary *TestSummary `json:"test_summary,omitempty" yaml:"test_summary"`
	Tests []*GoTestResult `json:"tests,omitempty" yaml:"tests"`
	//Markdown is a slice of markdowns to render
	MarkDown []*Markdown `json:"markdown" yaml:"markdown"`
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
	Passed bool
}

type TestSummary struct {
	TestCount int `json:"test_count,omitempty" yaml:"test_count"`
	PassedTests int `json:"passed_tests,omitempty" yaml:"passed_tests"`
	FailedTests int `json:"failed_tests,omitempty" yaml:"failed_tests"`
}

type Markdown struct {
	Source string `json:"source" yaml:"source"`
	Content string `json:"content" yaml:"content"`
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

//AddCommitInfo reads the commitsFile and appends to spec
func AddCommitInfo(spec *Specification, commits_file io.Reader) (*Specification,error) {
	bytes, err := ioutil.ReadAll(commits_file)
	if err != nil {
		return spec, err
	}

	err = json.Unmarshal(bytes, &spec.Commits)

	if err != nil {
		return &Specification{}, err
	}

	return spec, nil
}

func ProcessSourceToContent(mdReference *Markdown) error{
	resp := &MarkDownResponse{ remoteResource: mdReference}
	cont, err := resp.Read()

	if err != nil {
		return err
	}

	mdReference.Content = mdReference.Content + cont

	return nil
}

func GetTestResultsFromString(input string) ([]*GoTestResult,error) {
	gtrs := []*GoTestResult{}
	lines := strings.Split(input,"\n")

	for _, v := range lines{
		if strings.Contains(v,"---"){
			gtr := GoTestResult{}
			//This is a result line
			lineBytes := []byte(v)
			err := json.Unmarshal(lineBytes,&gtr)

			if err != nil {
				return gtrs, err
			}

			//Cleanup the time
			parseFormat := "2006-01-02T15:04:05"
			outputFormat := "2 Jan 2006 15:04:05"
			timeValue := strings.Split(gtr.Time,".")[0]

			t, err  := time.Parse(parseFormat,timeValue)

			if err != nil {
				return gtrs, err
			}

			gtr.Time = t.Format(outputFormat)

			gtrs = append(gtrs,&gtr)
		}
	}


	return gtrs,nil
}

func TestsByTag(tag string, tests []*GoTestResult) []*GoTestResult {
	var gtrs []*GoTestResult

	for _, v := range tests {
		if strings.ToLower(v.Test) == strings.ToLower(tag){
			gtrs = append(gtrs,v)
		}
	}

	return gtrs
}

type HTTPMarkDownSource struct {
	Source string
}

type S3MarkDownSource struct {
	Source string
}

type MarkDownResponse struct {
	remoteResource  MarkDownGetter
}

type MarkDownGetter interface{
	Get() (string, error)
}

type MarkDownResponseReader interface {
	Read() (string, error)
}

func (md *Markdown) Get() (string,error){
	resp, err := http.Get(md.Source)
	if err != nil || resp == nil {
		return "", err
	}

	if  resp.StatusCode != 200 {
		return "", fmt.Errorf("Invalid response code when attempting to acquire %s", md.Source)
	}

	bytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(bytes), nil
}


func (md *MarkDownResponse) Read() (string, error){
	return md.remoteResource.Get()
}

