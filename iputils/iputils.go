package iputils

import (
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"golang.org/x/net/context"
	"google.golang.org/grpc/peer"
)

const (
	localIPV4 = "127.0.0.1"
)

func regularIP(ip string) string {
	if ip == "::1" {
		return localIPV4
	}
	return ip
}

// GrpcGetRealIP get client ip from context of grpce
func GrpcGetRealIP(ctx context.Context) string {
	clientIP := metautils.ExtractIncoming(ctx).Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(metautils.ExtractIncoming(ctx).Get("X-Real-Ip"))
	}
	if clientIP == "" {
		client, ok := peer.FromContext(ctx)
		if ok && client.Addr != net.Addr(nil) {
			addSlice := strings.Split(client.Addr.String(), ":")
			if addSlice[0] == "[" {
				return localIPV4
			}
			clientIP = addSlice[0]
		}
	}
	return regularIP(clientIP)
}

// HttpGetRealIP get client ip from http request
func HttpGetRealIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = r.Header.Get("X-Real-Ip")
	}
	if clientIP == "" {
		clientIP = r.RemoteAddr
		addSlice := strings.Split(clientIP, ":")
		if addSlice[0] == "[" {
			return localIPV4
		}
		clientIP = addSlice[0]
	}
	return regularIP(clientIP)
}

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

// IsPrivateIP function
func IsPrivateIP(ip net.IP) bool {
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// GetOutboundIP functions
func GetOutboundIP() (ip net.IP, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer func() {
		_ = conn.Close()
	}()
	ip = conn.LocalAddr().(*net.UDPAddr).IP
	return
}

// LocalIPs return all non-loopback IPv4 addresses
func LocalIPv4s() ([]string, error) {
	mainIPV4 := ""
	ip, err := GetOutboundIP()
	if err == nil && ip.To4() != nil {
		mainIPV4 = ip.To4().String()
	}
	ips := make([]string, 0)
	if mainIPV4 != "" {
		ips = append(ips, mainIPV4)
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && ipnet.IP.IsGlobalUnicast() && ipnet.IP.To4() != nil {
			// remove 169.254
			ipV4 := ipnet.IP.String()
			if ipV4 == mainIPV4 {
				continue
			}
			ips = append(ips, ipV4)
		}
	}

	return ips, nil
}

// GetIPv4ByInterface return IPv4 address from a specific interface IPv4 addresses
func GetIPv4ByInterface(name string) ([]string, error) {
	var ips []string

	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}
