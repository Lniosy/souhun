# 原理

[English](principle.en.md)

搜魂不保存记忆。记忆在 Grok Bot 的云端。这个应用只做两件事：用本机已经登录的身份调用接口，以及在个人记忆没有直接接口时，请当前 Bot 代为读写。

Mac 界面在 `macapp/GrokMemory.swift`，Windows 界面在 `win/`。两边走同一套接口。`grok_client.py` 是早期实现，客户端版本号写死在文件里；两个应用都会读取本机 Grok Bot 的版本。

## 数据在哪

| 位置 | 有什么 | 没有什么 |
|---|---|---|
| Mac 钥匙串 `Grok Bot Safe Storage` / `Grok Bot Key` | 解开登录记录用的口令 | 记忆正文 |
| Mac `~/Library/Application Support/Grok Bot/sand-secrets.json` | 加密后的访问令牌和机器标识 | 记忆正文 |
| Windows `%APPDATA%\Grok Bot\Local State` | 解开登录记录用的系统密钥 | 记忆正文 |
| Windows `%APPDATA%\Grok Bot\sand-secrets.json` | 加密后的访问令牌和机器标识 | 记忆正文 |
| `~/.grokbot` | 本机设置 | 记忆库 |
| Grok Bot 云端 | Bot 列表、对话、记忆 | 搜魂自己的副本 |

Mac 上，登录记录是 Chromium 的 `v10` 封装。口令来自上面的钥匙串项，解密在 `macapp/crypt.c`。Windows 上，先用系统加密解开 `Local State` 里的密钥，再解开同一份 `v10` 记录。两边解开后都只留在这次运行的内存里。

令牌是一次登录会话。应用用它的过期时间判断「登录状态」。过期后要回到 Grok Bot 重新登录。

## 接口

列出 Bot、发消息、读对话，走 Connect 的 JSON 一元调用：

```
POST https://api2.cursor.sh/aiserver.v1.GrokBotService/<方法名>
```

请求带本机登录令牌、机器标识做成的校验值，以及桌面客户端的类型、版本和系统。校验值每次现算，不写入仓库。

应用实际用到的方法：

| 方法 | 作用 |
|---|---|
| `ListGrokBotAgents` | 左边的 Bot 列表 |
| `ListGrokBotUserBotMemories` | 尝试直接列出记忆。个人 Bot 通常被拒绝 |
| `SendGrokBotUserMessage` | 把拉取、写入或忘掉发给当前 Bot |
| `ListGrokBotTranscriptEntries` | 读这条对话，取出 Bot 的回复 |

直接列出记忆时，服务端只允许团队 Bot 的拥有者或管理者。个人 Bot 会返回拒绝。应用把这类拒绝当成「这条路不通」，改走对话，不把它显示成故障。

## 为什么写入要经过 Bot

个人记忆的增删在 Bot 运行时的工具里，不在桌面端能调用的方法列表里。

| 能力 | 桌面端 | Bot 运行时 |
|---|---|---|
| 列出 Bot、发消息、读对话 | 可以 | |
| 直接改一条个人记忆 | 没有这个方法 | `UpdateMemory` |
| 按内容回忆 | 没有完整列表 | `RecallMemory`，按相关度返回，有上限 |
| 团队记忆的列出、提交、审核 | 有方法，个人 Bot 会被拒绝 | 不从搜魂写入 |

所以拉取、写入、忘掉都是同一条路径：

```
选中的 Bot
    │  SendGrokBotUserMessage
    ▼
云端 Bot 调用自己的记忆工具
    │
    ▼
ListGrokBotTranscriptEntries
    只认这次请求之后、带有本次编号的回复
```

对话里会先出现「正在汇总」这类进度。应用忽略它，直到回复里出现这次的请求编号，或拉取时出现记忆列表。写入和忘掉要求 Bot 只回一个带 `ok` 的 JSON。编号对上并且 `ok` 为真，界面才算成功。

三分钟内没有合格回复，就停下来，提示到 Grok Bot 里看这条对话。

## 写入时两个选择怎么落地

界面上的「写到」和「属性」会写进发给 Bot 的原文：`scope` 是 `user` 或 `agent`，`tier` 是 `profile` 或 `log`。Bot 再换成自己的类型。对应关系见 [记忆怎么放](memory.md)。

忘掉不带新的位置和属性。它把这条的原文交给 Bot，请它从已保存的记忆里去掉，不要另写一条。

## 自检在查什么

启动时按这个顺序，前一项失败就不再往后装登录会话：

```
本机是否有 Grok Bot
        │  Mac 找 Grok Bot.app，Windows 找 Grok Bot.exe
登录文件是否存在
        │
能否解开令牌和机器标识
        │  Mac 用钥匙串，Windows 用系统加密
令牌是否过期
        │
ListGrokBotAgents 是否返回 Bot 列表
```

`--check` 只打印每项是否通过，不打印邮箱和令牌。
