// Copyright (c) 2016 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package virtcontainers

import (
	"fmt"
	"syscall"

	"github.com/mitchellh/mapstructure"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

// AgentType describes the type of guest agent a Sandbox should run.
type AgentType string

// ProcessListOptions contains the options used to list running
// processes inside the container
type ProcessListOptions struct {
	// Format describes the output format to list the running processes.
	// Formats are unrelated to ps(1) formats, only two formats can be specified:
	// "json" and "table"
	Format string

	// Args contains the list of arguments to run ps(1) command.
	// If Args is empty the agent will use "-ef" as options to ps(1).
	Args []string
}

// ProcessList represents the list of running processes inside the container
type ProcessList []byte

const (
	// NoopAgentType is the No-Op agent.
	NoopAgentType AgentType = "noop"

	// HyperstartAgent is the Hyper hyperstart agent.
	HyperstartAgent AgentType = "hyperstart"

	// KataContainersAgent is the Kata Containers agent.
	KataContainersAgent AgentType = "kata"

	// SocketTypeVSOCK is a VSOCK socket type for talking to an agent.
	SocketTypeVSOCK = "vsock"

	// SocketTypeUNIX is a UNIX socket type for talking to an agent.
	// It typically means the agent is living behind a host proxy.
	SocketTypeUNIX = "unix"
)

// Set sets an agent type based on the input string.
func (agentType *AgentType) Set(value string) error {
	switch value {
	case "noop":
		*agentType = NoopAgentType
		return nil
	case "hyperstart":
		*agentType = HyperstartAgent
		return nil
	case "kata":
		*agentType = KataContainersAgent
		return nil
	default:
		return fmt.Errorf("Unknown agent type %s", value)
	}
}

// String converts an agent type to a string.
func (agentType *AgentType) String() string {
	switch *agentType {
	case NoopAgentType:
		return string(NoopAgentType)
	case HyperstartAgent:
		return string(HyperstartAgent)
	case KataContainersAgent:
		return string(KataContainersAgent)
	default:
		return ""
	}
}

// newAgent returns an agent from an agent type.
func newAgent(agentType AgentType) agent {
	switch agentType {
	case NoopAgentType:
		return &noopAgent{}
	case HyperstartAgent:
		return &hyper{}
	case KataContainersAgent:
		return &kataAgent{}
	default:
		return &noopAgent{}
	}
}

// newAgentConfig returns an agent config from a generic SandboxConfig interface.
func newAgentConfig(config SandboxConfig) interface{} {
	switch config.AgentType {
	case NoopAgentType:
		return nil
	case HyperstartAgent:
		var hyperConfig HyperConfig
		err := mapstructure.Decode(config.AgentConfig, &hyperConfig)
		if err != nil {
			return err
		}
		return hyperConfig
	case KataContainersAgent:
		var kataAgentConfig KataAgentConfig
		err := mapstructure.Decode(config.AgentConfig, &kataAgentConfig)
		if err != nil {
			return err
		}
		return kataAgentConfig
	default:
		return nil
	}
}

// agent is the virtcontainers agent interface.
// Agents are running in the guest VM and handling
// communications between the host and guest.
type agent interface {
	// init is used to pass agent specific configuration to the agent implementation.
	// agent implementations also will typically start listening for agent events from
	// init().
	// After init() is called, agent implementations should be initialized and ready
	// to handle all other Agent interface methods.
	init(sandbox *Sandbox, config interface{}) error

	// capabilities should return a structure that specifies the capabilities
	// supported by the agent.
	capabilities() capabilities

	// check will check the agent liveness
	check() error

	// disconnect will disconnect the connection to the agent
	disconnect() error

	// createSandbox will tell the agent to perform necessary setup for a Sandbox.
	createSandbox(sandbox *Sandbox) error

	// exec will tell the agent to run a command in an already running container.
	exec(sandbox *Sandbox, c Container, cmd Cmd) (*Process, error)

	// startSandbox will tell the agent to start all containers related to the Sandbox.
	startSandbox(sandbox *Sandbox) error

	// stopSandbox will tell the agent to stop all containers related to the Sandbox.
	stopSandbox(sandbox *Sandbox) error

	// cleanup will clean the resources for sandbox
	cleanupSandbox(sandbox *Sandbox) error

	// createContainer will tell the agent to create a container related to a Sandbox.
	createContainer(sandbox *Sandbox, c *Container) (*Process, error)

	// startContainer will tell the agent to start a container related to a Sandbox.
	startContainer(sandbox *Sandbox, c *Container) error

	// stopContainer will tell the agent to stop a container related to a Sandbox.
	stopContainer(sandbox *Sandbox, c Container) error

	// signalProcess will tell the agent to send a signal to a
	// container or a process related to a Sandbox. If all is true, all processes in
	// the container will be sent the signal.
	signalProcess(c *Container, processID string, signal syscall.Signal, all bool) error

	// winsizeProcess will tell the agent to set a process' tty size
	winsizeProcess(c *Container, processID string, height, width uint32) error

	// writeProcessStdin will tell the agent to write a process stdin
	writeProcessStdin(c *Container, ProcessID string, data []byte) (int, error)

	// closeProcessStdin will tell the agent to close a process stdin
	closeProcessStdin(c *Container, ProcessID string) error

	// readProcessStdout will tell the agent to read a process stdout
	readProcessStdout(c *Container, processID string, data []byte) (int, error)

	// readProcessStderr will tell the agent to read a process stderr
	readProcessStderr(c *Container, processID string, data []byte) (int, error)

	// processListContainer will list the processes running inside the container
	processListContainer(sandbox *Sandbox, c Container, options ProcessListOptions) (ProcessList, error)

	// updateContainer will update the resources of a running container
	updateContainer(sandbox *Sandbox, c Container, resources specs.LinuxResources) error

	// waitProcess will wait for the exit code of a process
	waitProcess(c *Container, processID string) (int32, error)

	// onlineCPUMem will online CPUs and Memory inside the Sandbox.
	// This function should be called after hot adding vCPUs or Memory.
	// cpus specifies the number of CPUs that were added and the agent should online
	onlineCPUMem(cpus uint32) error

	// statsContainer will tell the agent to get stats from a container related to a Sandbox
	statsContainer(sandbox *Sandbox, c Container) (*ContainerStats, error)

	// pauseContainer will pause a container
	pauseContainer(sandbox *Sandbox, c Container) error

	// resumeContainer will resume a paused container
	resumeContainer(sandbox *Sandbox, c Container) error
}
