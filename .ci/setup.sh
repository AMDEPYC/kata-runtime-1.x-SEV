#!/bin/bash
#
# Copyright (c) 2018 Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

set -e

cidir=$(dirname "$0")
source "${cidir}/lib.sh"

clone_tests_repo

pushd "${tests_repo_dir}"
.ci/setup.sh
popd

echo "Setup virtcontainers environment"
chronic sudo -E PATH=$PATH bash -c "${cidir}/../virtcontainers/utils/virtcontainers-setup.sh"

echo "Install virtcontainers"
make -C "${cidir}/../virtcontainers" && chronic sudo make -C "${cidir}/../virtcontainers" install
