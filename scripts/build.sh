#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DIST_DIR="${PROJECT_ROOT}/dist"
APP_NAME="SpiritGen"
GUI_PKG="./cmd/spiritgui"

usage() {
  echo "Usage: $(basename "$0") [mac|win|all]"
  echo ""
  echo "  mac   macOS .app 번들 생성 (fyne CLI 필요)"
  echo "  win   Windows .exe 생성 (mingw-w64 크로스 컴파일)"
  echo "  all   mac + win 모두 빌드"
  echo ""
  echo "기본값: mac"
}

require_cmd() {
  if ! command -v "$1" &>/dev/null; then
    echo "❌ '$1' 명령을 찾을 수 없습니다. $2"
    exit 1
  fi
}

build_mac() {
  echo "▶ macOS 빌드 중..."
  require_cmd fyne "설치: go install fyne.io/fyne/v2/cmd/fyne@latest"

  mkdir -p "${DIST_DIR}"
  cd "${PROJECT_ROOT}"
  fyne package \
    -os darwin \
    -name "${APP_NAME}" \
    -sourceDir "${GUI_PKG}" \
    -output "${DIST_DIR}/${APP_NAME}.app"

  echo "✅ macOS 빌드 완료: ${DIST_DIR}/${APP_NAME}.app"
}

build_win() {
  echo "▶ Windows 빌드 중..."
  require_cmd x86_64-w64-mingw32-gcc "설치: brew install mingw-w64"

  mkdir -p "${DIST_DIR}"
  cd "${PROJECT_ROOT}"
  CGO_ENABLED=1 \
  GOOS=windows \
  GOARCH=amd64 \
  CC=x86_64-w64-mingw32-gcc \
    go build \
      -ldflags="-H windowsgui" \
      -o "${DIST_DIR}/spiritgui.exe" \
      "${GUI_PKG}"

  echo "✅ Windows 빌드 완료: ${DIST_DIR}/spiritgui.exe"
}

TARGET="${1:-mac}"

case "${TARGET}" in
  mac)
    build_mac
    ;;
  win)
    build_win
    ;;
  all)
    build_mac
    build_win
    ;;
  -h|--help|help)
    usage
    ;;
  *)
    echo "❌ 알 수 없는 타겟: ${TARGET}"
    usage
    exit 1
    ;;
esac
