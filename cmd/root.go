/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"anonymise/conector"
	"anonymise/model"
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"time"
)

var cfgFile string
var nonStringAnonymisationFieldError = errors.New("cannot anonymise anything but string field")
var anonymisationFieldNotFoundError = errors.New("cannot anonymise anything but string field")

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "anonymise",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("anonymise started")
		cfg, err := initConfig()
		if err != nil {
			logrus.Error(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		destinationCtxt, destinationCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer destinationCancel()

		data, err := getDataFromSource(ctx, cfg)
		if err != nil {
			logrus.Error(err)
		}

		data, err = anonymise(data, cfg)
		if err != nil {
			logrus.Error(err)
		}

		err = writeDataToDestination(destinationCtxt, cfg, data)
		if err != nil {
			logrus.Error(err)
		}

		logrus.Info("Successfully coping done")

	},
}

func anonymise(data []map[string]any, cfg *model.Config) ([]map[string]any, error) {
	if len(cfg.FieldToAnonymise) == 0 {
		return data, nil
	}
	anonymisedData := make([]map[string]any, cap(data))
	for idx, element := range data {
		anonymisedData[idx] = make(map[string]any)
		var foundField = false
		for key, val := range element {
			if key == cfg.FieldToAnonymise {
				foundField = true
				switch val.(type) {
				default:
					return nil, nonStringAnonymisationFieldError
				case string:
					anonymisedData[idx][key] = "Anonymised"
				}
			} else {
				anonymisedData[idx][key] = val
			}
		}
		if !foundField {
			return nil, anonymisationFieldNotFoundError
		}
	}
	return anonymisedData, nil
}

func getDataFromSource(ctx context.Context, cfg *model.Config) ([]map[string]any, error) {
	var sourceConnector conector.Connector
	if strings.Contains(cfg.SourceDatabaseDSN, "postgresql://") || strings.Contains(cfg.DestinationDatabaseDSN, "postgres://") {
		sourceConnector = &conector.PostgresConnection{}
	} else if strings.Contains(cfg.SourceDatabaseDSN, "mongodb://") || strings.Contains(cfg.DestinationDatabaseDSN, "mongodb+srv://") {
		sourceConnector = &conector.MongoConnector{
			DBName: cfg.SourceDatabaseName,
		}
	} else {
		return nil, errors.New(cfg.SourceDatabaseDSN + " is not supported")
	}
	err := sourceConnector.Connect(ctx, cfg.SourceDatabaseDSN)
	if err != nil {
		return nil, err
	}
	data, err := sourceConnector.GetData(ctx, cfg.SourceCollectionName)
	if err != nil {
		return nil, err
	}

	err = sourceConnector.Close(ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func writeDataToDestination(ctx context.Context, cfg *model.Config, data []map[string]any) error {
	var destinationConnector conector.Connector

	if strings.Contains(cfg.DestinationDatabaseDSN, "postgresql://") || strings.Contains(cfg.DestinationDatabaseDSN, "postgres://") {
		destinationConnector = &conector.PostgresConnection{}
	} else if strings.Contains(cfg.DestinationDatabaseDSN, "mongodb://") || strings.Contains(cfg.DestinationDatabaseDSN, "mongodb+srv://") {
		destinationConnector = &conector.MongoConnector{
			DBName: cfg.SourceDatabaseName,
		}
	} else {
		return errors.New(cfg.DestinationDatabaseDSN + " is not supported")
	}

	err := destinationConnector.Connect(ctx, cfg.DestinationDatabaseDSN)
	if err != nil {
		return err
	}

	err = destinationConnector.WriteData(ctx, cfg.DestinationCollectionName, data)
	if err != nil {
		return err
	}

	err = destinationConnector.Close(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is ./config.yaml)")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() (*model.Config, error) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Search config in home directory with name ".anonymise" (without extension).
		viper.AddConfigPath("./")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	var cfg model.Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil

}
