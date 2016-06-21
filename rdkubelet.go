package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

func check(e error) {
       if e != nil {
		panic(e)
        }
}

func rdKubelet(fp string, pods []string) map[string]string{
	file, err := os.Open(fp + "/kubelet.log")
	check(err)
	defer file.Close()

	mapPodKubelet := map[string]string{}

	scanner := bufio.NewScanner(file)
	check(scanner.Err())
		
	for _, pod := range(pods) {
		mapPodKubelet[pod] = ""
	}
	
	for scanner.Scan() {
		line := scanner.Text()
		for _, pod := range(pods) {
			if strings.Contains(line, pod) {
				a := []string{mapPodKubelet[pod], line}
				mapPodKubelet[pod] = strings.Join(a,"\n")
			}
		}
	}
	
	for pod, lines := range mapPodKubelet {
		fmt.Println("====================================")
		fmt.Println("Kubelet log for pod:", pod)
		fmt.Println(lines)
		fmt.Println("====================================")
	}

	return mapPodKubelet
}

