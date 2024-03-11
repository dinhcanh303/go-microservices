package kafka

import (
	"bytes"
	"net/http"
	"os"
)

type ConnectorConfig struct {
	ServerUrl       string
	DataType        string
	ConnectorConfig string
	Connector       string
}

func RegisterConnector(config ConnectorConfig) *http.Response {
	plan, _ := os.ReadFile(config.ConnectorConfig)
	response, err := http.Post(config.ServerUrl+"/connectors/", config.DataType, bytes.NewBuffer(plan))
	if err != nil {
		panic(err)
	}
	return response
}
func CheckConnector(config ConnectorConfig) {
	response, err := http.Get(config.ServerUrl + "/connectors/" + config.Connector)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		RegisterConnector(config)
	}
}
