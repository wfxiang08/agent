package g

import (
	"bytes"
	"fmt"
	"github.com/toolkits/file"
	"os/exec"
	"strings"
)

func GetCurrPluginVersion() string {
	if !Config().Plugin.Enabled {
		return "plugin not enabled"
	}

	pluginDir := Config().Plugin.Dir
	if !file.IsExist(pluginDir) {
		return "plugin dir not existent"
	}

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = pluginDir

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return fmt.Sprintf("Error:%s", err.Error())
	}

	// 如何管理Plugin呢?
	// 所有的Plugin放在指定的git目录中，然后通过Git的版本号来管理
	return strings.TrimSpace(out.String())
}
