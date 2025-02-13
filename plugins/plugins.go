package plugins

type Plugin struct {
	FilePath string
	MTime    int64
	Cycle    int
}

var (
	Plugins              = make(map[string]*Plugin)
	PluginsWithScheduler = make(map[string]*PluginScheduler)
)

func DelNoUsePlugins(newPlugins map[string]*Plugin) {
	for currKey, currPlugin := range Plugins {
		newPlugin, ok := newPlugins[currKey]
		// 新的plugins不包含旧的plugin, 或者修改时间变化了
		if !ok || currPlugin.MTime != newPlugin.MTime {
			// 删除Plugins中的元素，不会导致iterator失效?
			deletePlugin(currKey)
		}
	}
}

func AddNewPlugins(newPlugins map[string]*Plugin) {
	for fpath, newPlugin := range newPlugins {
		if _, ok := Plugins[fpath]; ok && newPlugin.MTime == Plugins[fpath].MTime {
			continue
		}

		// Plugin/Schedule
		Plugins[fpath] = newPlugin
		sch := NewPluginScheduler(newPlugin)
		PluginsWithScheduler[fpath] = sch
		sch.Schedule()
	}
}

func ClearAllPlugins() {
	for k := range Plugins {
		deletePlugin(k)
	}
}

func deletePlugin(key string) {
	v, ok := PluginsWithScheduler[key]
	if ok {
		v.Stop()
		delete(PluginsWithScheduler, key)
	}
	delete(Plugins, key)
}
