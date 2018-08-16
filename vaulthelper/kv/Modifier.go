package kv

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/vault/api"
)

// Set all paths that don't use environments to true
var noEnvironments = map[string]bool{
	"templates/": true,
	"cubbyhole/": true,
}

// Modifier maintains references to the active client and
// respective logical needed to write to the vault. Path
// can be changed to alter where in the vault the key,value
// pair is stored
type Modifier struct {
	client  *api.Client  // Client connected to vault
	logical *api.Logical // Logical used for read/write options
	Env     string       // Environment (local/dev/QA; Initialized to secrets)
}

// NewModifier Constructs a new modifier struct and connects to the vault
// @param token 	The access token needed to connect to the vault
// @param address	The address of the API endpoint for the server
// @return 			A pointer to the newly contstructed modifier object (Note: path set to default),
// 		   			Any errors generated in creating the client
func NewModifier(token string, address string) (*Modifier, error) {
	if len(address) == 0 {
		address = "http://127.0.0.1:8200" // Default address
	}
	httpClient, err := CreateHTTPClient()
	if err != nil {
		return nil, err
	}
	// Create client
	modClient, err := api.NewClient(&api.Config{
		Address: address, HttpClient: httpClient,
	})
	if err != nil {
		return nil, err
	}

	// Set access token and path for this modifier
	modClient.SetToken(token)

	// Return the modifier
	return &Modifier{client: modClient, logical: modClient.Logical(), Env: "secret"}, nil
}

// Writes the key,value pairs in data to the vault
//
// @param   data A set of key,value pairs to be written
//
// @return	Warnings (if any) generated from the vault,
//			errors generated by writing
func (m *Modifier) Write(path string, data map[string]interface{}) ([]string, error) {
	// Wrap data and send
	sendData := map[string]interface{}{"data": data}
	// Create full path
	pathBlocks := strings.SplitAfterN(path, "/", 2)
	if len(pathBlocks) == 1 {
		pathBlocks[0] += "/"
	}
	fullPath := pathBlocks[0] + "data/"
	if !noEnvironments[pathBlocks[0]] {
		fullPath += m.Env + "/"
	}
	if len(pathBlocks) > 1 {
		fullPath += pathBlocks[1]
	}
	Secret, err := m.logical.Write(fullPath, sendData)
	if Secret == nil { // No warnings
		return nil, err
	}
	return Secret.Warnings, err
}

// ReadData Reads the most recent data from the path referenced by this Modifier
// @return	A Secret pointer that contains key,value pairs and metadata
//			errors generated from reading
func (m *Modifier) ReadData(path string) (map[string]interface{}, error) {
	// Create full path
	pathBlocks := strings.SplitAfterN(path, "/", 2)
	fullPath := pathBlocks[0] + "data/"
	if !noEnvironments[pathBlocks[0]] {
		fullPath += m.Env + "/"
	}
	if len(pathBlocks) > 1 {
		fullPath += pathBlocks[1]
	}
	secret, err := m.logical.Read(fullPath)
	if secret == nil {
		return nil, err
	}
	if data, ok := secret.Data["data"].(map[string]interface{}); ok {
		return data, err
	}
	return nil, errors.New("Could not get data from vault response")

}

//ReadValue takes a path and a key and returns the corresponding value from the vault
func (m *Modifier) ReadValue(path string, key string) (string, error) {
	valueMap, err := m.ReadData(path)
	if err != nil {
		return "", err
	}
	//return value corresponding to the key
	if valueMap[key] != nil {
		if value, ok := valueMap[key].(string); ok {
			return value, nil
		} else if stringer, ok := valueMap[key].(fmt.GoStringer); ok {
			return stringer.GoString(), nil
		} else {
			return "", fmt.Errorf("Cannot convert value at %s to string", key)
		}
	}
	return "", fmt.Errorf("Key '%s' not found in '%s'", key, path)
}

// ReadMetadata Reads the Metadata from the path referenced by this Modifier
// @return	A Secret pointer that contains key,value pairs and metadata
//			errors generated from reading
func (m *Modifier) ReadMetadata(path string) (map[string]interface{}, error) {
	pathBlocks := strings.SplitAfterN(path, "/", 2)
	fullPath := pathBlocks[0] + "data/"
	if !noEnvironments[pathBlocks[0]] {
		fullPath += m.Env + "/"
	}
	fullPath += pathBlocks[1]
	secret, err := m.logical.Read(fullPath)
	if data, ok := secret.Data["metadata"].(map[string]interface{}); ok {
		return data, err
	}
	return nil, errors.New("Could not get metadata from vault response")
}

//List lists the paths underneath this one
func (m *Modifier) List(path string) (*api.Secret, error) {
	pathBlocks := strings.SplitAfterN(path, "/", 2)
	if len(pathBlocks) == 1 {
		pathBlocks[0] += "/"
	}

	fullPath := pathBlocks[0] + "metadata/"
	if !noEnvironments[pathBlocks[0]] {
		fullPath += m.Env + "/"
	}
	if len(pathBlocks) > 1 {
		fullPath += pathBlocks[1]
	}
	return m.logical.List(fullPath)
}

//AdjustValue adjusts the value at the given path/key by n
func (m *Modifier) AdjustValue(path string, key string, n int) ([]string, error) {
	// Get the existing data at the path
	oldData, err := m.ReadData(path)
	if err != nil {
		return nil, err
	}
	if oldData == nil { // Path has not been used yet, create an empty map
		oldData = make(map[string]interface{})
	}
	// Try to fetch the value with the given key, start empty values with 0
	if oldData[key] == nil {
		oldData[key] = "0"
	}
	// Convert from stored string value to int
	oldValue, err := strconv.Atoi(oldData[key].(string))
	if err != nil {
		return []string{"Could not convert value to int at: " + key}, err
	}
	newValue := strconv.Itoa(oldValue + n)
	oldData[key] = newValue
	return m.Write(path, oldData)
}