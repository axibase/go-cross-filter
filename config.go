/**
* Copyright 2015 Axibase Corporation or its affiliates. All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License").
* You may not use this file except in compliance with the License.
* A copy of the License is located at
*
* https://www.axibase.com/atsd/axibase-apache-2.0.pdf
*
* or in the "license" file accompanying this file. This file is distributed
* on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
* express or implied. See the License for the specific language governing
* permissions and limitations under the License.
*
** Created by Gregory Kutuzov on 26/10/15.
*
 */

package main

import (
	"bytes"
	"encoding/json"
	neturl "net/url"
	"strings"
	"time"
)

type Config struct {
	Url          Url
	Port         uint
	UpdatePeriod Duration
	TableConfigs []*TableConfig
}

type TableConfig struct {
	Name             	string
	UseEntityGroupFilter 	bool
	SqlQuery       		string
	EntityGroups     	[]string
	PortalConfigPath 	string
	ColumnsConfig    	[]map[string]interface{}
}

type ConfigDto struct {
	Url          Url            `json:"url"`
	Port         uint           `json:"port"`
	UpdatePeriod Duration       `json:"updatePeriod"`
	TableConfigs []*TableConfigDto `json:"tables"`
}

type TableConfigDto struct {
	Name             	string                   `json:"name"`
	UseEntityGroupFilter 	bool			 `json:"useEntityGroupFilter"`
	MultilineSqlQuery       []string                 `json:"sqlQuery"`
	EntityGroups     	[]string                 `json:"entityGroups"`
	PortalConfigPath 	string                   `json:"PortalConfigPath"`
	ColumnsConfig    	[]map[string]interface{} `json:"columns"`
}

func (self *Config) UnmarshalJSON(data []byte) error {
	var configDto ConfigDto
	var err = json.Unmarshal(data, &configDto)
	if err != nil {
		return err
	}

	var config Config
	config.Password = configDto.Password
	config.Port = configDto.Port
	config.UpdatePeriod = configDto.UpdatePeriod

	var tableConfigs = make([]*TableConfig, len(configDto.TableConfigs))
	for i := 0; i < len(configDto.TableConfigs); i++ {
		tableConfig := new(TableConfig)
		var tableConfigDto = configDto.TableConfigs[i]

		tableConfig.Name = tableConfigDto.Name
		tableConfig.UseEntityGroupFilter = tableConfigDto.UseEntityGroupFilter
		tableConfig.SqlQuery = strings.Join(tableConfigDto.MultilineSqlQuery, "\n")
		tableConfig.EntityGroups = tableConfigDto.EntityGroups
		tableConfig.PortalConfigPath = tableConfigDto.PortalConfigPath
		tableConfig.ColumnsConfig = tableConfigDto.ColumnsConfig

		tableConfigs[i] = tableConfig
	}

	config.TableConfigs = tableConfigs

	*self = config;
	return nil
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


