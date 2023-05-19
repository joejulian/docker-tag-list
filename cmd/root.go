/*
Copyright © 2023 Joe Julian <me@joejulian.name>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	repo "github.com/joejulian/docker-tag-list/pkg/repository"
)

var (
	cfgFile    string
	repository string
	output     string
	constraint string
	latest     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "docker-tag-list",
	Short: "Retrieves a list of Docker tag names",
	Long:  `docker-tag-list returns a list of Docker tag names for a given repository`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		log := log.Default()
		tags, err := repo.QueryImageRepositoryTags(repository)
		if err != nil {
			log.Fatalf("Error querying repository tags: %v", err)
		}

		if constraint != "" {
			t, err := filterTagsOnConstraints(tags, constraint)
			if err != nil {
				log.Fatalf("Error filtering tags on constraints %v", err)
			}
			tags = t
		}

		if latest {
			fmt.Println(findHighestSemverTag(tags))
			os.Exit(0)
		}

		var result string
		switch output {
		case "json":
			marshalled, err := json.Marshal(tags)
			if err != nil {
				log.Fatalf("Error marshalling tags: %v", err)
			}
			result = string(marshalled)
		default:
			result = fmt.Sprintf("tags: %s\n", strings.Join(tags, ", "))
		}
		fmt.Println(result)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.docker-tag-list.yaml)")

	rootCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "repository name to list tags from")
	rootCmd.MarkPersistentFlagRequired("repository")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output format [string, json]")
	rootCmd.PersistentFlags().StringVarP(&constraint, "constraint", "c", "", "filter on semver constraints, ie. '>= 1.2.3','~1.3', etc.")
	rootCmd.PersistentFlags().BoolVar(&latest, "latest", false, "return only the latest version. If constraints are specified, only the latest version that matches the constraints.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".docker-tag-list" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".docker-tag-list")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func filterTagsOnConstraints(tags []string, constraints string) ([]string, error) {
	var newTags []string
	for _, tag := range tags {
		c, err := semver.NewConstraint(constraints)
		if err != nil {
			return nil, err
		}
		v, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}
		if ok, _ := c.Validate(v); !ok {
			continue
		}
		newTags = append(newTags, tag)
	}
	return newTags, nil
}

func findHighestSemverTag(tags []string) string {
	// Find the highest semver tag
	var highestSemver *semver.Version
	for _, tag := range tags {
		version, err := semver.NewVersion(tag)
		if err != nil {
			continue
		}
		if highestSemver == nil || version.GreaterThan(highestSemver) {
			highestSemver = version
		}
	}

	return highestSemver.Original()
}
