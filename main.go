package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/brandur/modulir"
)

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Main
//
//
//
//////////////////////////////////////////////////////////////////////////////

func main() {
	rootCmd := &cobra.Command{
		Use:   "mutelight",
		Short: "Mutelight is a static site generator",
		Long: strings.TrimSpace(`
Mutelight is a static site generator for the Mutelight' blog.
See the product in action at https://mutelight.org.`),
	}

	buildCommand := &cobra.Command{
		Use:   "build",
		Short: "Run a single build loop",
		Long: strings.TrimSpace(`
Starts the build loop that watches for local changes and runs
when they're detected. A webserver is started on PORT (default
5009).`),
		Run: func(cmd *cobra.Command, args []string) {
			modulir.Build(getModulirConfig(), build)
		},
	}
	rootCmd.AddCommand(buildCommand)

	loopCommand := &cobra.Command{
		Use:   "loop",
		Short: "Start build and serve loop",
		Long: strings.TrimSpace(`
Runs the build loop one time and places the result in TARGET_DIR
(default ./public/).`),
		Run: func(cmd *cobra.Command, args []string) {
			modulir.BuildLoop(getModulirConfig(), build)
		},
	}
	rootCmd.AddCommand(loopCommand)

	// Make sure to seed the random number generator or else we'll end up with
	// the same random results for every build.
	rand.Seed(time.Now().UnixNano())

	if err := envdecode.Decode(&conf); err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding conf from env: %v", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command: %v", err)
		os.Exit(1)
	}
}

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Variables
//
//
//
//////////////////////////////////////////////////////////////////////////////

// Left as a global for now for the sake of convenience, but it's not used in
// very many places and can probably be refactored as a local if desired.
var conf Conf

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Types
//
//
//
//////////////////////////////////////////////////////////////////////////////

// Conf contains configuration information for the command. It's extracted from
// environment variables.
type Conf struct {
	// AbsoluteURL is the absolute URL where the compiled site will be hosted.
	// It's used for things like Atom feeds.
	AbsoluteURL string `env:"ABSOLUTE_URL,default=https://mutelight.org"`

	// Concurrency is the number of build Goroutines that will be used to
	// perform build work items.
	Concurrency int `env:"CONCURRENCY,default=30"`

	// Drafts is whether drafts of articles and fragments should be compiled
	// along with their published versions.
	//
	// Activating drafts also prompts the creation of a robots.txt to make sure
	// that drafts aren't inadvertently accessed by web crawlers.
	Drafts bool `env:"DRAFTS,default=false"`

	// SorgEnv is the environment to run the app with. Use "development" to
	// activate development features.
	MutelightEnv string `env:"MUTELIGHT_ENV,default=production"`

	// GoogleAnalyticsID is the account identifier for Google Analytics to use.
	GoogleAnalyticsID string `env:"GOOGLE_ANALYTICS_ID"`

	// NumAtomEntries is the number of entries to put in Atom feeds.
	NumAtomEntries int `env:"NUM_ATOM_ENTRIES,default=20"`

	// Port is the port on which to serve HTTP when looping in development.
	Port int `env:"PORT,default=5009"`

	// TargetDir is the target location where the site will be built to.
	TargetDir string `env:"TARGET_DIR,default=./public"`

	// Verbose is whether the program will print debug output as it's running.
	Verbose bool `env:"VERBOSE,default=false"`
}

//////////////////////////////////////////////////////////////////////////////
//
//
//
// Private
//
//
//
//////////////////////////////////////////////////////////////////////////////

const (
	mutelightEnvDevelopment = "development"
)

func getLog() modulir.LoggerInterface {
	log := logrus.New()

	if conf.Verbose {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}

// getModulirConfig interprets Conf to produce a configuration suitable to pass
// to a Modulir build loop.
func getModulirConfig() *modulir.Config {
	return &modulir.Config{
		Concurrency: conf.Concurrency,
		Log:         getLog(),
		LogColor:    terminal.IsTerminal(int(os.Stdout.Fd())),
		Port:        conf.Port,
		SourceDir:   ".",
		TargetDir:   conf.TargetDir,
		Websocket:   conf.MutelightEnv == mutelightEnvDevelopment,
	}
}
