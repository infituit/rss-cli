# Release Source Archive Verification

GitHub Releases for `rss-cli` provide an independent GPG detached signature for
the source archive. Use the project's published public key and the signature
file to verify that the downloaded source archive was not modified after it was
signed.

## Signing Key

The release signing public key is published in `release-signing-key.asc`:

| Item        | Value                                               |
| ----------- | --------------------------------------------------- |
| Key name    | `rss-cli`                                           |
| Fingerprint | `72FA A718 8D3A 1B51 6093 D0EC 7D7E E6A1 E403 3897` |
| Public key  | `release-signing-key.asc`                           |

Always compare the complete fingerprint with a trusted project announcement.
`Good signature` only means that the signature matches the archive; it does not
prove that the public key belongs to the `rss-cli` project.

## Verify A Source Archive

Replace `<tag>` with the release tag to verify, for example `v0.1.0-alpha.1`.

```bash
curl -fsSLo rss-cli-<tag>.tar.gz \
  https://github.com/infituit/rss-cli/archive/refs/tags/<tag>.tar.gz
curl -fsSLo rss-cli-<tag>.tar.gz.asc \
  https://github.com/infituit/rss-cli/releases/download/<tag>/rss-cli-<tag>.tar.gz.asc
curl -fsSLo release-signing-key.asc \
  https://raw.githubusercontent.com/infituit/rss-cli/main/docs/operation/release-signing-key.asc

gpg --import release-signing-key.asc
gpg --fingerprint "rss-cli"
gpg --verify rss-cli-<tag>.tar.gz.asc rss-cli-<tag>.tar.gz
```

The `gpg --fingerprint` output must contain the complete fingerprint published
above. A successful verification prints `Good signature`. Do not use the
archive if verification prints `BAD signature`, the filenames do not match, or
the fingerprint is not the trusted project fingerprint.

## Verification Scope

This signature verifies only the GitHub Release source archive and its matching
`.asc` file. It does not verify:

- compiled binaries;
- npm packages;
- Docker images;
- third-party dependencies;
- source code correctness or runtime security.
