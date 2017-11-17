// Package aws contains methods querying certain Azure, host-specific data.
// More documentation on retriving metadata information from the 
// Azure instance can be found at:
// https://docs.microsoft.com/en-us/azure/virtual-machines/windows/instance-metadata-service

package azure

const (
	INSTANCE_METADATA_URL = "http://169.254.169.254/metadata/instance?api-version=2017-04-02"
)

// Makes a request to the given url adding the required header data to
// query a Azure instance. 
// It returns the response text, status code and any error encountered.
func makeRequest(url string) (string, int, error) {

	// Construct an http Client with an acceptable timeout
	client := &http.Client{
		Timeout: time.Second * 10
	}

	// Create a new http GET request with the given URL
	req, err := http.NewRequest("GET", url, nil)

	// Set the headers required to make a GET request.
	req.Header.Set("Metadata", "true")

	// Make the request
	resp, err := client.Do(req) 

	if err != nil {
		return "", resp.StatusCode, err
	}
	defer rs.Body.Close()

	// Converts the response body bytes[] to a string
	// and returns it
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode, err
	}
	return string(body), resp.StatusCode, err
}
