// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package cosmos

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos"
	"github.com/cmendible/azqr/internal/scanners"
)

// CosmosDBScanner - Scanner for CosmosDB Databases
type CosmosDBScanner struct {
	config              *scanners.ScannerConfig
	diagnosticsSettings scanners.DiagnosticsSettings
	databasesClient     *armcosmos.DatabaseAccountsClient
	listDatabasesFunc   func(resourceGroupName string) ([]*armcosmos.DatabaseAccountGetResults, error)
}

// Init - Initializes the CosmosDBScanner
func (a *CosmosDBScanner) Init(config *scanners.ScannerConfig) error {
	a.config = config
	var err error
	a.databasesClient, err = armcosmos.NewDatabaseAccountsClient(config.SubscriptionID, config.Cred, config.ClientOptions)
	if err != nil {
		return err
	}
	a.diagnosticsSettings = scanners.DiagnosticsSettings{}
	err = a.diagnosticsSettings.Init(config)
	if err != nil {
		return err
	}
	return nil
}

// Scan - Scans all CosmosDB Databases in a Resource Group
func (c *CosmosDBScanner) Scan(resourceGroupName string, scanContext *scanners.ScanContext) ([]scanners.AzureServiceResult, error) {
	log.Printf("Scanning CosmosDB Databases in Resource Group %s", resourceGroupName)

	databases, err := c.listDatabases(resourceGroupName)
	if err != nil {
		return nil, err
	}
	engine := scanners.RuleEngine{}
	rules := c.GetRules()
	results := []scanners.AzureServiceResult{}

	for _, database := range databases {
		rr := engine.EvaluateRules(rules, database, scanContext)

		results = append(results, scanners.AzureServiceResult{
			SubscriptionID: c.config.SubscriptionID,
			ResourceGroup:  resourceGroupName,
			ServiceName:    *database.Name,
			Type:           *database.Type,
			Location:       *database.Location,
			Rules:          rr,
		})
	}
	return results, nil
}

func (c *CosmosDBScanner) listDatabases(resourceGroupName string) ([]*armcosmos.DatabaseAccountGetResults, error) {
	if c.listDatabasesFunc == nil {
		pager := c.databasesClient.NewListByResourceGroupPager(resourceGroupName, nil)

		domains := make([]*armcosmos.DatabaseAccountGetResults, 0)
		for pager.More() {
			resp, err := pager.NextPage(c.config.Ctx)
			if err != nil {
				return nil, err
			}
			domains = append(domains, resp.Value...)
		}
		return domains, nil
	}

	return c.listDatabasesFunc(resourceGroupName)
}
