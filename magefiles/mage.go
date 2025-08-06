//go:build mage

package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/magefile/mage/mg" // TODO: can we depend on these with them only being in go.tools.mod?
	"github.com/magefile/mage/sh"

	// mage:import
	_ "github.com/hslatman/magefiles/targets" // shared targets
)

var (
	// packageRegex is a regular expression to find a package name similar to github.com/openconfig/ygot
	packageRegex = regexp.MustCompile(`(?m)/.*github\.com/openconfig/ygot@[^/]+/genutil/names\.go`)
)

// Generate generates the code from YANG files by calling
// the ygot generator. The generator name is overridden in the
// generated code to be more informative.
func Generate(ctx context.Context) error {
	mg.Deps(Tools)

	args := []string{
		"tool",
		"-modfile=./.tools/go.mod",
		"github.com/openconfig/ygot/generator",
		"-path=yang",
		"-output_file=-",
		"-generate_simple_unions",
		"-package_name=mudyang",
		"-generate_fakeroot",
		"-fakeroot_name=mudfile",
		"yang/ietf-packet-fields@2019-03-04.yang",
		"yang/ietf-ethertypes@2019-03-04.yang",
		"yang/ietf-acldns.yang",
		"yang/ietf-access-control-list@2019-03-04.yang", // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-access-control-list@2019-03-04.yang
		"yang/ietf-inet-types@2024-10-21.yang",          // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-inet-types@2024-10-21.yang
		"yang/iana-tls-profile@2025-04-18.yang",         // NOTE: sourced from https://www.yangcatalog.org/all_modules/iana-tls-profile@2025-04-18.yang
		"yang/ietf-acl-tls@2025-04-18.yang",             // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-acl-tls@2025-04-18.yang
		"yang/iana-hash-algs.yang",                      // NOTE: sourced from https://raw.githubusercontent.com/YangModels/yang/3af23949e11a2acd2f36df1dc0afca73ffe118ac/experimental/ietf-extracted-YANG-modules/iana-hash-algs@2020-03-08.yang
		"yang/ietf-netconf-acm.yang",                    // NOTE: sourced from https://raw.githubusercontent.com/huawei/yang/855d2d384d49fea03872e75fcea4d40619cf3528/network-router/8.20.0/atn980b/ietf-netconf-acm.yang
		"yang/ietf-crypto-types@2021-09-14.yang",        // NOTE: sourced from https://yangcatalog.org/YANG-modules/
		"yang/ietf-mud-transparency@2023-10-10.yang",    // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-mud-transparency@2023-10-10.yang
		"yang/ietf-ol@2024-04-26.yang",                  // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-ol@2024-04-26.yang
		"yang/ietf-mud-tls@2025-04-18.yang",             // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-mud-tls@2025-04-18.yang
		"yang/ietf-mud@2019-01-28.yang",                 // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-mud@2019-01-28.yang
	}

	out, err := sh.Output("go", args...)
	if err != nil {
		return err
	}

	modVersion, err := sh.Output("go", "list", "-m", "github.com/openconfig/ygot")
	if err != nil {
		return err
	}

	moduleNameAndVersion := fmt.Sprintf("github.com/openconfig/ygot/generator@%s", strings.Split(modVersion, " ")[1])

	result := packageRegex.ReplaceAllString(out, moduleNameAndVersion)

	return os.WriteFile("mudyang.go", []byte(result), 0644)
}

// Tools ensures the tools get installed
func Tools() error {
	return sh.RunV("go", "mod", "tidy", "-modfile=./.tools/go.mod")
}

// Lint runs the linter
func Lint(ctx context.Context) error {
	mg.Deps(Tools)
	args := []string{"tool", "-modfile=./.tools/go.mod", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint", "run", "--config", ".golangci.yml"}
	return sh.RunV("go", args...)
}
