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
package common

import (
	"fmt"
	"sync"

	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/sirupsen/logrus"
)

type falconConf struct {
	Agent string `json:"agent"`
}

type metricConf struct {
	IgnoreMetrics []string `json:"ignoreMetrics"`
	EndPoint      string   `json:"endpoint"`
}

type logConf struct {
	Level string `json:"level"`
	Dir   string `json:"dir"`
}

//Conf 用户定义的配置
type Conf struct {
	Falcon       falconConf
	Metric       metricConf
	Log          logConf
	MetricFilter map[string]struct{}
	IsCrontab    bool
}

var (
	globalConf Conf                // 读取的配置项
	configLock = new(sync.RWMutex) // 加锁
)

// Config 返回全局配置
func Config() *Conf {
	configLock.Lock()
	defer configLock.Unlock()
	return &globalConf
}

func readConfigFile(configPath string) (err error) {
	if !fileExist(configPath) {
		return fmt.Errorf("config file %s is not exicted", configPath)
	}
	err = config.Load(file.NewSource(file.WithPath(configPath)))
	if err != nil {
		return err
	}
	return initGlobalConfig()
}

func initGlobalConfig() (err error) {
	//todo: 使用命令行配置项覆盖
	err = config.Get("falcon").Scan(&globalConf.Falcon)
	if err != nil {
		return err
	}
	err = config.Get("metric").Scan(&globalConf.Metric)
	if err != nil {
		return err
	}
	err = config.Get("log").Scan(&globalConf.Log)
	if err != nil {
		return err
	}
	return err
}

func initIgnoreMetrics() {
	metricFilter := make(map[string]struct{})
	for _, metric := range globalConf.Metric.IgnoreMetrics {
		metricFilter[metric] = struct{}{}
	}
	globalConf.MetricFilter = metricFilter
	Logger.WithFields(logrus.Fields{
		"ignore metrics": globalConf.Metric.IgnoreMetrics,
	}).Info("ignore metrics ")
}

func initConfig(configPath string, isCrontab bool) error {
	err := readConfigFile(configPath)
	if err != nil {
		return fmt.Errorf("unable to read Config file: %v", err)
	}
	globalConf.IsCrontab = isCrontab
	initIgnoreMetrics()
	return nil
}
