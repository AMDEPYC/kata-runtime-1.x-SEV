// Copyright (c) 2018 ARM Limited
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"fmt"
	"path/filepath"
	goruntime "runtime"
)

func getExpectedHostDetails(tmpdir string) (HostInfo, error) {
	type filesToCreate struct {
		file     string
		contents string
	}

	const expectedKernelVersion = "99.1"
	const expectedArch = goruntime.GOARCH

	expectedDistro := DistroInfo{
		Name:    "Foo",
		Version: "42",
	}

	expectedCPU := CPUInfo{
		Vendor: "0x41",
		Model:  "8",
	}

	expectedNormalizeCPU := CPUInfo{
		Vendor: "ARM Limited",
		Model:  "v8",
	}

	expectedHostDetails := HostInfo{
		Kernel:             expectedKernelVersion,
		Architecture:       expectedArch,
		Distro:             expectedDistro,
		CPU:                expectedNormalizeCPU,
		VMContainerCapable: true,
	}

	testProcCPUInfo := filepath.Join(tmpdir, "cpuinfo")
	testOSRelease := filepath.Join(tmpdir, "os-release")

	// XXX: This file is *NOT* created by this function on purpose
	// (to ensure the only file checked by the tests is
	// testOSRelease). osReleaseClr handling is tested in
	// utils_test.go.
	testOSReleaseClr := filepath.Join(tmpdir, "os-release-clr")

	testProcVersion := filepath.Join(tmpdir, "proc-version")

	// override
	procVersion = testProcVersion
	osRelease = testOSRelease
	osReleaseClr = testOSReleaseClr
	procCPUInfo = testProcCPUInfo

	procVersionContents := fmt.Sprintf("Linux version %s a b c",
		expectedKernelVersion)

	osReleaseContents := fmt.Sprintf(`
NAME="%s"
VERSION_ID="%s"
`, expectedDistro.Name, expectedDistro.Version)

	procCPUInfoContents := fmt.Sprintf(`
%s      : %s
%s      : %s
`,
		archCPUVendorField,
		expectedCPU.Vendor,
		archCPUModelField,
		expectedCPU.Model)

	data := []filesToCreate{
		{procVersion, procVersionContents},
		{osRelease, osReleaseContents},
		{procCPUInfo, procCPUInfoContents},
	}

	for _, d := range data {
		err := createFile(d.file, d.contents)
		if err != nil {
			return HostInfo{}, err
		}
	}

	return expectedHostDetails, nil
}
