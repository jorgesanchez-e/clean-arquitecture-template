package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/spf13/viper"
)

const (
	configPathETC      string = "/etc/example"
	configPathLocalETC string = "/usr/local/etc/example"
	configPathHome     string = "$HOME/.example"
	configPathLocal    string = "./example/"

	configFileName string = "config"
	configFileType string = "yaml"

	ErrReadConfigFile Error = "unable to read config file"
)

var onlyOnce sync.Once

type Error string

func (e Error) Error() string {
	return string(e)
}

type config struct {
	vp *viper.Viper
}

func New() (*config, error) {
	var cnf *config
	var err error

	onlyOnce.Do(func() {
		cnf = new(config)
		err = cnf.read()
	})

	return cnf, err
}

func (c *config) read() error {
	c.vp = viper.New()

	c.vp.SetConfigName(configFileName)
	c.vp.SetConfigType(configFileType)
	c.vp.AddConfigPath(configPathETC)
	c.vp.AddConfigPath(configPathLocalETC)
	c.vp.AddConfigPath(configPathHome)
	c.vp.AddConfigPath(configPathLocal)

	if err := c.vp.ReadInConfig(); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrReadConfigFile)
	}

	return nil
}

func (c *config) Find(node string) (io.Reader, error) {
	if c == nil {
		return nil, errors.New("config haven't been read")
	}

	n := c.vp.Get(node)
	if d, is := n.(map[string]interface{}); is {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(d); err != nil {
			return nil, err
		}

		return buf, nil
	}

	return nil, nil
}
