package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

type handler struct {
    client *azcosmos.Client
}

func New(client *azcosmos.Client) handler {
    return handler{client}
}