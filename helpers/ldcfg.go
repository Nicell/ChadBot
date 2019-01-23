package helpers

import (
	"encoding/json"
	"os"

	S "chad/structs"
)

// LdCFG loads config at config.json to the given config parameter
func LdCFG(cfg *S.Config) error {

	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)

	return nil
}
