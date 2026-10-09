<div align="center">

<img src="build/appicon.png" width="96" alt="VinPN logo">

# VinPN

**Encrypted DNS for Windows — with a recovery path you can actually trust.**

[![CI](https://github.com/namtit2932k/vinpn/actions/workflows/ci.yml/badge.svg)](https://github.com/namtit2932k/vinpn/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/namtit2932k/vinpn?include_prereleases)](https://github.com/namtit2932k/vinpn/releases)
[![Downloads](https://img.shields.io/github/downloads/namtit2932k/vinpn/total)](https://github.com/namtit2932k/vinpn/releases)
[![License: GPL v3](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)
![Platform](https://img.shields.io/badge/platform-Windows%2010%2F11%20x64-blue)

English · [Tiếng Việt](README.vi.md)

<a href="https://github.com/namtit2932k/vinpn/releases/latest/download/vinpn-amd64-installer.exe"><img src="https://img.shields.io/badge/%E2%AC%87%20Download-Windows%2010%2F11%20x64-2ea44f?style=for-the-badge" alt="Download VinPN for Windows"></a>

</div>

---

VinPN is a desktop DNS client for Windows. It runs a resolver on `127.0.0.1` / `::1`, points every network adapter at it, and forwards your queries over **DoH, DoT, DoQ or DNSCrypt** to the fastest healthy upstream. On networks that inspect and reset encrypted connections (DPI), it can also drive a bundled bypass engine — **zapret2** or **GoodbyeDPI** — and it can act as a local **HTTP/HTTPS/SOCKS proxy** for the machine and, on request, for your whole LAN.

Two things set it apart: **the inputs it trusts are signed**, and **the system it touches it always gives back**.

## What VinPN promises

1. **Your original DNS comes back.** Every adapter's DNS is snapshotted *before* the first change, and four independent layers restore it: a clean disconnect, a watchdog process, restore on next launch, and a logon recovery task. Pull the power cord mid-session and the next boot still ends with your old DNS.
2. **Nothing privileged is loaded on trust.** The server list, the strategy list and the DNSCrypt list are cryptographically verified (ed25519 / minisign) before use; the DPI engines are hash-pinned and re-verified before every start; anything a normal user could overwrite is staged into an admin-only directory before a task or the watchdog runs it elevated.
3. **No telemetry, ever.** No analytics, no crash reporting, no visited-domain logging. The optional query log lives in RAM only and is off by default; the on-disk log records event codes, not names.

## Quick start

1. Download and run the installer (it asks for administrator rights — DNS changes need them).
2. Press **Connect**. Simple mode is all most people need.
3. Done: queries now leave over encrypted DNS, and the status line tells you whether the leak check passed.

Advanced mode (toggle in the app) unlocks servers, DPI bypass, the proxy, rules, tools and the LAN DNS server.

## Capabilities

| Area | What you get |
| --- | --- |
| **Encrypted DNS** | DoH / DoT / DoQ / DNSCrypt upstreams via [dnsproxy](https://github.com/AdguardTeam/dnsproxy); parallel scanning rejects poisoned answers and remembers the best servers per network. |
| **DPI bypass** | Two engines: hash-pinned **zapret2** (recommended — split methods, QUIC, signed strategy list with auto-tune) and **GoodbyeDPI**; automatic fallback when antivirus blocks one. |
| **Local proxy** | HTTP, HTTPS and SOCKS4/5; can become the Windows system proxy or be shared with phones over Wi-Fi (QR code); names resolve through VinPN's DNS. TLS ClientHello fragmentation fixes SNI blocking without a driver. |
| **Rules** | Block, allow, fake-answer, fragment or route by domain, keyword, regexp or CIDR. Import hosts, AdBlock, dnsmasq, Unbound, RPZ, Clash, v2ray, sing-box or CIDR lists from a URL on a schedule. |
| **Home DNS server** | Plain DNS (port 53) or DoH for other devices on your Wi-Fi, with a QR setup page and an iOS configuration profile. |
| **Fake SNI** *(off by default)* | Send an allowed domain name to the network for sites behind fronting-friendly CDNs. The session certificate can sign **only** the domains you tick, lives in memory, and is removed on disconnect. |
| **Tools** | DNS lookup, port scanner, Cloudflare clean-IP scan, DNS stamp editor, settings backup/restore. |

## Trust and safety model

VinPN runs elevated, so its safety rules are explicit:

- **One admin-only zone.** `%ProgramData%\VinPN` holds everything that drives the system from the outside — the LAN CA key, `state.json` (the DNS/proxy snapshot) and the DPI engine binaries. The directory gets a protected DACL (SYSTEM + Administrators only), ownership is reclaimed from Administrators, and anything planted there by a normal user is removed at startup.
- **Tasks never run a user-writable image.** The logon recovery task and the watchdog re-launch the executable; if it lives somewhere a normal user can replace (portable mode, a checkout), it is first staged into the admin-only zone, and the staged copy carries the original data directory on its command line.
- **Guessable locks cannot cause outages.** The cross-process state lock times out instead of blocking forever; recovery proceeds without it rather than leaving your DNS on loopback.
- **Signed inputs.** Server list, strategy list, DNSCrypt resolvers and the Fake SNI presets are signature-checked; a bad signature is an error, not a warning.
- **Scoped output.** Firewall rules are allowlisted by name, certificates are removed only by their session prefix, restore only touches adapters that still point where VinPN left them, and update links open only `https://github.com/…`.

## How it works

```
 your apps
     │  standard DNS
     ▼
 VinPN local resolver 127.0.0.1:53 ──► encrypted upstream (DoH/DoT/DoQ/DNSCrypt)
     │                                        ▲
     │                                        └── per-network health scan,
     │                                            poisoned answers rejected
     ▼
 optional local proxy  ◄── rules (block / fake / fragment / upstream)
     │
     └──► DPI engine (zapret2 or GoodbyeDPI) when the network needs it
```

Connect snapshots DNS → installs the resolver → verifies the leak → starts the watchdog and the recovery task. Disconnect reverses everything in the same order.

## Privacy

- Query log: **off** by default, RAM only, cleared on exit.
- Log files: event codes and counts — no domains, no query names.
- Backups: your settings and lists; upstream passwords are DPAPI-protected and redacted from exports.
- Update checks: plain HTTPS GET to the GitHub releases API, every 6 hours, disable it in Settings.

## Known limitations

- Windows 10/11 x64 only; the DPI engines need WinDivert (a kernel driver).
- An app crash leaves a brief window where DNS is still loopback — that is exactly what the recovery layers cover.
- Portable mode needs its `data\` folder kept next to the executable.
- Fake SNI only helps behind CDNs that tolerate fronting; it decrypts nothing you did not explicitly add.

## Building from source

Requirements: Go 1.27+, Node 24+, [wails3](https://wails.io).

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.27

git clone https://github.com/namtit2932k/vinpn && cd vinpn
cd frontend && npm ci && cd ..

wails3 generate bindings -clean=true -ts -i   # TS bindings from the Go services
cd frontend && npm run build && cd ..         # produces frontend/dist (embedded)

wails3 build                                  # or: wails3 dev
```

Verify a change the way CI does:

```bash
go test ./...
go test -race ./internal/rules/... ./internal/proxy/... ./internal/app/...
cd frontend && npm test
golangci-lint run
```

## Repository layout

| Path | Contents |
| --- | --- |
| `internal/dnsserver`, `internal/sysdns` | Local resolver, adapter DNS snapshot/restore |
| `internal/proxy` | Local proxy, MITM for Fake SNI, TLS fragmentation |
| `internal/dpi` | Engine manager: hash pins, extraction, auto-tune |
| `internal/rules` | Rule language, list formats, fetcher with signature checks |
| `internal/servers`, `internal/upstreams` | Server list, stamps, upstream building |
| `internal/watchdog` | The four restore layers |
| `internal/startup` | Scheduled tasks and the staged executable |
| `internal/winutil`, `internal/certstore` | Windows plumbing: firewall, services, DPAPI, ACLs |
| `frontend/` | React UI (simple + advanced modes), i18n en/vi |
| `lists/` | Signed server, strategy and Fake SNI lists |

## Contributing and security

- Bugs and pull requests welcome — see [CONTRIBUTING.md](CONTRIBUTING.md) for the rules (tests first, no telemetry, no third-party process killing).
- Found a vulnerability? Please follow [SECURITY.md](SECURITY.md) — do not open a public issue. DNS leaks, restore failures and privilege escalation through VinPN are explicitly in scope.

## License

[GNU GPL v3.0 only](LICENSE). See [NOTICE](NOTICE) for bundled third-party components (zapret2, GoodbyeDPI, WinDivert, cygwin1.dll, dnsproxy, JetBrains Mono, Wails).

VinPN is a network tool, not a legal shield. Use it where you are allowed to, and follow your local law.
