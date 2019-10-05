package main

import (
	"flag"
	"fmt"
	"relparser/relfile"
)

var appID string

func init() {
	flag.StringVar(&appID, "-app-id", "", "Check a provisioning profile for the application identifier.")
}
func main() {
	flag.Parse()
	var cmd string = flag.Arg(0)

	switch cmd {
	case "check":
		var path string = flag.Arg(1)
		pp := relfile.MakeProvisioningProfile(path)
		ids := pp.GetValidIdentities()
		if len(ids) > 0 {
			fmt.Printf("%v,%v\n", ids[0].Sha1, ids[0].Name)
		}
	default:
		break
	}

}
