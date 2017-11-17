package main


import (
	"fmt"
	"net"

	"github.com/rupalkhilari/go-host-info/cloud/aws"
	"github.com/rupalkhilari/go-host-info/cloud/gcp"
)

func main() {
	if cname, err := net.LookupCNAME("www.google.com"); err != nil {
		fmt.Printf("lookup failed: %q\n", err)
	} else {
		fmt.Printf("lookup success: %s\n", cname)
	}
	fmt.Println(IsCloudInstance())
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
	}
	fmt.Println("This host is no where")
	return false
}
