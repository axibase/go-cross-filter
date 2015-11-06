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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func loadConfig(filename string) (*Config, error) {
	var config Config
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bs, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
func updateConfig(app *App) {
	for {
		err := watchFile("config.json")
		if err != nil {
			fmt.Println(err)
			continue
		}
		config, err := loadConfig("config.json")
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Start updating config...")
		fmt.Println("Stopping app...")
		app.Stop()
		fmt.Println("Initialising app...")
		app.Init(config)
	}
}

func main() {
	app := Instance()
	config, err := loadConfig("config.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(config)
	app.Init(config)
	go updateConfig(app)
	for {
		fmt.Println("Starting app...")
		app.Start()
	}
}
