package main

import (
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/json"
	"os"
)

type Config struct {
	Profile   string
	LocalKey  string
	ServerKey string
}

func ReadOrCreateConfig() Config {

	var configData Config

	config, err := os.OpenFile("./config", os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		log_error(err, 2)
	}
	info, err := config.Stat()
	if err != nil {
		log_error(err, 2)
	}

	// Empty config file, creating one
	if info.Size() == 0 {

		key := make([]byte, 32)

		_, err := rand.Read(key)
		if err != nil {
			log_error(err, 2)
		}

		var encryptionKey string = b64.StdEncoding.EncodeToString(key)
		configData = Config{
			Profile:   "default",
			LocalKey:  encryptionKey,
			ServerKey: "",
		}

		data, err := json.Marshal(configData)
		if err != nil {
			log_error(err, 2)
		}

		_, err = config.Write(data)

		if err != nil {
			log_error(err, 2)
		}

		return Config{
			Profile:  "default",
			LocalKey: encryptionKey,
		}

	} else {
		jsonParser := json.NewDecoder(config)
		jsonParser.Decode(&configData)
		return configData
	}
}
