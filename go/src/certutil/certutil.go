package certutil

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var logger = log.New(os.Stderr, "", 0)

const (
	AppleWWDRCAName           = "Apple Worldwide Developer Relations Certification Authority"
	AppleWWDRCACertificateURL = "https://developer.apple.com/certificationauthority/AppleWWDRCA.cer"
)

func InstallCertificate(cn string) {
	var (
		resp                *http.Response
		temp                *os.File
		tempName, sourceURL string
		err                 error
	)

	switch cn {
	case AppleWWDRCAName:
		tempName = "relax/AppleWWDRCA.cert"
		sourceURL = AppleWWDRCACertificateURL
	default:
		logger.Fatalf("\"%s\" isn't supported.", cn)
	}

	logger.Printf("Installing \"%s\"...", cn)
	if temp, err = ioutil.TempFile("", tempName); err != nil {
		logger.Fatalf("Failed to open temp file for \"%s\": %v", cn, err)
	}
	defer temp.Close()
	if resp, err = http.Get(sourceURL); err != nil {
		logger.Fatalf("Failed to download \"%s\": %v", cn, err)
	}
	defer resp.Body.Close()
	if _, err = io.Copy(temp, resp.Body); err != nil {
		logger.Fatalf("Failed to download \"%s\": %v", cn, err)
	}
	if _, err = exec.Command("/usr/bin/security", "add-certificates", temp.Name()).Output(); err != nil {
		logger.Fatalf("Failed to install \"%s\": %v", cn, err)
	}
}
