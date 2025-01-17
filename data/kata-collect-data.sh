#!/bin/bash
#
# Copyright (c) 2017-2018 Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0
#

typeset -r script_name=${0##*/}
typeset -r runtime_name="kata-runtime"
typeset -r runtime=$(command -v "$runtime_name" 2>/dev/null)
typeset -r issue_url="https://github.com/kata-containers/kata-containers/issues/new"
typeset -r script_version="1.1.0 (commit 292786890bd4792226c4fbf9b5268a92c27fe0f1-dirty)"

typeset -r unknown="unknown"

typeset -r osbuilder_file="/var/lib/osbuilder/osbuilder.yaml"

# Maximum number of errors to show for a single system component
# (such as runtime or proxy).
PROBLEM_LIMIT=${PROBLEM_LIMIT:-50}

# List of patterns used to detect problems in logfiles.
problem_pattern="("
problem_pattern+="\<abort|"
problem_pattern+="\<bug\>|"
problem_pattern+="\<cannot\>|"
problem_pattern+="\<catastrophic|"
problem_pattern+="\<could not\>|"
problem_pattern+="\<couldn\'t\>|"
problem_pattern+="\<critical|"
problem_pattern+="\<die\>|"
problem_pattern+="\<died\>|"
problem_pattern+="\<does.*not.*exist\>|"
problem_pattern+="\<dying\>|"
problem_pattern+="\<empty\>|"
problem_pattern+="\<erroneous|"
problem_pattern+="\<error|"
problem_pattern+="\<expected\>|"
problem_pattern+="\<fail|"
problem_pattern+="\<fatal|"
problem_pattern+="\<impossible\>|"
problem_pattern+="\<impossibly\>|"
problem_pattern+="\<incorrect|"
problem_pattern+="\<invalid\>|"
problem_pattern+="\<level=\"*error\"* |"
problem_pattern+="\<level=\"*fatal\"* |"
problem_pattern+="\<level=\"*panic\"* |"
problem_pattern+="\<level=\"*warning\"* |"
problem_pattern+="\<missing\>|"
problem_pattern+="\<need\>|"
problem_pattern+="\<no.*such.*file\>|"
problem_pattern+="\<not.*found\>|"
problem_pattern+="\<not.*supported\>|"
problem_pattern+="\<too many\>|"
problem_pattern+="\<unable\>|"
problem_pattern+="\<unavailable\>|"
problem_pattern+="\<unexpected|"
problem_pattern+="\<unknown\>|"
problem_pattern+="\<urgent|"
problem_pattern+="\<warn\>|"
problem_pattern+="\<warning\>|"
problem_pattern+="\<wrong\>"
problem_pattern+=")"

# List of patterns used to exclude messages that are not problems
problem_exclude_pattern="("
problem_exclude_pattern+="\<launching .* with:"
problem_exclude_pattern+=")"

usage()
{
	cat <<EOT
Usage: $script_name [options]

Summary: Collect data about an installation of Kata Containers.

Description: Run this script as root to obtain a markdown-formatted summary
  of the environment of the Kata Containers installation. The output of this script
  can be pasted directly into a github issue at the address below:

      $issue_url

Options:

 -h | --help    : show this usage summary.
 -v | --version : show program version.

EOT
}

version()
{
	cat <<EOT
$script_name version $script_version
EOT
}

die()
{
	local msg="$*"
	echo >&2 "ERROR: $script_name: $msg"
	exit 1
}

msg()
{
	local msg="$*"
	echo "$msg"
}

heading()
{
	local name="$*"
	echo -e "\n# $name\n"
}

subheading()
{
	local name="$*"
	echo -e "\n## $name\n"
}

separator()
{
	echo -e '\n---\n'
}

have_cmd()
{
	local cmd="$1"

	command -v "$cmd" &>/dev/null
	local ret=$?

	if [ $ret -eq 0 ]; then
		msg "Have \`$cmd\`"
	else
		msg "No \`$cmd\`"
	fi

	[ $ret -eq 0 ]
}

show_quoted_text()
{
	local text="$*"

	echo "\`\`\`"
	echo "$text"
	echo "\`\`\`"
}

run_cmd_and_show_quoted_output()
{
	local cmd="$*"

	local output

	msg "Output of \"\`$cmd\`\":"
	output=$(eval "$cmd" 2>&1)
	show_quoted_text "$output"
}

show_runtime_configs()
{
	local configs config

	heading "Runtime config files"
	
	configs=$($runtime --kata-show-default-config-paths)
	if [ $? -ne 0 ]; then
		version=$($runtime --version|tr '\n' ' ')
		die "failed to check config files - runtime is probably too old ($version)"
	fi

	subheading "Runtime default config files"

	show_quoted_text "$configs"

	# add in the standard defaults for good measure "just in case"
	configs+=" /etc/kata-containers/configuration.toml"
	configs+=" /usr/share/defaults/kata-containers/configuration.toml"
	configs+=" /usr/share/defaults/kata-containers/configuration.toml"
	configs+=" /etc/kata-containers/configuration.toml"

	# create a unique list of config files
	configs=$(echo $configs|tr ' ' '\n'|sort -u)

	subheading "Runtime config file contents"

	for config in $configs; do
		if [ -e "$config" ]; then
			run_cmd_and_show_quoted_output "cat \"$config\""
		else
			msg "Config file \`$config\` not found"
		fi
	done

	separator
}

show_log_details()
{
	heading "Logfiles"

	show_runtime_log_details
	show_proxy_log_details
	show_shim_log_details

	separator
}

show_runtime_log_details()
{
	subheading "Runtime logs"

	find_system_journal_problems "runtime" "kata-runtime" ""
}

find_system_journal_problems()
{
	local name="$1"
	local program="$2"
	local unit="$3"

	# select by identifier
	local selector="-t"

	local data_source="system journal"

	local problems=$(journalctl -q -o cat -a "$selector" "$program" |\
		grep "time=" |\
		egrep -i "$problem_pattern" |\
		egrep -iv "$problem_exclude_pattern" |\
		tail -n ${PROBLEM_LIMIT})

	if [ -n "$problems" ]; then
		msg "Recent $name problems found in $data_source:"
		show_quoted_text "$problems"
	else
		msg "No recent $name problems found in $data_source."
	fi
}

show_proxy_log_details()
{
	subheading "Proxy logs"

	find_system_journal_problems "proxy" "kata-proxy" ""
}

show_shim_log_details()
{
	subheading "Shim logs"

	find_system_journal_problems "shim" "kata-shim" ""
}

show_package_versions()
{
	heading "Packages"

	local pattern="("
	local project

	# CC 2.x, 3.0 and runv runtimes. They shouldn't be installed but let's
	# check anyway.
	pattern+="cc-oci-runtime"
	pattern+="cc-runtime"
	pattern+="runv"

	# core components
	for project in kata
	do
		pattern+="|${project}-proxy"
		pattern+="|${project}-runtime"
		pattern+="|${project}-shim"
		pattern+="|${project}-containers-image"
	done

	# assets
	pattern+="|linux-container"

	# hypervisor name prefix
	pattern+="|qemu-"

	pattern+=")"

	if have_cmd "dpkg"; then
		run_cmd_and_show_quoted_output "dpkg -l|egrep \"$pattern\""
	fi

	if have_cmd "rpm"; then
		run_cmd_and_show_quoted_output "rpm -qa|egrep \"$pattern\""
	fi

	separator
}

show_container_mgr_details()
{
	heading "Container manager details"

	if have_cmd "docker"; then
		subheading "Docker"
		run_cmd_and_show_quoted_output "docker version"
		run_cmd_and_show_quoted_output "docker info"
		run_cmd_and_show_quoted_output "systemctl show docker"
	fi

	if have_cmd "kubectl"; then
		subheading "Kubernetes"
		run_cmd_and_show_quoted_output "kubectl version"
		run_cmd_and_show_quoted_output "kubectl config view"
		run_cmd_and_show_quoted_output "systemctl show kubelet"

		if have_cmd "crio"; then
			run_cmd_and_show_quoted_output "crio --version"
			run_cmd_and_show_quoted_output "systemctl show crio"
		fi
	fi

	separator
}

show_meta()
{
	local date

	heading "Meta details"

	date=$(date '+%Y-%m-%d.%H:%M:%S.%N%z')
	msg "Running \`$script_name\` version \`$script_version\` at \`$date\`."

	separator
}

show_runtime()
{
	local cmd

	msg "Runtime is \`$runtime\`."

	cmd="kata-env"
	heading "\`$cmd\`"
	run_cmd_and_show_quoted_output "$runtime $cmd"

	separator
}

# Parameter 1: Path to disk image file.
# Returns: Details of the image, or "$unknown" on error.
get_image_details()
{
	local img="$1"

	[ -z "$img" ] && { echo "$unknown"; return;}
	[ -e "$img" ] || { echo "$unknown"; return;}

	local loop_device
	local partition_path
	local partitions
	local partition
	local count
	local mountpoint
	local contents
	local expected

	loop_device=$(loopmount_image "$img")
	if [ -z "$loop_device" ]; then
		echo "$unknown"
		return
	fi

	partitions=$(get_partitions "$loop_device")
	count=$(echo "$partitions"|wc -l)

	expected=1

	if [ "$count" -ne "$expected" ]; then
		release_device "$loop_device"
		echo "$unknown"
		return
	fi

	partition="$partitions"

	partition_path="/dev/${partition}"
	if [ ! -e "$partition_path" ]; then
		release_device "$loop_device"
		echo "$unknown"
		return
	fi

	mountpoint=$(mount_partition "$partition_path")

	contents=$(read_osbuilder_file "${mountpoint}")
	[ -z "$contents" ] && contents="$unknown"

	unmount_partition "$mountpoint"
	release_device "$loop_device"

	echo "$contents"
}

# Parameter 1: Path to the initrd file.
# Returns: Details of the initrd, or "$unknown" on error.
get_initrd_details()
{
	local initrd="$1"

	[ -z "$initrd" ] && { echo "$unknown"; return;}
	[ -e "$initrd" ] || { echo "$unknown"; return;}

	local file
	local relative_file=""
	local tmp

	file="${osbuilder_file}"

	# All files in the cpio archive are relative so remove leading slash
	relative_file=$(echo "$file"|sed 's!^/!!g')

	local tmpdir=$(mktemp -d)

	# Note: 'cpio --directory' seems to be non-portable, so cd(1) instead.
	(cd "$tmpdir" && gzip -dc "$initrd" | cpio \
		--extract \
		--make-directories \
		--no-absolute-filenames \
		$relative_file 2>/dev/null)

	contents=$(read_osbuilder_file "${tmpdir}")
	[ -z "$contents" ] && contents="$unknown"

	tmp="${tmpdir}/${file}"
	[ -d "$tmp" ] && rm -rf "$tmp"

	echo "$contents"
}

# Returns: Full path to the image file.
get_image_file()
{
	local cmd="kata-env"
	local cmdline="$runtime $cmd"

	local image=$(eval "$cmdline" 2>/dev/null |\
		grep -A 1 '^\[Image\]' |\
		egrep "\<Path\> =" |\
		awk '{print $3}' |\
		tr -d '"')

	echo "$image"
}

# Returns: Full path to the initrd file.
get_initrd_file()
{
	local cmd="kata-env"
	local cmdline="$runtime $cmd"

	local initrd=$(eval "$cmdline" 2>/dev/null |\
		grep -A 1 '^\[Initrd\]' |\
		egrep "\<Path\> =" |\
		awk '{print $3}' |\
		tr -d '"')

	echo "$initrd"
} 

# Parameter 1: Path to disk image file.
# Returns: Path to loop device.
loopmount_image()
{
	local img="$1"
	[ -n "$img" ] || die "need image file"

	local device_path

	losetup -fP "$img"

	device_path=$(losetup -j "$img" |\
		cut -d: -f1 |\
		sort -k1,1 |\
		tail -1)

	echo "$device_path"
}

# Parameter 1: Path to loop device.
# Returns: Partition names.
get_partitions()
{
	local device_path="$1"
	[ -n "$device_path" ] || die "need device path"

	local device
	local partitions

	device=${device_path/\/dev\//}

	partitions=$(lsblk -nli -o NAME "${device_path}" |\
		grep -v "^${device}$")

	echo "$partitions"
}

# Parameter 1: Path to disk partition device.
# Returns: Mountpoint.
mount_partition()
{
	local partition="$1"
	[ -n "$partition" ] || die "need partition"
	[ -e "$partition" ] || die "partition does not exist: $partition"

	local mountpoint

	mountpoint=$(mktemp -d)

	mount -oro,noload "$partition" "$mountpoint"

	echo "$mountpoint"
}

# Parameter 1: Mountpoint.
unmount_partition()
{
	local mountpoint="$1"
	[ -n "$mountpoint" ] || die "need mountpoint"
	[ -n "$mountpoint" ] || die "mountpoint does not exist: $mountpoint"

	umount "$mountpoint"
}

# Parameter 1: Loop device path.
release_device()
{
	local device="$1"
	[ -n "$device" ] || die "need device"
	[ -e "$device" ] || die "device does not exist: $device"

	losetup -d "$device"
}

# Retrieve details of the image containing
# the rootfs used to boot the virtual machine.
show_image_details()
{
	local image
	local details
	
	image=$(get_image_file)

	heading "Image details"

	if [ -n "$image" ]
	then
		details=$(get_image_details "$image")
		show_quoted_text "$details"
	else
		msg "No image"
	fi

	separator
}

# Retrieve details of the initrd containing
# the rootfs used to boot the virtual machine.
show_initrd_details()
{
	local initrd
	local details
	
	initrd=$(get_initrd_file)

	heading "Initrd details"

	if [ -n "$initrd" ]
	then
		details=$(get_initrd_details "$initrd")
		show_quoted_text "$details"
	else
		msg "No initrd"
	fi

	separator
}

read_osbuilder_file()
{
	local rootdir="$1"

	[ -n "$rootdir" ] || die "need root directory"

	local file="${rootdir}/${osbuilder_file}"

	[ ! -e "$file" ] && return

	cat "$file"
}

main()
{
	args=$(getopt \
		-n "$script_name" \
		-a \
		--options="dhv" \
		--longoptions="debug help version" \
		-- "$@")

	eval set -- "$args"
	[ $? -ne 0 ] && { usage && exit 1; }
	[ $# -eq 0 ] && { usage && exit 0; }

	while [ $# -gt 1 ]
	do
		case "$1" in
			-d|--debug)
				set -x
				;;

			-h|--help)
				usage && exit 0
				;;

			-v|--version)
				version && exit 0
				;;

			--)
				shift
				break
				;;
		esac
		shift
	done

	[ $(id -u) -eq 0 ] || die "Need to run as root"
	[ -n "$runtime" ] || die "cannot find runtime '$runtime_name'"

	show_meta
	show_runtime
	show_runtime_configs
	show_image_details
	show_initrd_details
	show_log_details
	show_container_mgr_details
	show_package_versions
}

main "$@"
