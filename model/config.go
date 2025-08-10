package model

type Config struct {
	SourceDatabaseDSN         string `yaml:"source_database_dsn" mapstructure:"source_database_dsn"`
	SourceCollectionName      string `yaml:"source_collection_name" mapstructure:"source_collection_name"`
	SourceDatabaseName        string `yaml:"source_database_name" mapstructure:"source_database_name"`
	DestinationDatabaseDSN    string `yaml:"destination_database_dsn" mapstructure:"destination_database_dsn"`
	DestinationCollectionName string `yaml:"destination_collection_name" mapstructure:"destination_collection_name"`
	DestinationDatabaseName   string `yaml:"destination_database_name" mapstructure:"destination_database_name"`
	FieldToAnonymise          string `yaml:"field_to_anonymise" mapstructure:"field_to_anonymise"`
}
