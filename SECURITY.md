# 安全说明

[English](SECURITY.en.md)

搜魂用你本机已经登录的 Grok Bot 身份去读、写记忆。它不另做账号，也不把登录信息写进这个仓库。

## 请不要公开这些东西

- 钥匙串里的口令、访问令牌、刷新令牌
- `sand-secrets.json` 的内容
- 记忆正文、账号邮箱、机器标识

问题复现只要写操作步骤和界面上的报错。令牌打码后再贴。

## 这个程序会碰什么

| 位置 | 用途 |
|---|---|
| Mac：`/Applications/Grok Bot.app` | 确认已安装，并读取版本号 |
| Windows：本机的 `Grok Bot.exe` | 确认已安装，并读取版本号 |
| `sand-secrets.json`（Mac 在 `~/Library/Application Support/Grok Bot/`，Windows 在 `%APPDATA%\Grok Bot\`） | 读出本机保存的登录记录 |
| Mac 钥匙串 `Grok Bot Safe Storage`，或 Windows 的系统加密 | 解开上面那份登录记录 |
| `api2.cursor.sh` 上的 Grok Bot 接口 | 列出 Bot、发送读写请求 |

登录记录只留在这次运行的内存里。仓库里没有示例密钥。

## 报告问题

仓库公开之后，用 GitHub 的私密安全通告。不要开公开 Issue 贴凭证。
