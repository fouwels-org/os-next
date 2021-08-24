// SPDX-FileCopyrightText: 2021 Belcan Advanced Solutions
// SPDX-FileCopyrightText: 2021 Kaelan Thijs Fouwels <kaelan.thijs@fouwels.com>
//
// SPDX-License-Identifier: Apache-2.0

package console

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os-next/init/config"
	"os-next/init/shell"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
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

		if checkAuthenticator(auth.Root, santext) {
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

func generateAuthenticator(text string) (string, error) {
	const _bcryptCost = 10

	bytes, err := bcrypt.GenerateFromPassword([]byte(text), _bcryptCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func checkAuthenticator(hash string, text string) bool {

	if hash == "" || text == "" {
		return false
	}

	bts, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		log.Printf("failed to decode base64 authenticator hash: %v", err)
		return false
	}

	if len(bts) == 0 {
		log.Printf("authenticator hash decoded to 0 length?")
		return false
	}

	err = bcrypt.CompareHashAndPassword(bts, []byte(text))

	//lint:ignore S1008 clearer flow being verbose
	if err != nil {
		return false
	}

	return true
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
