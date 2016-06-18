package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	type Testcase struct {
		Name      string `xml:"name,attr"`
		ClassName string `xml:"classname,attr"`
		Failure   string `xml:"failure"`
	}
	type Testsuite struct {
		TestCount int        `xml:"tests,attr"`
		FailCount int        `xml:"failures,attr"`
		Testcases []Testcase `xml:"testcase"`
	}

	testSuite := &Testsuite{}
	pwd, _ := os.Getwd()
	fp := filepath.Join(pwd, "10528-06_15_2016-15_31_15/artifacts/junit_03.xml")
	failures := map[string]string{}
	
	file, err := os.Open(fp)
	check(err)
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	check(err)
	
	err = xml.Unmarshal(data, testSuite)
	
	check(err)

	for _, tc := range testSuite.Testcases{
		if  tc.Failure != "" {
			failures[fmt.Sprintf("%v {%v}", tc.Name, tc.ClassName)] = tc.Failure
		}
	}
	fmt.Println(failures)

	fmt.Println("\nDone printing failures\n")
	for _, v := range failures {
		a := strings.SplitAfter(v, "\n")
		for _, line := range a {
			line = strings.TrimSpace(line)
			if (strings.HasPrefix(line, "pod")) {
				podInfo := strings.Split(line, "'")
				fmt.Println(podInfo[1])
			}
		}

		fmt.Println("")
	}
}



	
