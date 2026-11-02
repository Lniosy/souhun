# Memory

[中文](memory.md)

Souhun manages the memory Grok Bot stores in the cloud. It is not the chat log, and it is not the persona or the skills.

## Two choices

A write needs both.

| In the window | Meaning | Type the bot writes |
|---|---|---|
| Shared with the account · instruction | Every bot for this user will use it | `shared_profile` |
| Shared with the account · fact | Every bot may recall it when relevant | `shared_journal` |
| This bot only · instruction | Only this bot will use it | `profile` |
| This bot only · fact | Only this bot may recall it when relevant | `journal` |

The window still sends `scope` (`user` or `agent`) and `tier` (`profile` or `log`). The bot maps those to the types above. That map is the bot's own account of itself. The installed app does not contain a table to check it against.

Labels on each row:

| Label | Scope |
|---|---|
| Shared with the account | `user` |
| This bot only | `agent` |
| Team | `team` |
| Unlabeled | This reply did not give a scope |

## Which path a call takes

```
Souhun
  │
  ├─ List bots, confirm login     Call the Grok Bot API directly
  │
  └─ Pull / write / forget        Send one chat message to the selected bot
                                    The bot calls UpdateMemory or recall
                                    The result is read back from that chat
```

Personal memory has no desktop call that creates, updates, or deletes an item. `UpdateMemory` and `RecallMemory` are tools in the bot runtime.

So:

- Each write or forget leaves a request in that bot's chat
- A pull asks the bot to export. It is not a full-list call. One pull returns as many as the bot can recall, and the window says to pull again when the export is incomplete
- Team memory is not written here. Team review is separate, and a personal bot is refused by those calls

## What this memory is not

Persona, skills, routines, and chat history live somewhere else. Souhun does not change them.

A fact that was never marked shared can be recalled only by the bot that saved it. A team bot has its own summary, which can lag behind the personal memory shown here.
