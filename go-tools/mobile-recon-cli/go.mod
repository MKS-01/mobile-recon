module github.com/MKS-01/mobile-recon/go-tools/mobile-recon-cli

go 1.21

require (
	github.com/MKS-01/mobile-recon/go-tools/adb-toolkit v0.0.0
	github.com/MKS-01/mobile-recon/go-tools/apk-analyzer v0.0.0
	github.com/MKS-01/mobile-recon/go-tools/common v0.0.0
	github.com/MKS-01/mobile-recon/go-tools/ios-toolkit v0.0.0
	github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit v0.0.0
	github.com/fatih/color v1.18.0
	github.com/spf13/cobra v1.10.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/ulikunitz/xz v0.5.15 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

replace (
	github.com/MKS-01/mobile-recon/go-tools/adb-toolkit => ../adb-toolkit
	github.com/MKS-01/mobile-recon/go-tools/apk-analyzer => ../apk-analyzer
	github.com/MKS-01/mobile-recon/go-tools/common => ../common
	github.com/MKS-01/mobile-recon/go-tools/ios-toolkit => ../ios-toolkit
	github.com/MKS-01/mobile-recon/go-tools/nmap-toolkit => ../nmap-toolkit
)
