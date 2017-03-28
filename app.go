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
** Created by Gregory Kutuzov on 19/10/15.
*
 */

package main

import (
	table "github.com/axibase/go-cross-filter/model/table"
	web "github.com/axibase/go-cross-filter/web"
	neturl "net/url"
	"strings"
	"time"
)

var app *App
var defaultConfig *Config

func init() {
	url, _ := neturl.Parse("http://localhost:8088")
	defaultConfig = &Config{
		Url:          Url(*url),
		User:         "admin",
		Password:     "admin",
		UpdatePeriod: Duration(1 * time.Minute),
		Port:         8000,
		TableConfigs: []*TableConfig{},
	}

	app = &App{
		web:          web.NewWeb(),
		TableService: table.NewTableService(),
	}

	app.Init(defaultConfig)
}

type App struct {
	config       *Config
	web          *web.Web
	TableService *table.TableService

	isRunning bool
}

func (self *App) Config() Config {
	return *self.config
}

func Instance() *App {
	return app
}

func (self *App) Init(config *Config) {
	self.config = config
	tableConfigs := []*table.TableConfig{}
	for _, tableConfig := range config.TableConfigs {
		tableConfigs = append(tableConfigs, &table.TableConfig{
			Name:     tableConfig.Name,
			SqlQuery: strings.Join(tableConfig.MultilineSqlQuery, "\n"),
		})
	}
	self.TableService.Init(tableConfigs, neturl.URL(config.Url), config.User, config.Password)
	app.web.ResetHandlers()
	app.web.Register("/", &TableController{})
	app.web.Register("/", NewPortalController(neturl.URL(config.Url)))
}
func (self *App) Start() error {
	if self.isRunning {
		panic("Error: app is already running")
	}
	go self.TableService.StartUpdatingTables(time.Duration(self.config.UpdatePeriod))
	self.isRunning = true
	err := self.web.Serve(self.config.Port)
	if err != nil {
		return err
	}
	return nil
}
func (self *App) Stop() {
	if !self.isRunning {
		panic("Error: app is already stopping")
	}
	self.isRunning = false
	self.web.Stop()
	self.TableService.StopUpdatingService()
}

func (self *App) IsRunning() bool {
	return self.isRunning
}
