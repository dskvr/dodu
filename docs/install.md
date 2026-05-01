# Installing dodu

`dodu` is a single statically-linked Go binary. There are no runtime dependencies
beyond access to a Docker daemon socket.

## From a release archive

1. Download the archive matching your OS / arch from the
   [GitHub Releases page](https://github.com/tyutyutyu/dodu/releases).
2. Verify the checksum (recommended):

   ```bash
   sha256sum -c checksums.txt --ignore-missing
   ```

3. Extract and install:

   ```bash
   tar -xzf dodu_*_linux_amd64.tar.gz
   sudo install -m 0755 dodu /usr/local/bin/dodu
   dodu version
   ```

## From source (Go 1.23+)

```bash
git clone https://github.com/tyutyutyu/dodu.git
cd dodu
task build           # binary at ./bin/dodu
# or:
go install ./cmd/dodu@latest
```

## Permissions

`dodu` only needs read access to the Docker socket for normal operation:

- Linux: ensure your user is in the `docker` group, or set `DOCKER_HOST`.
- macOS: Docker Desktop's default socket works out of the box.

For destructive operations (`dodu prune --apply`), the same permissions are used.

## Read-only mode

To explicitly disable any destructive subcommand:

```bash
export DODU_READONLY=1
dodu prune --apply --yes        # → exits 3 with "read-only mode" error
```

## Uninstall

```bash
sudo rm /usr/local/bin/dodu
rm -rf ~/.cache/dodu ~/.local/state/dodu     # cache + audit log
```
