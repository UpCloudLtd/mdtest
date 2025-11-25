package testrun

import (
	"encoding/xml"
	"fmt"
	"os"
	"slices"

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

func readTestsuite(run RunResult) junitTestsuite {
	testsuite := junitTestsuite{
		Name:      run.Name,
		Tests:     len(run.TestResults),
		Time:      run.Finished.Sub(run.Started).Seconds(),
		Timestamp: run.Started.UTC().Format("2006-01-02T15:04:05.000000"),
	}

	var errors, failures int
	for _, result := range run.TestResults {
		testcase := junitTestcase{
			Classname: run.Name,
			Name:      result.Name,
			Time:      result.Finished.Sub(result.Started).Seconds(),
		}

		if result.Error != nil {
			if utils.IsCanceled(result.Error) {
				testcase.Failure = append(testcase.Failure, result.Error.Error())
			} else {
				testcase.Error = append(testcase.Error, result.Error.Error())
				errors++
			}
		}

		for i, step := range result.Results {
			if !step.Success && step.Error != nil {
				msg := step.Error.Error()
				if !slices.Contains(testcase.Failure, msg) {
					testcase.Failure = append(testcase.Failure, msg)
				}
			}

			output := step.Output
			if output == "" {
				output = "# No output"
			}
			testcase.SystemOut += fmt.Sprintf("# Step %d:\n%s", i+1, output)
		}

		if len(testcase.Failure) > 0 {
			failures++
		}

		testsuite.Testcases = append(testsuite.Testcases, testcase)
	}

	testsuite.Errors = errors
	testsuite.Failures = failures

	return testsuite
}

func junitReport(run RunResult, path string) error {
	testsuite := readTestsuite(run)

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
