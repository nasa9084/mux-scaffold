package command

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nasa9084/mux-scaffold/template"
)

const (
	initDescription = `mux-scaffold init command creates a new web application`
	initHelp        = `Usage:
        mux-scaffold init <path>

Description:
        %s

Example:
        mux-scaffold init foo
`
)

var dirsToCreate = []string{
	"db",
	"model",
}

// Init command
type Init struct {
	ProjectName  string
	ProjectDir   string
	DBNamePrefix string
	PackageName  string
}

// Exec to implement command interface
func (c *Init) Exec(args []string) int {
	if len(args) == 0 {
		c.Help()
		return 2
	}
	dir, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Printf(`given directory is not existing.`)
		return 2
	}

	wd, _ := os.Getwd()
	wd = filepath.ToSlash(wd)
	root := ""
	for _, p := range filepath.SplitList(os.Getenv("GOPATH")) {
		p = filepath.ToSlash(p)
		if strings.HasPrefix(strings.ToLower(wd), strings.ToLower(filepath.ToSlash(filepath.Join(p, "src"))+"/")) {
			root = wd[len(p+"/src/"):]
		}
	}

	c.ProjectName = filepath.Base(dir)
	c.ProjectDir = dir
	c.DBNamePrefix = filepath.Base(dir)
	c.PackageName = path.Join(root, c.ProjectName)
	dirsToCreate = append(dirsToCreate, filepath.Join("cmd/", c.ProjectName))
	c.createLayout()

	//c.installFiles("helpers")
	///c.installFiles("config")
	//c.installFiles("controllers")
	c.installFiles("db")
	c.installFiles(filepath.Join("cmd", "{{.ProjectName}}"))
	c.installFile("", "main.go.tmpl", c.ProjectName+".go")
	c.installFile("", "handlers.go.tmpl", "handlers.go")
	c.installFile("", "httputils.go.tmpl", "httputils.go")
	return 0
}

// Description to implement command interface
func (c *Init) Description() string {
	return initDescription
}

// Help to implement command interface
func (c *Init) Help() {
	fmt.Printf(initHelp, initDescription)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *Init) createLayout() {
	for _, dirName := range dirsToCreate {
		path := filepath.Join(c.ProjectDir, dirName)
		must(os.MkdirAll(path, 00755))
	}
}

func (c *Init) installFiles(dirName string) {
	helperFiles, err := filepath.Glob(template.TemplatePath(filepath.Join(dirName, "*.tmpl")))
	if err != nil {
		panic(err)
	}
	for _, templateFile := range helperFiles {
		outputFileName := filepath.Base(templateFile)
		outputFileName = strings.TrimRight(outputFileName, ".tmpl")
		c.installFile(dirName, templateFile, outputFileName)
	}
}

func (c *Init) installFile(dirName string, templateFile string, outputFileName string) {
	dirName = strings.Replace(dirName, "{{.ProjectName}}", c.ProjectName, -1)
	builder := template.NewBuilder(templateFile)
	builder.WriteToPath(filepath.Join(c.ProjectDir, dirName, outputFileName), c)
}
