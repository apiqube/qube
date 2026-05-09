package commands

import "github.com/apiqube/engine"

// mockPluginSchemas returns a synthetic plugin list for table-rendering tests.
func mockPluginSchemas() []engine.PluginSchema {
	return []engine.PluginSchema{
		{
			Name:         "http",
			Version:      "0.1.0",
			Protocols:    []engine.Protocol{engine.ProtocolHTTP, engine.ProtocolHTTPS},
			Capabilities: []string{"http"},
		},
		{
			Name:         "demo",
			Version:      "0.0.1",
			Protocols:    []engine.Protocol{"demo"},
			Capabilities: nil,
		},
	}
}
