package authenticator

import (
	"testing"
)

func Test_CreateGoogleAuthQRCodeData(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "CreateGoogleAuthQRCodeData", want: "otpauth://totp/issuer:account?secret=key&issuer=issuer"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			qr := CreateGoogleAuthQRCodeData("key", "account", "issuer")
			if qr != tt.want {
				t.Fatalf("[CreateGoogleAuthQRCodeData] qr: %s, want: %s", qr, tt.want)
			}
		})
	}
}

func Test_MakeGoogleAuthenticator(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "MakeGoogleAuthenticator", want: "215669"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := MakeGoogleAuthenticator("HXDMVJECJJWSRB3HWIZR4IFUGFTMXBOZ", 123123)
			if err != nil {
				t.Fatalf("[MakeGoogleAuthenticator] error: %s", err)
			}
			if code != tt.want {
				t.Fatalf("[MakeGoogleAuthenticator] code: %s, want: %s", code, tt.want)
			}
		})
	}
}
