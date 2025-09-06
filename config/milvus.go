package config

import (
	"context"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/spf13/viper"
)

var (
	MilvusClient client.Client
)

func InitMilvus() {
	ctx := context.Background()
	// 链接
	cli, err := client.NewClient(ctx, client.Config{
		Address:  viper.GetString("milvus.address"),
		Username: viper.GetString("milvus.username"),
		Password: viper.GetString("milvus.password"),
	})

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	version, err := cli.GetVersion(ctx)
	if err != nil {
		log.Fatalf("Failed to get version: %v", err)
	}

	log.Printf("Successfully connected to Milvus. Version: %v\n", version)

	MilvusClient = cli
}
