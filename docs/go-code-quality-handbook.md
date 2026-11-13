# Go Code Quality, Testing, Debugging & Security Handbook

**A practical toolkit for Go projects on Linux (Arch Linux focused)**  
**Research checked:** 8 November 2026  
**Primary goal:** catch formatting problems, compile errors, likely bugs, unsafe patterns, concurrency defects, dependency vulnerabilities, regressions, and performance problems before they reach users.

---

## 1. The important idea: no single tool can prove code is correct

A code-quality process is a collection of different checks. Each tool sees a different class of problem. The Go ecosystem contains many niche tools, and listing every one would create an unmaintainable catalogue; this handbook covers the main, actively useful tools and the specialist tools that add a distinct capability:

| Layer | Main question | Typical tools |
|---|---|---|
| Formatting | Is the code consistently formatted? | `gofmt`, `goimports`, `gofumpt` |
| Compilation and type checking | Does the code compile for the target? | `go build`, `go test` |
| Static analysis | Can suspicious code be identified without running it? | `go vet`, Staticcheck, `golangci-lint` |
| Behaviour testing | Does the code do what it is supposed to do? | `go test`, integration tests |
| Concurrency testing | Are goroutines accessing shared state safely? | Go race detector (`-race`) |
| Input robustness | Does the program survive unexpected or malformed input? | Go fuzzing |
| Security analysis | Are there risky coding patterns or known vulnerable dependencies? | Gosec, Govulncheck, Trivy, Gitleaks |
| Debugging | What is happening at a specific point in execution? | Delve, debugger in editor |
| Performance analysis | Where is CPU time or memory going? | benchmarks, `pprof`, `trace`, `benchstat` |
| Continuous integration | Do the checks run consistently on every change? | GitHub Actions or another CI system |

A clean scan does **not** prove that an algorithm is logically correct, that a service is secure, or that a multi-machine deployment will work. Tests still need meaningful expected results, and system-level behaviour still needs integration testing.

## 2. What I recommend using day to day

Don't turn every edit into a giant audit. Use a small group of tools continuously, then run the deeper checks at suitable points.

### Essential baseline

1. **`gofmt` / `go fmt`** — standard Go formatting.
2. **`goimports`** — formats code and fixes import lists.
3. **`go test ./...`** — executes all package tests.
4. **`go vet ./...`** — finds suspicious Go constructs.
5. **`golangci-lint run ./...`** — runs a coordinated set of linters.
6. **`govulncheck ./...`** — checks for known vulnerabilities in dependencies and the Go standard library that are relevant to your code.
7. **`go test -race ./...`** — checks for data races in the code paths exercised by tests.

### Add when useful

- **Gosec** for security-focused source-code checks.
- **Gitleaks** for secrets accidentally committed to the repository or present in the working tree.
- **Trivy** for repository, dependency, configuration, and container/root-filesystem scanning.
- **Go fuzzing** for parsers, configuration input, protocol messages, and other untrusted input.
- **Delve** when you need to stop execution and inspect variables, goroutines, and call stacks.
- **`pprof` and `benchstat`** when you have a real performance question.
- **NilAway** when nil-pointer risk is a recurring concern in a sufficiently large codebase.

**Important overlap:** the default `golangci-lint` set already includes `errcheck`, `govet`, `ineffassign`, Staticcheck, and `unused`. You can run `go vet` and `staticcheck` separately while learning or diagnosing a specific report, but you do not need to run every overlapping analyzer independently forever.

---

## 3. Prepare the Go tool environment

### 3.1 Check the toolchain already installed

Run these from a terminal:

```bash
go version
go env GOROOT GOPATH GOBIN GOTOOLCHAIN GOPROXY
```

As of the research date, the latest stable Go release listed by the official Go project is **Go 1.27.3**. That does not mean you must change your project's `go` directive just to use the newest compiler. Check the project's `go.mod` and test compatibility before changing its declared Go version.

Official references:

- Go downloads and releases: <https://go.dev/dl/>
- Go release history: <https://go.dev/doc/devel/release>
- Go command reference: <https://go.dev/doc/cmd>

### 3.2 Make sure Go-installed commands are on PATH

By default, `go install` puts command binaries in `$(go env GOPATH)/bin`, which is usually `$HOME/go/bin`, unless `GOBIN` is set.

Check the path:

```bash
 go env GOPATH
 go env GOBIN
```

For Bash or Zsh, add the following to `~/.bashrc` or `~/.zshrc` if `$HOME/go/bin` is not already on your PATH:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

Reload the shell configuration or open a new terminal, then verify:

```bash
command -v go
command -v gofmt
```

If you prefer a dedicated user-local binary directory, set `GOBIN="$HOME/.local/bin"` in your shell configuration and add that directory to PATH. Pick one location and use it consistently.

### 3.3 Arch Linux package-manager setup

For the Go compiler and basic development utilities:

```bash
sudo pacman -Syu
sudo pacman -S --needed go git
```

For Trivy, which is available in Arch's Extra repository at the time checked:

```bash
sudo pacman -S --needed trivy
```

Arch packages can lag behind upstream releases. For Go developer commands that you want to update independently, `go install module/path@version` is usually the most direct option. Avoid installing the same command from multiple places unless you know which copy your PATH will select:

```bash
command -v golangci-lint
which -a golangci-lint
```

### 3.4 Install the main Go tools

The following commands install the tools into your Go binary directory. `@latest` is convenient for an initial setup; for a repeatable team or CI environment, pin versions instead of allowing tools to change silently.

```bash
# Editor integration and import management
go install golang.org/x/tools/gopls@latest
go install golang.org/x/tools/cmd/goimports@latest

# Formatting alternative (optional; see the formatter section)
go install mvdan.cc/gofumpt@latest

# Linting and static analysis
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
go install honnef.co/go/tools/cmd/staticcheck@latest

# Security
go install golang.org/x/vuln/cmd/govulncheck@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Debugging and test output
go install github.com/go-delve/delve/cmd/dlv@latest
go install gotest.tools/gotestsum@latest

# Benchmark comparison utility
go install golang.org/x/perf/cmd/benchstat@latest
```

Then open a new terminal (or refresh your shell) and verify the binaries:

```bash
gopls version
goimports -h
gofumpt -version
golangci-lint version
staticcheck -version
govulncheck -version
gosec -version
dlv version
gotestsum --version
benchstat -h
```

**Optional separate error checker:**

```bash
go install github.com/kisielk/errcheck@latest
errcheck ./...
```

`errcheck` is already included in the default `golangci-lint` set, so install it separately only when you want to run or study that check directly.

**Optional nil-pointer analyzer:**

```bash
go install go.uber.org/nilaway/cmd/nilaway@latest
nilaway ./...
```

NilAway performs deeper nil-flow analysis. It can be more expensive and may require configuration for a large project; begin with your own packages and review its findings instead of assuming every report is a real defect.

### 3.5 Install Gitleaks

Gitleaks is distributed as official release binaries and through several package managers. On Arch, first check whether your configured repositories or AUR helper provide it. Otherwise:

1. Open the [official Gitleaks releases page](https://github.com/gitleaks/gitleaks/releases).
2. Download the archive matching Linux and your CPU architecture (for your x86-64 PC, choose the Linux x64/amd64 build).
3. Download the release checksum file as well and verify the archive with `sha256sum -c` using the matching published checksum.
4. Extract the archive and place the `gitleaks` executable in a directory on PATH, such as `~/.local/bin`.

Project and usage documentation: <https://github.com/gitleaks/gitleaks>

Verify installation:

```bash
gitleaks version
```

---

## 4. Formatting tools: make the code's appearance consistent

Formatting does not prove the code is correct. It removes avoidable visual differences so that code reviews focus on behaviour and design.

### 4.1 `gofmt` / `go fmt` — required baseline

**Purpose:** apply Go's standard formatting rules. It handles indentation, spacing, and layout in the canonical Go style.

**Why use it:** Go developers expect Go code to be `gofmt`-formatted. It removes subjective formatting debates and makes diffs easier to review.

No installation is needed; it comes with Go.

Format all packages in the current module:

```bash
go fmt ./...
```

Format one file directly:

```bash
gofmt -w path/to/file.go
```

List files that are not formatted:

```bash
find . -type f -name '*.go' -not -path './vendor/*' -print0 \
  | xargs -0 -r gofmt -l
```

If you want the formatter to print changed files and edit them, use `-w`. Without `-w`, `gofmt` normally prints formatted output rather than changing the original file.

### 4.2 `goimports` — formatting plus import cleanup

**Purpose:** formats Go source and adds/removes import statements as needed.

**Why use it:** it saves you from manually adding imports or leaving unused ones after refactoring.

Install:

```bash
go install golang.org/x/tools/cmd/goimports@latest
```

Use it on a file:

```bash
goimports -w internal/config/config.go
```

For a project-wide pass:

```bash
find . -type f -name '*.go' -not -path './vendor/*' -print0 \
  | xargs -0 -r goimports -w
```

### 4.3 `gofumpt` — stricter formatting (optional)

**Purpose:** enforces a stricter subset of Go formatting rules than `gofmt`.

**Why use it:** useful when you want a few additional consistency rules across a team. It is designed as a stricter, backwards-compatible formatter.

Install and run:

```bash
go install mvdan.cc/gofumpt@latest
gofumpt -w .
```

Use `gofmt` as the universal default unless you have a reason to enforce `gofumpt`. Don't configure different formatters to fight over the same files. Choose one canonical formatting policy and configure your editor to use it.

---

## 5. Static analysis: inspect code without running the program

Static analysis is useful for finding suspicious code patterns, likely mistakes, and maintainability issues. It cannot prove runtime behaviour is correct.

### 5.1 `go vet` — official Go checks

**Purpose:** examines Go source for suspicious constructs that compile but may not behave as intended.

Examples include incorrect `Printf` format arguments, suspicious calls, and other patterns covered by the vet analyzers.

**Why use it:** it is official, maintained alongside the Go toolchain, and should be part of a normal Go quality check.

No installation is required:

```bash
go vet ./...
```

To see the available analyzers and flags:

```bash
go tool vet help
```

A successful `go build` does not mean `go vet` will be clean. They check different things.

### 5.2 Staticcheck — high-value bug and code-quality checks

**Purpose:** performs a broad set of static checks for bugs, suspicious code, performance issues, simplifications, and style problems.

**Why use it:** it is a strong complement to `go vet`, with many focused checks designed to produce actionable findings.

Install:

```bash
go install honnef.co/go/tools/cmd/staticcheck@latest
```

Run:

```bash
staticcheck ./...
```

Staticcheck reports categories such as `SA` (static analysis), `S` (simple), `ST` (style), and `QF` (quick-fix opportunities). Read the explanation for each report before changing code; a simplification is not automatically a better design in every context.

### 5.3 `golangci-lint` — the main linting coordinator

**Purpose:** runs multiple Go linters through one command, supports configuration, caching, and CI integration.

**Why use it:** instead of remembering and invoking many individual linters, you get a repeatable central lint command.

Install:

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
```

Check version and available linters:

```bash
golangci-lint version
golangci-lint linters
```

Run against the entire module:

```bash
golangci-lint run ./...
```

The current v2 default set includes these important checks:

| Linter | What it checks |
|---|---|
| `errcheck` | Returned errors that are ignored |
| `govet` | Suspicious Go constructs (the Go vet analyzers) |
| `ineffassign` | Assignments whose values are never used |
| `staticcheck` | Staticcheck's bug, simplification, and style checks |
| `unused` | Unused constants, variables, functions, and types |

Start with the default configuration. It is better to have a manageable, trusted baseline than to enable every available linter and drown in noisy findings.

A small v2 configuration can add a couple of checks that suit service code. Create `.golangci.yml` in the repository root:

```yaml
version: "2"

linters:
  default: standard
  enable:
    - bodyclose
    - errorlint
    - gosec
```

- `bodyclose` checks that HTTP response bodies are appropriately closed.
- `errorlint` checks error-wrapping and error-inspection patterns.
- `gosec` adds security-focused source checks; a separate Gosec installation is not needed if you are happy using its integrated linter.

Run the configured set:

```bash
golangci-lint run ./...
```

To see the available configuration options for the exact installed version:

```bash
golangci-lint help linters
golangci-lint help formatters
```

**Do not enable every linter on day one.** Add checks when their purpose is clear, fix genuine findings, and document any justified exclusions. Avoid blanket exclusions that hide entire folders or categories of errors.

### 5.4 `errcheck` — ignored error returns

**Purpose:** identifies calls whose returned errors are not handled or explicitly discarded.

**Why it matters:** Go uses returned errors extensively. Accidentally ignoring an error from file writes, closes, network operations, or process execution can produce silent failures.

Install separately only if you want to run it directly:

```bash
go install github.com/kisielk/errcheck@latest
errcheck ./...
```

An explicit `_ = someCall()` tells the reader that the return value is being discarded, but it does not make that decision safe. Decide deliberately whether the error can be ignored.

### 5.5 Optional analyzers for specific concerns

| Tool | Install | Run | Useful when |
|---|---|---|---|
| NilAway | `go install go.uber.org/nilaway/cmd/nilaway@latest` | `nilaway ./...` | You want deeper analysis of potential nil-pointer panics. Begin with your own packages; the analyzer can be resource-intensive on large dependency graphs. |
| `gocyclo` | `go install github.com/fzipp/gocyclo/cmd/gocyclo@latest` | `gocyclo -over 15 .` | You want a rough cyclomatic-complexity report. |
| `revive` | `go install github.com/mgechev/revive@latest` | `revive ./...` | You want additional opinionated style and documentation rules. |
| `dupl` | `go install github.com/mibk/dupl@latest` | `dupl .` | You want to identify repeated code blocks. |

These are not mandatory gates. Complexity scores and duplication reports are signals for human review, not automatic proof that a function or package must be rewritten. A short function can still be confusing, and a long function can still be straightforward.

---

## 6. Testing: verify actual behaviour

### 6.1 `go test ./...` — the normal test command

**Purpose:** compiles packages and runs their tests.

**Why use it:** a linter cannot tell you whether your provisioning function creates the right container, whether a configuration is parsed correctly, or whether a network message produces the expected response. Tests verify expected behaviour.

Run all tests in the module:

```bash
go test ./...
```

Run one package:

```bash
go test ./internal/config
```

Run one test by name:

```bash
go test ./internal/config -run '^TestParseConfig$' -v
```

Run tests without reusing cached successful results and randomize test order:

```bash
go test -count=1 -shuffle=on ./...
```

- `-count=1` forces a fresh test execution instead of accepting a cached test result.
- `-shuffle=on` helps reveal tests that accidentally depend on execution order.

### 6.2 Organize Go tests

Go test files end in `_test.go` and normally live beside the code they test.

Example layout:

```text
internal/
  config/
    config.go
    config_test.go
  provision/
    provision.go
    provision_test.go
```

A small example:

```go
package config

import "testing"

func TestValidateRejectsEmptyHostname(t *testing.T) {
    cfg := Config{Hostname: ""}

    err := cfg.Validate()
    if err == nil {
        t.Fatal("Validate() error = nil, want an error for an empty hostname")
    }
}
```

Good tests should state the behaviour being protected, include meaningful input and expected output, and test both normal and failure cases. For validation code, test empty values, malformed values, boundaries, duplicate entries, and unexpected input.

### 6.3 `go test -race` — detect data races

**Purpose:** instruments the program to detect unsynchronized concurrent access to shared memory.

**Why use it:** a program can pass ordinary tests and still contain races that appear only when goroutines overlap. The race detector is especially relevant to network servers, background workers, shared state, and concurrent provisioning tasks.

Run:

```bash
go test -race ./...
```

A stronger routine test command is:

```bash
go test -race -count=1 -shuffle=on ./...
```

The race detector only finds races that occur during the execution being observed. It does not prove that all possible schedules are race-free. It also has noticeable CPU and memory overhead, so use it as a testing mode rather than assuming it must be enabled for every production build.

### 6.4 Coverage — find areas your tests never execute

**Purpose:** measures which statements are exercised by tests.

**Why use it:** low coverage can reveal packages or branches with no tests. Coverage is a map of test reach, not a score of correctness. A test can execute every line while asserting almost nothing.

Generate a coverage profile:

```bash
go test -coverprofile=coverage.out ./...
```

Show function-level coverage:

```bash
go tool cover -func=coverage.out
```

Create an HTML report:

```bash
go tool cover -html=coverage.out -o coverage.html
```

Open `coverage.html` in a browser to see which source lines were executed.

### 6.5 Fuzzing — automatically explore weird inputs

**Purpose:** repeatedly feeds generated inputs into a function to find crashes, panics, and unexpected behaviour. Go has native fuzzing in its standard toolchain; no third-party fuzzer is required to get started.

**Why it matters:** parsers and protocol handlers often fail on inputs that a human tester would not think to write.

Fuzzing is particularly useful for:

- JSON or YAML configuration parsing.
- Service discovery and protocol messages.
- Domain names, addresses, ports, and path validation.
- Data received over TCP or HTTP.
- Functions that split, decode, normalize, or validate user-controlled strings.

Example fuzz test:

```go
package config

import "testing"

func FuzzParseConfig(f *testing.F) {
    f.Add([]byte(`{"hostname":"node-1"}`))
    f.Add([]byte(`{}`))
    f.Add([]byte(`not-json`))

    f.Fuzz(func(t *testing.T, input []byte) {
        // The parser should return a result or an error; it must not panic.
        _, _ = Parse(input)
    })
}
```

Run a fuzz target for 30 seconds:

```bash
go test ./internal/config -fuzz '^FuzzParseConfig$' -fuzztime=30s
```

Run for a longer session:

```bash
go test ./internal/config -fuzz '^FuzzParseConfig$' -fuzztime=5m
```

When Go finds a failing input, it saves it to the fuzz corpus so the failure can be reproduced. Do not discard a discovered regression input without understanding and fixing the bug.

A fuzz target should be deterministic, fast, and isolated from persistent external side effects. Do not let a fuzzer create/delete real LXC containers or change the host firewall. Fuzz the parser or planning logic independently; test real provisioning in controlled integration tests.

### 6.6 Benchmarks — measure performance rather than guessing

**Purpose:** measures the time and allocations of a function or package under a repeatable workload.

**Why use it:** it helps establish whether a proposed optimization is actually faster or merely more complicated.

Run package benchmarks:

```bash
go test -bench=. -benchmem ./...
```

Run one benchmark repeatedly:

```bash
go test -bench '^BenchmarkParseConfig$' -benchmem -count=5 ./internal/config
```

Compare two benchmark result files with `benchstat`:

```bash
go test -bench=. -benchmem -count=10 ./internal/config > before.txt
# Make the change, then run the same benchmark:
go test -bench=. -benchmem -count=10 ./internal/config > after.txt
benchstat before.txt after.txt
```

Benchmarks should use realistic input and comparable machine conditions. Do not optimize based on one noisy run.

### 6.7 `gotestsum` — more readable test output

**Purpose:** runs tests via `go test -json` and presents a human-friendly summary. It can also emit JUnit XML for CI systems.

**Why use it:** helpful when the project has many packages and you want a clear summary of failures and slow tests.

Install and run:

```bash
go install gotest.tools/gotestsum@latest
gotestsum -- ./...
```

Run with race detection:

```bash
gotestsum -- -race ./...
```

Generate a JUnit report:

```bash
gotestsum --junitfile test-results.xml -- ./...
```

It is optional. `go test` remains the actual test engine; `gotestsum` improves reporting and CI integration.

---

## 7. Security tools: inspect application code, dependencies, secrets, and deployment artefacts

These tools inspect different layers. Running one does not make the others unnecessary.

### 7.1 Govulncheck — known vulnerabilities relevant to your Go code

**Purpose:** checks dependencies and the Go standard library against the Go vulnerability database. It analyses whether vulnerable functions are reachable from your code, reducing noise compared with reporting every vulnerable module version without context.

**Why use it:** an old dependency is not automatically exploitable in every program; call-path information helps prioritize the issues that actually affect your application.

Install and run:

```bash
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...
```

Run this regularly and before releases. Review each finding, update the affected dependency or toolchain where possible, then rerun the scan.

### 7.2 Gosec — security-focused source-code analysis

**Purpose:** scans Go source and program representations for common security weaknesses. Its checks cover issues such as unsafe use, hard-coded credentials, weak cryptography, unsafe file permissions, path handling, and injection-related patterns.

**Why use it:** it can flag risky source patterns even if the code compiles and the dependencies have no known CVEs.

Install and run:

```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec ./...
```

Write a machine-readable report:

```bash
gosec -fmt json -out gosec-report.json ./...
```

Or produce SARIF for compatible code-scanning systems:

```bash
gosec -fmt sarif -out gosec-report.sarif ./...
```

Gosec findings need human review. A scanner cannot always understand why a particular path, command, permission, or cryptographic choice is safe in your context.

If using the `gosec` linter inside `golangci-lint`, avoid running both configurations in a way that creates duplicate reports. Standalone Gosec remains useful when you want its dedicated report formats or separate security job.

### 7.3 Gitleaks — find credentials and secrets in Git

**Purpose:** scans Git history or a directory for likely secrets such as access tokens, API keys, and passwords.

**Why use it:** accidentally committed credentials can remain in repository history even after the line is deleted. Removing the visible line does not invalidate a leaked credential.

Examples using the current command style:

```bash
# Scan Git history
gitleaks git . --redact

# Scan files in a directory
gitleaks dir . --redact
```

Before sharing a report, use redaction options so the secret itself is not copied into logs or chat. If a real credential is exposed, revoke or rotate it immediately; deleting the file is not enough.

### 7.4 Trivy — repository and container supply-chain scanning

**Purpose:** scans source repositories and deployment artefacts for known vulnerabilities, misconfigurations, secrets, and license information. It can also inspect container images and root filesystems.

**Why use it:** useful when a Go application is only one part of a larger system containing OS packages, configuration files, containers, and deployment assets.

Install on Arch (package available in Extra at the time checked):

```bash
sudo pacman -S --needed trivy
```

Scan the repository:

```bash
trivy repo .
```

Scan a filesystem tree:

```bash
trivy fs .
```

Scan a container image:

```bash
trivy image IMAGE_NAME:TAG
```

Scan a container/root filesystem directory:

```bash
trivy rootfs /path/to/rootfs
```

For an LXC-based project, Trivy can be useful for inspecting the Debian root filesystem or deployment files. It does not validate that the container networking, firewall, DNS, NGINX, or provisioning behaviour is correct; those need runtime and integration tests.

### 7.5 Which security tool should you run?

| Concern | Best starting tool |
|---|---|
| A vulnerable Go module or standard-library function | `govulncheck` |
| Risky Go source-code pattern | Gosec or the `gosec` linter |
| A token/password committed to Git | Gitleaks |
| Dependencies, secrets, configuration, OS/container content | Trivy |

A sensible baseline is **Govulncheck + Gosec + a secret scanner**. Add Trivy when you begin scanning the LXC image/root filesystem and deployment configuration as well.

---

## 8. Debugging: inspect what the program is doing

### 8.1 `gopls` — editor intelligence and live diagnostics

**Purpose:** the official Go language server. It provides completion, navigation, diagnostics, refactoring, references, and other editor features through LSP.

**Why use it:** it catches many problems while you are editing instead of waiting until you manually run a command.

Install:

```bash
go install golang.org/x/tools/gopls@latest
```

Usually the editor starts and manages `gopls` automatically once its Go/LSP integration is enabled. You normally do not run it manually. Check the Go extension or editor's language-server settings if diagnostics are missing.

### 8.2 Delve (`dlv`) — source-level debugger

**Purpose:** lets you pause a Go program, set breakpoints, inspect variables, move through execution, and inspect goroutines and call stacks.

**Why use it:** logs tell you what you chose to print. A debugger lets you inspect program state at the exact point where something goes wrong.

Install:

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

Check:

```bash
dlv version
```

Debug a main package (replace the path with your command package):

```bash
dlv debug ./cmd/your-program
```

Debug tests in a package:

```bash
dlv test ./internal/provision
```

Useful commands inside Delve:

| Command | Meaning |
|---|---|
| `break main.main` | Set a breakpoint at a function |
| `continue` | Continue execution until the next breakpoint |
| `next` | Execute the next source line without stepping into a call |
| `step` | Step into a function call |
| `print variableName` | Inspect a variable |
| `locals` | Display local variables |
| `goroutines` | List goroutines |
| `stack` | Show the current call stack |
| `exit` | Leave the debugger |

For day-to-day work, configure your editor to use Delve so that breakpoints and variable inspection are integrated into the UI.

---

## 9. Performance diagnostics: measure CPU, memory, and blocking

Do not optimize code just because it looks complicated. First reproduce a performance issue, measure it, find the hot path, and then compare before and after.

### 9.1 `pprof` — CPU and memory profiling

**Purpose:** collects and displays profiles that identify where the application spends CPU time or allocates memory.

**Why use it:** it helps identify actual hot functions and memory-heavy operations rather than guessing from source code.

Collect profiles from tests (run against a focused package):

```bash
go test -cpuprofile=cpu.prof -memprofile=mem.prof ./internal/config
```

Inspect the CPU profile:

```bash
go tool pprof cpu.prof
```

Inside the pprof prompt, useful commands include:

```text
top
list FunctionName
web
```

You can also render an interactive web view (where supported):

```bash
go tool pprof -http=:8080 cpu.prof
```

For a running HTTP server, Go provides `net/http/pprof`. Only expose profiling endpoints on a trusted local or protected interface; do not accidentally publish them as a public service endpoint.

### 9.2 `go tool trace` — execution timeline

**Purpose:** records a runtime trace showing scheduling, goroutines, blocking, network activity, and other runtime events.

**Why use it:** it is useful when a server appears slow because goroutines are blocked, work is queued, or scheduling is not behaving as expected.

Collect a trace from tests:

```bash
go test -trace=trace.out ./internal/provision
```

Open it:

```bash
go tool trace trace.out
```

### 9.3 `benchstat` — compare benchmark results

**Purpose:** statistically summarizes benchmark outputs and makes before/after performance comparisons easier to interpret.

**Why use it:** it reduces the temptation to trust one noisy benchmark run.

Install:

```bash
go install golang.org/x/perf/cmd/benchstat@latest
```

Then compare saved result files:

```bash
benchstat before.txt after.txt
```

---

## 10. Dependency and module hygiene

Go modules already provide several useful integrity and consistency checks. Use them rather than adding third-party tools for every small task.

### 10.1 `go mod tidy`

**Purpose:** adds missing module requirements and removes requirements that are no longer needed by the module's packages.

Run:

```bash
go mod tidy
```

Review the changes to `go.mod` and `go.sum` and commit them when appropriate. Do not run it blindly and ignore a surprising dependency change.

### 10.2 `go mod verify`

**Purpose:** checks that downloaded module contents have not been modified since they were downloaded and recorded.

Run:

```bash
go mod verify
```

### 10.3 Other useful module commands

```bash
# Download module dependencies
go mod download

# List dependency graph
go mod graph

# Explain why a module is present
go mod why example.com/module

# Show available module updates
go list -m -u all
```

Do not update every dependency automatically just because updates exist. Review release notes, compatibility, and security advisories, then test the result.

### 10.4 Tool versions and reproducible development

`go install tool@latest` is convenient for installing your own machine's tools, but the meaning of `latest` changes over time. For consistent local and CI results, pin the versions used by the project.

Go 1.24 and newer support tool dependencies through `go.mod`. For example, a project can record a tool dependency with:

```bash
go get -tool golang.org/x/tools/cmd/goimports
```

Then invoke the recorded tool with:

```bash
go tool goimports -w path/to/file.go
```

Check the installed tool's documentation and the module's `go.mod` for its compatibility requirements. The benefit is that the project records its developer-tool dependency instead of relying entirely on whichever global version happens to be installed.

---

## 11. A practical quality command set for your project

Run these from the repository root. They form a good progression from fast checks to deeper checks.

### Fast local check

```bash
go fmt ./...
go test ./...
go vet ./...
golangci-lint run ./...
```

### Before a commit or pull request

```bash
go test -count=1 -shuffle=on ./...
go test -race -count=1 ./...
golangci-lint run ./...
govulncheck ./...
gosec ./...
gitleaks git .
```

### Before a release

```bash
go mod tidy
go mod verify
go build ./...
go test -count=1 -shuffle=on ./...
go test -race -count=1 ./...
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
golangci-lint run ./...
govulncheck ./...
gosec ./...
gitleaks git .
trivy repo .
```

The release checklist should also include target-platform builds, integration tests, container/root-filesystem scans, and a manual review of security-sensitive changes.

---

## 12. A Makefile to make the checks easy to remember

Create a `Makefile` at the repository root. The commands below assume the tools in the baseline section are installed.

```makefile
.PHONY: fmt fmt-check build test test-race vet lint vuln security quality coverage

fmt:
	go fmt ./...

fmt-check:
	@test -z "$$(find . -type f -name '*.go' -not -path './vendor/*' -print0 | xargs -0 -r gofmt -l)" || \
	  (echo "Some Go files are not gofmt-formatted. Run 'make fmt'."; exit 1)

build:
	go build ./...

test:
	go test -count=1 -shuffle=on ./...

test-race:
	go test -race -count=1 -shuffle=on ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

security:
	gosec ./...
	gitleaks git .

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

quality: fmt-check build vet lint test test-race vuln
```

Use it like this:

```bash
make fmt       # automatically format
make test      # regular tests
make test-race # race detector
make lint      # static linting
make security  # source and secret scans
make quality   # combined quality gate
make coverage  # produce coverage.html
```

The `quality` target intentionally does not run the slower whole-repository Trivy scan or fuzzing. Run those at a release/checkpoint cadence. You can add them to a separate `release-check` target when the repository and CI pipeline are ready.

---

## 13. Continuous integration: repeat checks on every change

A tool installed only on a developer's laptop can be forgotten. CI makes the baseline repeatable and catches changes that were not tested locally.

A simple GitHub Actions workflow can run the core checks. Create `.github/workflows/go-quality.yml`:

```yaml
name: Go quality

on:
  push:
  pull_request:

permissions:
  contents: read

jobs:
  quality:
    runs-on: ubuntu-latest
    steps:
      - name: Check out source
        uses: actions/checkout@v7

      - name: Set up Go
        uses: actions/setup-go@v7
        with:
          go-version-file: go.mod
          cache: true

      - name: Format check
        run: |
          test -z "$(find . -type f -name '*.go' -not -path './vendor/*' -print0 | xargs -0 -r gofmt -l)

      - name: Build
        run: go build ./...

      - name: Test
        run: go test -count=1 -shuffle=on ./...

      - name: Race detector
        run: go test -race -count=1 ./...

      - name: Vet
        run: go vet ./...

      - name: Install golangci-lint
        run: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

      - name: Lint
        run: golangci-lint run ./...

      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest

      - name: Vulnerability check
        run: govulncheck ./...
```

For a stable production workflow, pin external action versions and Go tool versions to reviewed releases. The example uses current major action tags and `@latest` for readability; pinning makes CI more reproducible and avoids surprise tool changes. Add Gosec, Gitleaks, and Trivy as separate security steps when you decide how to handle their reports and failure thresholds.

---

## 14. Testing advice specific to infrastructure and provisioning code

For infrastructure software, a program can pass every static scan and still damage a host if its provisioning logic is wrong. Build tests around the boundaries between your Go code and the operating system.

### 14.1 Keep planning logic separate from host side effects

For operations such as creating an LXC container, changing firewall rules, writing system files, or starting services, separate:

1. **Decision/planning logic** — computes what should happen from the current state and desired configuration.
2. **Execution layer** — actually invokes `lxc`, firewall tools, filesystem operations, and service managers.
3. **Verification layer** — checks whether the operation succeeded and whether the resulting state is correct.

The decision layer should be testable without root privileges or a real container. Give it inputs and verify the commands/actions it intends to perform.

### 14.2 Test external commands safely

Use `exec.CommandContext` with an executable and separate argument values rather than building a shell command by concatenating strings. This makes argument boundaries explicit and reduces shell-injection risk.

Tests can replace the command runner with a fake implementation that records:

- Which executable was requested.
- Which arguments were passed.
- Which environment variables were set.
- What exit code, stdout, and stderr were returned.
- Whether the code handled failures, timeouts, and cancellation correctly.

Then use a small number of controlled integration tests to verify the real LXC/firewall/service interactions inside a disposable VM or test host.

### 14.3 Test failure paths, not just successful provisioning

For each important operation, test:

- The command succeeds.
- The executable is missing.
- The command returns a non-zero exit code.
- The command times out or is cancelled.
- The command returns unexpected output.
- The target file already exists or has invalid permissions.
- The container or service is already in the desired state.
- A later stage fails after an earlier stage succeeded.
- A retry does not corrupt or duplicate the previous state.

For a bootstrap pipeline, verify the difference between a warning-level failure that can be recovered from and a fatal failure that must stop the pipeline. Assert that errors propagate to the correct caller and that cleanup or rollback is attempted when necessary.

### 14.4 Treat network input as untrusted

For DNS/control-plane messages, configuration, and service-to-service communication, fuzz parsers and validators independently of the networking stack. Also write integration tests for framing, timeouts, malformed messages, invalid identity, authorization, and partial reads/writes. TCP is a byte stream; a single `Read` is not guaranteed to return one complete application message.

### 14.5 Test on the target platform

If the supported deployment target is Linux amd64, ensure CI at least builds and tests on Linux. Cross-compiling can catch platform-specific compilation issues, but it does not replace running integration tests on the actual target environment.

---

## 15. What each tool cannot do

| Tool | It is good at | It does not guarantee |
|---|---|---|
| `gofmt` / `goimports` | Consistent source formatting and imports | Correct logic, safe architecture |
| `go build` | Compilation and type checking | Runtime correctness |
| `go vet` | Suspicious constructs covered by analyzers | Detection of every bug |
| Staticcheck | Many common static defects and simplifications | Complete program understanding |
| `golangci-lint` | Coordinated linting | That enabled linters cover every risk |
| `go test` | Behaviour covered by assertions | Correctness of untested paths |
| `-race` | Races that actually occur during instrumented execution | Absence of all races |
| Fuzzing | Bugs exposed by generated input | All possible inputs or business correctness |
| Coverage | Which code was executed | Whether tests asserted the right thing |
| Govulncheck | Known vulnerabilities reachable through analysed call paths | Unknown vulnerabilities or unsafe application logic |
| Gosec | Common source-level security patterns | A full security audit or threat model |
| Gitleaks | Likely secrets in files and Git history | Every possible secret format or leaked credential outside the repo |
| Trivy | Known issues in supported repository/container artefacts | Correct runtime configuration or secure network design |
| Delve | Runtime state and control flow while debugging | Automatic bug discovery |
| pprof / trace | Performance and runtime behaviour under observed workload | Performance under every production workload |

The strongest process combines these tools with careful design, readable code, targeted tests, code review, and realistic integration testing.

---

## 16. A sensible adoption order

Don't install and configure everything in one giant leap. Build the workflow in layers:

**Stage 1 — start immediately**

- Verify Go and PATH.
- Use `go fmt ./...`.
- Add meaningful `_test.go` tests.
- Run `go test ./...` and `go vet ./...`.

**Stage 2 — consistent static analysis**

- Install `golangci-lint`.
- Run its default set first.
- Fix real findings and add a small `.golangci.yml` only when necessary.
- Add `goimports` or choose `gofumpt` as the formatter policy.

**Stage 3 — reliability and security**

- Add `go test -race ./...`.
- Add `govulncheck ./...`.
- Add Gosec or its integrated linter.
- Add Gitleaks for repository history and staged changes.

**Stage 4 — deeper testing and deployment assurance**

- Add fuzz tests for parsers and message validation.
- Add integration tests in a disposable VM/container environment.
- Add Trivy scans for the repository and LXC root filesystem.
- Add CI so the checks run on every push and pull request.

**Stage 5 — investigate actual performance problems**

- Write benchmarks for the hot path.
- Use `benchstat` for comparisons.
- Use `pprof` and `go tool trace` to locate the actual bottleneck.

For a growing Go infrastructure project, this is a strong quality system without excessive ceremony. Start with the baseline, keep the checks passing, and add specialist tools when the project reaches the point where they provide measurable value.

---

## 17. Official documentation and project links

These are the primary references used for the tool descriptions and commands. Tool interfaces and installation instructions can change, so consult their official docs when updating the toolchain.

### Go toolchain and testing

- Go downloads: <https://go.dev/dl/>
- Go release history: <https://go.dev/doc/devel/release>
- Go command documentation: <https://go.dev/doc/cmd>
- Go testing package: <https://pkg.go.dev/testing>
- Go fuzzing: <https://go.dev/doc/security/fuzz/>
- Go race detector: <https://go.dev/doc/articles/race_detector>
- Go diagnostics, profiling, tracing, and debugging: <https://go.dev/doc/diagnostics>
- Go pprof HTTP handlers: <https://pkg.go.dev/net/http/pprof>
- Go module dependency management: <https://go.dev/doc/modules/managing-dependencies>
- Go vulnerability management: <https://go.dev/doc/security/vuln/>
- Govulncheck tutorial: <https://go.dev/doc/tutorial/govulncheck>

### Static analysis and formatting

- Go tools (`gopls`, `goimports`, and related tools): <https://go.googlesource.com/tools>
- gopls: <https://go.dev/gopls/>
- Staticcheck: <https://staticcheck.dev/docs/>
- golangci-lint: <https://golangci-lint.run/docs/welcome/quick-start/>
- golangci-lint formatters: <https://golangci-lint.run/docs/formatters/>
- gofumpt: <https://github.com/mvdan/gofumpt>
- errcheck: <https://github.com/kisielk/errcheck>
- NilAway: <https://github.com/uber-go/nilaway>
- gocyclo: <https://github.com/fzipp/gocyclo>
- revive: <https://github.com/mgechev/revive>
- dupl: <https://github.com/mibk/dupl>

### Security, debugging, and CI

- Gosec: <https://github.com/securego/gosec>
- Gitleaks: <https://github.com/gitleaks/gitleaks>
- Trivy repository scanning: <https://www.trivy.dev/docs/latest/guide/target/repository/>
- Trivy filesystem scanning: <https://trivy.dev/docs/latest/guide/target/filesystem/>
- Delve installation and docs: <https://github.com/go-delve/delve>
- gotestsum: <https://github.com/gotestyourself/gotestsum>
- GitHub Actions setup-go: <https://github.com/actions/setup-go>
- Arch Linux Go package: <https://archlinux.org/packages/extra/x86_64/go/>
- Arch Linux Trivy package: <https://archlinux.org/packages/extra/x86_64/trivy/>

---

**End of handbook.** The recommended starting point is sections 3, 4, 5.3, 6.1, and 7.1. The remaining sections are your reference library as the codebase grows.
