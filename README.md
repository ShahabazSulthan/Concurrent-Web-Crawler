# 🕷️ Concurrent Web Crawler in Go

This is a simple and efficient concurrent web crawler written in Golang. It reads a list of URLs from a file (`urls.txt`), fetches each URL concurrently with a timeout, extracts the HTML title, and logs the results.

## 🚀 Features

- Concurrent crawling with configurable concurrency
- Graceful handling of timeouts using Go's `context`
- Deduplication of URLs
- Title extraction from HTML
- Logging of HTTP status and errors
- Easy-to-read, beginner-friendly Go code

---

## 📂 Project Structure

```
.
├── main.go         # Entry point and core logic
├── urls.txt        # List of URLs to crawl
└── README.md       # Project documentation
```

---

## 🛠️ Requirements

- Go 1.18 or higher

---

## 📦 Installation

1. **Clone the repo**
   ```bash
   git clone https://github.com/your-username/go-web-crawler.git
   cd go-web-crawler
   ```

2. **Create `urls.txt`**
   Add one URL per line. Example:
   ```
   https://golang.org
   https://example.com
   https://google.com
   ```

3. **Run the program**
   ```bash
   go run main.go
   ```

---

## ⚙️ Configuration

You can modify the following constants in `main.go` to tweak behavior:

```go
const (
    MaxConcurrency = 5                   // Maximum number of concurrent HTTP requests
    RequestTimeout = 5 * time.Second     // Timeout per request
    URLFile        = "urls.txt"          // File with list of URLs
)
```

---

## 📤 Output Example

```
2025/04/10 12:34:56 [OK] https://golang.org (200) -> The Go Programming Language
2025/04/10 12:34:57 [ERROR] https://broken-link.com: Get "https://broken-link.com": dial tcp ...
2025/04/10 12:34:58 ✅ Done in 2.134123s
```

---

## 🧠 How It Works

- Reads URLs from `urls.txt`
- Deduplicates them using a `map`
- Spawns a goroutine per URL up to a maximum concurrency
- Each request has a timeout using `context.WithTimeout`
- Parses HTML using `golang.org/x/net/html` and extracts the `<title>`
- Collects results and prints them after all are done

---

## 🔧 Dependencies

- [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html) - for HTML parsing

Install using:

```bash
go get golang.org/x/net/html
```

---

## 📄 License

MIT License - do whatever you want, just don't forget to have fun 🐹

---

Let me know if you'd like this crawler to:
- Save output to a CSV or JSON file
- Crawl internal links recursively
- Add retry logic or proxy support

Happy crawling! 🕸️
