package autograder

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Labs struct {
	Labs []Lab
}

type Lab struct {
	name        string
	LabTestCase []Testcase
}

type Testcase struct {
	Type     string
	Expected []Expected
}

type Expected struct {
	Values   []string
	Points   float32
	Feedback string
}

func (c *Labs) getConf(path string) *Labs {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatal("Unmarshal: ", err)
	}

	return c
}

func (c *Labs) GetNumLabs() {
	fmt.Println("this is a test")

	var a Labs
	a.getConf("/Users/ninjamian/go/src/autograder/test_case.yaml")

	for _, element := range a.Labs {
		fmt.Println(element.name)
	}
}
