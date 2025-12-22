// cmd/e2e_test.go

package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestCreateProjectViaCLI(t *testing.T) {
	// Setup temporary directory for test
	tempDir, err := os.MkdirTemp("", "goback-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	projectName := "testproject"
	outputDir := filepath.Join(tempDir, projectName)

	// Mock cobra command
	testCmd := &cobra.Command{}
	testCmd.Flags().String("framework", "", "")
	testCmd.Flags().String("database", "", "")
	testCmd.Flags().String("tool", "", "")
	testCmd.Flags().String("architecture", "", "")
	testCmd.Flags().String("output", "", "")
	testCmd.Flags().String("module", "", "")
	testCmd.Flags().Bool("devops", false, "")
	testCmd.Flags().StringSlice("devops-tools", []string{}, "")

	testCmd.Flags().Set("framework", "gin")
	testCmd.Flags().Set("database", "postgresql")
	testCmd.Flags().Set("tool", "sqlx")
	testCmd.Flags().Set("architecture", "simple")
	testCmd.Flags().Set("output", outputDir)
	testCmd.Flags().Set("module", "github.com/test/project")
	testCmd.Flags().Set("devops", "true")
	testCmd.Flags().Set("devops-tools", "helm")

	// WARNING: createProjectViaCLI calls os.Exit(1) on failure.
	createProjectViaCLI(testCmd, []string{projectName})

	// Verify files exist
	expectedFiles := []string{
		"go.mod",
		"cmd/api/main.go",
		"internal/database/connection.go", // sqlx connection
		"devops/helm/Chart.yaml", // helm
		"devops/helm/values.yaml",
		"devops/helm/templates/deployment.yaml",
		"db/migrations", // should exist for sqlx
	}

	for _, relPath := range expectedFiles {
		fullPath := filepath.Join(outputDir, relPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", relPath)
		}
	}
}
