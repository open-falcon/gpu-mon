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
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// 日志文件
var (
	commonLogName = "monitor.log"
	Logger        = logrus.New()
)

// 创建log日志路径，返回创建目录及创建的相关信息
func createLogPath(logDirPath string, logName string) (logPath string, err error) {

	if fileExist(logDirPath) {
		return filepath.Join(logDirPath, logName), err
	}
	err = os.MkdirAll(logDirPath, os.ModeDir)
	if err != nil {
		logDirPath = "."
	}
	logPath = filepath.Join(logDirPath, logName)
	return logPath, err
}

// default Warn Level
func setLogLevel(logLevel string, Log *logrus.Logger) error {
	switch logLevel {
	case "Warn":
		Log.SetLevel(logrus.WarnLevel)
	case "Error":
		Log.SetLevel(logrus.ErrorLevel)
	case "Debug":
		Log.SetLevel(logrus.DebugLevel)
	default:
		Log.SetLevel(logrus.WarnLevel)
		return fmt.Errorf("input logLevel [%s] is not supported, only support %s",
			logLevel,
			"Warn/Error/Debug")
	}
	return nil
}

func setLogOut(logDirPath, logName string, log *logrus.Logger) {
	logPath, _ := createLogPath(logDirPath, logName)
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Out = os.Stderr
	} else {
		log.Out = file
	}
}

func setLogFormat(log *logrus.Logger) {
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})
}

func configLogger(logDirPath, logLevel, logName string, log *logrus.Logger) {
	setLogOut(logDirPath, logName, log)
	setLogFormat(log)
	err := setLogLevel(logLevel, log)
	if err != nil {
		log.Printf("set log level error: %v", err)
	}
}

// 初始化日志对象
func initLoggor(logDirPath string, logLevel string) {
	configLogger(logDirPath, logLevel, commonLogName, Logger)
}
