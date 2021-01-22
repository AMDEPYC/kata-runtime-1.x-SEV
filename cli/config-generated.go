// WARNING: This file is auto-generated - DO NOT EDIT!
//
// Note that some variables are "var" to allow them to be modified
// by the tests.
package main

import (
	"fmt"
)

// name is the name of the runtime
const name = "kata-runtime"

// name of the project
const project = "Kata Containers"

// prefix used to denote non-standard CLI commands and options.
const projectPrefix = "kata"

// systemdUnitName is the systemd(1) target used to launch the agent.
const systemdUnitName = "kata-containers.target"

// original URL for this project
const projectURL = "https://github.com/kata-containers"

// commit is the git commit the runtime is compiled from.
var commit = "292786890bd4792226c4fbf9b5268a92c27fe0f1-dirty"

// version is the runtime version.
var version = "1.1.0"

// project-specific command names
var envCmd = fmt.Sprintf("%s-env", projectPrefix)
var checkCmd = fmt.Sprintf("%s-check", projectPrefix)

// project-specific option names
var configFilePathOption = fmt.Sprintf("%s-config", projectPrefix)
var showConfigPathsOption = fmt.Sprintf("%s-show-default-config-paths", projectPrefix)

var defaultHypervisorPath = "/usr//bin/qemu-lite-system-x86_64"
var defaultImagePath = "/usr//share/kata-containers/kata-containers.img"
var defaultKernelPath = "/usr//share/kata-containers/vmlinuz.container"
var defaultInitrdPath = "/usr//share/kata-containers/kata-containers-initrd.img"
var defaultFirmwarePath = ""
var defaultMachineAccelerators = ""
var defaultShimPath = "/usr//libexec/kata-containers/kata-shim"

const defaultKernelParams = ""
const defaultMachineType = "pc"
const defaultRootDirectory = "/var/run/kata-containers"

const defaultVCPUCount uint32 = 1
const defaultMaxVCPUCount uint32 = 0
const defaultMemSize uint32 = 2048 // MiB
const defaultBridgesCount uint32 = 1
const defaultInterNetworkingModel = "macvtap"
const defaultDisableBlockDeviceUse bool = false
const defaultBlockDeviceDriver = "virtio-scsi"
const defaultEnableIOThreads bool = false
const defaultEnableMemPrealloc bool = false
const defaultEnableHugePages bool = false
const defaultEnableSwap bool = false
const defaultEnableDebug bool = false
const defaultDisableNestingChecks bool = false
const defaultMsize9p uint32 = 8192

// Default config file used by stateless systems.
var defaultRuntimeConfiguration = "/usr/share/defaults/kata-containers/configuration.toml"

// Alternate config file that takes precedence over
// defaultRuntimeConfiguration.
var defaultSysConfRuntimeConfiguration = "/etc/kata-containers/configuration.toml"

var defaultProxyPath = "/usr//libexec/kata-containers/kata-proxy"
