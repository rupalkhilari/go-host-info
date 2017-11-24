package main


import (
	"fmt"
	"net"
	"os"
	"errors"

	"github.com/rupalkhilari/go-host-info/cloud/aws"
	"github.com/rupalkhilari/go-host-info/cloud/gcp"
	"github.com/rupalkhilari/go-host-info/cloud/azure"
)

type CloudProvider int
const (
	GCP CloudProvider = iota
	AWS
	AZURE
	UNKNOWN
)

func main() {
	if cname, err := net.LookupCNAME("www.google.com"); err != nil {
		fmt.Printf("lookup failed: %q\n", err)
	} else {
		fmt.Printf("lookup success: %s\n", cname)
	}
	fmt.Println(IsCloudInstance())
	if IsCloudInstance() == true {
		provider := DetermineCurrentCloud()
		switch provider {
		case GCP:
			RunGCPCloudFuncs()
		case AWS:
			RunAWSCloudFuncs()
		case AZURE:
			RunAzureCloudFuncs()
		default:
			fmt.Println("No implementation available")
		}
	}


	/// Printing out the Go-host info.
	fmt.Println(GetHostname())

	// Printing some static IP information
	fmt.Println(GetCName())
	fmt.Println(GetHostInfo())
	fmt.Println(ExternalIP())
	LookupHost()


}

// Determines if the host is a cloud instance.
func IsCloudInstance() bool {
	// Ping the internal URLs of GCP/AWS/Azure to test.
	if aws.HasMetadataHost() == true {
		fmt.Println("This is an aws host")
		return true
	} else if gcp.HasMetadataHost() == true {
		fmt.Println("This is a gcp host")
		return true
	} else if azure.HasMetadataHost() == true {
		fmt.Println("This is an azure host")
		return true
	} else {
		fmt.Println("This host is not on any known cloud provider")
	}
	return false
}

func DetermineCurrentCloud() CloudProvider {
	if aws.HasMetadataHost() == true {
		return AWS
	} else if gcp.HasMetadataHost() == true {
		return GCP
	}
	return UNKNOWN
}
// Runs all the cloud host info to retrieve details.

func RunGCPCloudFuncs() {
	gcp.FQDN()
	gcp.PublicHostname()
	gcp.Hostname()
	gcp.LocalIPAddress()
	gcp.PublicIPAddress()
	gcp.Id()
	gcp.Zone()
	gcp.Type()
	gcp.IsPreemptible()
	gcp.Tags()
}

func RunAWSCloudFuncs() {
	aws.FQDN()
	aws.PublicHostname()
	aws.Hostname()
	aws.LocalIPAddress()
	aws.PublicIPAddress()
	aws.Id()
	aws.Zone()
	aws.Type()
	aws.ImageId()
}

func RunAzureCloudFuncs() {
	azure.PublicHostname()
	azure.Hostname()
	azure.LocalIPAddress()
	azure.PublicIPAddress()
	azure.Id()
	azure.Zone()
	azure.Type()
	azure.Offer()
	azure.Publisher()
	azure.Version()
	azure.SKU()
}

func GetHostname() (string, error) {
	return os.Hostname()
}

func GetCName() (string, error) {
	hostname, err := GetHostname()
	if err != nil {
		return "", err
	}
	return net.LookupCNAME(hostname)
}

func LookupHost() {
	hostname, err := GetHostname()
	if err != nil {
		return
	}
	addrs, err := net.LookupHost(hostname)


	if err != nil {
	    fmt.Printf("Error: %v\n", err)
	}

	for _, a := range addrs {
	    fmt.Println(a)
	}
}
func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func GetHostInfo() ([]string, error) {
	hostname, err := GetHostname()
	if err != nil {
		return []string{}, err
	}
	return net.LookupHost(hostname)
}
