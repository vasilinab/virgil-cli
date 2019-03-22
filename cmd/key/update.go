/*
 * Copyright (C) 2015-2019 Virgil Security Inc.
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are
 * met:
 *
 *     (1) Redistributions of source code must retain the above copyright
 *     notice, this list of conditions and the following disclaimer.
 *
 *     (2) Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in
 *     the documentation and/or other materials provided with the
 *     distribution.
 *
 *     (3) Neither the name of the copyright holder nor the names of its
 *     contributors may be used to endorse or promote products derived from
 *     this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE AUTHOR ''AS IS'' AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING
 * IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * Lead Maintainer: Virgil Security Inc. <support@virgilsecurity.com>
 */

package key

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/VirgilSecurity/virgil-cli/utils"

	"github.com/VirgilSecurity/virgil-cli/models"

	"github.com/VirgilSecurity/virgil-cli/client"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

func Update(vcli *client.VirgilHttpClient) *cli.Command {
	return &cli.Command{
		Name:      "update",
		Aliases:   []string{"u"},
		ArgsUsage: "api_key_id",
		Usage:     "Update existing api-key",
		Flags:     []cli.Flag{&cli.StringFlag{Name: "app_id"}},
		Action: func(context *cli.Context) (err error) {

			appID := context.String("app_id")
			if appID == "" {
				appID, _ := utils.LoadAppID()
				if appID == "" {
					return errors.New("Please, specify app_id (flag --app_id)")
				}
			} else {
				utils.SaveAppID(appID)
			}

			if context.NArg() < 1 {
				return errors.New("Invalid number of arguments. Please, specify api-key id")
			}

			err = UpdateFunc(context.Args().First(), vcli)

			if err != nil {
				return err
			}

			fmt.Println("Key successfully updated")
			return nil
		},
	}
}

func UpdateFunc(apiKeyID string, vcli *client.VirgilHttpClient) (err error) {

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter new api-key name:")
	name := ""
	for name == "" {
		scanner.Scan()
		name = scanner.Text()
		if name == "" {
			fmt.Printf("name can't be empty")
			fmt.Println("Enter new api-key name:")
		}
	}

	req := &models.UpdateAccessKeyRequest{Name: name}

	token, err := utils.LoadAccessTokenOrLogin(vcli)

	if err != nil {
		return err
	}

	for err == nil {
		_, _, vErr := vcli.Send(http.MethodPut, token, "access_keys/"+apiKeyID, req, nil)
		if vErr == nil {
			break
		}

		token, err = utils.CheckRetry(vErr, vcli)
	}

	return err
}