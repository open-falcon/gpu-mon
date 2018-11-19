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
package cfg

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

//ToDo :逻辑比较乱
func Test_setLogLevel(t *testing.T) {
	type args struct {
		logLevel string
		Log      *logrus.Logger
	}
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantLevel logrus.Level
	}{
		{"test1", args{"Warn", CommonLogger}, false, logrus.WarnLevel},
		{"test2", args{"Error", CommonLogger}, false, logrus.ErrorLevel},
		{"test3", args{"Debug", CommonLogger}, false, logrus.DebugLevel},
		{"test4", args{"wrongLevel", CommonLogger}, true, logrus.WarnLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setLogLevel(tt.args.logLevel, tt.args.Log)
			getLevel := CommonLogger.GetLevel()
			if (err != nil) != tt.wantErr {
				t.Errorf("setLogLevel() \nerror = %v,\nwantErr = %v", err, tt.wantErr)
			}

			if getLevel != tt.wantLevel {
				t.Errorf("setLogLevel() \ngetlevel = %v, \nwantLevel = %v", getLevel, tt.wantLevel)
			}

		})
	}
}

func Test_createLogPath(t *testing.T) {
	const (
		message1 = "logDirPath is not exit, and create it failed"
		message2 = "logDirPath is not exit, created it"
	)
	type args struct {
		logDirPath string
		logName    string
	}
	tests := []struct {
		name        string
		args        args
		wantLogPath string
		wantMessage string
	}{
		{"test1", args{"./testData", "test.log"}, "testData/test.log", ""},
		{"test2", args{"./testData/testData1", "test.log"}, "testData/testData1/test.log", message2},
		{"test3", args{"", "test.log"}, "", message1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLogPath, gotMessage := createLogPath(tt.args.logDirPath, tt.args.logName)
			if gotLogPath != tt.wantLogPath {
				t.Errorf("createLogPath() gotLogPath = %v, want %v", gotLogPath, tt.wantLogPath)
			}
			if gotMessage != tt.wantMessage {
				t.Errorf("createLogPath() gotMessage = %v, want %v", gotMessage, tt.wantMessage)
			}
			if fileExist("./testData/testData1") {
				os.RemoveAll("./testData/testData1")
			}
		})
	}

}
