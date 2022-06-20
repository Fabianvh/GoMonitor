package gossm

import (
	"encoding/json"

	"github.com/ssimunic/gossm/validate"
)

type Config struct {
	Servers  Servers   `json:"servers"`
	Settings *Settings `json:"settings"`
}

//NewConfig geeft een pointer terug naar config welke is aangemaakt door de json config file
//Deze word gevalideerd in validate
func NewConfig(jsonData []byte) *Config {
	config := &Config{}
	err := json.Unmarshal(jsonData, config)
	if err != nil {
		panic("error parsing json configuration data")
	}
	if err := validate.ValidateAll(config); err != nil {
		panic(err)
	}
	return config
}
