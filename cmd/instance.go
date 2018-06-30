// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"terraform-playground/deployer"

	"github.com/spf13/cobra"
)

var instanceProviders []string
var instancePrivateKey string
var instancePublicKey string
var instanceCount int
var regionAws []string
var regionDo []string
var regionAzure []string
var regionGoogle []string
var numberInput string

var instance = &cobra.Command{
	Use:   "instance",
	Short: "instance parent command",
	Long:  `Domain Front Command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'instance --help' for usage.")
	},
}

var instanceDeploy = &cobra.Command{
	//TODO: need to trim spaces
	Use:   "deploy",
	Short: "deploys an instance",
	Long:  `deploys an instance`,
	Args: func(cmd *cobra.Command, args []string) error {
		deployer.InitializeTerraformFiles()
		if deployer.ProviderCheck(instanceProviders) {
			return nil
		}
		return fmt.Errorf("invalid providers specified: %s", instanceProviders)
	},
	Run: func(cmd *cobra.Command, args []string) {

		marshalledOutput := deployer.TerraformOutputMarshaller()
		masterList := deployer.InstanceDeploy(instanceProviders, regionAws, regionDo, regionAzure, regionGoogle, instanceCount, instancePrivateKey, instancePublicKey, marshalledOutput)

		mainFile := deployer.CreateMasterFile(masterList)

		deployer.CreateTerraformMain(mainFile)
	},
}

var instanceDestroy = &cobra.Command{
	Use:   "destroy",
	Short: "destroy",
	Long:  `destroys an instance`,
	Args: func(cmd *cobra.Command, args []string) error {
		if !deployer.IsValidNumberInput(numberInput) {
			return fmt.Errorf("invalid formatting specified: %s", numberInput)
		}
		numsToDestroy := deployer.ExpandNumberInput(numberInput)
		largestInstanceNumToDestroy := deployer.FindLargestNumber(numsToDestroy)

		//get the number of instances actually available in state
		marshalledOutput := deployer.TerraformOutputMarshaller()
		for i := range marshalledOutput.Master.ProviderValues.AWSProvider.Instances {
			awsCount := awsCount + marshalledOutput.Master.ProviderValues.AWSProvider.Instances[i].Config.Count
		}
		for i := range marshalledOutput.Master.ProviderValues.DoProvider.Instances {
			doCount := doCount + marshalledOutput.Master.ProviderValues.DoProvider.Instances[i].Config.Count
		}
		for i := range marshalledOutput.Master.ProviderValues.GoogleProvider.Instances {
			googleCount := googleCount + marshalledOutput.Master.ProviderValues.GoogleProvider.Instances[i].Config.Count
		}
		for i := range marshalledOutput.Master.ProviderValues.AzureProvider.Instances {
			azureCount := azureCount + marshalledOutput.Master.ProviderValues.AzureProvider.Instances[i].Config.Count
		}

		//make sure the largestInstanceNumToDestroy is not bigger than totalInstancesAvailable
		if awsCount < largestInstanceNumToDestroy {
			return error("The number you entered is too big. Try running `list` to see the number of instances you have.")
		}

		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {
		numsToDelete := deployer.ExpandNumberInput(numberInput)
		IPIDList := deployer.GenerateIPIDList()

		//check numstodelete agaisnt ipidlist
		var IDsToDelete []string

		for key, value := IPIDList {
			IDsToDelete = IDsToDelete + value
		}

		for _, id := IDsToDelete {
			//do a terraform state list id = '', and then do terraform destroy target 'name from list' target 'another name'
		}
		//delete those based on array nums
	},
}

var instanceList = &cobra.Command{
	Use:   "list",
	Short: "list instances",
	Long:  `list instances`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list of active instances: ", deployer.GenerateIPIDList())
	},
}

var instanceInfo = &cobra.Command{
	Use:   "info",
	Short: "info",
	Long:  `provides information on specific instance`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Info Called")
	},
}

func init() {
	rootCmd.AddCommand(instance)
	instance.AddCommand(instanceDeploy, instanceDestroy, instanceInfo, instanceList)

	instanceDeploy.PersistentFlags().StringSliceVarP(&instanceProviders, "providers", "p", nil, "list of providers to enter")
	instanceDeploy.MarkPersistentFlagRequired("providers")

	instanceDeploy.PersistentFlags().IntVarP(&instanceCount, "count", "c", 0, "number of instances to deploy")
	instanceDeploy.MarkPersistentFlagRequired("count")

	instanceDeploy.PersistentFlags().StringVarP(&instancePrivateKey, "privatekey", "v", "", "full path to private key to connect to instances")
	instanceDeploy.MarkPersistentFlagRequired("privatekey")

	instanceDeploy.PersistentFlags().StringVarP(&instancePublicKey, "publickey", "b", "", "full path to public key corresponding to the private key")
	instanceDeploy.MarkPersistentFlagRequired("publickey")

	// instanceDestroy.PersistentFlags().IntVarP(&numberInput, "input", "i", 0, "number of instances to destroy")
	// instanceDestroy.MarkPersistentFlagRequired("input")

	//TODO: default all regions
	rootCmd.PersistentFlags().StringSliceVar(&regionAws, "region-aws", []string{"us-east-1", "us-west-2"}, "list of regions for aws. ex: us-east-1,us-west-2,ap-northeast-1")
	rootCmd.PersistentFlags().StringSliceVar(&regionDo, "region-do", []string{"AMS2", "SFO2"}, "list of regions for digital ocean. ex: AMS2,SFO2,NYC1")
	rootCmd.PersistentFlags().StringSliceVar(&regionAzure, "region-azure", []string{"westus", "centralus"}, "list of regions for azure. ex: centralus, eastus, westus")
	rootCmd.PersistentFlags().StringSliceVar(&regionGoogle, "region-google", []string{"us-west1", "us-east1"}, "list of regions for google. ex: us-east1, us-west1, us-central1")

}
