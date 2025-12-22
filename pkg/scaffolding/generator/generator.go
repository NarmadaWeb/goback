// pkg/scaffolding/generator/generator.go

package generator

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/NarmadaWeb/goback/pkg/config"
	"github.com/NarmadaWeb/goback/pkg/scaffolding"
	"github.com/iancoleman/strcase"
)

const (
	archSimple    = "simple"
	archClean     = "clean"
	archDDD       = "ddd"
	archHexagonal = "hexagonal"
)

const (
	pathConfig     = "config/framework.go"
	pathMigrate    = "pkg/migrate/migrate.go"
	pathDatabase   = "infrastructure/database/connection.go"
	pathRoutes     = "interfaces/routes/routes.go"
	pathHandlers   = "interfaces/handlers/handlers.go"
	pathMiddleware = "interfaces/middleware/middleware.go"
	templatesDir   = "templates"
	connectionTmpl = "connection.go.tmpl"
	frameworksDir  = "frameworks"
	devopsDir      = "devops"
	helmDir        = "helm"
	ansibleDir     = "ansible"
)

// TemplateGenerator handles project generation from templates
type TemplateGenerator struct {
	Config           *config.ProjectConfig
	OutputDir        string
	progressCallback func(step int, message string)
	errorCallback    func(step int, err error)
	currentStep      int
	totalSteps       int
}

// NewTemplateGenerator creates a new template generator
func NewTemplateGenerator(cfg *config.ProjectConfig) *TemplateGenerator {
	return &TemplateGenerator{
		Config:     cfg,
		OutputDir:  cfg.OutputDir,
		totalSteps: 7,
	}
}

// SetProgressCallback sets the progress callback
func (tg *TemplateGenerator) SetProgressCallback(callback func(step int, message string)) {
	tg.progressCallback = callback
}

// SetErrorCallback sets the error callback
func (tg *TemplateGenerator) SetErrorCallback(callback func(step int, err error)) {
	tg.errorCallback = callback
}

// Generate generates the project structure and files
func (tg *TemplateGenerator) Generate() error {
	steps := []struct {
		name    string
		handler func() error
	}{
		{"Validating configuration", tg.validateConfiguration},
		{"Generating base files", tg.generateBaseFiles},
		{"Generating framework files", tg.generateFrameworkFiles},
		{"Generating database config", tg.generateDatabaseConfig},
		{"Generating Tool files", tg.generateToolFiles},
		{"Generating architecture files", tg.generateArchitectureFiles},
		{"Generating DevOps files", tg.generateDevOpsFiles},
	}

	for i, step := range steps {
		tg.currentStep = i
		tg.reportProgress(i, fmt.Sprintf("Step %d/%d: %s", i+1, len(steps), step.name))

		if err := step.handler(); err != nil {
			tg.reportError(i, err)
			return fmt.Errorf("step %d (%s) failed: %w", i+1, step.name, err)
		}
	}

	tg.reportProgress(len(steps), "Project generation completed successfully!")
	return nil
}

// generateFileFromTemplate is the main helper function for processing templates.
// It reads a template file, creates the destination directory if it doesn't exist,
// executes the template with the config data, and writes the result.
func (tg *TemplateGenerator) generateFileFromTemplate(destPath, templatePath string, delims ...string) error {
	if destPath == "" {
		return fmt.Errorf("empty destination path for template %s", templatePath)
	}

	// Remove .tmpl extension from destination path
	destPath = strings.TrimSuffix(destPath, ".tmpl")

	fullDestPath := filepath.Join(tg.OutputDir, destPath)
	// All template paths are now relative to the embedded `templates` directory
	fullTemplatePath := filepath.ToSlash(filepath.Join(templatesDir, templatePath))

	// Read template content from embedded FS
	templateContent, err := scaffolding.Templates.ReadFile(fullTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read embedded template %s: %w", fullTemplatePath, err)
	}

	// Create destination directory if it doesn't exist
	if err2 := os.MkdirAll(filepath.Dir(fullDestPath), 0755); err2 != nil {
		return fmt.Errorf("failed to create directory for %s: %w", fullDestPath, err2)
	}

	// Create destination file
	outputFile, err := os.Create(fullDestPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", fullDestPath, err)
	}
	defer outputFile.Close()

	// Custom template functions
	funcMap := template.FuncMap{
		"title":      strings.ToTitle,
		"toTitle":    strings.ToTitle,
		"snakeCase":  strcase.ToSnake,
		"kebabCase":  strcase.ToKebab,
		"upper":      strings.ToUpper,
		"lower":      strings.ToLower,
		"replaceAll": strings.ReplaceAll,
		"b64enc":     func(s string) string { return base64.StdEncoding.EncodeToString([]byte(s)) },
		"default": func(val string, def string) string {
			if val == "" {
				return def
			}
			return val
		},
	}

	// Parse and execute template
	tmpl := template.New(filepath.Base(templatePath)).Funcs(funcMap)
	if len(delims) == 2 {
		tmpl = tmpl.Delims(delims[0], delims[1])
	}
	parsedTmpl, err := tmpl.Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Use tg.Config directly so the template can access .Architecture.String(), etc.
	if err := parsedTmpl.Execute(outputFile, tg.Config); err != nil {
		return fmt.Errorf("failed to execute template %s: %w", templatePath, err)
	}

	return nil
}

func (tg *TemplateGenerator) getDestinationPath(fileType string) string {
	arch := strings.ToLower(string(tg.Config.Architecture))

	// Default to simple architecture paths
	paths := map[string]string{
		"config":     "internal/config/framework.go",
		"database":   "internal/database/connection.go",
		"routes":     "internal/routes/routes.go",
		"handlers":   "internal/handlers/handlers.go",
		"middleware": "internal/middleware/middleware.go",
		"models":     "internal/models/base_model.go",
		"validator":  "internal/utils/validator.go",
		"migrate":    "internal/migrate/migrate.go",
	}

	// gocritic (ifElseChain) fix: Changed if-else to switch.
	switch arch {
	case archDDD:
		paths["config"] = pathConfig
		paths["database"] = pathDatabase
		paths["routes"] = pathRoutes
		paths["handlers"] = pathHandlers
		paths["middleware"] = pathMiddleware
		paths["migrate"] = pathMigrate
	case archClean:
		paths["config"] = pathConfig
		paths["database"] = pathDatabase
		paths["routes"] = pathRoutes
		paths["handlers"] = pathHandlers
		paths["middleware"] = pathMiddleware
		paths["migrate"] = pathMigrate
	case archHexagonal:
		paths["config"] = pathConfig
		paths["database"] = "adapters/secondary/database/connection.go"
		paths["routes"] = "adapters/primary/http/routes.go"
		paths["handlers"] = "adapters/primary/http/handlers.go"
		paths["middleware"] = "adapters/primary/http/middleware.go"
		paths["migrate"] = pathMigrate
	}

	return paths[fileType]
}

func (tg *TemplateGenerator) validateConfiguration() error {
	validationErrors := config.ValidateProjectConfig(tg.Config)
	if len(validationErrors) > 0 {
		// Join the errors into a single string to return as an error
		return fmt.Errorf("configuration validation failed: %s", strings.Join(validationErrors, ", "))
	}
	return nil
}

// generateBaseFiles generates the base project files.
func (tg *TemplateGenerator) generateBaseFiles() error {
	baseTemplates := map[string]string{
		"go.mod":             "base/go.mod.tmpl",
		".gitignore":         "base/gitignore.tmpl",
		"README.md":          "base/README.md.tmpl",
		"Makefile":           "base/Makefile.tmpl",
		".env":               "base/.env.tmpl",
		".env.example":       "base/env.example.tmpl",
		"Dockerfile":         "base/Dockerfile.tmpl",
		"docker-compose.yml": "base/docker-compose.yml.tmpl",
	}

	for dest, src := range baseTemplates {
		if err := tg.generateFileFromTemplate(dest, src); err != nil {
			return err
		}
	}

	// Conditionally generate validator
	arch := strings.ToLower(string(tg.Config.Architecture))
	if arch == archSimple || arch == "" {
		destPath := tg.getDestinationPath("validator")
		templatePath := "base/internal/utils/validator.go.tmpl"
		if err := tg.generateFileFromTemplate(destPath, templatePath); err != nil {
			return err
		}
	}

	return nil
}

// generateFrameworkFiles generates the framework-specific files.
func (tg *TemplateGenerator) generateFrameworkFiles() error {
	framework := string(tg.Config.Framework)
	if framework == "" {
		return nil // No framework selected
	}

	frameworkDir := filepath.ToSlash(filepath.Join(templatesDir, frameworksDir, framework))
	globPath := filepath.Join(frameworkDir, "*.tmpl")

	files, err := fs.Glob(scaffolding.Templates, globPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		templatePath := strings.TrimPrefix(file, templatesDir+string(filepath.Separator))
		frameworkDirInTmpl := filepath.Join(frameworksDir, framework)

		var destPath string
		switch filepath.Base(templatePath) {
		case "main.go.tmpl":
			destPath = "cmd/api/main.go"
		case "routes.go.tmpl":
			destPath = tg.getDestinationPath("routes")
		case "config.go.tmpl":
			destPath = tg.getDestinationPath("config")
		case "handlers.go.tmpl":
			destPath = tg.getDestinationPath("handlers")
		case "middleware.go.tmpl":
			destPath = tg.getDestinationPath("middleware")
		default:
			destPath = strings.TrimPrefix(templatePath, frameworkDirInTmpl+string(filepath.Separator))
		}

		if err := tg.generateFileFromTemplate(destPath, templatePath); err != nil {
			return fmt.Errorf("failed to generate framework file from %s: %w", templatePath, err)
		}
	}
	return nil
}

// generateDatabaseConfig generates the database configuration files.
func (tg *TemplateGenerator) generateDatabaseConfig() error {
	tool := strings.ToLower(tg.Config.Tool.String())
	if tool == "" {
		return nil // No tool selected
	}

	templatePath := filepath.Join("databases", tool, connectionTmpl)
	destPath := tg.getDestinationPath("database")
	fullTemplatePath := filepath.ToSlash(filepath.Join(templatesDir, templatePath))

	// Check if the template file exists in embedded FS
	if _, err := scaffolding.Templates.Open(fullTemplatePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// Fallback to database type if tool-specific connection doesn't exist
			dbType := strings.ToLower(tg.Config.Database.String())
			templatePath = filepath.Join("databases", dbType, connectionTmpl)
			fullTemplatePath = filepath.ToSlash(filepath.Join(templatesDir, templatePath))
			if _, err2 := scaffolding.Templates.Open(fullTemplatePath); err2 != nil {
				if errors.Is(err2, fs.ErrNotExist) {
					return nil // Ignore if no suitable template is found
				}
				return err2 // Return other unexpected errors
			}
		} else {
			return err // Other error
		}
	}

	return tg.generateFileFromTemplate(destPath, templatePath)
}

// generateToolFiles generates the Tool-specific files.
func (tg *TemplateGenerator) generateToolFiles() error {
	tool := strings.ToLower(tg.Config.Tool.String())
	if tool == "" {
		return nil
	}

	templateRootDir := filepath.ToSlash(filepath.Join(templatesDir, "tools", tool))
	if _, err := scaffolding.Templates.Open(templateRootDir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil // No tool-specific files to generate
		}
		return err
	}

	return fs.WalkDir(scaffolding.Templates, templateRootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Skip connection files as they are handled by generateDatabaseConfig
		if filepath.Base(path) == connectionTmpl {
			return nil
		}

		relPath, err := filepath.Rel(templateRootDir, path)
		if err != nil {
			return err
		}

		slashedRelPath := filepath.ToSlash(relPath)
		var destPath string

		// Determine destination based on file path pattern and architecture
		arch := strings.ToLower(string(tg.Config.Architecture))

		switch {
		// Handle migrations folder (e.g. db/migrations or just migrations)
		case strings.HasPrefix(slashedRelPath, "db/migrations/") || strings.HasPrefix(slashedRelPath, "migrations/"):
			var filename string
			if strings.HasPrefix(slashedRelPath, "db/migrations/") {
				filename = strings.TrimPrefix(slashedRelPath, "db/migrations/")
			} else {
				filename = strings.TrimPrefix(slashedRelPath, "migrations/")
			}

			switch arch {
			case "simple":
				destPath = filepath.Join("db/migrations", filename)
			case "clean", "ddd":
				destPath = filepath.Join("infrastructure/database/migrations", filename)
			case "hexagonal":
				destPath = filepath.Join("secondary/database/migrations", filename)
			default:
				destPath = filepath.Join("migrations", filename)
			}

		// Handle queries folder
		case strings.HasPrefix(slashedRelPath, "db/queries/") || strings.HasPrefix(slashedRelPath, "queries/"):
			var filename string
			if strings.HasPrefix(slashedRelPath, "db/queries/") {
				filename = strings.TrimPrefix(slashedRelPath, "db/queries/")
			} else {
				filename = strings.TrimPrefix(slashedRelPath, "queries/")
			}

			switch arch {
			case "simple":
				destPath = filepath.Join("db/queries", filename)
			case "clean", "ddd":
				destPath = filepath.Join("infrastructure/database/query", filename)
			case "hexagonal":
				destPath = filepath.Join("secondary/database/query", filename)
			default:
				destPath = filepath.Join("queries", filename)
			}

		// Specific file handling
		case filepath.Base(path) == "model.go.tmpl":
			destPath = tg.getDestinationPath("models")
		case filepath.Base(path) == "sqlc.yaml.tmpl":
			destPath = "sqlc.yaml"
		case strings.HasSuffix(path, "migrate.go.tmpl"):
			destPath = tg.getDestinationPath("migrate")

		default:
			// Fallback: maintain relative structure if not matched
			// Use the filename if it's at the root of the tool template
			// destPath remains empty to skip
		}

		if destPath == "" {
			// Skip files we don't know how to place
			return nil
		}

		templatePath := strings.TrimPrefix(path, templatesDir+string(filepath.Separator))
		return tg.generateFileFromTemplate(destPath, templatePath)
	})
}

// generateArchitectureFiles generates the architecture-specific files recursively.
func (tg *TemplateGenerator) generateArchitectureFiles() error {
	architecture := string(tg.Config.Architecture)
	if architecture == "" {
		return nil
	}

	templateRootDir := filepath.ToSlash(filepath.Join(templatesDir, "architectures", architecture))

	return fs.WalkDir(scaffolding.Templates, templateRootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		relPath, err := filepath.Rel(templateRootDir, path)
		if err != nil {
			return err
		}

		destPath := relPath
		templatePath := strings.TrimPrefix(path, templatesDir+string(filepath.Separator))

		return tg.generateFileFromTemplate(destPath, templatePath)
	})
}

// generateDevOpsFiles generates the DevOps-specific files recursively.
func (tg *TemplateGenerator) generateDevOpsFiles() error {
	if !tg.Config.DevOps.Enabled {
		return nil
	}

	for _, tool := range tg.Config.DevOps.Tools {
		toolName := strings.ToLower(tool)
		if toolName == helmDir {
			if err := tg.generateHelmChart(); err != nil {
				return fmt.Errorf("failed to generate files for DevOps tool %s: %w", toolName, err)
			}
			continue
		}

		templateRootDir := filepath.ToSlash(filepath.Join(templatesDir, devopsDir, toolName))
		if _, err := scaffolding.Templates.Open(templateRootDir); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return err
		}

		err := fs.WalkDir(scaffolding.Templates, templateRootDir, func(path string, d fs.DirEntry, err_ error) error {
			if err_ != nil {
				return err_
			}
			if d.IsDir() || (!strings.HasSuffix(path, ".tmpl") && !strings.HasSuffix(path, ".tf")) {
				// Note: Added .tf check just in case, though they likely have .tmpl
				return nil
			}

			relPath, err := filepath.Rel(templateRootDir, path)
			if err != nil {
				return err
			}

			destPath := filepath.Join(devopsDir, toolName, relPath)
			templatePath := strings.TrimPrefix(path, templatesDir+string(filepath.Separator))

			if toolName == ansibleDir {
				return tg.generateFileFromTemplate(destPath, templatePath, "<<", ">>")
			}
			return tg.generateFileFromTemplate(destPath, templatePath)
		})

		if err != nil {
			return fmt.Errorf("failed to generate files for DevOps tool %s: %w", toolName, err)
		}
	}
	return nil
}

func (tg *TemplateGenerator) generateHelmChart() error {
	chartFS := scaffolding.Templates
	chartRoot := filepath.Join(templatesDir, devopsDir, helmDir)

	// Walk the embedded chart directory
	err := fs.WalkDir(chartFS, chartRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(chartRoot, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}

		// Skip if directory, but make sure to create the dir in output
		destPath := filepath.Join(tg.OutputDir, devopsDir, helmDir, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			return nil
		}

		// Handle files
		templatePath := strings.TrimPrefix(path, templatesDir+string(filepath.Separator))
		baseName := filepath.Base(path)

		// 1. Chart.yaml.tmpl and values.yaml.tmpl -> Process as Go Template
		if baseName == "Chart.yaml.tmpl" || baseName == "values.yaml.tmpl" {
			return tg.generateFileFromTemplate(filepath.Join(devopsDir, helmDir, relPath), templatePath)
		}

		// 2. All other files -> Copy, perform string replacement for {{.ProjectName}}, and strip .tmpl
		// Read content
		data, err := fs.ReadFile(chartFS, path)
		if err != nil {
			return err
		}
		content := string(data)

		// Replace {{.ProjectName}}
		content = strings.ReplaceAll(content, "{{.ProjectName}}", tg.Config.ProjectName)
		// Handle potential spaced version {{ .ProjectName }} just in case
		content = strings.ReplaceAll(content, "{{ .ProjectName }}", tg.Config.ProjectName)

		// Determine output filename (strip .tmpl)
		finalDestPath := destPath
		if strings.HasSuffix(destPath, ".tmpl") {
			finalDestPath = strings.TrimSuffix(destPath, ".tmpl")
		}

		// Create parent directory if needed (should have been created by IsDir check but just in case)
		if err := os.MkdirAll(filepath.Dir(finalDestPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory for %s: %w", finalDestPath, err)
		}

		// Write content
		if err := os.WriteFile(finalDestPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %w", finalDestPath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scaffold helm chart: %w", err)
	}

	return nil
}

// reportProgress sends a progress update via the callback.
func (tg *TemplateGenerator) reportProgress(step int, message string) {
	if tg.progressCallback != nil {
		tg.progressCallback(step, message)
	}
}

func (tg *TemplateGenerator) reportError(step int, err error) {
	if tg.errorCallback != nil {
		tg.errorCallback(step, err)
	}
}
