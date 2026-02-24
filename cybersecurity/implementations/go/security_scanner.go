package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// SecurityScanner performs various security scans
type SecurityScanner struct {
	timeout time.Duration
	threads int
}

// PortScanResult represents the result of a port scan
type PortScanResult struct {
	Port   int
	Open   bool
	Banner string
}

// TLSInfo contains TLS/SSL information
type TLSInfo struct {
	Version     string
	CipherSuite string
	Certificate struct {
		Subject  string
		Issuer   string
		NotAfter time.Time
	}
}

// VulnerabilityResult represents a potential vulnerability
type VulnerabilityResult struct {
	Type        string
	Severity    string
	Description string
	Port        int
	Details     string
}

// NewSecurityScanner creates a new security scanner instance
func NewSecurityScanner(timeout time.Duration, threads int) *SecurityScanner {
	return &SecurityScanner{
		timeout: timeout,
		threads: threads,
	}
}

// PortScan performs a TCP port scan on the target
func (s *SecurityScanner) PortScan(target string, ports []int) []PortScanResult {
	var results []PortScanResult
	var mutex sync.Mutex
	var wg sync.WaitGroup

	// Create a channel to limit concurrent connections
	semaphore := make(chan struct{}, s.threads)

	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			semaphore <- struct{}{} // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			result := s.scanPort(target, p)

			mutex.Lock()
			results = append(results, result)
			mutex.Unlock()
		}(port)
	}

	wg.Wait()

	// Sort results by port number
	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	return results
}

// scanPort scans a single port
func (s *SecurityScanner) scanPort(target string, port int) PortScanResult {
	address := fmt.Sprintf("%s:%d", target, port)
	conn, err := net.DialTimeout("tcp", address, s.timeout)

	result := PortScanResult{Port: port, Open: false}

	if err != nil {
		return result
	}

	result.Open = true
	defer conn.Close()

	// Try to grab banner
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err == nil && n > 0 {
		result.Banner = strings.TrimSpace(string(buffer[:n]))
	}

	return result
}

// HTTPHeaderScan checks for security headers
func (s *SecurityScanner) HTTPHeaderScan(url string) map[string]string {
	client := &http.Client{
		Timeout: s.timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return map[string]string{"error": err.Error()}
	}
	defer resp.Body.Close()

	headers := make(map[string]string)

	// Check for important security headers
	securityHeaders := []string{
		"Strict-Transport-Security",
		"Content-Security-Policy",
		"X-Frame-Options",
		"X-Content-Type-Options",
		"X-XSS-Protection",
		"Referrer-Policy",
		"Permissions-Policy",
	}

	for _, header := range securityHeaders {
		if value := resp.Header.Get(header); value != "" {
			headers[header] = value
		} else {
			headers[header] = "MISSING"
		}
	}

	// Add server information
	headers["Server"] = resp.Header.Get("Server")
	headers["X-Powered-By"] = resp.Header.Get("X-Powered-By")

	return headers
}

// TLSScan analyzes TLS/SSL configuration
func (s *SecurityScanner) TLSScan(target string, port int) (*TLSInfo, error) {
	address := fmt.Sprintf("%s:%d", target, port)

	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: s.timeout},
		"tcp",
		address,
		&tls.Config{InsecureSkipVerify: true},
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	state := conn.ConnectionState()

	tlsInfo := &TLSInfo{}

	// Get TLS version
	switch state.Version {
	case tls.VersionTLS10:
		tlsInfo.Version = "TLS 1.0"
	case tls.VersionTLS11:
		tlsInfo.Version = "TLS 1.1"
	case tls.VersionTLS12:
		tlsInfo.Version = "TLS 1.2"
	case tls.VersionTLS13:
		tlsInfo.Version = "TLS 1.3"
	default:
		tlsInfo.Version = "Unknown"
	}

	// Get cipher suite
	tlsInfo.CipherSuite = tls.CipherSuiteName(state.CipherSuite)

	// Get certificate information
	if len(state.PeerCertificates) > 0 {
		cert := state.PeerCertificates[0]
		tlsInfo.Certificate.Subject = cert.Subject.String()
		tlsInfo.Certificate.Issuer = cert.Issuer.String()
		tlsInfo.Certificate.NotAfter = cert.NotAfter
	}

	return tlsInfo, nil
}

// VulnerabilityCheck performs basic vulnerability checks
func (s *SecurityScanner) VulnerabilityCheck(target string, portResults []PortScanResult) []VulnerabilityResult {
	var vulnerabilities []VulnerabilityResult

	for _, result := range portResults {
		if !result.Open {
			continue
		}

		// Check for common vulnerable services
		switch result.Port {
		case 21: // FTP
			if strings.Contains(strings.ToLower(result.Banner), "vsftpd 2.3.4") {
				vulnerabilities = append(vulnerabilities, VulnerabilityResult{
					Type:        "Backdoor",
					Severity:    "Critical",
					Description: "vsftpd 2.3.4 backdoor vulnerability",
					Port:        result.Port,
					Details:     "This version contains a backdoor that allows remote code execution",
				})
			}
		case 22: // SSH
			if strings.Contains(strings.ToLower(result.Banner), "openssh") {
				// Check for old versions (simplified check)
				if strings.Contains(result.Banner, "OpenSSH_6.") || strings.Contains(result.Banner, "OpenSSH_5.") {
					vulnerabilities = append(vulnerabilities, VulnerabilityResult{
						Type:        "Outdated Software",
						Severity:    "Medium",
						Description: "Outdated OpenSSH version detected",
						Port:        result.Port,
						Details:     fmt.Sprintf("Banner: %s", result.Banner),
					})
				}
			}
		case 23: // Telnet
			vulnerabilities = append(vulnerabilities, VulnerabilityResult{
				Type:        "Insecure Protocol",
				Severity:    "High",
				Description: "Telnet service detected (unencrypted)",
				Port:        result.Port,
				Details:     "Telnet transmits data in plaintext, including passwords",
			})
		case 80, 8080: // HTTP
			vulnerabilities = append(vulnerabilities, VulnerabilityResult{
				Type:        "Insecure Protocol",
				Severity:    "Medium",
				Description: "HTTP service detected (unencrypted)",
				Port:        result.Port,
				Details:     "Consider using HTTPS instead of HTTP",
			})
		case 139, 445: // SMB
			vulnerabilities = append(vulnerabilities, VulnerabilityResult{
				Type:        "Information Disclosure",
				Severity:    "Medium",
				Description: "SMB service detected",
				Port:        result.Port,
				Details:     "SMB services can leak system information",
			})
		case 1433: // SQL Server
			vulnerabilities = append(vulnerabilities, VulnerabilityResult{
				Type:        "Database Exposure",
				Severity:    "High",
				Description: "SQL Server detected on default port",
				Port:        result.Port,
				Details:     "Database services should not be directly accessible from external networks",
			})
		case 3389: // RDP
			vulnerabilities = append(vulnerabilities, VulnerabilityResult{
				Type:        "Remote Access",
				Severity:    "High",
				Description: "RDP service detected",
				Port:        result.Port,
				Details:     "RDP should be protected with VPN or strong authentication",
			})
		}

		// Check for default/common ports that might indicate services
		if result.Port > 1024 && result.Port < 65535 && result.Banner == "" {
			vulnerabilities = append(vulnerabilities, VulnerabilityResult{
				Type:        "Unknown Service",
				Severity:    "Low",
				Description: "Unknown service on non-standard port",
				Port:        result.Port,
				Details:     "Service fingerprinting recommended",
			})
		}
	}

	return vulnerabilities
}

// GenerateReport creates a comprehensive security report
func (s *SecurityScanner) GenerateReport(target string, portResults []PortScanResult, headers map[string]string, tlsInfo *TLSInfo, vulnerabilities []VulnerabilityResult) {
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("SECURITY SCAN REPORT FOR: %s\n", target)
	fmt.Printf(strings.Repeat("=", 60) + "\n\n")

	// Port scan results
	fmt.Println("PORT SCAN RESULTS:")
	fmt.Println(strings.Repeat("-", 40))
	openPorts := 0
	for _, result := range portResults {
		if result.Open {
			openPorts++
			fmt.Printf("Port %d: OPEN", result.Port)
			if result.Banner != "" {
				fmt.Printf(" - %s", result.Banner)
			}
			fmt.Println()
		}
	}
	fmt.Printf("Total open ports: %d\n\n", openPorts)

	// HTTP headers analysis
	if len(headers) > 0 && headers["error"] == "" {
		fmt.Println("HTTP SECURITY HEADERS:")
		fmt.Println(strings.Repeat("-", 40))
		for header, value := range headers {
			status := "✓"
			if value == "MISSING" {
				status = "✗"
			}
			fmt.Printf("%s %s: %s\n", status, header, value)
		}
		fmt.Println()
	}

	// TLS information
	if tlsInfo != nil {
		fmt.Println("TLS/SSL INFORMATION:")
		fmt.Println(strings.Repeat("-", 40))
		fmt.Printf("Version: %s\n", tlsInfo.Version)
		fmt.Printf("Cipher Suite: %s\n", tlsInfo.CipherSuite)
		if tlsInfo.Certificate.Subject != "" {
			fmt.Printf("Certificate Subject: %s\n", tlsInfo.Certificate.Subject)
			fmt.Printf("Certificate Expires: %s\n", tlsInfo.Certificate.NotAfter.Format("2006-01-02"))
		}
		fmt.Println()
	}

	// Vulnerabilities
	if len(vulnerabilities) > 0 {
		fmt.Println("POTENTIAL VULNERABILITIES:")
		fmt.Println(strings.Repeat("-", 40))

		// Group by severity
		severityOrder := []string{"Critical", "High", "Medium", "Low"}
		for _, severity := range severityOrder {
			found := false
			for _, vuln := range vulnerabilities {
				if vuln.Severity == severity {
					if !found {
						fmt.Printf("\n%s SEVERITY:\n", strings.ToUpper(severity))
						found = true
					}
					fmt.Printf("  • %s (Port %d)\n", vuln.Description, vuln.Port)
					fmt.Printf("    Type: %s\n", vuln.Type)
					fmt.Printf("    Details: %s\n", vuln.Details)
					fmt.Println()
				}
			}
		}
	} else {
		fmt.Println("No obvious vulnerabilities detected.\n")
	}

	// Recommendations
	fmt.Println("RECOMMENDATIONS:")
	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("• Close unnecessary open ports")
	fmt.Println("• Implement proper security headers for web services")
	fmt.Println("• Use strong TLS configurations (TLS 1.2+)")
	fmt.Println("• Regular security updates and patches")
	fmt.Println("• Implement network segmentation")
	fmt.Println("• Use intrusion detection systems")
	fmt.Println("• Regular vulnerability assessments")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run security_scanner.go <target>")
		fmt.Println("Example: go run security_scanner.go example.com")
		os.Exit(1)
	}

	target := os.Args[1]

	// Initialize scanner
	scanner := NewSecurityScanner(3*time.Second, 50)

	fmt.Printf("Starting security scan of %s...\n", target)

	// Common ports to scan
	commonPorts := []int{
		21, 22, 23, 25, 53, 80, 110, 111, 135, 139, 143, 443, 993, 995,
		1723, 3306, 3389, 5432, 5900, 8080, 8443,
	}

	// Perform port scan
	fmt.Println("Performing port scan...")
	portResults := scanner.PortScan(target, commonPorts)

	// HTTP header scan (if web server detected)
	var headers map[string]string
	for _, result := range portResults {
		if result.Open && (result.Port == 80 || result.Port == 443 || result.Port == 8080 || result.Port == 8443) {
			fmt.Println("Analyzing HTTP security headers...")
			protocol := "http"
			if result.Port == 443 || result.Port == 8443 {
				protocol = "https"
			}
			url := fmt.Sprintf("%s://%s:%d", protocol, target, result.Port)
			headers = scanner.HTTPHeaderScan(url)
			break
		}
	}

	// TLS scan (if HTTPS detected)
	var tlsInfo *TLSInfo
	for _, result := range portResults {
		if result.Open && (result.Port == 443 || result.Port == 8443) {
			fmt.Println("Analyzing TLS configuration...")
			info, err := scanner.TLSScan(target, result.Port)
			if err == nil {
				tlsInfo = info
			}
			break
		}
	}

	// Vulnerability check
	fmt.Println("Checking for potential vulnerabilities...")
	vulnerabilities := scanner.VulnerabilityCheck(target, portResults)

	// Generate report
	scanner.GenerateReport(target, portResults, headers, tlsInfo, vulnerabilities)
}
