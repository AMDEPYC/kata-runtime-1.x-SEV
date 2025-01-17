// Copyright 2015 The rkt Authors
// Copyright (c) 2016 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
//

package virtcontainers

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

func TestBindMountInvalidSourceSymlink(t *testing.T) {
	source := filepath.Join(testDir, "fooFile")
	os.Remove(source)

	err := bindMount(source, "", false)
	if err == nil {
		t.Fatal()
	}
}

func TestBindMountFailingMount(t *testing.T) {
	source := filepath.Join(testDir, "fooLink")
	fakeSource := filepath.Join(testDir, "fooFile")
	os.Remove(source)
	os.Remove(fakeSource)

	_, err := os.OpenFile(fakeSource, os.O_CREATE, mountPerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Symlink(fakeSource, source)
	if err != nil {
		t.Fatal(err)
	}

	err = bindMount(source, "", false)
	if err == nil {
		t.Fatal()
	}
}

func TestBindMountSuccessful(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip(testDisabledAsNonRoot)
	}

	source := filepath.Join(testDir, "fooDirSrc")
	dest := filepath.Join(testDir, "fooDirDest")
	syscall.Unmount(dest, 0)
	os.Remove(source)
	os.Remove(dest)

	err := os.MkdirAll(source, mountPerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(dest, mountPerm)
	if err != nil {
		t.Fatal(err)
	}

	err = bindMount(source, dest, false)
	if err != nil {
		t.Fatal(err)
	}

	syscall.Unmount(dest, 0)
}

func TestBindMountReadonlySuccessful(t *testing.T) {
	if os.Geteuid() != 0 {
		t.Skip(testDisabledAsNonRoot)
	}

	source := filepath.Join(testDir, "fooDirSrc")
	dest := filepath.Join(testDir, "fooDirDest")
	syscall.Unmount(dest, 0)
	os.Remove(source)
	os.Remove(dest)

	err := os.MkdirAll(source, mountPerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.MkdirAll(dest, mountPerm)
	if err != nil {
		t.Fatal(err)
	}

	err = bindMount(source, dest, true)
	if err != nil {
		t.Fatal(err)
	}

	defer syscall.Unmount(dest, 0)

	// should not be able to create file in read-only mount
	destFile := filepath.Join(dest, "foo")
	_, err = os.OpenFile(destFile, os.O_CREATE, mountPerm)
	if err == nil {
		t.Fatal(err)
	}
}

func TestEnsureDestinationExistsNonExistingSource(t *testing.T) {
	err := ensureDestinationExists("", "")
	if err == nil {
		t.Fatal()
	}
}

func TestEnsureDestinationExistsWrongParentDir(t *testing.T) {
	source := filepath.Join(testDir, "fooFile")
	dest := filepath.Join(source, "fooDest")
	os.Remove(source)
	os.Remove(dest)

	_, err := os.OpenFile(source, os.O_CREATE, mountPerm)
	if err != nil {
		t.Fatal(err)
	}

	err = ensureDestinationExists(source, dest)
	if err == nil {
		t.Fatal()
	}
}

func TestEnsureDestinationExistsSuccessfulSrcDir(t *testing.T) {
	source := filepath.Join(testDir, "fooDirSrc")
	dest := filepath.Join(testDir, "fooDirDest")
	os.Remove(source)
	os.Remove(dest)

	err := os.MkdirAll(source, mountPerm)
	if err != nil {
		t.Fatal(err)
	}

	err = ensureDestinationExists(source, dest)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEnsureDestinationExistsSuccessfulSrcFile(t *testing.T) {
	source := filepath.Join(testDir, "fooDirSrc")
	dest := filepath.Join(testDir, "fooDirDest")
	os.Remove(source)
	os.Remove(dest)

	_, err := os.OpenFile(source, os.O_CREATE, mountPerm)
	if err != nil {
		t.Fatal(err)
	}

	err = ensureDestinationExists(source, dest)
	if err != nil {
		t.Fatal(err)
	}
}
