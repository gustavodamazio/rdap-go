# 🚀 Optimized Go RDAP Lookup (parallel, fast exit on first success)

### A fast and efficient IP lookup tool written in **Go** that queries multiple **RDAP servers in parallel**, returning standardized JSON information about the IP address owner, country, and block details.

> **Built for speed, reliability, and clean output.**  
> Supports both **IPv4** and **IPv6**.

---

## 🚀 Features

- **Parallel RDAP queries** to major RIRs (ARIN, RIPE, LACNIC, APNIC, AFRINIC)
- Returns consistent **JSON output**, including:
  - IP
  - IP Block Range
  - Organization Name
  - Holder
  - Country
  - City (if available)
- Automatically stops querying when first valid response is received
- Timeout handling (default 5s)
- Lightweight and dependency-free (pure Go)

---

## 📥 Installation

```bash
git clone https://github.com/gustavodamazio/rdap-go.git
cd rdap-go
go build -o rdaplookup main.go
```

---

## 🛠️ Usage

```bash
./rdaplookup <IP-ADDRESS>
```

**Example:**

```bash
./rdaplookup 8.8.8.8
```

**Output:**

```json
{
  "ip": "8.8.8.8",
  "start_address": "8.8.8.0",
  "end_address": "8.8.8.255",
  "organization": "Google LLC",
  "country": "US",
  "city": "",
  "holder": "GOGL"
}
```

---

## 📄 Example

```bash
./rdaplookup 2804:10c:1234::1
```

---

## 🌍 Supported RDAP Servers

- ARIN (North America)
- RIPE NCC (Europe)
- LACNIC (Latin America)
- APNIC (Asia Pacific)
- AFRINIC (Africa)

---

## ⚙️ Configuration Ideas

- Add local cache (optional)
- Export output to a file or pipe to another tool
- Dockerize for easy deployment

---

## 🤝 Contributions

Contributions, issues, and feature requests are welcome!  
Feel free to open a pull request or submit issues.

---

## 📄 License

This project is **MIT Licensed** — do whatever you want, just give credit!
