// SPDX-FileCopyrightText: 2025 Margo
//
// SPDX-License-Identifier: MIT
//
// SPDX-FileContributor: Michael Adler <michael.adler@siemens.com>

package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
)

func verifyGPGSignature(pubKey io.Reader, signedFile, signatureFile string) error {
	log.Println("Verifying signature of", signedFile)

	keyring, err := openpgp.ReadArmoredKeyRing(pubKey)
	if err != nil {
		return err
	}

	signature, err := os.Open(signatureFile)
	if err != nil {
		return err
	}
	defer signature.Close()

	signed, err := os.Open(signedFile)
	if err != nil {
		return err
	}
	defer signed.Close()

	if _, err := openpgp.CheckDetachedSignature(keyring, signed, signature, nil); err != nil {
		return fmt.Errorf("signature verification failed: %v", err)
	}
	log.Println("Signature verified succesfully")
	return nil
}
