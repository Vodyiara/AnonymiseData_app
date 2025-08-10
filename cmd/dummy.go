/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"anonymise/model"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/spf13/cobra"
)

// dummyCmd represents the dummy command
var dummyCmd = &cobra.Command{
	Use:   "dummy",
	Short: "Add dummy config file",
	Long:  `Create test config file to be filled in by user`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dummy called")
		config := &model.Config{
			SourceDatabaseDSN:         "DSN connect to DB",
			SourceCollectionName:      "collection name to data from DB",
			SourceDatabaseName:        "db name to take data from",
			DestinationDatabaseDSN:    "DS to connect destination DB",
			DestinationCollectionName: "collection name to put data to",
			DestinationDatabaseName:   "db name to put data to",
			FieldToAnonymise:          "mane of the column to anonymise",
		}

		yamlStr, err := yaml.Marshal(config)
		if err != nil {
			fmt.Println(err)
			return
		}

		configFile, err := os.Create("dummy-config.yaml")
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = configFile.WriteString(string(yamlStr))
		if err != nil {
			fmt.Println(err)
			return
		}

		err = configFile.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("dummy config file succesfuly created")

	},
}

func init() {
	rootCmd.AddCommand(dummyCmd)

}
