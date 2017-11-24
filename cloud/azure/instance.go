// Package aws contains methods querying certain Azure, host-specific data.
// More documentation on retriving metadata information from the 
// Azure instance can be found at:
// https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service

package azure

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"io/ioutil"
)

const (
	INSTANCE_METADATA_URL = "http://169.254.169.254/metadata/instance"
	API_VERSION_PARAMETER = "api-version=2017-04-02"
)

// Makes a request to the given url adding the required header data to
// query a Azure instance. 
// It returns the response text, status code and any error encountered.
func makeRequest(url string) (string, int, error) {

	// Construct an http Client with an acceptable timeout
	client := http.Client{Timeout: time.Second * 5}

	// Create a new http GET request with the given URL
	req, err := http.NewRequest("GET", url, nil)

	// Set the headers required to make a GET request.
	req.Header.Set("Metadata", "true")

	// Make the request
	resp, err := client.Do(req) 

	if err != nil {
		return "", 0, err
	}

	if resp != nil {
		defer resp.Body.Close()

		// Converts the response body bytes[] to a string
		// and returns it
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", resp.StatusCode, err
		}
		return string(body), resp.StatusCode, err
	} else {
		return "", 0, err
	}
}

func HasMetadataHost() bool {
	respText, statusCode, err := makeRequest(
		fmt.Sprintf("%s/?%s", INSTANCE_METADATA_URL, API_VERSION_PARAMETER))

	if err != nil {
		return false
	}
	if statusCode == http.StatusNotFound {
		return false
	}
	fmt.Printf("Got HasMetadataHost %s", respText)
	return true
}

func GetInstanceData() {
	respText, _, _ := makeRequest(
		fmt.Sprintf("%s/?%s", INSTANCE_METADATA_URL, API_VERSION_PARAMETER))

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(respText), &data); err != nil {
		fmt.Println(err)
	}
	fmt.Println(respText)

	/*
	if _, ok := data["imageId"]; ok {
		fmt.Printf("Got %s", data["imageId"].(string))
		return data["imageId"].(string), nil
	}
	*/
}

// Gets the instance-ID of the host
func Id() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/compute/vmId?%s&format=text", 
		INSTANCE_METADATA_URL, API_VERSION_PARAMETER))
	if err != nil {
		return "", err
	}
	fmt.Printf("Got Id %s", respText)
	return respText, nil
}
