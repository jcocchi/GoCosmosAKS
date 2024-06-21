package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

func connectCosmosClient() *azcosmos.Client {
	clientOptions := azcosmos.ClientOptions{
		EnableContentResponseOnWrite: true,
	}

	cosmosDbEndpoint := os.Getenv("ENDPOINT")
	cosmosDbKey := os.Getenv("PRIMARY_KEY")

	cred, err := azcosmos.NewKeyCredential(cosmosDbKey)
	handle(err)
	client, err := azcosmos.NewClientWithKey(cosmosDbEndpoint, cred, &clientOptions)
	handle(err)

	fmt.Printf("Connected to account\t%s", cosmosDbEndpoint)

	return client
}

func getContainer(client azcosmos.Client) *azcosmos.ContainerClient {
	dbName := os.Getenv("DATABASE_NAME")
	containerName := os.Getenv("CONTAINER_NAME")

	container, err := client.NewContainer(dbName, containerName)
	handle(err)

	fmt.Printf("Connected to container:\t%s", container.ID())

	return container
}

func getAlbumByIdFromCosmos(container azcosmos.ContainerClient, id string) *album {
	context := context.TODO()

	// Container is partitioned by id so pass it for both id and pk
	itemResponse, err := container.ReadItem(context, azcosmos.NewPartitionKeyString(id), id, nil)
	if err != nil {
		var responseErr *azcore.ResponseError
		errors.As(err, &responseErr)

		if responseErr.StatusCode != 404 {
			// TODO: need to add a real error response here. So far not handling any error other than 404
			handle(responseErr)
		} else {
			return nil
		}
	}

	var myAlbum album
	err = json.Unmarshal(itemResponse.Value, &myAlbum)
	handle(err)

	fmt.Printf("Read Item:\t%s", myAlbum)

	return &myAlbum
}

func handle(err error) {
	if err != nil {
		log.Printf(err.Error())
	}
}
