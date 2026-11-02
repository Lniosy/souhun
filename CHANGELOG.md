# 更新记录

本文件记录已经做进应用里的变化。版本号对应 `macapp/Info.plist` 的 `CFBundleShortVersionString`。

## 1.0.0

- Mac 应用「搜魂」：选择 Bot、拉取记忆、写入、忘掉
- 启动时检查 Grok Bot、登录文件、钥匙串、登录是否过期、接口是否连通
- 写入时选择「账号共用 / 仅此 Bot」和「普通事实 / 指令」
- 每条记忆标出归属
- 个人记忆通过当前 Bot 的记忆工具读写；团队记忆不在这里写入
- 原理说明见 docs/principle.md
- Windows 64 位版本在 win/，界面与 Mac 版一致
- 文档同时说明 Mac 和 Windows，并提供英文版
