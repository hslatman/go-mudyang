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
	"github.com/hslatman/magefiles/targets" // shared targets

	"github.com/hslatman/magefiles/release"
)

var (
	// packageRegex is a regular expression to find a package name similar to github.com/openconfig/ygot
	packageRegex = regexp.MustCompile(`(?m)/.*github\.com/openconfig/ygot@[^/]+/genutil/names\.go`)
)

// Generate generates the code from YANG files by calling
// the ygot generator. The generator name is overridden in the
// generated code to be more informative.
func Generate(ctx context.Context) error {
	mg.Deps(targets.Tools)

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
		"yang/ietf-packet-fields@2019-03-04.yang",       // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-packet-fields@2019-03-04.yang
		"yang/ietf-ethertypes@2019-03-04.yang",          // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-ethertypes@2019-03-04.yang
		"yang/ietf-acldns@2019-01-28.yang",              // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-acldns@2019-01-28.yang
		"yang/ietf-access-control-list@2019-03-04.yang", // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-access-control-list@2019-03-04.yang
		"yang/ietf-inet-types@2025-12-22.yang",          // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-inet-types@2025-12-22.yang
		"yang/iana-tls-profile@2025-04-18.yang",         // NOTE: sourced from https://www.yangcatalog.org/all_modules/iana-tls-profile@2025-04-18.yang
		"yang/ietf-acl-tls@2025-04-18.yang",             // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-acl-tls@2025-04-18.yang
		"yang/iana-hash-algs@2020-03-08.yang",           // NOTE: sourced from https://www.yangcatalog.org/all_modules/iana-hash-algs@2020-03-08.yang
		"yang/ietf-netconf-acm@2018-02-14.yang",         // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-netconf-acm@2018-02-14.yang
		"yang/ietf-crypto-types@2024-10-10.yang",        // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-crypto-types@2024-10-10.yang
		"yang/ietf-mud-transparency@2023-10-10.yang",    // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-mud-transparency@2023-10-10.yang
		"yang/ietf-owner-license@2026-06-28.yang",       // NOTE: sourced from https://www.yangcatalog.org/all_modules/ietf-owner-license@2026-06-28.yang
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

// Lint runs the linter
func Lint(ctx context.Context) error {
	mg.Deps(targets.Tools)
	args := []string{"tool", "-modfile=./.tools/go.mod", "github.com/golangci/golangci-lint/v2/cmd/golangci-lint", "run", "--config", ".golangci.yml"}
	return sh.RunV("go", args...)
}

func Check(ctx context.Context) error {
	mg.Deps(release.Tools)

	return release.Check()
}
