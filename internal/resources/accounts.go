package resources

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Account struct
type Account struct {
	method   string
	metadata map[string]interface{}
}

func (a *Account) objectName() string {
	const obName = "Account"
	return obName
}

// newAccount writes that new Customer instance
func NewAccount(metadata map[string]interface{}) (*Account, error) {
	rawMethod, ok := metadata["method"]
	if !ok {
		return nil, errors.New("missing required parameters: method")
	}
	method, ok := rawMethod.(string)
	if !ok {
		return nil, errors.New("failed to convert interface{} to string")
	}
	return &Account{
		method:   method,
		metadata: metadata,
	}, nil
}

// getMetadata mold customer get metadata
func (a *Account) getMetadata() (map[string]interface{}, error) {
	idIF, ok := a.metadata["id"]
	if !ok {
		// TODO: key が account_id の場合のための暫定対応
		idIF, ok = a.metadata["account_id"]
		if !ok {
			return nil, fmt.Errorf("account_id or id is required")
		}
	}
	pathParam, ok := idIF.(string)
	if !ok {
		return nil, errors.New("failed to convert interface{} to string")
	}
	return buildMetadata(a.method, a.objectName(), pathParam, nil, "", "account_get"), nil
}

// postMetadata mold customer post metadata
func (a *Account) postMetadata() (map[string]interface{}, error) {
	data := map[string]interface{}{}
	dataIF, exists := a.metadata["data"]
	if exists && dataIF != nil {
		if _, ok := dataIF.(map[string]interface{}); !ok {
			return nil, errors.New("failed to convert interface{} to map[string]inteeface{}")
		}
		data, _ = dataIF.(map[string]interface{})
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON marshal error: %v", err)
	}
	body := string(bytes)
	return buildMetadata(a.method, a.objectName(), "", nil, body, "account_post"), nil
}

// updateMetadata mold customer update metadata
func (a *Account) updateMetadata() (map[string]interface{}, error) {
	dataIF, exists := a.metadata["data"]
	if !exists {
		return nil, fmt.Errorf("data is required")
	}
	data, _ := dataIF.(map[string]interface{})
	idIF, ok := a.metadata["id"]
	if !ok {
		idIF, ok = a.metadata["account_id"]
		if !ok {
			return buildMetadata(a.method, a.objectName(), "", nil, "", "account_put"), nil
		}
	}
	id := idIF.(string)
	data["Id"] = id
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON marshal error: %v", err)
	}
	body := string(bytes)
	return buildMetadata(a.method, a.objectName(), "", nil, body, "account_put"), nil
}

// BuildMetadata
func (a *Account) BuildMetadata() (map[string]interface{}, error) {
	switch a.method {
	case "get":
		return a.getMetadata()
	case "post":
		return a.postMetadata()
	case "put":
		return a.updateMetadata()
	}
	return nil, fmt.Errorf("invalid method: %s", a.method)
}

func buildMetadata(method, object, pathParam string, queryParams map[string]string, body string, connectionKey string) map[string]interface{} {
	metadata := map[string]interface{}{
		"method":         method,
		"object":         object,
		"connection_key": connectionKey,
	}
	if len(pathParam) > 0 {
		metadata["path_param"] = pathParam
	}
	if queryParams != nil {
		metadata["query_params"] = queryParams
	}
	if body != "" {
		metadata["body"] = body
	}
	return metadata
}
