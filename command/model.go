package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nasa9084/mux-scaffold/template"

	"bitbucket.org/pkg/inflect"
)

const (
	modelDescription = `mux-scaffold model command creates a new model with the given fields`
	modelHelp        = `Usage:
        mux-scaffold model <model name> <field name>:<field type> ...

Description:
        %s

Example:
        mux-scaffold model Post Title:string Body:string
`
)

// Field represents a field of struct
type Field struct {
	FieldType string
	Misc      string
}

// Model command
type Model struct {
	PackageName        string
	ModelName          string
	ModelNamePlural    string
	InstanceName       string
	InstanceNamePlural string
	TemplateName       string
	Fields             map[string]Field
}

// Exec to implements command interface
func (c *Model) Exec(args []string) int {
	if len(args) < 2 {
		c.Help()
		os.Exit(2)
	}
	c.ModelName = inflect.Titleize(args[0])
	c.ModelNamePlural = inflect.Pluralize(c.ModelName)
	c.InstanceName = inflect.CamelizeDownFirst(c.ModelName)
	c.InstanceNamePlural = inflect.Pluralize(c.InstanceName)

	c.PackageName = template.PackageName()

	c.Fields = processFields(args[1:])

	outPath := filepath.Join("model", inflect.Underscore(c.ModelName)+".go")

	builder := template.NewBuilder(filepath.Join("model", "model.go.tmpl"))
	builder.WriteToPath(outPath, c)

	outPath = filepath.Join("db", inflect.Underscore(c.ModelName)+".go")
	builder = template.NewBuilder("db.go.tmp")
	builder.WriteToPath(outPath, c)
	return 0
}

// Description to implements command interface
func (c *Model) Description() string {
	return modelDescription
}

// Help to implements command interface
func (c *Model) Help() {
	fmt.Printf(modelHelp, modelDescription)
}

func processFields(args []string) map[string]Field {
	fields := map[string]Field{}
	for _, arg := range args {
		fieldNameAndType := strings.SplitN(arg, ":", 2)
		key := inflect.Titleize(fieldNameAndType[0])
		name, misc := findFieldType(fieldNameAndType[1])
		field := Field{
			FieldType: name,
			Misc:      misc,
		}
		fields[key] = field
	}
	return fields
}

func findFieldType(name string) (string, string) {
	misc := ""
	switch name {
	case "test":
		name = "string"
	case "float", "double":
		name = "float64"
	case "boolean":
		name = "bool"
	case "integer":
		name = "int"
	case "decimal":
		name = "int64"
	case "time", "date", "datetime":
		name = "int64"
		misc = "time"
	}
	return name, misc
}
