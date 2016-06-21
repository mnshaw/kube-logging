package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
)

func rdKubeAPI(fp string, mapPodKubelet map[string]string) {
	file, err := os.Open(fp + "/kube-apiserver.log")
	check(err)
	defer file.Close()

	for pod, kubelet := range(mapPodKubelet) {
		splitKubelet := strings.Split(kubelet, "\n")
		
	}