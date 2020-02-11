package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	projectvalidator "github.com/metrumresearchgroup/projectValidator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const VERSION string = "v0.0.1"

var pvgoCmd = &cobra.Command{
	Use:   "pvgo",
	Short: "PVGO is a project validation tool for any language",
	Long: `Define your stories, tests and tags in a json file and tell pvgo where to look
for test output. Tests should follow the structure of TestSummary in this package. PVGO will then
stitch located tests into each story and provide a unified output, tying your stories to their respective tests.`,
	Run: func(cmd *cobra.Command, args []string) {
		fs := afero.NewOsFs()
		//Maps by tag all of the test results
		testMap := make(map[string][]projectvalidator.TestResult)
		//Codify viper details into the struct
		var vc projectvalidator.ValidationConfiguration

		err := viper.Unmarshal(&vc)

		if err != nil {
			log.Fatal("Unable to marshal Viper configuration details into struct ")
		}

		file, err := fs.Open(vc.ScenarioFile)

		if err != nil {
			log.Fatalf("Unable to open the specification file at %s",vc.ScenarioFile)
		}

		spec, err := projectvalidator.NewSpecification(file)
		println(spec.Scope)

		if err != nil {
			log.Fatalf("Scenario file is listed as %s, but accessing it yielded the following error: %s", vc.ScenarioFile,err)
		}


		if ok, _ := afero.Exists(fs,vc.TestsDirectory); ! ok {
			log.Fatalf("%s is designated as the test outputs directory, but it doesn't appear to exist", vc.TestsDirectory)
		}

		//var tests []projectvalidator.TestResult

		contents, err := afero.ReadDir(fs,vc.TestsDirectory)

		if err != nil {
			log.Fatalf("Unable to access or read contents of test file directory at %s", vc.TestsDirectory)
		}

		for _, v := range contents {
			var tr projectvalidator.TestResult

			file, err := fs.Open(filepath.Join(vc.TestsDirectory,v.Name()))

			if err != nil {
				log.Errorf("Test file listed at %s could not be opened. Details are: %s",filepath.Join(vc.TestsDirectory,v.Name()),err)
				continue
			}

			tbytes, err := ioutil.ReadAll(file)

			if err != nil {
				log.Errorf("Attempting to read the file failed. %s", err)
				continue
			}

			err = json.Unmarshal(tbytes,&tr)

			if err != nil {
				log.Errorf("Unable to marshal the contents! Possibly invalid JSON. %s", err)
				continue
			}

			for _, tag := range tr.Tags{
				//Does the key exist
				if val, ok := testMap[tag]; !ok {
					testMap[tag] = []projectvalidator.TestResult{
						tr,
					}
				} else {
					//If it does. does this test exist in it?
					matched := false
					for _, t := range val {
						if t.Name == tr.Name {
							matched = true
						}
					}

					//If not, add it to the key
					if ! matched {
						val = append(val,tr)
					}
				}
			}
		}

		println("Loaded")

		//Now stitch the tests from the map into the stories

		for k, s := range spec.Stories {
			locals := s
			for _, tag := range s.Tags {
				if tests, ok := testMap[tag]; ok {
					locals.Tests = append(locals.Tests,tests...)
				}
			}
			spec.Stories[k] = locals
		}

		//Summarize
		println("Beginning summarization")

		for k, s := range spec.Stories {
			locals := s
			for _, test := range locals.Tests{
				locals.Summary.TestCount ++
				if test.Passed {
					locals.Summary.PassedTests++
				} else {
					locals.Summary.FailedTests++
				}
			}
			spec.Stories[k] = locals
		}

		//pretty, _ := json.MarshalIndent(spec," ", "    ")
		//
		//println(string(pretty))

		//Try to generate the template
		t, err := template.New("specification.md.t").ParseFiles("../../specification.md.t")

		buf := new(bytes.Buffer)

		err = t.Execute(buf,spec)



		println(buf.String())
	},
}

func Execute() {
	if err := pvgoCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init(){
	//subcommands
	const versionIdentifier string = "version"

	//Flags
	const scenarioFileIdentifier string = "scenarioFile"
	pvgoCmd.Flags().StringP(scenarioFileIdentifier,"s","scenario.json","Specify the path to a JSON file containing your validation scenario")
	viper.BindPFlag(scenarioFileIdentifier, pvgoCmd.Flags().Lookup(scenarioFileIdentifier))

	const testsDirectoryIdentifier string = "testsDirectory"
	pvgoCmd.Flags().StringP(testsDirectoryIdentifier,"t","tests", "Specify a directory in which to look for the JSON results of all specified tests")
	viper.BindPFlag(testsDirectoryIdentifier, pvgoCmd.Flags().Lookup(testsDirectoryIdentifier))
}