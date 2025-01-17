// Copyright (c) 2016 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package virtcontainers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/kata-containers/runtime/virtcontainers/device/api"
	"github.com/kata-containers/runtime/virtcontainers/device/config"
	"github.com/kata-containers/runtime/virtcontainers/device/drivers"
)

// sandboxResource is an int representing a sandbox resource type.
//
// Note that some are specific to the sandbox itself and others can apply to
// sandboxes and containers.
type sandboxResource int

const (
	// configFileType represents a configuration file type
	configFileType sandboxResource = iota

	// stateFileType represents a state file type
	stateFileType

	// networkFileType represents a network file type (sandbox only)
	networkFileType

	// hypervisorFileType represents a hypervisor file type (sandbox only)
	hypervisorFileType

	// agentFileType represents an agent file type (sandbox only)
	agentFileType

	// processFileType represents a process file type
	processFileType

	// lockFileType represents a lock file type (sandbox only)
	lockFileType

	// mountsFileType represents a mount file type
	mountsFileType

	// devicesFileType represents a device file type
	devicesFileType
)

// configFile is the file name used for every JSON sandbox configuration.
const configFile = "config.json"

// stateFile is the file name storing a sandbox state.
const stateFile = "state.json"

// networkFile is the file name storing a sandbox network.
const networkFile = "network.json"

// hypervisorFile is the file name storing a hypervisor's state.
const hypervisorFile = "hypervisor.json"

// agentFile is the file name storing an agent's state.
const agentFile = "agent.json"

// processFile is the file name storing a container process.
const processFile = "process.json"

// lockFile is the file name locking the usage of a sandbox.
const lockFileName = "lock"

const mountsFile = "mounts.json"

// devicesFile is the file name storing a container's devices.
const devicesFile = "devices.json"

// dirMode is the permission bits used for creating a directory
const dirMode = os.FileMode(0750) | os.ModeDir

// storagePathSuffix is the suffix used for all storage paths
//
// Note: this very brief path represents "virtcontainers sandboxes". It is as
// terse as possible to minimise path length.
const storagePathSuffix = "/vc/sbs"

// configStoragePath is the sandbox configuration directory.
// It will contain one config.json file for each created sandbox.
var configStoragePath = filepath.Join("/var/lib", storagePathSuffix)

// runStoragePath is the sandbox runtime directory.
// It will contain one state.json and one lock file for each created sandbox.
var runStoragePath = filepath.Join("/run", storagePathSuffix)

// resourceStorage is the virtcontainers resources (configuration, state, etc...)
// storage interface.
// The default resource storage implementation is filesystem.
type resourceStorage interface {
	// Create all resources for a sandbox
	createAllResources(sandbox *Sandbox) error

	// Resources URIs functions return both the URI
	// for the actual resource and the URI base.
	containerURI(sandboxID, containerID string, resource sandboxResource) (string, string, error)
	sandboxURI(sandboxID string, resource sandboxResource) (string, string, error)

	// Sandbox resources
	storeSandboxResource(sandboxID string, resource sandboxResource, data interface{}) error
	deleteSandboxResources(sandboxID string, resources []sandboxResource) error
	fetchSandboxConfig(sandboxID string) (SandboxConfig, error)
	fetchSandboxState(sandboxID string) (State, error)
	fetchSandboxNetwork(sandboxID string) (NetworkNamespace, error)
	storeSandboxNetwork(sandboxID string, networkNS NetworkNamespace) error

	// Hypervisor resources
	fetchHypervisorState(sandboxID string, state interface{}) error
	storeHypervisorState(sandboxID string, state interface{}) error

	// Agent resources
	fetchAgentState(sandboxID string, state interface{}) error
	storeAgentState(sandboxID string, state interface{}) error

	// Container resources
	storeContainerResource(sandboxID, containerID string, resource sandboxResource, data interface{}) error
	deleteContainerResources(sandboxID, containerID string, resources []sandboxResource) error
	fetchContainerConfig(sandboxID, containerID string) (ContainerConfig, error)
	fetchContainerState(sandboxID, containerID string) (State, error)
	fetchContainerProcess(sandboxID, containerID string) (Process, error)
	storeContainerProcess(sandboxID, containerID string, process Process) error
	fetchContainerMounts(sandboxID, containerID string) ([]Mount, error)
	storeContainerMounts(sandboxID, containerID string, mounts []Mount) error
	fetchContainerDevices(sandboxID, containerID string) ([]api.Device, error)
	storeContainerDevices(sandboxID, containerID string, devices []api.Device) error
}

// filesystem is a resourceStorage interface implementation for a local filesystem.
type filesystem struct {
}

// Logger returns a logrus logger appropriate for logging filesystem messages
func (fs *filesystem) Logger() *logrus.Entry {
	return virtLog.WithField("subsystem", "filesystem")
}

func (fs *filesystem) createAllResources(sandbox *Sandbox) (err error) {
	for _, resource := range []sandboxResource{stateFileType, configFileType} {
		_, path, _ := fs.sandboxURI(sandbox.id, resource)
		err = os.MkdirAll(path, dirMode)
		if err != nil {
			return err
		}
	}

	for _, container := range sandbox.containers {
		for _, resource := range []sandboxResource{stateFileType, configFileType} {
			_, path, _ := fs.containerURI(sandbox.id, container.id, resource)
			err = os.MkdirAll(path, dirMode)
			if err != nil {
				fs.deleteSandboxResources(sandbox.id, nil)
				return err
			}
		}
	}

	sandboxlockFile, _, err := fs.sandboxURI(sandbox.id, lockFileType)
	if err != nil {
		fs.deleteSandboxResources(sandbox.id, nil)
		return err
	}

	_, err = os.Stat(sandboxlockFile)
	if err != nil {
		lockFile, err := os.Create(sandboxlockFile)
		if err != nil {
			fs.deleteSandboxResources(sandbox.id, nil)
			return err
		}
		lockFile.Close()
	}

	return nil
}

func (fs *filesystem) storeFile(file string, data interface{}) error {
	if file == "" {
		return errNeedFile
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	jsonOut, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Could not marshall data: %s", err)
	}
	f.Write(jsonOut)

	return nil
}

// TypedDevice is used as an intermediate representation for marshalling
// and unmarshalling Device implementations.
type TypedDevice struct {
	Type string

	// Data is assigned the Device object.
	// This being declared as RawMessage prevents it from being  marshalled/unmarshalled.
	// We do that explicitly depending on Type.
	Data json.RawMessage
}

// storeDeviceFile is used to provide custom marshalling for Device objects.
// Device is first marshalled into TypedDevice to include the type
// of the Device object.
func (fs *filesystem) storeDeviceFile(file string, data interface{}) error {
	if file == "" {
		return errNeedFile
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	devices, ok := data.([]api.Device)
	if !ok {
		return fmt.Errorf("Incorrect data type received, Expected []Device")
	}

	var typedDevices []TypedDevice
	for _, d := range devices {
		tempJSON, _ := json.Marshal(d)
		typedDevice := TypedDevice{
			Type: string(d.DeviceType()),
			Data: tempJSON,
		}
		typedDevices = append(typedDevices, typedDevice)
	}

	jsonOut, err := json.Marshal(typedDevices)
	if err != nil {
		return fmt.Errorf("Could not marshal devices: %s", err)
	}

	if _, err := f.Write(jsonOut); err != nil {
		return err
	}

	return nil
}

func (fs *filesystem) fetchFile(file string, resource sandboxResource, data interface{}) error {
	if file == "" {
		return errNeedFile
	}

	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	switch resource {
	case devicesFileType:
		devices, ok := data.(*[]api.Device)
		if !ok {
			return fmt.Errorf("Could not cast %v into *[]Device type", data)
		}

		return fs.fetchDeviceFile(fileData, devices)
	}

	return json.Unmarshal(fileData, data)
}

// fetchDeviceFile is used for custom unmarshalling of device interface objects.
func (fs *filesystem) fetchDeviceFile(fileData []byte, devices *[]api.Device) error {
	var typedDevices []TypedDevice
	if err := json.Unmarshal(fileData, &typedDevices); err != nil {
		return err
	}

	var tempDevices []api.Device
	for _, d := range typedDevices {
		l := fs.Logger().WithField("device-type", d.Type)
		l.Info("Device type found")

		switch d.Type {
		case string(config.DeviceVFIO):
			// TODO: remove dependency of drivers package
			var device drivers.VFIODevice
			if err := json.Unmarshal(d.Data, &device); err != nil {
				return err
			}
			tempDevices = append(tempDevices, &device)
			l.Infof("VFIO device unmarshalled [%v]", device)

		case string(config.DeviceBlock):
			// TODO: remove dependency of drivers package
			var device drivers.BlockDevice
			if err := json.Unmarshal(d.Data, &device); err != nil {
				return err
			}
			tempDevices = append(tempDevices, &device)
			l.Infof("Block Device unmarshalled [%v]", device)
		case string(config.DeviceGeneric):
			// TODO: remove dependency of drivers package
			var device drivers.GenericDevice
			if err := json.Unmarshal(d.Data, &device); err != nil {
				return err
			}
			tempDevices = append(tempDevices, &device)
			l.Infof("Generic device unmarshalled [%v]", device)

		default:
			return fmt.Errorf("Unknown device type, could not unmarshal")
		}
	}

	*devices = tempDevices
	return nil
}

// resourceNeedsContainerID determines if the specified
// sandboxResource needs a containerID. Since some sandboxResources can
// be used for both sandboxes and containers, it is necessary to specify
// whether the resource is being used in a sandbox-specific context using
// the sandboxSpecific parameter.
func resourceNeedsContainerID(sandboxSpecific bool, resource sandboxResource) bool {

	switch resource {
	case lockFileType, networkFileType, hypervisorFileType, agentFileType:
		// sandbox-specific resources
		return false
	default:
		return !sandboxSpecific
	}
}

func resourceDir(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource) (string, error) {
	if sandboxID == "" {
		return "", errNeedSandboxID
	}

	if resourceNeedsContainerID(sandboxSpecific, resource) == true && containerID == "" {
		return "", errNeedContainerID
	}

	var path string

	switch resource {
	case configFileType:
		path = configStoragePath
		break
	case stateFileType, networkFileType, processFileType, lockFileType, mountsFileType, devicesFileType, hypervisorFileType, agentFileType:
		path = runStoragePath
		break
	default:
		return "", errInvalidResource
	}

	dirPath := filepath.Join(path, sandboxID, containerID)

	return dirPath, nil
}

// If sandboxSpecific is true, the resource is being applied for an empty
// sandbox (meaning containerID may be blank).
// Note that this function defers determining if containerID can be
// blank to resourceDIR()
func (fs *filesystem) resourceURI(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource) (string, string, error) {
	if sandboxID == "" {
		return "", "", errNeedSandboxID
	}

	var filename string

	dirPath, err := resourceDir(sandboxSpecific, sandboxID, containerID, resource)
	if err != nil {
		return "", "", err
	}

	switch resource {
	case configFileType:
		filename = configFile
		break
	case stateFileType:
		filename = stateFile
	case networkFileType:
		filename = networkFile
	case hypervisorFileType:
		filename = hypervisorFile
	case agentFileType:
		filename = agentFile
	case processFileType:
		filename = processFile
	case lockFileType:
		filename = lockFileName
		break
	case mountsFileType:
		filename = mountsFile
		break
	case devicesFileType:
		filename = devicesFile
		break
	default:
		return "", "", errInvalidResource
	}

	filePath := filepath.Join(dirPath, filename)

	return filePath, dirPath, nil
}

func (fs *filesystem) containerURI(sandboxID, containerID string, resource sandboxResource) (string, string, error) {
	if sandboxID == "" {
		return "", "", errNeedSandboxID
	}

	if containerID == "" {
		return "", "", errNeedContainerID
	}

	return fs.resourceURI(false, sandboxID, containerID, resource)
}

func (fs *filesystem) sandboxURI(sandboxID string, resource sandboxResource) (string, string, error) {
	return fs.resourceURI(true, sandboxID, "", resource)
}

// commonResourceChecks performs basic checks common to both setting and
// getting a sandboxResource.
func (fs *filesystem) commonResourceChecks(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource) error {
	if sandboxID == "" {
		return errNeedSandboxID
	}

	if resourceNeedsContainerID(sandboxSpecific, resource) == true && containerID == "" {
		return errNeedContainerID
	}

	switch resource {
	case configFileType:
	case stateFileType:
	case networkFileType:
	case hypervisorFileType:
	case agentFileType:
	case processFileType:
	case mountsFileType:
	case devicesFileType:
	default:
		return errInvalidResource
	}

	return nil
}

func (fs *filesystem) storeSandboxAndContainerConfigResource(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource, file interface{}) error {
	if resource != configFileType {
		return errInvalidResource
	}

	configFile, _, err := fs.resourceURI(sandboxSpecific, sandboxID, containerID, configFileType)
	if err != nil {
		return err
	}

	return fs.storeFile(configFile, file)
}

func (fs *filesystem) storeStateResource(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource, file interface{}) error {
	if resource != stateFileType {
		return errInvalidResource
	}

	stateFile, _, err := fs.resourceURI(sandboxSpecific, sandboxID, containerID, stateFileType)
	if err != nil {
		return err
	}

	return fs.storeFile(stateFile, file)
}

func (fs *filesystem) storeNetworkResource(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource, file interface{}) error {
	if resource != networkFileType {
		return errInvalidResource
	}

	// sandbox only resource
	networkFile, _, err := fs.resourceURI(true, sandboxID, containerID, networkFileType)
	if err != nil {
		return err
	}

	return fs.storeFile(networkFile, file)
}

func (fs *filesystem) storeProcessResource(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource, file interface{}) error {
	if resource != processFileType {
		return errInvalidResource
	}

	processFile, _, err := fs.resourceURI(sandboxSpecific, sandboxID, containerID, processFileType)
	if err != nil {
		return err
	}

	return fs.storeFile(processFile, file)
}

func (fs *filesystem) storeMountResource(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource, file interface{}) error {
	if resource != mountsFileType {
		return errInvalidResource
	}

	mountsFile, _, err := fs.resourceURI(sandboxSpecific, sandboxID, containerID, mountsFileType)
	if err != nil {
		return err
	}

	return fs.storeFile(mountsFile, file)
}

func (fs *filesystem) storeDeviceResource(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource, file interface{}) error {
	if resource != devicesFileType {
		return errInvalidResource
	}

	devicesFile, _, err := fs.resourceURI(sandboxSpecific, sandboxID, containerID, devicesFileType)
	if err != nil {
		return err
	}

	return fs.storeDeviceFile(devicesFile, file)
}

func (fs *filesystem) storeResource(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource, data interface{}) error {
	if err := fs.commonResourceChecks(sandboxSpecific, sandboxID, containerID, resource); err != nil {
		return err
	}

	switch file := data.(type) {
	case SandboxConfig, ContainerConfig:
		return fs.storeSandboxAndContainerConfigResource(sandboxSpecific, sandboxID, containerID, resource, file)

	case State:
		return fs.storeStateResource(sandboxSpecific, sandboxID, containerID, resource, file)

	case NetworkNamespace:
		return fs.storeNetworkResource(sandboxSpecific, sandboxID, containerID, resource, file)

	case Process:
		return fs.storeProcessResource(sandboxSpecific, sandboxID, containerID, resource, file)

	case []Mount:
		return fs.storeMountResource(sandboxSpecific, sandboxID, containerID, resource, file)

	case []api.Device:
		return fs.storeDeviceResource(sandboxSpecific, sandboxID, containerID, resource, file)

	default:
		return fmt.Errorf("Invalid resource data type")
	}
}

func (fs *filesystem) fetchResource(sandboxSpecific bool, sandboxID, containerID string, resource sandboxResource, data interface{}) error {
	if err := fs.commonResourceChecks(sandboxSpecific, sandboxID, containerID, resource); err != nil {
		return err
	}

	path, _, err := fs.resourceURI(sandboxSpecific, sandboxID, containerID, resource)
	if err != nil {
		return err
	}

	return fs.fetchFile(path, resource, data)
}

func (fs *filesystem) storeSandboxResource(sandboxID string, resource sandboxResource, data interface{}) error {
	return fs.storeResource(true, sandboxID, "", resource, data)
}

func (fs *filesystem) fetchSandboxConfig(sandboxID string) (SandboxConfig, error) {
	var sandboxConfig SandboxConfig

	if err := fs.fetchResource(true, sandboxID, "", configFileType, &sandboxConfig); err != nil {
		return SandboxConfig{}, err
	}

	return sandboxConfig, nil
}

func (fs *filesystem) fetchSandboxState(sandboxID string) (State, error) {
	var state State

	if err := fs.fetchResource(true, sandboxID, "", stateFileType, &state); err != nil {
		return State{}, err
	}

	return state, nil
}

func (fs *filesystem) fetchSandboxNetwork(sandboxID string) (NetworkNamespace, error) {
	var networkNS NetworkNamespace

	if err := fs.fetchResource(true, sandboxID, "", networkFileType, &networkNS); err != nil {
		return NetworkNamespace{}, err
	}

	return networkNS, nil
}

func (fs *filesystem) fetchHypervisorState(sandboxID string, state interface{}) error {
	return fs.fetchResource(true, sandboxID, "", hypervisorFileType, state)
}

func (fs *filesystem) fetchAgentState(sandboxID string, state interface{}) error {
	return fs.fetchResource(true, sandboxID, "", agentFileType, state)
}

func (fs *filesystem) storeSandboxNetwork(sandboxID string, networkNS NetworkNamespace) error {
	return fs.storeSandboxResource(sandboxID, networkFileType, networkNS)
}

func (fs *filesystem) storeHypervisorState(sandboxID string, state interface{}) error {
	hypervisorFile, _, err := fs.resourceURI(true, sandboxID, "", hypervisorFileType)
	if err != nil {
		return err
	}

	return fs.storeFile(hypervisorFile, state)
}

func (fs *filesystem) storeAgentState(sandboxID string, state interface{}) error {
	agentFile, _, err := fs.resourceURI(true, sandboxID, "", agentFileType)
	if err != nil {
		return err
	}

	return fs.storeFile(agentFile, state)
}

func (fs *filesystem) deleteSandboxResources(sandboxID string, resources []sandboxResource) error {
	if resources == nil {
		resources = []sandboxResource{configFileType, stateFileType}
	}

	for _, resource := range resources {
		_, dir, err := fs.sandboxURI(sandboxID, resource)
		if err != nil {
			return err
		}

		err = os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	return nil
}

func (fs *filesystem) storeContainerResource(sandboxID, containerID string, resource sandboxResource, data interface{}) error {
	if sandboxID == "" {
		return errNeedSandboxID
	}

	if containerID == "" {
		return errNeedContainerID
	}

	return fs.storeResource(false, sandboxID, containerID, resource, data)
}

func (fs *filesystem) fetchContainerConfig(sandboxID, containerID string) (ContainerConfig, error) {
	var config ContainerConfig

	if err := fs.fetchResource(false, sandboxID, containerID, configFileType, &config); err != nil {
		return ContainerConfig{}, err
	}

	return config, nil
}

func (fs *filesystem) fetchContainerState(sandboxID, containerID string) (State, error) {
	var state State

	if err := fs.fetchResource(false, sandboxID, containerID, stateFileType, &state); err != nil {
		return State{}, err
	}

	return state, nil
}

func (fs *filesystem) fetchContainerProcess(sandboxID, containerID string) (Process, error) {
	var process Process

	if err := fs.fetchResource(false, sandboxID, containerID, processFileType, &process); err != nil {
		return Process{}, err
	}

	return process, nil
}

func (fs *filesystem) storeContainerProcess(sandboxID, containerID string, process Process) error {
	return fs.storeContainerResource(sandboxID, containerID, processFileType, process)
}

func (fs *filesystem) fetchContainerMounts(sandboxID, containerID string) ([]Mount, error) {
	var mounts []Mount

	if err := fs.fetchResource(false, sandboxID, containerID, mountsFileType, &mounts); err != nil {
		return []Mount{}, err
	}

	return mounts, nil
}

func (fs *filesystem) fetchContainerDevices(sandboxID, containerID string) ([]api.Device, error) {
	var devices []api.Device

	if err := fs.fetchResource(false, sandboxID, containerID, devicesFileType, &devices); err != nil {
		return []api.Device{}, err
	}

	return devices, nil
}

func (fs *filesystem) storeContainerMounts(sandboxID, containerID string, mounts []Mount) error {
	return fs.storeContainerResource(sandboxID, containerID, mountsFileType, mounts)
}

func (fs *filesystem) storeContainerDevices(sandboxID, containerID string, devices []api.Device) error {
	return fs.storeContainerResource(sandboxID, containerID, devicesFileType, devices)
}

func (fs *filesystem) deleteContainerResources(sandboxID, containerID string, resources []sandboxResource) error {
	if resources == nil {
		resources = []sandboxResource{configFileType, stateFileType}
	}

	for _, resource := range resources {
		_, dir, err := fs.sandboxURI(sandboxID, resource)
		if err != nil {
			return err
		}

		containerDir := filepath.Join(dir, containerID, "/")

		err = os.RemoveAll(containerDir)
		if err != nil {
			return err
		}
	}

	return nil
}
