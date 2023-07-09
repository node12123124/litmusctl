/*
Copyright © 2021 The LitmusChaos Authors

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
package disconnect

import (
	"fmt"
	"os"
	"strings"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use: "chaos-delegate",
	Short: `Disconnect a Chaos Infrastructure
	Example:
	#disconnect a Chaos Delegate
	litmusctl disconnect chaos-delegate c520650e-7cb6-474c-b0f0-4df07b2b025b --project-id=c520650e-7cb6-474c-b0f0-4df07b2b025b

	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
	Run: func(cmd *cobra.Command, args []string) {

		// Fetch user credentials
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		// Handle blank input for project ID
		if projectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&projectID)

			if projectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		var agentID string
		if len(args) == 0 {
			utils.White_B.Print("\nEnter the Agent ID: ")
			fmt.Scanln(&agentID)
		} else {
			agentID = args[0]
		}
		// Handle blank input for agent ID
		if agentID == "" {
			utils.Red.Println("⛔ Chaos Delegate ID can't be empty!!")
			os.Exit(1)
		}

		// Perform authorization
		userDetails, err := apis.GetProjectDetails(credentials)
		utils.PrintError(err)
		var editAccess = false
		var project apis.Project
		for _, p := range userDetails.Data.Projects {
			if p.ID == projectID {
				project = p
			}
		}
		for _, member := range project.Members {
			if (member.UserID == userDetails.Data.ID) && (member.Role == "Owner" || member.Role == "Editor") {
				editAccess = true
			}
		}
		if !editAccess {
			utils.Red.Println("⛔ User doesn't have edit access to the project!!")
			os.Exit(1)
		}

		// Make API call
		var agentIDs []*string
		agentIDs = append(agentIDs, &agentID)
		disconnectedInfra, err := apis.DisconnectInfra(projectID, agentIDs, credentials)
		if err != nil {
			utils.Red.Println("\n❌ Error in disconnecting Chaos Infrastructure: ", err.Error())
			os.Exit(1)
		}

		if strings.Contains(disconnectedInfra.Data.Message, "infra deleted successfully") {
			utils.White_B.Println("\n🚀 Chaos Infrastructure successfully disconnected.")
		} else {
			utils.White_B.Println("\n❌ Failed to disconnect Chaos Infrastructure. Please check if the ID is correct or not.")
		}
	},
}

func init() {
	DisconnectCmd.AddCommand(agentCmd)

	agentCmd.Flags().String("project-id", "", "Set the project-id to disconnect Chaos Infrastructure for the particular project. To see the projects, apply litmusctl get projects")
}
