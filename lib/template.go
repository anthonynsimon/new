package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
	"gopkg.in/yaml.v2"
)

var (
	// configFilename is the root template config to look for when rendering a template
	configFilename = ".new.yml"
)

// templateCfg is the parsed template config
type templateCfg struct {
	Version     string
	Description string
	Params      []templateParam
}

// templateParam represents each user defined parameter in the config
type templateParam struct {
	Name     string
	Required bool
	Kind     string
	Prompt   string
	Enum     []string
}

// TemplateContext is a map containing the resolved template parameters and rendering context
type TemplateContext struct {
	Params map[string]interface{}
}

// localTemplate represents all metadata required for a template rendering
type localTemplate struct {
	config          *templateCfg
	sourcePath      string
	destinationPath string
	ctx             *TemplateContext
}

// Template contains the metadata for rendering a template from a source into a destination path
type Template interface {
	SourcePath() string
	DestinationPath() string
	Context() *TemplateContext
	Resolve() error
	Render() error
}

// NewTemplate creates a new template
func NewTemplate(templatePath, destinationPath string) Template {
	return &localTemplate{
		config:          nil,
		destinationPath: destinationPath,
		sourcePath:      templatePath,
		ctx:             nil,
	}
}

// Resolve loads a template config and resolves its context
func (t *localTemplate) Resolve() error {
	config, err := loadTemplateConfig(t.sourcePath)
	if err != nil {
		return err
	}

	t.config = config

	ctx, err := buildContext(config)
	if err != nil {
		return err
	}

	t.ctx = ctx

	return nil
}

// Context returns the resolved template context
func (t *localTemplate) Context() *TemplateContext {
	return t.ctx
}

// SourcePath returns the source template path
func (t *localTemplate) SourcePath() string {
	return t.sourcePath
}

// DestinationPath returns the destination path for the template to be rendered
func (t *localTemplate) DestinationPath() string {
	return t.destinationPath
}

// Render resolves a template if not resolved yet and renders it into a given path
func (t *localTemplate) Render() error {
	if t.config == nil || t.ctx == nil {
		if err := t.Resolve(); err != nil {
			return err
		}
	}

	if err := os.MkdirAll(t.destinationPath, os.ModePerm); err != nil {
		return err
	}

	return filepath.Walk(t.sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == t.sourcePath || filepath.Base(path) == configFilename {
			return nil
		}

		filename, err := renderFilename(path, t.ctx)
		if err != nil {
			return err
		}

		dst, err := resolveDestinationPath(t.sourcePath, t.destinationPath, filename)
		if err != nil {
			return err
		}

		fmt.Printf("Rendering %s\n", dst)

		if info.IsDir() {
			return os.MkdirAll(dst, os.ModePerm)
		}

		return renderFile(path, dst, t.ctx)
	})
}

func renderFilename(path string, ctx *TemplateContext) (string, error) {
	filenameTemplate, err := template.New("filename").Parse(path)
	if err != nil {
		return "", err
	}
	var templatedFilename bytes.Buffer
	if err := filenameTemplate.Execute(&templatedFilename, ctx); err != nil {
		return "", err
	}
	return templatedFilename.String(), nil
}

// resolveDestinationPath resolves the destination path for a template file
func resolveDestinationPath(sourceTemplatePath, destinationPath, targetFilePath string) (string, error) {
	relativeToDestination := strings.Replace(targetFilePath, sourceTemplatePath, "", 1)
	if strings.HasPrefix(relativeToDestination, "/") {
		relativeToDestination = relativeToDestination[1:]
	}
	destination := filepath.Join(destinationPath, relativeToDestination)
	return destination, nil
}

// renderFile renders a file from a source path into the destination path given a rendering context
func renderFile(source, destination string, ctx *TemplateContext) error {
	template, err := template.ParseFiles(source)
	if err != nil {
		return nil
	}

	fout, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer fout.Close()

	return template.Execute(fout, ctx)
}

// loadTemplateConfig reads and parses a template config at the root of the template path
func loadTemplateConfig(templatePath string) (*templateCfg, error) {
	swissConfig, err := ioutil.ReadFile(filepath.Join(templatePath, configFilename))
	if err != nil {
		return nil, err
	}

	config := templateCfg{}
	if err := yaml.Unmarshal(swissConfig, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// buildContext creates a new rendering context given the user defined parameters
func buildContext(config *templateCfg) (*TemplateContext, error) {
	if config.Description != "" {
		fmt.Printf("\n%s\n\n", config.Description)
	}

	ctx := TemplateContext{
		Params: make(map[string]interface{}),
	}
	for _, param := range config.Params {
		switch param.Kind {
		case "enum":
			prompt := promptui.Select{
				Label: param.Prompt,
				Items: param.Enum,
			}

			_, result, err := prompt.Run()
			if err != nil {
				return nil, err
			}

			ctx.Params[param.Name] = result
		default:
			prompt := promptui.Prompt{
				Label: param.Prompt,
				Validate: func(value string) error {
					if param.Required && len(value) == 0 {
						return fmt.Errorf("%s is required", param.Name)
					}
					return nil
				},
			}

			result, err := prompt.Run()
			if err != nil {
				return nil, err
			}

			ctx.Params[param.Name] = result
		}
	}
	return &ctx, nil
}
