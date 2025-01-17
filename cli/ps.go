// Copyright (c) 2014,2015,2016 Docker, Inc.
// Copyright (c) 2017 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"fmt"

	vc "github.com/kata-containers/runtime/virtcontainers"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var psCLICommand = cli.Command{
	Name:      "ps",
	Usage:     "ps displays the processes running inside a container",
	ArgsUsage: `<container-id> [ps options]`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "format, f",
			Value: "table",
			Usage: `select one of: ` + formatOptions,
		},
	},
	Action: func(context *cli.Context) error {
		if context.Args().Present() == false {
			return fmt.Errorf("Missing container ID, should at least provide one")
		}

		var args []string
		if len(context.Args()) > 1 {
			// [1:] is to remove container_id:
			// context.Args(): [container_id ps_arg1 ps_arg2 ...]
			// args:           [ps_arg1 ps_arg2 ...]
			args = context.Args()[1:]
		}

		return ps(context.Args().First(), context.String("format"), args)
	},
	SkipArgReorder: true,
}

func ps(containerID, format string, args []string) error {
	if containerID == "" {
		return fmt.Errorf("Missing container ID")
	}

	// Checks the MUST and MUST NOT from OCI runtime specification
	status, sandboxID, err := getExistingContainerInfo(containerID)
	if err != nil {
		return err
	}

	containerID = status.ID

	kataLog = kataLog.WithFields(logrus.Fields{
		"container": containerID,
		"sandbox":   sandboxID,
	})

	// container MUST be running
	if status.State.State != vc.StateRunning {
		return fmt.Errorf("Container %s is not running", containerID)
	}

	var options vc.ProcessListOptions

	options.Args = args
	if len(options.Args) == 0 {
		options.Args = []string{"-ef"}
	}

	options.Format = format

	msg, err := vci.ProcessListContainer(containerID, sandboxID, options)
	if err != nil {
		return err
	}

	fmt.Print(string(msg))

	return nil
}
