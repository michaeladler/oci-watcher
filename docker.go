// SPDX-FileCopyrightText: 2025 Margo
//
// SPDX-License-Identifier: MIT
//
// SPDX-FileContributor: Michael Adler <michael.adler@siemens.com>

package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/client"
)

func uploadToDocker(filePath string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	response, err := cli.ImageLoad(ctx, file, true)
	if err != nil {
		return fmt.Errorf("failed to load image into Docker: %w", err)
	}
	defer response.Body.Close()

	// show response
	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	return nil
}
