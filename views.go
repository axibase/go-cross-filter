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
	"net/http"
	"text/template"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	ts "github.com/axibase/go-cross-filter/model/table"
)

func PortalView(query string, w http.ResponseWriter) {
	t, _ := template.ParseFiles("templates/view.tpl")
	t.Execute(w, struct {
		Query string
	}{query})

}

type Col struct {
	Name  string
	Value string
}

func ConfigView(table string, cols []*Col, w http.ResponseWriter) {
	portalConfigPath := ""
	for _, config := range app.Config().TableConfigs {
		if config.Name == table {
			portalConfigPath = config.PortalConfigPath
		}
	}
	if portalConfigPath == "" {
		http.NotFound(w, nil)
		return
	}
	t, err := template.ParseFiles(portalConfigPath)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	colMap := map[string]string{}
	for _, c := range cols {
		colMap[c.Name] = c.Value
	}
	t.Execute(w, colMap)
}

func EntityGroupListView(table string, w http.ResponseWriter) {
	entityGroupsNames := []string{}
	for _, tableConfig := range Instance().Config().TableConfigs {
		if table == tableConfig.Name {
			entityGroupsNames = tableConfig.EntityGroups
		}
	}
	if len(entityGroupsNames) == 0 {
		entityGroupsNames = Instance().TableService.GetEntityGroups()
	}
	resp, _ := json.Marshal(entityGroupsNames)
	w.Header()["Content-Type"] = append(w.Header()["Content-Type"], "application/json")
	w.Write(resp)
}

func IndexView(w http.ResponseWriter) {
	t, err := ioutil.ReadFile("templates/index.tpl")
	if err != nil {
		fmt.Println(err)
		http.NotFound(w, nil)
		return
	}
	w.Write(t)
}

type Table struct {
	Name                 string                   `json:"name"`
	UpdateTime           time.Time                `json:"updateTime"`
	UseEntityGroupFilter bool                     `json:"useEntityGroupFilter"`
	Columns              []map[string]interface{} `json:"columns"`
	Rows                 [][]string               `json:"rows"`
}

func TableView(table, entityGroup string, w http.ResponseWriter) {
	if table, ok := Instance().TableService.GetTable(table); ok {
		var tableConfig *TableConfig
		resultTable := table
		for _, config := range Instance().Config().TableConfigs {
			if config.Name == table.Name {
				tableConfig = config
			}
		}

		if entityGroup != "" {
			entities := Instance().TableService.GetGroupEntities(entityGroup)
			resultTable = filterEntities(entities, table)
		}
		json, err := json.Marshal(convertToTable(resultTable, tableConfig))
		if err != nil {
			panic(err)
		}
		w.Write(json)
	} else {
		http.NotFound(w, nil)
		return
	}
}
func convertToTable(table *ts.Table, tableConfig *TableConfig) *Table {
	newTable := &Table{
		Name:                 table.Name,
		UseEntityGroupFilter: tableConfig.UseEntityGroupFilter,
		UpdateTime:           table.UpdateTime,
		Rows:                 table.Rows,
	}
	for _, col := range table.Columns {
		column := map[string]interface{}{}

		for _, colConfig := range tableConfig.ColumnsConfig {
			if colConfig["name"] == col.Label {
				for key, value := range colConfig {
					if key != "name" && key != "label" {
						column[key] = value
					}
				}
			}
		}
		column["name"] = col.Name
		column["label"] = col.Label
		column["numeric"] = col.Numeric

		newTable.Columns = append(newTable.Columns, column)
	}
	return newTable
}
func filterEntities(entities []string, table *ts.Table) *ts.Table {
	entityColIndex := 0
	for i, col := range table.Columns {
		if col.Name == "entity" {
			entityColIndex = i
		}
	}
	newTable := &ts.Table{
		Name:       table.Name,
		Columns:    table.Columns,
		UpdateTime: table.UpdateTime,
		Rows:       [][]string{},
	}
	for _, row := range table.Rows {
		for _, entity := range entities {
			if row[entityColIndex] == entity {
				newTable.Rows = append(newTable.Rows, row)
			}
		}
	}
	return newTable
}

func TableListView(w http.ResponseWriter) {
	tableNames := Instance().TableService.GetTablesList()
	resp, _ := json.Marshal(tableNames)
	w.Header()["Content-Type"] = append(w.Header()["Content-Type"], "application/json")
	w.Write(resp)
}
