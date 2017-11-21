// Package gcp contains methods querying certain GCP, host-specific data.
// More documentation on retriving metadata information from the 
// GCP instance can be found at:
// https://cloud.google.com/compute/docs/storing-retrieving-metadata

package gcp

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	INSTANCE_METADATA_URL = "http://metadata.google.internal/computeMetadata/v1/instance"
)

// Makes a request to the given url adding the required header data to
// query a GCP instance. 
// It returns the response text, status code and any error encountered.
func makeRequest(url string) (string, int, error) {

	// Construct an http Client with an acceptable timeout
	client := http.Client{Timeout: time.Second * 5}

	// Create a new http GET request with the given URL
	req, err := http.NewRequest("GET", url, nil)

	// Set the headers required to make a GET request.
	req.Header.Set("Metadata-Flavor", "Google")

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

// Make a dummy request to the metadata server to check if it is reachable.
// This server is only accessible from GCP instances, hence
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

// Gets the FQDN of the instance.
func FQDN() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/hostname", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Println("Got %s", respText)
	return respText, err
}

// Returns data of the form "x.x.x.x.bc.googleusercontent.com"
func PublicHostname() (string, error) {
	// perform a reverse NSLookup on the IP address.
	address, err := PublicIPAddress()
	if err != nil {
		return "", err
	}

	names, err := net.LookupAddr(address)
	if len(names) == 0 {
		return "", err
	}
	fmt.Println("Got %s", names[0])
	return names[0], nil
}

// Gets the local hostname 
func Hostname() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/name", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Printf("Got %s", respText)
	return respText, nil
}

// Gets the Public/External IP Address 
func PublicIPAddress() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/network-interfaces/0/access-configs/0/external-ip",INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Printf("Got %s", respText)
	return respText, nil
}

// Gets the Local IP Address of the host.
func LocalIPAddress() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/network-interfaces/0/ip",
		INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Printf("Got %s", respText)
	return respText, nil
}

// Gets the instance-ID of the host
func Id() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/id", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Printf("Got %s", respText)
	return respText, nil
}

// Gets the zone name of the host
func Zone() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/zone", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Printf("Got %s", respText)
	return respText, nil
}

// Gets the base machine type of the host.
func Type() (string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/machine-type", INSTANCE_METADATA_URL))
	if err != nil {
		return "", err
	}
	fmt.Printf("Got %s", respText)
	return respText, nil
}

/* ========================== INFORMATION UNIQUE TO GCP ==================== */

// Gets the Pre-emptable setting
func IsPreemptible() (bool, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/scheduling/preemptible", INSTANCE_METADATA_URL))
	if err != nil {
		return false, err
	}
	fmt.Printf("Got %s", respText)
	if respText == "TRUE" {
		return true, nil
	}
	return false, nil
}

// Gets the tags attached to this instance
func Tags() ([]string, error) {
	respText, _, err := makeRequest(fmt.Sprintf("%s/tags?alt=text", INSTANCE_METADATA_URL))

	if err != nil {
		return []string{}, err
	}
	fmt.Printf("Got %s", respText)
	// Parse lines
	tags := strings.Split(respText, "\n")
	return tags, nil
}
