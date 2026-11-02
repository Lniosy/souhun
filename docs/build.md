# 编译

在 `macapp/` 里执行 `./build.sh`。脚本只生成 arm64 应用，输出在 `macapp/dist/搜魂.app`。

## 机器要求

| 项 | 要求 |
|---|---|
| 芯片 | Apple 芯片。脚本写死 `-arch arm64` |
| 系统 | 声明最低 macOS 14。当前用本机 SDK 15 编译 |
| 工具 | Xcode 命令行工具。不需要完整的 Xcode.app |
| 签名 | ad-hoc（`codesign -s -`）。本机打开即可，没有开发者证书 |

## Swift 模块冲突

只装了命令行工具时，`SwiftBridging` 会在两份 modulemap 里各定义一次，`swiftc` 因此失败。`build.sh` 用 `-vfsoverlay` 把 `bridging.modulemap` 换成空文件，并加上 `-parse-as-library`，这样 `@main` 才能作为程序入口。

不要删这段变通，除非机器上已经改用不会重复定义该模块的工具链。

## 图标

源图是 `macapp/AppIcon.png`（1024）。脚本会生成 `AppIcon.icns`，并把它和 PNG 一起放进应用资源。界面侧栏读的是 PNG。图的四角在文件里是方的，窗口里用连续圆角裁切。

## 自检

```bash
macapp/dist/搜魂.app/Contents/MacOS/GrokMemory --check
```

只打印 `app`、`secrets`、`keychain`、`login`、`api`、`ready`、`bots`、`login_named`。不打印邮箱和令牌。
