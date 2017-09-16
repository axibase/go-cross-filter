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
** Created by Gregory Kutuzov on 20/10/15.
*
 */

package table

import (
	"encoding/json"
	"fmt"
	neturl "net/url"
	"sort"
	"strings"
	"time"

	atsdHttp "github.com/axibase/atsd-api-go/http"
)

var (
	NUMERIC_TYPES = []string{
		"short",
		"integer",
		"long",
		"float",
		"double",
	}
)

type TableConfig struct {
	Name         string
	SqlQuery     string
	EntityGroups []string
}

type TableService struct {
	tables           map[string]*Table
	entitiesPerGroup map[string][]string
	tableConfigs     []*TableConfig
	client           *atsdHttp.Client
	entityGroups     []string
	stop             chan bool
	isUpdating       bool
}

func NewTableService() *TableService {
	return &TableService{
		stop:             make(chan bool),
		tables:           map[string]*Table{},
		entitiesPerGroup: map[string][]string{},
	}
}

func (self *TableService) Init(tableConfigs []*TableConfig, url neturl.URL) {
	self.tableConfigs = tableConfigs
	self.client = atsdHttp.New(url, true)
	for _, table := range self.tables {
		if !Contains(table.Name, self.tableConfigs) {
			delete(self.tables, table.Name)
		}
	}
}
func Contains(tableName string, tableConfigs []*TableConfig) bool {
	for _, config := range tableConfigs {
		if config.Name == tableName {
			return true
		}
	}
	return false
}

func (self *TableService) StartUpdatingTables(UpdatePeriod time.Duration) {
	if self.isUpdating {
		panic("Error: service is already running")
	}
	self.isUpdating = true
	for {
		self.entitiesPerGroup = map[string][]string{}
		fmt.Println("Updating tables...")
		for _, tableConfig := range self.tableConfigs {
			entityGroups, err := self.loadEntityGroups()
			if err == nil {
				self.entityGroups = entityGroups
			} else {
				fmt.Println(err)
			}
			if len(tableConfig.EntityGroups) > 0 {
				entityGroups = tableConfig.EntityGroups
			}
			for _, group := range entityGroups {
				self.entitiesPerGroup[group] = self.loadEntitiesForGroup(group)
			}
			table, err := self.loadTable(tableConfig)
			if err == nil && len(table.Rows) > 0 {
				self.tables[tableConfig.Name] = table
				fmt.Println("Table", tableConfig.Name, "ok")
			} else if err != nil {
				fmt.Println("Error: ", err)
			}
		}
		select {
		case <-time.After(UpdatePeriod):
		case <-self.stop:
			return
		}
	}
}

func (self *TableService) loadEntitiesForGroup(entityGroup string) []string {
	entities, err := self.client.EntityGroups.EntitiesList(entityGroup, "", nil, 0)
	if err != nil {
		fmt.Println(err)
	}
	entityNames := []string{}
	for _, entity := range entities {
		entityNames = append(entityNames, entity.Name())
	}
	return entityNames
}

func (self *TableService) StopUpdatingService() {
	if !self.isUpdating {
		panic("Error service has already stopped")
	}
	self.stop <- true
	self.isUpdating = false
}

func (self *TableService) loadEntityGroups() ([]string, error) {
	entityGroupsNames := []string{}
	entityGroups, err := self.client.EntityGroups.List("", nil, 0)
	if err != nil {
		return nil, err
	}
	for _, entityGroup := range entityGroups {
		entityGroupsNames = append(entityGroupsNames, entityGroup.Name)
	}
	return entityGroupsNames, nil
}
func (self *TableService) loadTable(tableConfig *TableConfig) (*Table, error) {
	fmt.Println("Start load table")
	responseTable, err := self.client.SQL.Query(tableConfig.SqlQuery)
	fmt.Println("Stop load table")
	if err != nil {
		return nil, err
	}
	table := &Table{
		Name:       tableConfig.Name,
		UpdateTime: time.Now(),
		Columns:    []*Column{},
		Rows:       [][]string{},
	}
	if meta, ok := responseTable.Metadata.(map[string]interface{}); ok {
		if schema, ok := meta["tableSchema"].(map[string]interface{}); ok {
			if columns, ok := schema["columns"].([]interface{}); ok {
				for _, responseCol := range columns {
					if c, ok := responseCol.(map[string]interface{}); ok {
						table.Columns = append(table.Columns, &Column{
							Name:    c["name"].(string),
							Label:   c["titles"].(string),
							Numeric: stringInSlice(c["datatype"].(string), NUMERIC_TYPES),
						})
					} else {
						fmt.Println("Col object not found")
					}
				}
			} else {
				fmt.Println("Columns not found")
			}
		} else {
			fmt.Println("Schema not found")
		}
	} else {
		fmt.Println("Meta not found")
	}

	for _, responseRow := range responseTable.Data {
		row := []string{}
		for _, responseCell := range responseRow {
			switch c := responseCell.(type) {
			case string:
				row = append(row, c)
			case json.Number:
				row = append(row, c.String())
			case nil:
				row = append(row, "null")
			default:
				panic("Error undefined type")
			}
		}
		table.Rows = append(table.Rows, row)
	}
	return table, nil
}

func (self *TableService) GetTable(name string) (*Table, bool) {
	table, ok := self.tables[name]
	return table, ok
}
func (self *TableService) GetTablesList() []string {
	tables := []string{}
	for table := range self.tables {
		tables = append(tables, table)
	}
	sort.Strings(tables)
	return tables
}
func (self *TableService) GetEntityGroups() []string {
	return self.entityGroups
}
func (self *TableService) GetGroupEntities(entityGroup string) []string {
	return self.entitiesPerGroup[entityGroup]
}
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}
