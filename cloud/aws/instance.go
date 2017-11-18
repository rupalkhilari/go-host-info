// Package aws contains methods querying certain AWS, host-specific data.
// More documentation on retriving metadata information from the 
// AWS instance can be found at:
// http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html

package aws

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	INSTANCE_METADATA_URL = "http://169.254.169.254/latest/meta-data/"
	INSTANCE_DYNAMIC_DATA_URL = "http://169.254.169.254/latest/dynamic/"
)

// Makes a request to the given url adding the required header data to
// query a AWS instance. 
// It returns the response text, status code and any error encountered.
func makeRequest(url string) (string, int, error) {

	// Construct an http Client with an acceptable timeout
	client := &http.Client{Timeout: time.Second * 5}

	// Create a new http GET request with the given URL
	req, err := http.NewRequest("GET", url, nil)

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

/* ========================== COMMON METHODS ============================ */

// Ping the metadata server to check if it is reachable.
// This server is only accessible from a GCP instances, hence
// acts as a good indication of whether the host running this
// function is a GCP instance.
func HasMetadataHost() bool {
	respText, statusCode, err := makeRequest(INSTANCE_METADATA_URL)
	if err != nil {
		return false
	}
	if statusCode == http.StatusNotFound {
		return false
	}
	fmt.Println("Got %s", respText)
	return true
}

// Gets the FQDN or local hostname of the instance.
func FQDN() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/local-hostname", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Println("Got %s", respText)
	return respText, err
}

// Returns data of the form "x.x.x.x.bc.googleusercontent.com"
func PublicHostname() (string, error) {

	respText, _, err := makeRequest(fmt.Sprintf("%s/public-hostname", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Println("Got %s", respText)
	return respText, err
}

// Gets the local hostname 
func Hostname() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/hostname", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Println("Got %s", respText)
	return respText, nil
}

// Gets the Public/External IP Address 
func PublicIPAddress() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/public-ipv4", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Println("Got %s", respText)
	return respText, nil
}

// Gets the Local IP Address of the host.
func LocalIPAddress() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/local-ipv4", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Println("Got %s", respText)
	return respText, nil
}

// Gets the instance-ID of the host
func Id() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/instance-id", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Println("Got %s", respText)
	return respText, nil
}

// Gets the zone name of the host
func Zone() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/availability-zone", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Println("Got %s", respText)
	return respText, nil
}

// Gets the base machine type of the instance
func Type() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/instance-type", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Println("Got %s", respText)
	return respText, nil	
}

/* ========================== INFORMATION UNIQUE TO AWS ======================= */

// Get the VM's image ID
func ImageId() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/instance-identity/document", INSTANCE_DYNAMIC_DATA_URL))

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(respText), &data); err != nil {
		return "", err
	}

	if _, ok := data["imageId"]; ok {
		return data["imageId"].(string), nil
	}
	
	return "", err
}

