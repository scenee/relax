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

	in = "mobileprovisions/AdHoc.mobileprovision"
	pp = newProvisioningProfile(in)
	if pp.TeamName != "Relax" {
		t.Errorf("TeamName failed")
	}
	if pp.TeamID() != "AAAAAAAAAA" {
		t.Errorf("TeamID failed %v", pp.TeamIdentifiers)
	}
	if pp.Name != "Sample Ad Hoc" {
		t.Errorf("Name failed")
	}
	if pp.Entitlements.GetTaskAllow {
		t.Errorf("failed")
	}
	if pp.ProvisioningType() != ProvisioningTypeAdHoc {
		t.Errorf("ProvisioningType failed")
	}
	if pp.CertificateType() != CertificateTypeDistribution {
		t.Errorf("CertificateType failed")
	}

	in = "mobileprovisions/Development.mobileprovision"
	pp = newProvisioningProfile(in)
	if pp.TeamName != "Relax" {
		t.Errorf("TeamName failed")
	}
	if pp.TeamID() != "AAAAAAAAAA" {
		t.Errorf("TeamID failed %v", pp.TeamIdentifiers)
	}
	if pp.Name != "Sample Development" {
		t.Errorf("Name failed")
	}
	if pp.Entitlements.GetTaskAllow == false {
		t.Errorf("failed %v", pp.Entitlements.GetTaskAllow)
	}
	if pp.ProvisioningType() != ProvisioningTypeDevelopment {
		t.Errorf("ProvisioningType failed")
	}
	if pp.CertificateType() != CertificateTypeDeveloper {
		t.Errorf("CertificateType failed")
	}

	in = "mobileprovisions/Enterprise.mobileprovision"
	pp = newProvisioningProfile(in)
	if pp.ProvisioningType() != ProvisioningTypeEnterprise {
		t.Errorf("ProvisioningType failed")
	}
	if pp.CertificateType() != CertificateTypeDistribution {
		t.Errorf("CertificateType failed")
	}

	in = "mobileprovisions/AppStore.mobileprovision"
	pp = newProvisioningProfile(in)
	if pp.ProvisioningType() != ProvisioningTypeAppStore {
		t.Errorf("ProvisioningType failed")
	}
	if pp.CertificateType() != CertificateTypeDistribution {
		t.Errorf("CertificateType failed")
	}
}

func TestFindProvisioningProfile(t *testing.T) {
	ClearCache()
	infos := FindProvisioningProfile("Relax", "")

	if len(infos) != 2 {
		t.Errorf("Failed %v", infos)
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
