# Security

[中文](SECURITY.md)

Souhun reads and writes memory with the Grok Bot login already on this computer. It does not create another account, and it does not commit login material to this repository.

## Do not publish

- The keychain password, access token, or refresh token
- The contents of `sand-secrets.json`
- Memory text, the account email, or the machine id

A bug report needs the steps and the error shown in the window. Redact any token before pasting it.

## What this program touches

| Place | Why |
|---|---|
| Mac: `/Applications/Grok Bot.app` | Confirm it is installed and read the version |
| Windows: the local `Grok Bot.exe` | Confirm it is installed and read the version |
| `sand-secrets.json` (Mac: `~/Library/Application Support/Grok Bot/`, Windows: `%APPDATA%\Grok Bot\`) | Read the saved login |
| Mac keychain `Grok Bot Safe Storage`, or Windows encryption | Unlock that login |
| The Grok Bot API on `api2.cursor.sh` | List bots and send read or write requests |

The unlocked login stays in memory for that run. The repository has no sample secrets.

## Reporting

After the repository is public, use a private GitHub security advisory. Do not open a public issue that contains credentials.
