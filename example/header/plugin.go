package main

import (
	authapi "github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/ext-auth-plugins/example/header/api"
)

func main() {}

var _ authapi.ExtauthPlugin = new(api.RequiredHeaderPlugin)

//noinspection GoUnusedGlobalVariable
var Plugin api.RequiredHeaderPlugin
