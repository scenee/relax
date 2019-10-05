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

func TestGetValidIdentities(t *testing.T) {
	infos := FindProvisioningProfile("Relax AdHoc", "")
	pp := infos[0].Pp
	ids := pp.GetValidIdentities()

	if len(ids) == 0 {
		t.Errorf("GetValidIdentities failed: not found identity for %v", pp.Name)
		return
	}
	var exp string

	id := ids[0]

	exp = "00D0F760D573CFAEBE09DB0E3E62B4F251999973"
	if id.Sha1 != exp {
		t.Errorf("GetValidIdentities failed: %v is not equal to %v", id.Sha1, exp)
	}
	exp = "iPhone Distribution: Shin Yamamoto (J3D7L9FHSS)"
	if id.Name != exp {
		t.Errorf("GetValidIdentities failed: %v is not equal to %v", id.Name, exp)
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
