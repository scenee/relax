package relfile

import (
	"fmt"
	"os"
	"testing"
)

func TestNewProvisioningProfile(t *testing.T) {
	var (
		in string
		pp *ProvisioningProfile
	)

	ClearCache()

	in = "../../../../sample/certificates/Relax_AdHoc.mobileprovision"
	pp = NewProvisioningProfile(in)
	if pp.TeamName != "Shin Yamamoto" {
		t.Errorf("TeamName failed")
	}
	if pp.TeamID() != "J3D7L9FHSS" {
		t.Errorf("TeamID failed %v", pp.TeamIdentifiers)
	}
	if pp.Name != "Relax AdHoc" {
		t.Errorf("Name failed")
	}
	if pp.Entitlements.GetTaskAllow {
		t.Errorf("failed")
	}

	in = "../../../../sample/certificates/Relax_Development.mobileprovision"
	pp = NewProvisioningProfile(in)
	if pp.TeamName != "Shin Yamamoto" {
		t.Errorf("TeamName failed")
	}
	if pp.TeamID() != "J3D7L9FHSS" {
		t.Errorf("TeamID failed %v", pp.TeamIdentifiers)
	}
	if pp.Name != "Relax Development" {
		t.Errorf("Name failed")
	}
	if pp.Entitlements.GetTaskAllow == false {
		t.Errorf("failed %v", pp.Entitlements.GetTaskAllow)
	}

}

func TestFindProvisioningProfile(t *testing.T) {
	ClearCache()
	infos := FindProvisioningProfile("Relax", "")
	if len(infos) != 2 {
		t.Errorf("Failed")
	}
	for _, info := range infos {
		fmt.Println(info.Name, info.Pp.Name)
		if _, err := os.Stat(info.Name); err != nil {
			t.Errorf("Not Found %v", info.Name)
		}
	}
}

func BenchmarkFindPP(b *testing.B) {
	ClearCache()
	_ = FindProvisioningProfile("Relax", "")
}

func BenchmarkFindPPFast(b *testing.B) {
	ClearCache()
	_ = FindProvisioningProfile("Relax", "")
	b.ResetTimer()
	_ = FindProvisioningProfile("Relax", "")
}
