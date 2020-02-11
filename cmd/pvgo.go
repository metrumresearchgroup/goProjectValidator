package cmd

import (
	"bytes"
	"fmt"
	projectvalidator "github.com/metrumresearchgroup/projectValidator"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		var locatedTests []projectvalidator.GoTestResult
		fs := afero.NewOsFs()
		//Maps by tag all of the test results
		//testMap := make(map[string][]projectvalidator.TestResult)
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

		if err != nil {
			log.Fatalf("%s",err)
		}

		//Populate content for the specifications
		for _, v := range spec.MarkDown{
			err := projectvalidator.ProcessSourceToContent(v)
			if err != nil {
				log.Fatal(err)
			}
		}


		if ok, _ := afero.Exists(fs,vc.TestsDirectory); ! ok {
			log.Fatalf("%s is designated as the test outputs directory, but it doesn't appear to exist", vc.TestsDirectory)
		}


		contents, err := afero.ReadDir(fs,vc.TestsDirectory)

		if err != nil {
			log.Fatalf("Unable to access or read contents of test file directory at %s", vc.TestsDirectory)
		}

		for _, v := range contents {
			if strings.Contains(v.Name(),".json"){
				//If it's a json file. Let's read all the contents in
				file, err := os.Open(filepath.Join(vc.TestsDirectory,v.Name()))
				if err != nil {
					log.Fatalf("Failure to open file %s", v.Name())
				}

				bytes, err := ioutil.ReadAll(file)
				if err != nil {
					log.Fatal("Unable to read contents of file")
				}

				file.Close()
				fileContents := string(bytes)

				//Deal with miscellaneous newlines created by go test json output.
				newContents := strings.ReplaceAll(fileContents,`\n"`,`"`)

				tests, err := projectvalidator.GetTestResultsFromString(newContents)

				if err != nil {
					log.Fatalf("An error occurred trying to extract go testing details from the located file. %s", err)
				}

				locatedTests = append(locatedTests,tests...)
			}
		}

		//Now we have a listing of tests. Build map of tags we care about
		for _, v := range spec.Stories{
			for _, m := range v.MarkDown {
				err := projectvalidator.ProcessSourceToContent(m)
				if err != nil {
					log.Fatal(err)
				}
			}
			var storytests []projectvalidator.GoTestResult
			for _, t := range v.Tags {
				storytests = append(storytests,projectvalidator.TestsByTag(t,locatedTests)...)
			}

			v.Tests = storytests
		}

		log.Info("We've collected tests now!")

		specOutput, err := MarkDownFromScenario("../../specification.md.t",spec)

		if err != nil {
			log.Fatal(err)
		}

		println(specOutput)

		testingOutput, err := MarkDownFromScenario("../../testing.md.t",spec)
		if err != nil {
			log.Fatal(err)
		}

		println(testingOutput)
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


func MarkDownFromScenario(markdownfile string, spec *projectvalidator.Specification) (string,error) {
	t, err := template.New(filepath.Base(markdownfile)).ParseFiles(markdownfile)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf,spec)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
