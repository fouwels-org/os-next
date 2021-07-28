// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package console

import (
	"bufio"
	"crypto"
	_ "crypto/sha256"
	"fmt"
	"init/config"
	"init/shell"
	"log"
	"os"
	"strings"
	"time"
)

// Start runtime console
func Start(auth config.Authenticators) error {

	err := login(auth)
	if err != nil {
		return err
	}

	err = bash()
	if err != nil {
		return err
	}
	return nil
}

func login(auth config.Authenticators) error {

	success := false
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}

	for !success {
		fmt.Printf("enter authenticator for shell\n> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}

		santext := strings.TrimSuffix(text, "\n")

		if checkPasswordHash(auth.Root, santext) {
			success = true
			log.Printf("user succeeded to authenticate")
		} else {
			fmt.Printf("authenticator incorrect\n")
			log.Printf("user failed to authenticate")
			time.Sleep(2 * time.Second)
		}
	}

	return nil
}

func checkPasswordHash(hash string, text string) bool {

	hasher := crypto.SHA256.New()
	_, err := hasher.Write([]byte(text))
	if err != nil {
		log.Printf("failed to write hash for login: %v", err)
		return false
	}
	textHash := fmt.Sprintf("%x", hasher.Sum(nil))

	if textHash == hash {
		return true
	} else {
		return false
	}
}

func bash() error {
	commands := []shell.Command{
		{Executable: shell.Ash, Arguments: []string{}},
	}

	err := shell.Executor.ExecuteInteractive(commands)
	if err != nil {
		return err
	}

	return nil
}
