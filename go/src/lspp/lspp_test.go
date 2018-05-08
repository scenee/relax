package main

import (
	"os"
	"testing"
)

func TestFindProvisioningProfile(t *testing.T) {
	for _, out := range findProvisioningProfile("Relax", "", true, false) {
		if _, err := os.Stat(out); err != nil {
			t.Error(err)
		}
	}
}
