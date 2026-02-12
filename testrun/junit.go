package testrun

import (
	"encoding/xml"
	"fmt"
	"os"
	"slices"

	"github.com/UpCloudLtd/mdtest/testcase"
	"github.com/UpCloudLtd/mdtest/utils"
)

type junitTestcase struct {
	Classname string   `xml:"classname,attr,omitempty"`
	Name      string   `xml:"name,attr"`
	Failure   []string `xml:"failure,omitempty"`
	Error     []string `xml:"error,omitempty"`
	SystemOut string   `xml:"system-out,omitempty"`
	SystemErr string   `xml:"system-err,omitempty"`
	Time      float64  `xml:"time,attr,omitempty"`
}

type junitTestsuite struct {
	XMLName xml.Name `xml:"testsuite"`

	Name      string  `xml:"name,attr"`
	Tests     int     `xml:"tests,attr"`
	Failures  int     `xml:"failures,attr"`
	Errors    int     `xml:"errors,attr"`
	Skipped   int     `xml:"skipped,attr"`
	Time      float64 `xml:"time,attr,omitempty"`
	Timestamp string  `xml:"timestamp,attr,omitempty"`

	Testcases []junitTestcase `xml:"testcase"`
}

func handleTestStep(i int, step testcase.StepResult, params RunParameters, tcase *junitTestcase) {
	if step.Status == testcase.StepStatusFailure && step.Error != nil {
		msg := step.Error.Error()
		if !slices.Contains(tcase.Failure, msg) {
			tcase.Failure = append(tcase.Failure, msg)
		}
	}

	output := step.Output
	if output == "" {
		output = "# No output\n"
	}
	if step.Status == testcase.StepStatusSkipped {
		output = "# Skipped\n"
	}

	for _, msg := range step.Warnings {
		if params.WarningsAsErrors {
			if !slices.Contains(tcase.Failure, msg) {
				tcase.Failure = append(tcase.Failure, msg)
			}
		} else {
			output += fmt.Sprintf("# Warning: %s\n", msg)
		}
	}

	tcase.SystemOut += fmt.Sprintf("# Step %d:\n%s", i+1, output)
}

func readTestsuite(run RunResult, params RunParameters) junitTestsuite {
	testsuite := junitTestsuite{
		Name:      run.Name,
		Tests:     len(run.TestResults),
		Time:      run.Finished.Sub(run.Started).Seconds(),
		Timestamp: run.Started.UTC().Format("2006-01-02T15:04:05.000000"),
	}

	var errors, failures int
	for _, result := range run.TestResults {
		tcase := junitTestcase{
			Classname: run.Name,
			Name:      result.Name,
			Time:      result.Finished.Sub(result.Started).Seconds(),
		}

		if result.Error != nil {
			if utils.IsCanceled(result.Error) {
				tcase.Failure = append(tcase.Failure, result.Error.Error())
			} else {
				tcase.Error = append(tcase.Error, result.Error.Error())
				errors++
			}
		}

		for i, step := range result.Results {
			handleTestStep(i, step, params, &tcase)
		}

		if len(tcase.Failure) > 0 || params.WarningsAsErrors && result.HasWarnings() {
			failures++
		}

		testsuite.Testcases = append(testsuite.Testcases, tcase)
	}

	testsuite.Errors = errors
	testsuite.Failures = failures

	return testsuite
}

func junitReport(run RunResult, params RunParameters) error {
	path := params.JUnitXML
	testsuite := readTestsuite(run, params)

	data, err := xml.MarshalIndent(testsuite, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JUnit XML: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create JUnit XML file: %w", err)
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write JUnit XML file: %w", err)
	}
	return nil
}
