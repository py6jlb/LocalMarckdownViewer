package main

import (
	"encoding/json"
	"io/ioutil"
)

type config struct {
	ListenAddress string `json:"listen_address"`
	ListenPort    int    `json:"listen_port"`
	DataPath      string `json:"data_path"`
	IndexFileName string `json:"index_file_name"`
}

func initConfig() (c *config, err error) {
	var file []byte
	if file, err = ioutil.ReadFile("config.json"); err != nil {
		return nil, err
	}
	c = new(config)
	if err = json.Unmarshal(file, c); err != nil {
		return nil, err
	}
	return c, nil
}
