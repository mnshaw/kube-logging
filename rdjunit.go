package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
)

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


func main() {
	pwd, _ := os.Getwd()
	
	images, err := filepath.Glob(os.Args[1] + "/artifacts/tmp*")
	check(err)

	mapPodError := map[string]string{}
	mapPodCont := map[string]string{}
	

	if len(images) != 0 {
		for _, fp := range images {
			junits, err := filepath.Glob(fp + "/junit*")
			check(err)
			for _, ju := range junits {
				pods := getFailedPods(ju, mapPodError, mapPodCont)
				if len(pods) > 0 {
					mapPodKubelet := rdKubelet(fp, pods)
					rdKubeAPI(fp, mapPodKubelet)
				}
			}
		}
	} else {
		moreFiles := true
		for i := 1; moreFiles; i++ {	
			fp := fmt.Sprintf("%s/artifacts/junit_%02d.xml", os.Args[1], i)
			_, err := os.Stat(filepath.Join(pwd, fp));
			if err == nil {
				getFailedPods(fp, mapPodError, mapPodCont)
			} else {
				moreFiles = false
			}
		}
	}

	fmt.Println("mapPodError")
	fmt.Println(mapPodError)
	fmt.Println("")
	fmt.Println("mapPodCont")
	fmt.Println(mapPodCont)
}


func getFailedPods(fp string, mapPodError map[string]string, mapPodCont map[string]string) ([]string){
	testSuite := &Testsuite{}
	
	failures := map[string]string{}
	
	file, err := os.Open(fp)
	check(err)
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	check(err)
	
	err = xml.Unmarshal(data, testSuite)

	if testSuite.FailCount == 0 {
		// no pods failed in this file
		return nil
	}

	check(err)

	for _, tc := range testSuite.Testcases{
		if  tc.Failure != "" {
			failures[fmt.Sprintf("%v {%v}", tc.Name, tc.ClassName)] = tc.Failure
		}
	}
	
	fmt.Println("Failed Pods in file", fp)
	podNames := make([]string, 0)
	for _, v := range failures {
		lines := strings.SplitAfter(v, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if (strings.HasPrefix(line, "pod")) {
				podInfo := strings.Split(line, "'")
				podName := podInfo[1]
				fmt.Println(podName)
				
				//map the pod to the error
				mapPodError[podName] = line

				//map the pod to the container
				cont := strings.Split(strings.Split(line, "ContainerID:")[1], "}")
				mapPodCont[podName] = cont[0]
				podNames = append(podNames, podName)
				
			}
		}
	}

	return podNames
}
	
