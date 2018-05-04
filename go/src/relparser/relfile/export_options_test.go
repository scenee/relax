package relfile

import (
	"testing"
)

func TestEncodeExportOptionsPlist(t *testing.T) {
	// fmt.Println("export options:", d.ExportOptions)
	var opts *ExportOptions

	opts = &ExportOptions{}
	opts.SetProvisioningProfiles("com.scenee.Sample", "J3D7L9FHSS", "Relax AdHoc")

	if opts.Method != ProvisioningTypeAdHoc {
		t.Fatalf("Invalid 'method'")
	}
	t.Logf("Method = %v", opts.Method)

	if opts.SigningCertificate != CertificateTypeDistribution {
		t.Fatalf("Invalid 'signingCertificate'")
	}
	t.Logf("SigningCertificate = %v", opts.SigningCertificate)

	opts = &ExportOptions{}
	opts.SetProvisioningProfiles("com.scenee.Sample", "J3D7L9FHSS", "Relax Development")

	if opts.Method != ProvisioningTypeDevelopment {
		t.Fatalf("Invalid 'method'")
	}
	t.Logf("Method = %v", opts.Method)

	if opts.SigningCertificate != CertificateTypeDeveloper {
		t.Fatalf("Invalid 'signingCertificate'")
	}
	t.Logf("SigningCertificate = %v", opts.SigningCertificate)

}
