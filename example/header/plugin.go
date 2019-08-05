package main

import (
	"github.com/solo-io/ext-auth-plugins/api"
	api2 "github.com/solo-io/ext-auth-plugins/example/header/api"
)

func main() {}

var _ api.ExtauthPlugin = new(api2.RequiredHeaderPlugin)

//noinspection GoUnusedGlobalVariable
var Plugin api2.RequiredHeaderPlugin
