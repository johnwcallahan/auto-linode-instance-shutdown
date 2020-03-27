package main

import (
	"context"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"

	"log"
	"net/http"
	"os"
)

func main() {
	apiKey, ok := os.LookupEnv("LINODE_TOKEN")
	if !ok {
		log.Fatal("Could not find LINODE_TOKEN, please assert it is set.")
	}
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{
			Source: tokenSource,
		},
	}

	linodeClient := linodego.NewClient(oauth2Client)
	linodeClient.SetDebug(false)

	res, err := linodeClient.ListInstances(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	shutdownCount := 0

	for _, instance := range res {
		// Persisted Linodes should have a "persist" or "secure" tag.
		shouldPersist := contains(instance.Tags, "persist") || contains(instance.Tags, "secure")
		// Very brittle, but the best we can do.
		isLKEInstance := instance.Label[0:3] == "lke"

		shouldShutDown := !shouldPersist && !isLKEInstance && instance.Status == "running"

		if shouldShutDown {
			err := linodeClient.ShutdownInstance(context.Background(), instance.ID)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Shutting down instance \"%s\" | ID: %d\n", instance.Label, instance.ID)
			shutdownCount++
		}
	}

	log.Printf("Shut down %d of %d instances", shutdownCount, len(res))
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
