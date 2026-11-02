#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")"
SDK="$(xcrun --show-sdk-path)"
ROOT="$(pwd)"
rm -rf build dist
mkdir -p build dist
printf '// placeholder so SwiftBridging is only defined once\n' > build/empty-bridging.modulemap
cat > build/overlay.yaml << EOF
{
  "version": 0,
  "roots": [
    {
      "name": "/Library/Developer/CommandLineTools/usr/include/swift/bridging.modulemap",
      "type": "file",
      "external-contents": "$ROOT/build/empty-bridging.modulemap"
    }
  ]
}
EOF
clang -c crypt.c -o build/crypt.o -isysroot "$SDK" -arch arm64
xcrun swiftc -parse-as-library -vfsoverlay "$ROOT/build/overlay.yaml" -swift-version 5 -O \
  -framework AppKit -framework SwiftUI \
  -o build/GrokMemory GrokMemory.swift build/crypt.o
APP="dist/搜魂.app"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
cp build/GrokMemory "$APP/Contents/MacOS/GrokMemory"
cp Info.plist "$APP/Contents/Info.plist"
printf 'APPL????' > "$APP/Contents/PkgInfo"
ICONSET="build/AppIcon.iconset"
mkdir -p "$ICONSET"
sips -s format png AppIcon.png --out build/icon-1024.png >/dev/null
for spec in "16:icon_16x16.png" "32:icon_16x16@2x.png" "32:icon_32x32.png" "64:icon_32x32@2x.png" "128:icon_128x128.png" "256:icon_128x128@2x.png" "256:icon_256x256.png" "512:icon_256x256@2x.png" "512:icon_512x512.png" "1024:icon_512x512@2x.png"; do
  size="${spec%%:*}"
  name="${spec##*:}"
  sips -z "$size" "$size" build/icon-1024.png --out "$ICONSET/$name" >/dev/null
done
iconutil -c icns "$ICONSET" -o "$APP/Contents/Resources/AppIcon.icns"
cp AppIcon.png "$APP/Contents/Resources/AppIcon.png"
codesign --force --sign - "$APP"
echo "built $APP"
