// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package st

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/cmendible/azqr/internal/scanners"
)

// StorageScanner - Scanner for Storage
type StorageScanner struct {
	config              *scanners.ScannerConfig
	diagnosticsSettings scanners.DiagnosticsSettings
	storageClient       *armstorage.AccountsClient
	listStorageFunc     func(resourceGroupName string) ([]*armstorage.Account, error)
}

// Init - Initializes the StorageScanner
func (c *StorageScanner) Init(config *scanners.ScannerConfig) error {
	c.config = config
	var err error

	c.storageClient, err = armstorage.NewAccountsClient(config.SubscriptionID, config.Cred, config.ClientOptions)
	if err != nil {
		return err
	}
	c.diagnosticsSettings = scanners.DiagnosticsSettings{}
	err = c.diagnosticsSettings.Init(config)
	if err != nil {
		return err
	}
	return nil
}

// Scan - Scans all Storage in a Resource Group
func (c *StorageScanner) Scan(resourceGroupName string, scanContext *scanners.ScanContext) ([]scanners.AzureServiceResult, error) {
	log.Printf("Scanning Storage in Resource Group %s", resourceGroupName)

	storage, err := c.listStorage(resourceGroupName)
	if err != nil {
		return nil, err
	}
	engine := scanners.RuleEngine{}
	rules := c.GetRules()
	results := []scanners.AzureServiceResult{}

	for _, storage := range storage {
		rr := engine.EvaluateRules(rules, storage, scanContext)

		results = append(results, scanners.AzureServiceResult{
			SubscriptionID: c.config.SubscriptionID,
			ResourceGroup:  resourceGroupName,
			ServiceName:    *storage.Name,
			Type:           *storage.Type,
			Location:       *storage.Location,
			Rules:          rr,
		})
	}
	return results, nil
}

func (c *StorageScanner) listStorage(resourceGroupName string) ([]*armstorage.Account, error) {
	if c.listStorageFunc == nil {
		pager := c.storageClient.NewListByResourceGroupPager(resourceGroupName, nil)

		staccounts := make([]*armstorage.Account, 0)
		for pager.More() {
			resp, err := pager.NextPage(c.config.Ctx)
			if err != nil {
				return nil, err
			}
			staccounts = append(staccounts, resp.Value...)
		}
		return staccounts, nil
	}

	return c.listStorageFunc(resourceGroupName)
}
