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

type Config struct {
	Url          Url            `json:"url"`
	User         string         `json:"user"`
	Password     string         `json:"password"`
	Port         uint           `json:"port"`
	UpdatePeriod Duration       `json:"updatePeriod"`
	TableConfigs []*TableConfig `json:"tables"`
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

type TableConfig struct {
	Name             string                   `json:"name"`
	SqlQuery         string                   `json:"sqlQuery"`
	EntityGroups     []string                 `json:"entityGroups"`
	PortalConfigPath string                   `json:"PortalConfigPath"`
	ColumnsConfig    []map[string]interface{} `json:"columns"`
}
