package common

import (
	"fmt"
)

//InitCommon 初始化 config 和 log
func InitCommon(configPath string, isCrontab bool) error {
	err := initConfig(configPath, isCrontab)
	if err != nil {
		return fmt.Errorf("unable to initialize config: %v", err)
	}
	initLoggor(globalConf.Log.Dir, globalConf.Log.Level)
	return nil
}
