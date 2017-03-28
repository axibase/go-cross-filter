package main

import (
	"bytes"
	"encoding/json"
	neturl "net/url"
	"strings"
	"time"
)

type JsonConfig struct {
	Url          Url            `json:"url"`
	Port         uint           `json:"port"`
	UpdatePeriod Duration       `json:"updatePeriod"`
	TableConfigs []*JsonTableConfig `json:"tables"`
}

type JsonTableConfig struct {
	Name             	string                   `json:"name"`
	UseEntityGroupFilter 	bool			 `json:"useEntityGroupFilter"`
	MultilineSqlQuery       []string                 `json:"sqlQuery"`
	EntityGroups     	[]string                 `json:"entityGroups"`
	PortalConfigPath 	string                   `json:"PortalConfigPath"`
	ColumnsConfig    	[]map[string]interface{} `json:"columns"`
}

type Url neturl.URL

func (self *Url) UnmarshalJSON(data []byte) error {
	var url string
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	if err := dec.Decode(&url); err != nil {
		return err
	}
	url1, err := neturl.Parse(url)
	*self = Url(*url1)
	return err
}

type Duration time.Duration

func (self *Duration) UnmarshalJSON(input []byte) error {
	str := strings.Trim(string(input), "\"")
	dur, err := time.ParseDuration(str)
	if err != nil {
		return err
	}
	*self = Duration(dur)
	return nil
}