package conector

import "context"

type Connector interface {
	Connect(ctx context.Context, dbDsn string) error
	GetData(ctx context.Context, collectionName string) ([]map[string]any, error)
	WriteData(ctx context.Context, collectionName string, data []map[string]any) error
	Close(ctx context.Context) error
}
