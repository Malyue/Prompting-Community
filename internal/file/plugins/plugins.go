package plugins

import "go.uber.org/zap"

type Plugin interface {
	Flag() bool
	Name() string
	New() interface{}
	Health()
	Close()
}

var Plugins = make(map[string]Plugin)

func RegisterPlugin(plugin Plugin) {
	Plugins[plugin.Name()] = plugin
}

func NewPlugins() {
	for _, p := range Plugins {
		if !p.Flag() {
			continue
		}
		zap.L().Info("Init " + p.Name())
		p.New()
		zap.L().Info("HealthCheck ... " + p.Name())
		p.Health()
		zap.L().Info("Success Init . " + p.Name())
	}
}

func ClosePlugins() {
	for _, p := range Plugins {
		if !p.Flag() {
			continue
		}
		p.Close()
		zap.L().Info("Closed" + p.Name())
	}
}
