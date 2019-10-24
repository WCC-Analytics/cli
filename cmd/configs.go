/*
Copyright © 2019 Doppler <support@doppler.com>

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
	"doppler-cli/api"
	configuration "doppler-cli/config"
	"doppler-cli/utils"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type configsResponse struct {
	Variables map[string]interface{}
	Success   bool
}

var configsCmd = &cobra.Command{
	Use:   "configs",
	Short: "List configs",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag := utils.GetBoolFlag(cmd, "json")
		localConfig := configuration.LocalConfig(cmd)

		_, configs := api.GetAPIConfigs(cmd, localConfig.Key.Value, localConfig.Project.Value)

		printConfigsInfo(configs, jsonFlag)
	},
}

var configsGetCmd = &cobra.Command{
	Use:   "get [config]",
	Short: "Get info for a config",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag := utils.GetBoolFlag(cmd, "json")
		localConfig := configuration.LocalConfig(cmd)

		config := localConfig.Config.Value
		if len(args) > 0 {
			config = args[0]
		}

		_, configInfo := api.GetAPIConfig(cmd, localConfig.Key.Value, localConfig.Project.Value, config)

		printConfigInfo(configInfo, jsonFlag)
	},
}

var configsCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a config",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag := utils.GetBoolFlag(cmd, "json")
		silent := utils.GetBoolFlag(cmd, "silent")
		defaults := utils.GetBoolFlag(cmd, "defaults")
		environment := cmd.Flag("environment").Value.String()

		name := cmd.Flag("name").Value.String()
		if len(args) > 0 {
			name = args[0]
		}

		localConfig := configuration.LocalConfig(cmd)
		_, info := api.CreateAPIConfig(cmd, localConfig.Key.Value, localConfig.Project.Value, name, environment, defaults)

		if !silent {
			printConfigInfo(info, jsonFlag)
		}
	},
}

var configsDeleteCmd = &cobra.Command{
	Use:   "delete [config]",
	Short: "Delete a config",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO prompt user with a confirmation before proceeding (and add a --yes flag to skip it)
		jsonFlag := utils.GetBoolFlag(cmd, "json")
		silent := utils.GetBoolFlag(cmd, "silent")
		localConfig := configuration.LocalConfig(cmd)

		config := localConfig.Config.Value
		if len(args) > 0 {
			config = args[0]
		}

		api.DeleteAPIConfig(cmd, localConfig.Key.Value, localConfig.Project.Value, config)

		// fetch and display configs
		if !silent {
			_, configs := api.GetAPIConfigs(cmd, localConfig.Key.Value, localConfig.Project.Value)
			printConfigsInfo(configs, jsonFlag)
		}
	},
}

var configsUpdateCmd = &cobra.Command{
	Use:   "update [config]",
	Short: "Update a config",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFlag := utils.GetBoolFlag(cmd, "json")
		silent := utils.GetBoolFlag(cmd, "silent")
		name := cmd.Flag("name").Value.String()
		localConfig := configuration.LocalConfig(cmd)

		config := localConfig.Config.Value
		if len(args) > 0 {
			config = args[0]
		}

		_, info := api.UpdateAPIConfig(cmd, localConfig.Key.Value, localConfig.Project.Value, config, name)

		if !silent {
			printConfigInfo(info, jsonFlag)
		}
	},
}

func init() {
	configsCmd.Flags().String("project", "", "doppler project (e.g. backend)")
	configsCmd.Flags().Bool("json", false, "output json")

	configsGetCmd.Flags().String("project", "", "doppler project (e.g. backend)")
	configsGetCmd.Flags().String("config", "", "doppler config (e.g. dev)")
	configsGetCmd.Flags().Bool("json", false, "output json")
	configsCmd.AddCommand(configsGetCmd)

	configsCreateCmd.Flags().String("project", "", "doppler project (e.g. backend)")
	configsCreateCmd.Flags().String("name", "", "config name")
	configsCreateCmd.Flags().String("environment", "", "config environment")
	configsCreateCmd.Flags().Bool("defaults", true, "populate config with environment's default secrets")
	configsCreateCmd.Flags().Bool("json", false, "output json")
	configsCreateCmd.Flags().Bool("silent", false, "don't output the response")
	configsCreateCmd.MarkFlagRequired("environment")
	configsCmd.AddCommand(configsCreateCmd)

	configsUpdateCmd.Flags().String("project", "", "doppler project (e.g. backend)")
	configsUpdateCmd.Flags().String("config", "", "doppler config (e.g. dev)")
	configsUpdateCmd.Flags().String("name", "", "config name")
	configsUpdateCmd.Flags().Bool("json", false, "output json")
	configsUpdateCmd.Flags().Bool("silent", false, "don't output the response")
	configsUpdateCmd.MarkFlagRequired("name")
	configsCmd.AddCommand(configsUpdateCmd)

	configsDeleteCmd.Flags().String("project", "", "doppler project (e.g. backend)")
	configsDeleteCmd.Flags().String("config", "", "doppler config (e.g. dev)")
	configsDeleteCmd.Flags().Bool("json", false, "output json")
	configsDeleteCmd.Flags().Bool("silent", false, "don't output the response")
	configsCmd.AddCommand(configsDeleteCmd)

	rootCmd.AddCommand(configsCmd)
}

func printConfigsInfo(info []api.ConfigInfo, jsonFlag bool) {
	if jsonFlag {
		resp, err := json.Marshal(info)
		if err != nil {
			utils.Err(err)
		}

		fmt.Println(string(resp))
		return
	}

	var rows [][]string
	for _, configInfo := range info {
		rows = append(rows, []string{configInfo.Name, strings.Join(configInfo.MissingVariables, ", "), configInfo.DeployedAt, configInfo.CreatedAt, configInfo.Environment, configInfo.Project})
	}
	utils.PrintTable([]string{"name", "missing_variables", "deployed_at", "created_at", "stage", "project"}, rows)
}

func printConfigInfo(info api.ConfigInfo, jsonFlag bool) {
	if jsonFlag {
		resp, err := json.Marshal(info)
		if err != nil {
			utils.Err(err)
		}

		fmt.Println(string(resp))
		return
	}

	rows := [][]string{{info.Name, strings.Join(info.MissingVariables, ", "), info.DeployedAt, info.CreatedAt, info.Environment, info.Project}}
	utils.PrintTable([]string{"name", "missing_variables", "deployed_at", "created_at", "stage", "project"}, rows)
}
