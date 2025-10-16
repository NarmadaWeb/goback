// cmd/root.go

package cmd

import (
	"fmt"
	"os"

	"github.com/NarmadaWeb/goback/internal/tui"
	"github.com/NarmadaWeb/goback/pkg/config"
	"github.com/NarmadaWeb/goback/pkg/scaffolding/generator"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goback",
	Short: "GoBack - TUI Backend Project Scaffolding Tool",
	Long: `GoBack is a TUI (Terminal User Interface) built with Bubble Tea
to make it easier for backend developers to create backend projects with various
choices of frameworks, databases, tools, architectures, and DevOps tools.

Use without arguments to open the interactive TUI, or use
subcommands for direct CLI operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior: start TUI
		startTUI()
	},
}

// tuiCmd starts the interactive TUI
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Start interactive TUI interface",
	Long:  `Starts the interactive TUI to create a backend project.`,
	Run: func(cmd *cobra.Command, args []string) {
		startTUI()
	},
}

// newCmd creates a new project via CLI
var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new backend project",
	Long:  `Creates a new backend project with the configuration specified via flags.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		createProjectViaCLI(cmd, args)
	},
}

// listCmd lists available templates and options
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available frameworks, databases, and architectures",
	Long:  `Lists the available frameworks, databases, tools, and architectures.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoBack offers the following options for your project:")

		// Frameworks
		fmt.Println("\n🏗️ Frameworks:")
		printChoiceList(config.GetValidFrameworks())

		// Databases
		fmt.Println("\n🗄️ Databases:")
		printChoiceList(config.GetValidDatabases())

		// Tools
		fmt.Println("\n🔗 Tools:")
		printChoiceList(config.GetValidTools())

		// Architectures
		fmt.Println("\n🏛️ Architectures:")
		printChoiceList(config.GetValidArchitectures())

		// DevOps Tools
		fmt.Println("\n🚀 DevOps Tools:")
		printChoiceList(config.GetValidDevOpsTools())
	},
}

// A helper interface to work with different choice types
type choice interface {
	String() string
	Description() string
}

// A wrapper for string slices to satisfy the choice interface
type stringChoice string

func (s stringChoice) String() string {
	return string(s)
}
func (s stringChoice) Description() string {
	return config.GetDevOpsToolDescription(string(s))
}

func printChoiceList[T any](items []T) {
	for _, item := range items {
		var choiceItem choice
		switch v := any(item).(type) {
		case config.FrameworkChoice:
			choiceItem = v
		case config.DatabaseChoice:
			choiceItem = v
		case config.ToolChoice:
			choiceItem = v
		case config.ArchitectureChoice:
			choiceItem = v
		case string:
			choiceItem = stringChoice(v)
		default:
			continue
		}
		// Print with padding for alignment
		fmt.Printf("  - %-25s %s\n", choiceItem.String(), choiceItem.Description())
	}
}

// configCmd manages configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage GoBack configuration",
	Long:  `Manages GoBack configuration, including default settings and preferences.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetConfig()
		fmt.Printf("Configuration file: %s\n", viper.ConfigFileUsed())
		fmt.Printf("Default output directory: %s\n", cfg.DefaultOutputDir)
		fmt.Printf("Default module prefix: %s\n", cfg.DefaultModulePrefix)
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]
		viper.Set(key, value)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("Error saving configuration: %v\n", err)
			return
		}
		fmt.Printf("Set %s = %s\n", key, value)
	},
}

// versionCmd shows version information
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("GoBack v0.1.7")
		fmt.Println("TUI Backend Project Scaffolding Tool")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goback.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")

	// Add subcommands
	rootCmd.AddCommand(tuiCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)

	// Config subcommands
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)

	// New command flags
	newCmd.Flags().StringP("framework", "f", "", "Framework to use (fiber, gin, chi, echo)")
	newCmd.Flags().StringP("database", "d", "", "Database to use (postgresql, mysql, sqlite)")
	newCmd.Flags().StringP("tool", "t", "", "Tool to use (sqlx, sqlc)")
	newCmd.Flags().StringP("architecture", "a", "", "Architecture pattern (simple, ddd, clean, hexagonal)")
	newCmd.Flags().StringP("output", "O", "", "Output directory")
	newCmd.Flags().StringP("module", "m", "", "Go module path")
	newCmd.Flags().Bool("devops", false, "Include DevOps configurations")
	newCmd.Flags().StringSlice("devops-tools", []string{},
		"DevOps tools to include (helm, terraform, ansible)")

	// Bind flags to viper
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("framework", newCmd.Flags().Lookup("framework"))
	_ = viper.BindPFlag("database", newCmd.Flags().Lookup("database"))
	_ = viper.BindPFlag("tool", newCmd.Flags().Lookup("tool"))
	_ = viper.BindPFlag("architecture", newCmd.Flags().Lookup("architecture"))
	_ = viper.BindPFlag("output", newCmd.Flags().Lookup("output"))
	_ = viper.BindPFlag("module", newCmd.Flags().Lookup("module"))
	_ = viper.BindPFlag("devops", newCmd.Flags().Lookup("devops"))
	_ = viper.BindPFlag("devops-tools", newCmd.Flags().Lookup("devops-tools"))
}

// initConfig reads in config file and ENV variables.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".goback" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".goback")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}

	// Initialize default configuration
	config.InitDefaults()
}

// startTUI initializes and runs the TUI application
func startTUI() {
	model := tui.NewMainModel()
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v\n", err)
		os.Exit(1)
	}
}

// createProjectViaCLI creates a project using CLI flags
func createProjectViaCLI(cmd *cobra.Command, args []string) {
	var projectName string
	if len(args) > 0 {
		projectName = args[0]
	} else {
		fmt.Println("Error: project name is required")
		fmt.Println("Usage: goback new [project-name]")
		os.Exit(1)
	}

	// Get flags
	framework, _ := cmd.Flags().GetString("framework")
	database, _ := cmd.Flags().GetString("database")
	tool, _ := cmd.Flags().GetString("tool")
	architecture, _ := cmd.Flags().GetString("architecture")
	output, _ := cmd.Flags().GetString("output")
	module, _ := cmd.Flags().GetString("module")
	devops, _ := cmd.Flags().GetBool("devops")
	devopsTools, _ := cmd.Flags().GetStringSlice("devops-tools")

	// Set module and output dir if not provided
	if module == "" {
		module = "github.com/user/" + projectName
	}
	if output == "" {
		output = "./" + projectName
	}

	// Create configuration
	cfg := &config.ProjectConfig{
		ProjectName:  projectName,
		ModulePath:   module,
		Description:  fmt.Sprintf("%s backend API", projectName),
		OutputDir:    output,
		Framework:    config.FrameworkChoice(framework),
		Database:     config.DatabaseChoice(database),
		Tool:         config.ToolChoice(tool),
		Architecture: config.ArchitectureChoice(architecture),
		DevOps: config.DevOpsConfig{
			Enabled: devops,
			Tools:   devopsTools,
		},
	}

	// Validate configuration
	if validationErrors := config.ValidateProjectConfig(cfg); len(validationErrors) > 0 {
		fmt.Println("❌ Configuration validation failed:")
		for _, err := range validationErrors {
			fmt.Printf("  - %s\n", err)
		}
		fmt.Println("\nPlease provide all required flags: --framework, --database, --tool, --architecture")
		os.Exit(1)
	}

	// Generate project
	fmt.Printf("Creating project '%s'...\n", projectName)
	gen := generator.NewTemplateGenerator(cfg)

	gen.SetProgressCallback(func(step int, message string) {
		fmt.Printf("  %s\n", message)
	})

	if err := gen.Generate(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Project '%s' created successfully!\n", projectName)
	fmt.Printf("Next steps:\n")
	fmt.Printf("  cd %s\n", output)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go run main.go\n")
}
