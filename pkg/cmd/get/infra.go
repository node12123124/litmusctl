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
package get

import (
	"fmt"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"os"
	"text/tabwriter"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// InfraCmd represents the agents command
var InfraCmd = &cobra.Command{
	Use:   "chaos-Infrastructures",
	Short: "Display list of Chaos Infrastructures within the project",
	Long:  `Display list of Chaos Infrastructures within the project`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if projectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&projectID)

			for projectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		infras, err := apis.GetInfraList(credentials, projectID, models.ListInfraRequest{})
		utils.PrintError(err)

		output, err := cmd.Flags().GetString("output")
		utils.PrintError(err)

		switch output {
		case "json":
			utils.PrintInJsonFormat(infras.Data)

		case "yaml":
			utils.PrintInYamlFormat(infras.Data)

		case "":

			writer := tabwriter.NewWriter(os.Stdout, 4, 8, 1, '\t', 0)
			utils.White_B.Fprintln(writer, "CHAOS INFRASTRUCTURE ID \tCHAOS INFRASTRUCTURE NAME\tSTATUS\t")

			for _, infra := range infras.Data.ListInfraDetails.Infras {
				var status string
				if infra.IsActive {
					status = "ACTIVE"
				} else {
					status = "INACTIVE"
				}

				//var isRegistered string
				//if agent.IsRemoved {
				//	isRegistered = "REGISTERED"
				//} else {
				//	isRegistered = "NOT REGISTERED"
				//}
				utils.White.Fprintln(writer, infra.InfraID+"\t"+infra.Name+"\t"+status+"\t")
				//+isRegistered+"\t"
			}
			writer.Flush()
		}
	},
}

func init() {
	GetCmd.AddCommand(InfraCmd)

	InfraCmd.Flags().String("project-id", "", "Set the project-id. To retrieve projects. Apply `litmusctl get projects`")

	InfraCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
