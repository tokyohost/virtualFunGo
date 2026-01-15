#!/usr/bin/env bash
set -e

APP_NAME="virtual-fun-go"
REPO_URL="https://github.com/tokyohost/virtualFunGo.git"
BRANCH="master"
WORKDIR="/opt/${APP_NAME}"
INSTALL_BIN="/usr/local/bin/${APP_NAME}"
SERVICE_FILE="${APP_NAME}.service"
MIN_GO_VERSION="1.24"

echo "======================================"
echo " Installing ${APP_NAME} (Go >= ${MIN_GO_VERSION})"
echo "======================================"

# ---------- root ----------
if [[ $EUID -ne 0 ]]; then
  echo "[ERROR] Please run as root"
  exit 1
fi

# ---------- detect OS ----------
if [ ! -f /etc/os-release ]; then
  echo "[ERROR] Cannot detect OS"
  exit 1
fi
. /etc/os-release
echo "[INFO] OS: $ID"

# ---------- install git ----------
if ! command -v git >/dev/null 2>&1; then
  echo "[INFO] Installing git..."
  case "$ID" in
    debian|ubuntu)
      apt update
      apt install -y git curl tar
      ;;
    centos|rhel|rocky|almalinux)
      yum install -y git curl tar
      ;;
    *)
      echo "[ERROR] Unsupported OS"
      exit 1
      ;;
  esac
fi

# ---------- fix git HTTP2 issue ----------
git config --global http.version HTTP/1.1
git config --global http.postBuffer 524288000

# ---------- install Go ----------
install_go_binary() {
  GO_TAR="go1.24.11.linux-amd64.tar.gz"
  DOWNLOAD_URL="https://mirrors.tuna.tsinghua.edu.cn/go/${GO_TAR}"

  echo "[INFO] Downloading Go binary ${DOWNLOAD_URL}"
  curl -L --http1.1 -k "$DOWNLOAD_URL"
  rm -rf /usr/local/go
  tar -C /usr/local -xzf "$GO_TAR"
  rm "$GO_TAR"

  export PATH="/usr/local/go/bin:$PATH"
}

check_go_version() {
  if ! command -v go >/dev/null 2>&1; then
    return 1
  fi
  CURRENT=$(go version | awk '{print $3}' | sed 's/go//')
  # compare version
  if [[ $(printf '%s\n' "$MIN_GO_VERSION" "$CURRENT" | sort -V | head -n1) != "$MIN_GO_VERSION" ]]; then
    return 0
  else
    return 1
  fi
}

if ! check_go_version; then
  echo "[INFO] Installing Go >= ${MIN_GO_VERSION}..."
  install_go_binary
else
  echo "[INFO] Found Go $(go version)"
fi

# ---------- clone repository ----------
echo "[INFO] Cloning repository..."
rm -rf "$WORKDIR"
mkdir -p "$WORKDIR"

GIT_HTTP_VERSION=HTTP/1.1 git clone --depth=1 -b "$BRANCH" "$REPO_URL" "$WORKDIR"

# ---------- build ----------
echo "[INFO] Building ${APP_NAME}..."
cd "$WORKDIR"
go build -o "$APP_NAME"

# ---------- install binary ----------
echo "[INFO] Installing binary to ${INSTALL_BIN}"
install -m 0755 "$APP_NAME" "$INSTALL_BIN"

# ---------- install systemd service ----------
if [ ! -f "$WORKDIR/${SERVICE_FILE}" ]; then
  echo "[ERROR] ${SERVICE_FILE} not found in repo"
  exit 1
fi

echo "[INFO] Installing systemd service"
install -m 0644 "$WORKDIR/${SERVICE_FILE}" "/etc/systemd/system/${SERVICE_FILE}"

systemctl daemon-reexec
systemctl daemon-reload
systemctl enable "${APP_NAME}.service"

# ---------- start service ----------
echo "[INFO] Starting service"
systemctl restart "${APP_NAME}.service"

echo
echo "=========== SERVICE STATUS ==========="
systemctl status "${APP_NAME}.service" --no-pager || true
echo "===================================="

echo
echo "[SUCCESS] ${APP_NAME} installed successfully"
