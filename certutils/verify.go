package certutils

import (
	"errors"
	"fmt"
)

// CertNameToSignatureKey .
func CertNameToSignatureKey(certName string) string {
	return "cert_" + certName + "_sig"
}

// QueryPublicKeySignature .
func QueryPublicKeySignature(certName string, secureOption *SecureOption) (string, error) {
	if certName == "" || secureOption == nil || len(secureOption.CertSignatures) <= 0 {
		return "", errors.New("invalid input parameters")
	}
	if signature, ok := secureOption.CertSignatures[CertNameToSignatureKey(certName)]; ok {
		return signature, nil
	}
	return "", fmt.Errorf("cert %v no signature record", certName)
}

// VerifyCertPublicKey .
func VerifyCertPublicKey(publicKey interface{}, certName string, secureOption *SecureOption) (ok bool, err error) {
	publicKeyHash, err := GetPublicKeyHash(publicKey)
	if err != nil {
		err = fmt.Errorf("GetPublicKeyHash failed: %v", err)
		return
	}
	signature, err := QueryPublicKeySignature(certName, secureOption)
	if err != nil {
		err = fmt.Errorf("QueryPublicKeySignature failed: %v", err)
		return
	}

	if publicKeyHash != signature {
		return
	}
	ok = true
	return
}
