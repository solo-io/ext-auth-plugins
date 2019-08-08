package main

import (
	authapi "github.com/solo-io/ext-auth-plugins/api"
	pluginapi "github.com/solo-io/ext-auth-plugins/examples/header/api"
)

func main() {}

var _ authapi.ExtAuthPlugin = new(pluginapi.RequiredHeaderPlugin)

//noinspection GoUnusedGlobalVariable
var Plugin pluginapi.RequiredHeaderPlugin
