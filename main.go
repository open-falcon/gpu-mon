/*
 * Copyright 2018 Xiaomi, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/open-falcon/gpu-mon/fetch"
	"github.com/open-falcon/gpu-mon/send"
)

func main() {
	flag.Parse()
	args := os.Args
	if len(args) == 1 || string(args[1][0]) != "-" {
		fmt.Println("improper use, no input parameters")
		showHelpMessage()
		return
	}
	if isUsage() {
		showMessage()
		return
	}
	//默认不采用crontab模式
	err := cfg.InitCommon(configPath, isContab)
	if err != nil {
		fmt.Printf("Initial configuration failed, %v", err)
		return
	}

	data, err := fetch.Data()
	if err != nil {
		cfg.CommonLogger.Error(err)
		return
	}

	metaDataList := send.BuildMetaDatas(data)
	err = send.Data(metaDataList)
	if err != nil {
		cfg.CommonLogger.Error(err)
		return

	}
}
