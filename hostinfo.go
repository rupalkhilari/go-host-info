package main


import (
	"fmt"
	"net"

	"github.com/rupalkhilari/go-host-info/cloud/aws"
	"github.com/rupalkhilari/go-host-info/cloud/gcp"
)

type CloudProvider int
const (
	GCP CloudProvider = iota
	AWS
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
		default:
			fmt.Println("No implementation available")
		}
	}

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

}