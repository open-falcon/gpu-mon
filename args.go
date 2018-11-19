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
)

// 版本信息
var (
	Version   string
	BuildTime string
	GoVersion string
	GitCommit string
)

// flag 变量
var (
	isShowVersion bool
	isShowHelp    bool
	isContab      bool
	configPath    string
)

func init() {
	flag.BoolVar(&isShowVersion, "v", false, "show version")
	flag.BoolVar(&isShowVersion, "version", false, "show version")
	flag.BoolVar(&isShowHelp, "h", false, "show help message")
	flag.BoolVar(&isShowHelp, "help", false, "show help message")
	flag.StringVar(&configPath, "c", "", "configure file path")
	flag.BoolVar(&isContab, "o", false, "Output monitoring information to the stdout")
}

// 版本信息
func showVersion() {
	fmt.Printf("Version:%6s\nGit commit:%6s\nGo version:%6s\nBuild time:%6s\n",
		Version, GitCommit, GoVersion, BuildTime)
}

// 帮助信息
func showHelpMessage() {
	fmt.Printf("%s\n\t-o\t%s\n\t-c\t%s\n\t-v, --version\t%s\n\t-h, --help\t%s\n%s\n",
		"[Usage] ./gpu-mon [-o] -c cfg.json",
		"Output to the screen",
		"Specify configuration file",
		"Show version message",
		"Show help message",
		"Note: You need to make sure the nv-hostengine process is running",
	)
}

func showMessage() {
	if isShowVersion {
		showVersion()
	}
	if isShowHelp {
		showHelpMessage()
	}
}

func isUsage() bool {
	if configPath != "" {
		return false
	}
	return true
}
