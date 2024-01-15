package main

import (
	"errors"
	keychain "github.com/keybase/go-keychain"
	touchid "github.com/lox/go-touchid"
	"github.com/spf13/cobra"
	"github.com/xlzd/gotp"
	"log"
	"os/exec"
)

const serviceName = "totp"

var copyToClipboard bool

func authenticate() error {
	ok, err := touchid.Authenticate("access totp code")
	if err != nil {
		log.Fatal(err)
		return err
	}

	if !ok {
		log.Fatal("Failed to authenticate")
	}
	return nil
}

func writeClipboard(text string) error {
	copyCmd := exec.Command("pbcopy")
	in, err := copyCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := copyCmd.Start(); err != nil {
		return err
	}
	if _, err := in.Write([]byte(text)); err != nil {
		return err
	}
	if err := in.Close(); err != nil {
		return err
	}
	return copyCmd.Wait()
}

func addItem(name, secret string) error {
	// Store it to the keychain
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(serviceName)
	item.SetAccount(name)
	item.SetLabel(name)
	item.SetData([]byte(secret))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenPasscodeSetThisDeviceOnly)
	return keychain.AddItem(item)
}

func delItem(name string) error {
	// Delete an item
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(serviceName)
	query.SetAccount(name)
	query.SetMatchLimit(keychain.MatchLimitOne)
	return keychain.DeleteItem(query)
}

func getItem(name string) (results []keychain.QueryResult, err error) {
	// Query an item
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(serviceName)
	query.SetAccount(name)
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)
	return keychain.QueryItem(query)
}

func queryItems() (results []keychain.QueryResult, err error) {
	// Query items
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(serviceName)
	query.SetMatchLimit(keychain.MatchLimitAll)
	query.SetReturnAttributes(true)
	return keychain.QueryItem(query)
}

func main() {
	log.SetFlags(0)
	var cmdGen = &cobra.Command{
		Use:   "gen <secret>",
		Short: "generate totp code from secret",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			secret := args[0]
			if secret == "" {
				return errors.New("No secret was given")
			}
			log.Printf(gotp.NewDefaultTOTP(string(secret)).Now())
			return nil
		},
	}

	var cmdAdd = &cobra.Command{
		Use:   "add <name> <secret>",
		Short: "Manually add a secret to the macOS keychain",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			secret := args[1]
			err := addItem(name, secret)
			if err != nil {
				return err
			}
			log.Printf("secret successfully registered as \"%v\".\n", name)
			return nil
		},
	}

	var cmdList = &cobra.Command{
		Use:   "ls",
		Short: "List all registered TOTP codes",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			results, err := queryItems()
			if err != nil {
				return err
			}
			for _, r := range results {
				log.Printf(r.Account)
			}
			return nil
		},
	}

	var cmdDelete = &cobra.Command{
		Use:   "del <name>",
		Short: "Delete a TOTP code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			authenticate()
			name := args[0]
			err := delItem(name)
			if err != nil {
				return err
			}
			log.Printf("Successfully deleted \"%v\".\n", name)
			return nil
		},
	}

	var rootCmd = &cobra.Command{
		Use:   "totp <name>",
		Short: "totp secrets manager",
		Args:  cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				authenticate()
				name := args[0]
				results, err := getItem(name)
				if err != nil {
					return err
				}
				if len(results) != 1 {
					return errors.New("Given name is not found")
				}
				r := results[0]
				// Generate a TOTP code
				totp := gotp.NewDefaultTOTP(string(r.Data)).Now()
				log.Printf(totp)
				if copyToClipboard {
					err := writeClipboard(totp)
					if err != nil {
						return err
					}
				}
			} else {
				cmd.Help()
			}
			return nil
		},
	}

	rootCmd.Flags().BoolVarP(
		&copyToClipboard,
		"copy",
		"c",
		false,
		"copy to clipboard",
	)

	rootCmd.AddCommand(cmdGen, cmdAdd, cmdList, cmdDelete)
	
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
	}
}
