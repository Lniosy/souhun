# How it works

[中文](principle.md)

Souhun does not keep a copy of memory. Memory stays in Grok Bot's cloud. The app does two things: call the API with the login already on this computer, and, when personal memory has no direct call, ask the current bot to read or write it.

The Mac UI is `macapp/GrokMemory.swift`. The Windows UI is `win/`. Both use the same API. `grok_client.py` is an early client with a fixed version string. Both apps read the Grok Bot version installed on the machine.

## Where the data is

| Place | What is there | What is not |
|---|---|---|
| Mac keychain `Grok Bot Safe Storage` / `Grok Bot Key` | Password that unlocks the saved login | Memory text |
| Mac `~/Library/Application Support/Grok Bot/sand-secrets.json` | Encrypted access token and machine id | Memory text |
| Windows `%APPDATA%\Grok Bot\Local State` | System key that unlocks the saved login | Memory text |
| Windows `%APPDATA%\Grok Bot\sand-secrets.json` | Encrypted access token and machine id | Memory text |
| `~/.grokbot` | Local settings | A memory database |
| Grok Bot cloud | Bot list, chats, memory | A Souhun copy |

On Mac the login blob is Chromium `v10`. The password comes from the keychain item above. Decryption is in `macapp/crypt.c`. On Windows the app first unlocks the key in `Local State` with Windows encryption, then unlocks the same `v10` blob. Either way the result stays in memory for that run.

The token is one login session. The app uses its expiry for the login check. After it expires, sign in again in Grok Bot.

## API

Listing bots, sending a message, and reading the transcript are Connect JSON unary calls:

```
POST https://api2.cursor.sh/aiserver.v1.GrokBotService/<Method>
```

The request carries the local login token, a checksum built from the machine id, and the desktop client type, version, and OS. The checksum is computed per request and is not stored in the repository.

Methods the app actually calls:

| Method | Role |
|---|---|
| `ListGrokBotAgents` | The bot list |
| `ListGrokBotUserBotMemories` | Try to list memories directly. Personal bots are usually refused |
| `SendGrokBotUserMessage` | Send a pull, write, or forget to the current bot |
| `ListGrokBotTranscriptEntries` | Read that chat and take the bot's reply |

A direct list is allowed only for an owner or manager of a team bot. A personal bot is refused. The app treats that refusal as "this path is closed" and uses the chat instead. It does not show the refusal as a failure.

## Why a write goes through the bot

Adding and removing personal memory is a tool inside the bot runtime. It is not in the method list the desktop client can call.

| Capability | Desktop | Bot runtime |
|---|---|---|
| List bots, send a message, read the chat | Yes | |
| Change one personal memory directly | No such method | `UpdateMemory` |
| Recall by content | No full list | `RecallMemory`, ranked, with a limit |
| List, submit, or review team memory | Methods exist; a personal bot is refused | Souhun does not write these |

Pull, write, and forget share one path:

```
Selected bot
    │  SendGrokBotUserMessage
    ▼
The cloud bot calls its own memory tool
    │
    ▼
ListGrokBotTranscriptEntries
    Only a reply after this request, carrying this request id
```

The chat first shows progress. The app ignores that until the reply contains this request id, or, for a pull, a memory list. A write or forget must come back as JSON with `ok`. The UI counts it as done only when the id matches and `ok` is true.

If no qualifying reply arrives within three minutes, the app stops and tells you to look at that chat in Grok Bot.

## How the two choices are sent

"Where" and "kind" are written into the message: `scope` is `user` or `agent`, `tier` is `profile` or `log`. The bot maps those to its own types. The map is in [Memory](memory.en.md).

Forget does not send a new place or kind. It sends the original text and asks the bot to remove that saved item, not to write another one.

## What the checks look at

The launch checks stop at the first failure. They do not build a session past that point:

```
Is Grok Bot installed?
        │  Mac looks for Grok Bot.app, Windows for Grok Bot.exe
Does the login file exist?
        │
Can the token and machine id be unlocked?
        │  Mac uses the keychain, Windows uses system encryption
Is the token unexpired?
        │
Does ListGrokBotAgents return a bot list?
```

`--check` prints pass or fail for each item. It does not print an email address or a token.
