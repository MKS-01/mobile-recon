module github.com/MKS-01/mobile-recon/go-tools/mobile-recon-cli

go 1.21

require (
	github.com/MKS-01/mobile-recon/go-tools/common v0.0.0
	github.com/fatih/color v1.18.0
	github.com/manifoldco/promptui v0.9.0
	github.com/spf13/cobra v1.8.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

replace github.com/MKS-01/mobile-recon/go-tools/common => ../common
