package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bahelit/confirmerator/database"
	"github.com/bahelit/pester"
)

var (
	mongoClient *mongo.Client
	cmcKey      string
)

// NOTE can make a request to CMC every 5 minutes without breaching the free limit.
const (
	cmcAPIKey = "CMC_API_KEY"
)

func init() {
	var statusOK bool

	cmcKey, statusOK = os.LookupEnv(cmcAPIKey)
	if !statusOK {
		cmcKey = ""
	}
}

func main() {
	var err error

	client := pester.New()
	req := createCMCRequest()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))

	mongoClient, err = database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to mongodb", err)
	}
	defer func() {
		err := mongoClient.Disconnect(context.Background())
		if err != nil {
			log.Printf("ERROR: failed to disconnect from mongo: %v", err)
		}
	}()
}

func createCMCRequest() (req *http.Request) {
	req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		log.Printf("failed to create http request error: %v", err)
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "1500")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", cmcKey)
	req.URL.RawQuery = q.Encode()

	return req
}
