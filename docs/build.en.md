# Build

[中文](build.md)

From `macapp/`, run `./build.sh`. The script builds an arm64 app at `macapp/dist/搜魂.app`.

## Mac

| Item | Requirement |
|---|---|
| CPU | Apple silicon. The script pins `-arch arm64` |
| OS | Minimum declared is macOS 14. This machine builds with the macOS 15 SDK |
| Tools | Xcode Command Line Tools. A full Xcode.app is not required |
| Signature | Ad-hoc (`codesign -s -`). It opens on this Mac. There is no developer certificate |

## Swift module clash

With only the command-line tools, `SwiftBridging` is defined in two module maps, and `swiftc` fails. `build.sh` uses `-vfsoverlay` to replace `bridging.modulemap` with an empty file, and adds `-parse-as-library`, so `@main` can be the entry point.

Leave that workaround in place unless the toolchain on the machine no longer defines the module twice.

## Icon

The source is `macapp/AppIcon.png` (1024). The script builds `AppIcon.icns` and copies both into the app resources. The sidebar reads the PNG. The file itself has square corners. The window clips them to a continuous rounded rect.

## Windows

The Windows app is in `win/`, written in Go. The window matches the Mac app: sidebar, bot list, pull, the two write choices, and forget. The saved login is unlocked with Windows encryption, not the Mac keychain.

On a machine that can cross-compile:

```bash
cd win
./build.sh
```

The output is `win/dist/搜魂.exe`, 64-bit Windows only. The computer needs Microsoft Edge WebView2, and a Grok Bot that is installed and signed in. Check:

```bat
搜魂.exe --check
```

The keys match the Mac check: `app`, `secrets`, `keychain`, `login`, `api`, `ready`, `bots`, `login_named`. The name `keychain` is the same check. On Windows it reads system encryption.

## Check

```bash
macapp/dist/搜魂.app/Contents/MacOS/GrokMemory --check
```

It prints only `app`, `secrets`, `keychain`, `login`, `api`, `ready`, `bots`, `login_named`. It does not print an email address or a token.
