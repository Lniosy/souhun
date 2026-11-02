# Contributing

[中文](CONTRIBUTING.md)

Keep a change on the problem in front of you. The Mac window follows `macapp/GrokMemory.swift`. The Windows window follows `win/ui.html`. Keep the wording aligned.

## Environment

- Apple silicon Mac, for the Mac app
- macOS 14 or later. This machine uses the macOS 15 SDK
- Xcode Command Line Tools: `xcode-select --install`
- Grok Bot installed and signed in

A full Xcode.app is not required. `macapp/build.sh` includes the Swift module workaround. See [docs/build.en.md](docs/build.en.md).

## Build

```bash
cd macapp
./build.sh
open "dist/搜魂.app"
```

On Windows, run `./build.sh` in `win/`. The output is `dist/搜魂.exe`. It needs Go 1.26 and can be cross-compiled. See [docs/build.en.md](docs/build.en.md).

The check prints status only, not an email address:

```bash
"dist/搜魂.app/Contents/MacOS/GrokMemory" --check
```

## Before you commit

- Do not commit `macapp/build/`, `macapp/dist/`, `win/dist/`, or `.DS_Store`
- Do not commit tokens, a keychain export, memory text, or an account email
- Say what a person can see, and how you checked it
- Screenshots use placeholder bot names

`grok_client.py` and `server.py` served the early local page. That page is retired. Put new behavior in the apps. If the Python write path still uses the old scope, say so, so the two copies do not drift in silence.
