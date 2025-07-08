package main

import "goweb/cmd"

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Provide the accessToken as a Bearer token: 'Bearer {accessToken}'
func main() {
	cmd.Execute()
}
