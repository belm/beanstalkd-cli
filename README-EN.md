# Beanstalkd CLI

[🇨🇳 Chinese](README.md) | [🇺🇸 English](README-EN.md)

An all-in-one toolkit for Beanstalkd featuring a powerful command-line interface, a modern web dashboard, and a comprehensive automated test suite.

## ✨ Highlights

- 🎯 **Full coverage** – every Beanstalkd protocol command is supported
- 🎨 **Two interfaces** – CLI and web UI work seamlessly together
- 🧪 **Battle-tested** – 31 automated test cases with 100 % pass rate
- 📚 **Rich documentation** – bilingual guides with detailed workflows
- 🚀 **Ready to run** – build once, use everywhere

## 📦 Quick Start

### Install the CLI

#### Recommended: `go install`

```bash
# Install the CLI
go install github.com/belm/beanstalkd-cli@latest

# Install the web service (optional)
go install github.com/belm/beanstalkd-cli/beanstalkd-web@latest
```

After installation:

- All binaries are placed in `$(go env GOPATH)/bin`; ensure this directory is in your `PATH`
- The CLI binary is named `beanstalkd-cli`
- The web service binary defaults to `beanstalkd-web`; static assets are embedded, so no manual copying is required

#### Build from source

```bash
# Clone and enter the project
git clone <repository-url>
cd beanstalkd-cli

# Download dependencies
go mod tidy

# Build the binary
go build -o beanstalkd-cli

# Or run directly
go run main.go
```

### Launch the Web UI

#### Using the installed binary
```bash
# Default: connect to 127.0.0.1:11300 and listen on 8080
beanstalkd-web

# Specify Beanstalkd host and web port
beanstalkd-web -beanstalkd 192.168.1.100:11300 -port 9090
```

#### Using source/scripts
```bash
cd beanstalkd-web
./start.sh
# or
go run server.go
```

Visit **http://localhost:8080** once the server is up.

### Run the Tests

```bash
make test         # run all tests
make coverage     # generate coverage report
```

## 🛠️ CLI Overview

- Colorful output and table formatting for improved readability
- Global flags: `-H / --host`, `-p / --port`, `-t / --tube`
- Command families: job lifecycle (put, reserve, delete, release, bury, kick, touch),
  inspection (peek*, stats*), and tube management (list/use/watch/ignore)
- Supports priority, delay, and TTR tuning out of the box

See in-depth examples in the Chinese README or run `./beanstalkd-cli <command> --help`.

## 🌐 Web UI Dashboard

- Tailwind CSS based responsive layout with live statistics every 10 seconds
- Tabs for tube overview, job metrics, operations center, and server stats
- Operations center consolidates insert, reserve, delete, and kick actions
- Configurable via command-line flags or environment variables (`BEANSTALKD_HOST`, `WEB_PORT`)

### Configuration Priority
`CLI flags` ➜ `Environment variables` ➜ `Defaults (127.0.0.1:11300, 8080)`

### API Endpoints
- `GET /api/stats`
- `GET /api/tubes`
- `GET /api/tubes/{name}/stats`
- `POST /api/put`
- `POST /api/reserve`
- `POST /api/delete`
- `POST /api/kick`

Screenshots are available in the main README’s “🖼️ 界面预览” section.

## 🧪 Test Suite

- 7 test files covering connection handling, job lifecycle, peeking, tube management, statistics, integration flows, and benchmarking
- `make test-verbose`, `make bench`, and `make coverage` ready for CI pipelines
- Tests are prepared for environments where Beanstalkd runs on `127.0.0.1:11300`

## 📁 Project Structure

```
beanstalkd-cli/
├── cmd/            # CLI command implementations
├── tests/          # Automated tests and helpers
├── beanstalkd-web/ # Web dashboard (Go server + static assets)
├── main.go         # CLI entrypoint
├── Makefile        # Convenience commands
└── README*.md      # Documentation (bilingual)
```

## 🚀 Deployment Tips

- **Development** – run the CLI locally and start the web UI with `./start.sh`
- **Testing** – use environment variables to point at staging Beanstalkd instances
- **Production** – run the CLI with explicit host/port flags and host the web UI behind HTTPS or a reverse proxy
- **Docker** – build the binary in one stage or run the web server with `BEANSTALKD_HOST` injected at runtime

## 🔗 Useful Links

- [Beanstalkd Official Site](https://beanstalkd.github.io/)
- [Protocol Specification](https://github.com/beanstalkd/beanstalkd/blob/master/doc/protocol.txt)
- [go-beanstalk Client](https://github.com/beanstalkd/go-beanstalk)

## 🤝 Contributions

Issues and pull requests are welcome! Please make sure tests pass before submitting.

## 📄 License

MIT License
