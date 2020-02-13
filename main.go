package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type templateModel struct {
	ProjectPath string
	ProjectName string
}

const templateRepoURL = "https://github.com/lon9/go-%s-api-starter-template.git"

func main() {
	var (
		projectPath string
		outputDir   string
		framework   string
	)

	flag.StringVar(&projectPath, "p", "github.com/lon9/awesomeproject", "path of your awesome project from GOPATH")
	flag.StringVar(&outputDir, "o", "", "path you want to set up project. Default is base of p option")
	flag.StringVar(&framework, "f", "echo", "framework you use (e.g. echo or gin)")
	flag.Parse()

	model := &templateModel{
		ProjectPath: projectPath,
		ProjectName: filepath.Base(projectPath),
	}

	if outputDir == "" {
		outputDir = model.ProjectName
	}

	fmt.Println(model.ProjectName)

	fmt.Println("Cloning template repository")
	tmpDirPath, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDirPath)
	if err := exec.Command("git", "clone", fmt.Sprintf(templateRepoURL, framework), tmpDirPath).Run(); err != nil {
		panic(err)
	}

	fmt.Println("Generating repository")
	if err := filepath.Walk(tmpDirPath, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(p, ".git") {
			return nil
		}
		outputPath := filepath.Join(outputDir, strings.ReplaceAll(p, tmpDirPath, ""))
		fmt.Printf("Generating %s\n", outputPath)

		t := template.Must(template.ParseFiles(p))

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			return err
		}
		f, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer f.Close()

		if err := t.Execute(f, model); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(err)
	}
	fmt.Println("Finished to create repository")
}
