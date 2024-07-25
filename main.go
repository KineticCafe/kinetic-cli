package main

import "github.com/KineticCommerce/kinetic-cli/cmd"

var BuildTime = "unset"

func main() {
	cmd.Execute(BuildTime)
}
