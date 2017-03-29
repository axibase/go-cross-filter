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
	"encoding/json"
	"strings"
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

func (self *Config) UnmarshalJSON(data []byte) error {
	var jsonConfig JsonConfig
	var err = json.Unmarshal(data, &jsonConfig)
	if err != nil {
		return err
	}

	var config Config
	config.Url = jsonConfig.Url
	config.Port = jsonConfig.Port
	config.UpdatePeriod = jsonConfig.UpdatePeriod

	var tableConfigs = make([]*TableConfig, len(jsonConfig.TableConfigs))
	for i := 0; i < len(jsonConfig.TableConfigs); i++ {
		tableConfig := new(TableConfig)
		var jsonTableConfig = jsonConfig.TableConfigs[i]

		tableConfig.Name = jsonTableConfig.Name
		tableConfig.UseEntityGroupFilter = jsonTableConfig.UseEntityGroupFilter
		tableConfig.SqlQuery = strings.Join(jsonTableConfig.MultilineSqlQuery, "\n")
		tableConfig.EntityGroups = jsonTableConfig.EntityGroups
		tableConfig.PortalConfigPath = jsonTableConfig.PortalConfigPath
		tableConfig.ColumnsConfig = jsonTableConfig.ColumnsConfig

		tableConfigs[i] = tableConfig
	}

	config.TableConfigs = tableConfigs

	*self = config;
	return nil
}


