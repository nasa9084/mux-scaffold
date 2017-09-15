package template

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"bitbucket.org/pkg/inflect"
)

type state struct {
	n int
}

func (s *state) Set(n int) int {
	s.n = n
	return n
}

func (s *state) Inc() int {
	s.n++
	return s.n
}

var s state
var (
	funcMap = template.FuncMap{
		"Pluralize":  inflect.Pluralize,
		"Underscore": inflect.Underscore,
		"ToUpper":    strings.ToUpper,
		"ToLower":    strings.ToLower,
		"set":        s.Set,
		"inc":        s.Inc,
		"ret": func(fieldType string) string {
			switch fieldType {
			case
				"int",
				"float64",
				"bool":
				return ", _"
			}
			return ""
		},
		"conf": func(origin string, fieldType string) string {
			switch fieldType {
			case "int":
				return "strconv.Atoi(" + origin + ")"
			case "float64":
				return "strconf.ParseFloat(" + origin + ", 64)"
			case "bool":
				return "strconv.ParseBool(" + origin + ")"
			}
			return origin
		},
		"mapUnderscore": func(strs []string) []string {
			ret := []string{}
			for _, s := range strs {
				ret = append(ret, inflect.Underscore(s))
			}
			return ret
		},
		"prefix": func(prefix, str string) string {
			return prefix + str
		},
		"suffix": func(str, suffix string) string {
			return str + suffix
		},
		"mapPrefix": func(strs []string, prefix string) []string {
			ret := []string{}
			for _, s := range strs {
				ret = append(ret, prefix + s)
			}
			return ret
		},
		"mapSuffix": func(strs []string, suffix string) []string {
			ret := []string{}
			for _, s := range strs {
				ret = append(ret, s + suffix)
			}
			return ret
		},
		"join": strings.Join,
		"joinfptr": func(prefix string, strs []string, sep string) string {
			ret := []string{}
			for _, s := range strs {
				ret = append(ret, "&" + prefix + "." + s)
			}
			return strings.Join(ret, ", ")
		},
	}
)

// Builder is Template builder
type Builder struct {
	TemplateName string
	TemplatePath string
}

// NewBuilder returns a new Buildler obj
func NewBuilder(templatePath string) *Builder {
	if !filepath.IsAbs(templatePath) {
		templatePath = TemplatePath(templatePath)
	}

	templateName := filepath.Base(templatePath)
	builder := &Builder{
		TemplateName: templateName,
		TemplatePath: templatePath,
	}

	return builder
}

// Template returns template the builder pointing
func (builder *Builder) Template() *template.Template {
	contents := LoadTemplateFromFile(builder.TemplatePath)
	tmpl := template.Must(template.New(builder.TemplateName).Funcs(funcMap).Parse(contents))

	return tmpl
}

func (builder *Builder) Write(writer io.Writer, data interface{}) {
	tmpl := builder.Template()
	err := tmpl.Execute(writer, data)
	if err != nil {
		panic(err)
	}
}

// WriteToPath write a file using template
func (builder *Builder) WriteToPath(outputPath string, data interface{}) {
	printAction("green+h:black", "create", outputPath)
	if _, err := os.Stat(outputPath); err == nil {
		printAction("red:hblack", "skip", outputPath)
		return
	}

	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	builder.Write(f, data)
}

// InsertAfterToPath insert something to given path
func (builder *Builder) InsertAfterToPath(outputPath, after string, data interface{}) {
	printAction("cyan+h:black", "insert", outputPath)

	newFilePath := outputPath + ".new"

	f, err := os.Open(outputPath)
	if err != nil {
		panic(err)
	}

	outputFile, err := os.Create(newFilePath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(f)
	writer := bufio.NewWriter(outputFile)

	for scanner.Scan() {
		line := scanner.Text()

		writer.WriteString(line + "\n")
		if strings.HasPrefix(line, after) {
			builder.Write(writer, data)
		}
	}

	writer.Flush()
	outputFile.Close()
	f.Close()

	os.Rename(newFilePath, outputPath)
}
