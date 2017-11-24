// Package aws contains methods querying certain Azure, host-specific data.
// More documentation on retriving metadata information from the 
// Azure instance can be found at:
// https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service

package azure

import (
	"fmt"
	"net"
	"net/http"
	"time"
	"io/ioutil"
)

const (
	INSTANCE_METADATA_URL string = "http://169.254.169.254/metadata/instance"
	API_VERSION string = "api-version=2017-04-02"
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

/* ========================== COMMON METHODS ============================ */

// Ping the metadata server to check if it is reachable.
// This server is only accessible from a AWS instances, hence
// acts as a good indication of whether the host running this
// function is a AWS instance.
// The URL returns a JSON structure with compute and network info for the instance.
// {"compute":{"location":"westus",
//				"name":"WinResServ16",
// 				"offer":"WindowsServer",
// 				"osType":"Windows",
//				"platformFaultDomain":"0",
//				"platformUpdateDomain":"0",
//				"publisher":"MicrosoftWindowsServer",
//				"sku":"2016-Datacenter",
//				"version":"2016.127.20171017",
//				"vmId":"9xxxxxxd-exxb-xxxd-9xxx-xxxxxxxxxxc0",
//				"vmSize":"Standard_A1"},
//	"network":{"interface":
//				[{"ipv4":{"ipAddress":
//							[{"privateIpAddress":"10.0.0.5",
//					  		  "publicIpAddress":"xx.xxx.xxx.xx"}],
//				  		  "subnet":
//							[{"address":"10.0.0.0",
//							  "prefix":"24"}]
//						  },
//				  "ipv6":{"ipAddress":
//							[]
//						  },
//				  "macAddress":"0012345678AC"
//				}]
//			  }
// }
func HasMetadataHost() bool {
	respText, statusCode, err := makeRequest(
		fmt.Sprintf("%s/?%s", INSTANCE_METADATA_URL, API_VERSION))

	if err != nil {
		return false
	}
	if statusCode == http.StatusNotFound {
		return false
	}
	fmt.Printf("Got HasMetadataHost %s", respText)
	return true
}

// Gets the public hostname of the public IP if any.
// NOTE: Azure instances do not seem to have a public hostname.
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
	fmt.Printf("Got public hostname %s", names[0])
	return names[0], nil
}

// Gets the local hostname 
func Hostname() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/compute/name?%s&format=text",
			INSTANCE_METADATA_URL,
			API_VERSION,
		),
	)
	if err != nil {
		return "", err
	}
	fmt.Printf("Got Hostname %s", respText)
	return respText, nil
}

// Gets the Public/External IP Address 
func PublicIPAddress() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/network/interface/0/ipv4/ipAddress/0/publicIpAddress?%s&format=text",
			INSTANCE_METADATA_URL,
			API_VERSION,
		),
	)
	if err != nil {
		return "", err
	}
	fmt.Printf("Got PublicIPAddress %s", respText)
	return respText, nil
}



// Gets the Local IP Address of the host.
func LocalIPAddress() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/network/interface/0/ipv4/ipAddress/0/privateIpAddress?%s&format=text",
			INSTANCE_METADATA_URL,
			API_VERSION,
		),
	)
	if err != nil {
		return "", err
	}
	fmt.Println("Got LocalIPAddress %s", respText)
	return respText, nil
}

// ** Gets the instance-ID of the host
func Id() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/compute/vmId?%s&format=text", 
			INSTANCE_METADATA_URL,
			API_VERSION,
		),
	)

	if err != nil {
		return "", err
	}
	fmt.Printf("Got Id %s\n", respText)
	return respText, nil
}

// ** Gets the zone name of the host
func Zone() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/compute/location?%s&format=text",
			INSTANCE_METADATA_URL,
			API_VERSION,
		),
	)
	if err != nil {
		return "", err
	}
	fmt.Printf("Got Zone %s", respText)
	return respText, nil
}

// Gets the base machine type of the instance
func Type() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/compute/vmSize?%s&format=text",
			INSTANCE_METADATA_URL,
			API_VERSION,
		),
	)

	if err != nil {
		return "", err
	}
	fmt.Printf("Got Type %s", respText)
	return respText, nil	
}

/* ========================== INFORMATION UNIQUE TO AWS ======================= */
// Gets the offer information of the VM image
// This contains information only for images in the gallery.
func Offer() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/compute/offer?%s&format=text",
			INSTANCE_METADATA_URL,
			API_VERSION,
		),
	)
	if err != nil {
		return "", err
	}
	fmt.Printf("Got Offer %s", respText)
	return respText, nil
}

// Gets the name of the publisher of the VM image.
func Publisher() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/compute/publisher?%s&format=text",
			INSTANCE_METADATA_URL,
			API_VERSION,
		),
	)
	if err != nil {
		return "", err
	}
	fmt.Printf("Got Publisher %s", respText)
	return respText, nil
}

// Gets the Stock Keeping Unit (SKU) of the VM image.
func SKU() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/compute/sku?%s&format=text",
			INSTANCE_METADATA_URL,
			API_VERSION,
		),
	)
	if err != nil {
		return "", err
	}
	fmt.Printf("Got SKU %s", respText)
	return respText, nil
}

// Gets the version of the VM image.
func Version() (string, error) {
	respText, _, err := makeRequest(
		fmt.Sprintf("%s/compute/version?%s&format=text",
			INSTANCE_METADATA_URL,
			API_VERSION,
		),
	)
	if err != nil {
		return "", err
	}
	fmt.Printf("Got Version %s", respText)
	return respText, nil
}
