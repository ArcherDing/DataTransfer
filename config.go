package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	configfile string
	RemoteAddr string
	RemotePort string
	LocaleAddr string
	LocalePort string
}

func NewConfig(configfile string) *Config {
	return &Config{
		configfile: configfile,
	}
}

func (self *Config) LoadConfig() error {
	file, err := os.Open(self.configfile)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	defer file.Close()

	dec := json.NewDecoder(file)
	err = dec.Decode(&self)
	if err != nil {
		log.Printf(err.Error())
		return err
	}
	return nil
}

func (self *Config) DumpConfig() {
	fmt.Printf("%v", self)
}
