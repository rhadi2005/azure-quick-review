// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package azqr

import (
	"github.com/cmendible/azqr/internal/scanners"
	"github.com/cmendible/azqr/internal/scanners/ci"
	"github.com/spf13/cobra"
)

func init() {
	scanCmd.AddCommand(ciCmd)
}

var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "Scan Azure Container Instances",
	Long:  "Scan Azure Container Instances",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		serviceScanners := []scanners.IAzureScanner{
			&ci.ContainerInstanceScanner{},
		}

		scan(cmd, serviceScanners)
	},
}
