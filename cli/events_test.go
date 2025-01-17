// Copyright (c) 2018 Huawei Corporation.
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"flag"
	"os"
	"testing"
	"time"

	vc "github.com/kata-containers/runtime/virtcontainers"
	vcAnnotations "github.com/kata-containers/runtime/virtcontainers/pkg/annotations"
	"github.com/kata-containers/runtime/virtcontainers/pkg/vcmock"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func TestEventsCliAction(t *testing.T) {
	assert := assert.New(t)

	// get Action function
	actionFunc, ok := eventsCLICommand.Action.(func(ctx *cli.Context) error)
	flagSet := flag.NewFlagSet("events", flag.ContinueOnError)

	// create a new fake context
	ctx := cli.NewContext(&cli.App{}, flagSet, nil)
	assert.True(ok)

	err := actionFunc(ctx)
	assert.Error(err, "Missing container ID")
}

func TestEventsCLIFailure(t *testing.T) {
	assert := assert.New(t)

	flagSet := flag.NewFlagSet("events", flag.ContinueOnError)
	ctx := cli.NewContext(&cli.App{}, flagSet, nil)

	actionFunc, ok := eventsCLICommand.Action.(func(ctx *cli.Context) error)
	assert.True(ok)

	// missing container ID
	err := actionFunc(ctx)
	assert.Error(err)

	// interval is negative
	flagSet.Parse([]string{testContainerID})
	flagSet.Duration("interval", (-1)*time.Second, "")
	ctx = cli.NewContext(&cli.App{}, flagSet, nil)
	err = actionFunc(ctx)
	assert.Error(err)

	// interval is zero
	flagSet = flag.NewFlagSet("events", flag.ContinueOnError)
	flagSet.Parse([]string{testContainerID})
	flagSet.Duration("interval", 0*time.Second, "")
	ctx = cli.NewContext(&cli.App{}, flagSet, nil)
	err = actionFunc(ctx)
	assert.Error(err)

	// not running
	sandbox := &vcmock.Sandbox{
		MockID: testContainerID,
	}

	sandbox.MockContainers = []*vcmock.Container{
		{
			MockID:      sandbox.ID(),
			MockSandbox: sandbox,
		},
	}

	testingImpl.StatusContainerFunc = func(sandboxID, containerID string) (vc.ContainerStatus, error) {
		return vc.ContainerStatus{
			ID: sandbox.ID(),
			Annotations: map[string]string{
				vcAnnotations.ContainerTypeKey: string(vc.PodContainer),
			},
		}, nil
	}

	defer func() {
		testingImpl.StatusContainerFunc = nil
	}()
	err = actionFunc(ctx)
	assert.Error(err)
}

func TestEventsCLISuccessful(t *testing.T) {
	assert := assert.New(t)

	sandbox := &vcmock.Sandbox{
		MockID: testContainerID,
	}

	sandbox.MockContainers = []*vcmock.Container{
		{
			MockID:      sandbox.ID(),
			MockSandbox: sandbox,
		},
	}

	testingImpl.StatusContainerFunc = func(sandboxID, containerID string) (vc.ContainerStatus, error) {
		return vc.ContainerStatus{
			ID: sandbox.ID(),
			Annotations: map[string]string{
				vcAnnotations.ContainerTypeKey: string(vc.PodContainer),
			},
			State: vc.State{
				State: vc.StateRunning,
			},
		}, nil
	}

	testingImpl.StatsContainerFunc = func(sandboxID, containerID string) (vc.ContainerStats, error) {
		return vc.ContainerStats{}, nil
	}

	defer func() {
		testingImpl.StatusContainerFunc = nil
		testingImpl.StatsContainerFunc = nil
	}()

	path, err := createTempContainerIDMapping(sandbox.ID(), sandbox.ID())
	assert.NoError(err)
	defer os.RemoveAll(path)

	actionFunc, ok := eventsCLICommand.Action.(func(ctx *cli.Context) error)
	assert.True(ok)

	flagSet := flag.NewFlagSet("events", flag.ContinueOnError)
	flagSet.Parse([]string{testContainerID})
	flagSet.Duration("interval", 5*time.Second, "")
	flagSet.Bool("stats", true, "")
	ctx := cli.NewContext(&cli.App{}, flagSet, nil)
	err = actionFunc(ctx)
	assert.NoError(err)
}
