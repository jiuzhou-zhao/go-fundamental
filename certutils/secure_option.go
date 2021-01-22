package certutils

// SecureOption class
type SecureOption struct {
	VerifyClient bool `yaml:"verify_client,omitempty" json:"verify_client,omitempty"`
	VerifyServer bool `yaml:"verify_server,omitempty" json:"verify_server,omitempty"`
	// 服务端是否启用了TLS
	ServerWithTLS bool `yaml:"enable_tls,omitempty" json:"server_with_tls,omitempty"`

	ServerCertFile string `yaml:"server_cert_file,omitempty" json:"server_cert_file,omitempty"`
	ServerKeyFile  string `yaml:"server_key_file,omitempty" json:"server_key_file,omitempty"`
	ClientCertFile string `yaml:"client_cert_file,omitempty" json:"client_cert_file,omitempty"`
	ClientKeyFile  string `yaml:"client_key_file,omitempty" json:"client_key_file,omitempty"`

	// RootCAFiles 验证客户端和服务端都是用的ca文件列表
	RootCAFiles []string `yaml:"root_ca_files,omitempty" json:"root_ca_files,omitempty"`
	// ClientCAFiles 服务端用来验证客户端证书的ca文件列表
	ClientCAFiles []string `yaml:"client_ca_files,omitempty" json:"client_ca_files,omitempty"`
	// ServerCAFiles 客户端用来验证服务端的ca证书列表
	ServerCAFiles []string `yaml:"server_ca_files,omitempty" json:"server_ca_files,omitempty"`
	// 是客戶端用来验证服务端证书名字的，对服务端无效
	ServerName string `yaml:"server_name,omitempty" json:"server_name,omitempty"`

	// CertSignatures .
	CertSignatures map[string]string `yaml:"cert_signatures,omitempty" json:"cert_signatures,omitempty"`
}
