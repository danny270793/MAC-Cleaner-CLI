set -e

REPO="danny270793/MAC-Cleaner-CLI"
BINARY_NAME="maccleaner"
INSTALL_DIR="$HOME/.local/bin"

if [ "$(uname -s)" != "Darwin" ]; then
  echo "$BINARY_NAME only supports macOS" >&2
  exit 1
fi

case "$(uname -m)" in
  arm64) arch="arm64" ;;
  x86_64) arch="amd64" ;;
  *)
    echo "unsupported architecture: $(uname -m)" >&2
    exit 1
    ;;
esac

tag=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
if [ -z "$tag" ]; then
  echo "failed to resolve the latest release" >&2
  exit 1
fi
version="${tag#v}"

archive="${BINARY_NAME}_${version}_darwin_${arch}.tar.gz"
workdir=$(mktemp -d)
trap 'rm -rf "$workdir"' EXIT

echo "downloading $archive ($tag)..."
curl -fsSL "https://github.com/$REPO/releases/download/$tag/$archive" -o "$workdir/$archive"

echo "verifying checksum..."
curl -fsSL "https://github.com/$REPO/releases/download/$tag/checksums.txt" -o "$workdir/checksums.txt"
(cd "$workdir" && grep "$archive" checksums.txt | shasum -a 256 -c -)

tar -xzf "$workdir/$archive" -C "$workdir"

mkdir -p "$INSTALL_DIR"
mv "$workdir/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
chmod +x "$INSTALL_DIR/$BINARY_NAME"

echo "installed $BINARY_NAME $version to $INSTALL_DIR/$BINARY_NAME"

case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "$INSTALL_DIR is not on your PATH. Add this to your shell profile:"
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    ;;
esac
