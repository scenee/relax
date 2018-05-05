package relfile

import (
	"testing"
)

func TestEncodeExportOptionsPlist(t *testing.T) {
	// fmt.Println("export options:", d.ExportOptions)
	var opts *ExportOptions

	opts = &ExportOptions{}
	opts.SetProvisioningProfiles("Relax AdHoc", "com.scenee.SampleApp")

	if opts.ProvisioningProfiles["com.scenee.SampleApp"] != "Relax AdHoc" {
		t.Fatalf("Invalid 'provisioningProfiles'")
	}

	if opts.Method != ProvisioningTypeAdHoc {
		t.Fatalf("Invalid 'method'")
	}
	t.Logf("Method = %v", opts.Method)

	if opts.SigningCertificate != CertificateTypeDistribution {
		t.Fatalf("Invalid 'signingCertificate'")
	}
	t.Logf("SigningCertificate = %v", opts.SigningCertificate)

	opts = &ExportOptions{}
	opts.SetProvisioningProfiles("Relax Development", "com.scenee.SampleApp")

	if opts.ProvisioningProfiles["com.scenee.SampleApp"] != "Relax Development" {
		t.Fatalf("Invalid 'provisioningProfiles'")
	}
	if opts.TeamID != "J3D7L9FHSS" {
		t.Fatalf("Invalid 'teamID'")
	}

	if opts.Method != ProvisioningTypeDevelopment {
		t.Fatalf("Invalid 'method'")
	}
	t.Logf("Method = %v", opts.Method)

	if opts.SigningCertificate != CertificateTypeDeveloper {
		t.Fatalf("Invalid 'signingCertificate'")
	}
	t.Logf("SigningCertificate = %v", opts.SigningCertificate)

}
