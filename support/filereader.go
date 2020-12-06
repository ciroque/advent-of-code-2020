package support

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func ReadFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(fmt.Errorf("failed to open '%s', %v", filename, err))
	}

	return strings.Split(string(content), "\n")
}
