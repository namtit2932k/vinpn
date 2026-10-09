<div align="center">

<img src="build/appicon.png" width="96" alt="Logo VinPN">

# VinPN

**DNS mã hoá cho Windows — với đường phục hồi bạn thực sự tin được.**

[![CI](https://github.com/sickyturtlez/vinpn/actions/workflows/ci.yml/badge.svg)](https://github.com/sickyturtlez/vinpn/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/sickyturtlez/vinpn?include_prereleases)](https://github.com/sickyturtlez/vinpn/releases)
[![Downloads](https://img.shields.io/github/downloads/sickyturtlez/vinpn/total)](https://github.com/sickyturtlez/vinpn/releases)
[![License: GPL v3](https://img.shields.io/badge/license-GPL--3.0-blue.svg)](LICENSE)
![Platform](https://img.shields.io/badge/platform-Windows%2010%2F11%20x64-blue)

[English](README.md) · Tiếng Việt

<a href="https://github.com/sickyturtlez/vinpn/releases/latest/download/vinpn-amd64-installer.exe"><img src="https://img.shields.io/badge/%E2%AC%87%20Download-Windows%2010%2F11%20x64-2ea44f?style=for-the-badge" alt="Tải VinPN cho Windows"></a>

</div>

---

VinPN là ứng dụng DNS cho Windows. Nó chạy một trình phân giải tại `127.0.0.1` / `::1`, chỉ mọi adapter mạng trỏ tới nó, và chuyển truy vấn của bạn qua **DoH, DoT, DoQ hoặc DNSCrypt** tới upstream khoẻ nhất. Ở mạng can thiệp kết nối mã hoá (DPI), nó còn điều khiển engine vượt chặn kèm theo — **zapret2** hoặc **GoodbyeDPI** — và có thể đóng vai **proxy cục bộ HTTP/HTTPS/SOCKS** cho máy tính, cũng như cho cả Wi-Fi nếu bạn bật.

Hai điểm làm nên khác biệt: **đầu vào đều được ký**, và **hệ thống đã động tới thì luôn trả lại nguyên vẹn**.

## VinPN cam kết

1. **DNS gốc luôn trở lại.** Mọi adapter được chụp ảnh DNS *trước khi* thay đổi, và bốn lớp độc lập phục hồi lại: ngắt kết nối sạch, tiến trình watchdog, khôi phục ở lần mở sau, và tác vụ chạy khi đăng nhập. Rút điện giữa phiên thì lần boot sau vẫn kết thúc bằng DNS cũ của bạn.
2. **Không gì đặc biệt mà lại nạp theo niềm tin.** Danh sách máy chủ, danh sách strategy và danh sách DNSCrypt đều được xác thực chữ ký (ed25519 / minisign) trước khi dùng; engine DPI được ghim hash và kiểm tra lại trước mỗi lần chạy; bất cứ file nào người dùng thường có thể ghi đều được sao chép vào thư mục admin-only trước khi tác vụ hay watchdog chạy nó với quyền cao.
3. **Không đo lường từ xa, tuyệt đối.** Không analytics, không báo lỗi crash, không ghi tên miền đã truy cập. Nhật ký truy vấn (nếu bật) chỉ nằm trong RAM và mặc định tắt; nhật ký trên đĩa chỉ ghi mã sự kiện.

## Bắt đầu nhanh

1. Tải và chạy bộ cài (có yêu cầu quyền quản trị viên — đổi DNS thì phải có).
2. Bấm **Connect**. Chế độ Đơn giản là đủ với hầu hết mọi người.
3. Xong: truy vấn giờ đi ra bằng DNS mã hoá, thanh trạng thái báo luôn phép kiểm tra rò rỉ có qua không.

Chế độ Nâng cao (bật trong ứng dụng) mở ra máy chủ, vượt DPI, proxy, rules, công cụ và DNS server cho mạng LAN.

## Năng lực

| Lĩnh vực | Bạn nhận được gì |
| --- | --- |
| **DNS mã hoá** | Upstream DoH / DoT / DoQ / DNSCrypt qua [dnsproxy](https://github.com/AdguardTeam/dnsproxy); quét song song, loại câu trả lời bị đầu độc, ghi nhớ máy chủ tốt nhất theo từng mạng. |
| **Vượt DPI** | Hai engine: **zapret2** (khuyến nghị — nhiều cách chia gói, QUIC, danh sách strategy có chữ ký kèm auto-tune) và **GoodbyeDPI**; tự dự phòng khi antivirus chặn cái này. |
| **Proxy cục bộ** | HTTP, HTTPS và SOCKS4/5; có thể trở thành system proxy của Windows hoặc chia cho điện thoại qua Wi-Fi (mã QR); tên miền phân giải qua DNS của VinPN. Phân mảnh TLS ClientHello sửa chặn SNI mà không cần driver. |
| **Rules** | Chặn, cho phép, trả lời giả, phân mảnh hoặc định tuyến theo domain, keyword, regexp hoặc CIDR. Nhập danh sách hosts, AdBlock, dnsmasq, Unbound, RPZ, Clash, v2ray, sing-box hoặc CIDR từ URL theo lịch. |
| **DNS server cho nhà** | DNS thường (cổng 53) hoặc DoH cho thiết bị khác trên Wi-Fi, kèm trang cài đặt bằng mã QR và cấu hình iOS. |
| **Fake SNI** *(mặc định tắt)* | Gửi một tên miền được phép tới mạng cho các site cho phép domain fronting. Chứng chỉ phiên chỉ ký được **đúng** các domain bạn tick, sống trong bộ nhớ, và gỡ ra khi ngắt kết nối. |
| **Công cụ** | Tra cứu DNS, quét cổng, quét IP sạch Cloudflare, soạn DNS stamp, sao lưu/khôi phục cài đặt. |

## Mô hình tin cậy

VinPN chạy với quyền cao, nên các quy tắc an toàn của nó rất rõ:

- **Một vùng admin-only.** `%ProgramData%\VinPN` giữ tất cả thứ điều khiển hệ thống từ bên ngoài — khóa LAN CA, `state.json` (ảnh chụp DNS/proxy) và các nhịp engine DPI. Thư mục được gắn DACL bảo vệ (chỉ SYSTEM + Administrators), quyền sở hữu được lấy lại về Administrators, và bất cứ thứ gì người dùng thường cấy vào đó đều bị xoá khi khởi động.
- **Tác vụ không bao giờ chạy một file người dùng ghi được.** Tác vụ khôi phục khi đăng nhập và watchdog khởi động lại chính file executable; nếu nó nằm ở nơi người dùng thường thay được (chế độ portable, bản checkout), nó được sao chép vào vùng admin-only trước, và bản sao mang theo thư mục dữ liệu gốc trên dòng lệnh.
- **Khóa đoán được không được phép gây mất mạng.** Khóa liên tiến trình có timeout thay vì chặn vô hạn; việc phục hồi chạy không cần khóa còn hơn là để DNS kẹt ở loopback.
- **Đầu vào có chữ ký.** Danh sách máy chủ, strategy, DNSCrypt và preset Fake SNI đều kiểm tra chữ ký; chữ ký sai là lỗi, không phải cảnh báo.
- **Đầu ra có phạm vi.** Quy tắc firewall được cho phép theo tên, chứng chỉ chỉ gỡ theo prefix phiên, phục hồi chỉ đụng adapter vẫn trỏ đúng chỗ VinPN để lại, và liên kết cập nhật chỉ mở `https://github.com/…`.

## Cách hoạt động

```
 ứng dụng của bạn
     │  DNS thường
     ▼
 trình phân giải cục bộ VinPN 127.0.0.1:53 ──► upstream mã hoá (DoH/DoT/DoQ/DNSCrypt)
     │                                              ▲
     │                                              └── quét sức khoẻ theo mạng,
     │                                                  loại câu trả lời bị đầu độc
     ▼
 proxy cục bộ (tuỳ chọn)  ◄── rules (chặn / giả / phân mảnh / upstream)
     │
     └──► engine DPI (zapret2 hoặc GoodbyeDPI) khi mạng cần
```

Connect: chụp ảnh DNS → cài trình phân giải → kiểm tra rò rỉ → bật watchdog và tác vụ phục hồi. Disconnect: đảo ngược mọi thứ theo cùng thứ tự.

## Quyền riêng tư

- Nhật ký truy vấn: **tắt** mặc định, chỉ RAM, xoá khi thoát.
- Tệp nhật ký: mã và số lượng sự kiện — không tên miền, không tên truy vấn.
- Bản sao lưu: cài đặt và danh sách của bạn; mật khẩu upstream được bảo vệ bằng DPAPI và che khi xuất.
- Kiểm tra cập nhật: GET HTTPS thuần tới GitHub releases API, mỗi 6 giờ, tắt được trong Cài đặt.

## Giới hạn đã biết

- Chỉ Windows 10/11 x64; engine DPI cần WinDivert (driver kernel).
- Ứng dụng crash để lại khoảng DNS còn trỏ loopback — chính thứ bốn lớp phục hồi cover.
- Chế độ portable phải giữ thư mục `data\` cạnh file thực thi.
- Fake SNI chỉ hữu ích sau các CDN chấp nhận fronting; nó không giải mã thứ gì bạn chưa chủ động thêm.

## Build từ nguồn

Yêu cầu: Go 1.27+, Node 24+, [wails3](https://wails.io).

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0-beta.27

git clone https://github.com/sickyturtlez/vinpn && cd vinpn
cd frontend && npm ci && cd ..

wails3 generate bindings -clean=true -ts -i   # sinh bindings TS từ service Go
cd frontend && npm run build && cd ..         # tạo frontend/dist (để nhúng)

wails3 build                                  # hoặc: wails3 dev
```

Kiểm tra thay đổi theo cách CI làm:

```bash
go test ./...
go test -race ./internal/rules/... ./internal/proxy/... ./internal/app/...
cd frontend && npm test
golangci-lint run
```

## Cấu trúc thư mục

| Đường dẫn | Nội dung |
| --- | --- |
| `internal/dnsserver`, `internal/sysdns` | Resolver cục bộ, chụp/phục hồi DNS adapter |
| `internal/proxy` | Proxy cục bộ, MITM cho Fake SNI, phân mảnh TLS |
| `internal/dpi` | Quản lý engine: ghim hash, giải nén, auto-tune |
| `internal/rules` | Ngôn ngữ rule, định dạng list, fetcher có xác thực chữ ký |
| `internal/servers`, `internal/upstreams` | Danh sách máy chủ, stamp, dựng upstream |
| `internal/watchdog` | Bốn lớp phục hồi |
| `internal/startup` | Tác vụ định kỳ và file thực thi đã staging |
| `internal/winutil`, `internal/certstore` | Plumbing Windows: firewall, dịch vụ, DPAPI, ACL |
| `frontend/` | Giao diện React (chế độ Đơn giản + Nâng cao), i18n en/vi |
| `lists/` | Các danh sách máy chủ, strategy và Fake SNI đã ký |

## Đóng góp và bảo mật

- Chào mừng báo lỗi và PR — xem [CONTRIBUTING.md](CONTRIBUTING.md) (test trước, không đo lường người dùng, không giết tiến trình bên thứ ba).
- Tìm ra lỗ hổng? Làm theo [SECURITY.md](SECURITY.md) — đừng mở issue công khai. Rò rỉ DNS, lỗi phục hồi và leo đặc quyền qua VinPN nằm trong phạm vi.

## Giấy phép

[GNU GPL v3.0 only](LICENSE). Xem [NOTICE](NOTICE) cho các thành phần bên thứ ba (zapret2, GoodbyeDPI, WinDivert, cygwin1.dll, dnsproxy, JetBrains Mono, Wails).

VinPN là công cụ mạng, không phải khiên pháp lý. Hãy dùng ở nơi bạn được phép và tuân theo luật địa phương.
