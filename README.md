totp secret manager, provide a cli interface to generate totp codes.

store your secrets in the macOS keychain and  protected by your macos touch id.

```
totp secrets manager

Usage:
  totp <name> [flags]
  totp [command]

Available Commands:
  add         Manually add a secret to the macOS keychain
  del         Delete a TOTP code
  gen         generate totp code from secret
  help        Help about any command
  ls          List all registered TOTP codes

Flags:
  -c, --copy   copy to clipboard
  -h, --help   help for totp

Use "totp

```