# 参与修改

[English](CONTRIBUTING.en.md)

改动保持在当前问题上。Mac 界面以 `macapp/GrokMemory.swift` 为准，Windows 界面以 `win/ui.html` 为准，两边的文案保持一致。

## 环境

- Apple 芯片的 Mac
- macOS 14 或更新版本，本机用的是 15 的 SDK
- Xcode 命令行工具：`xcode-select --install`
- 已安装并登录 Grok Bot

没有完整的 Xcode.app 也能编译。`macapp/build.sh` 里带了 Swift 模块冲突的变通，说明在 [docs/build.md](docs/build.md)。

## 构建

```bash
cd macapp
./build.sh
open "dist/搜魂.app"
```

Windows 在 `win/` 执行 `./build.sh`，得到 `dist/搜魂.exe`。需要 Go 1.26，可以在别的系统上交叉编译。说明在 [docs/build.md](docs/build.md)。

自检只打印状态，不打印邮箱：

```bash
"dist/搜魂.app/Contents/MacOS/GrokMemory" --check
```

## 提交前

- 不要提交 `macapp/build/`、`macapp/dist/`、`.DS_Store`
- 不要提交令牌、钥匙串导出、记忆正文或账号邮箱
- 说明用户能看见的变化，以及你怎么验证的

`grok_client.py` 和 `server.py` 是早期本机网页用的。网页已经停用。新行为请改 Swift 应用；如果 Python 里的写入范围还是旧的，在说明里写一句，避免两边各写各的。
