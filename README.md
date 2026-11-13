# Yathahi

**Unity through interoperability.**

Yathahi is an infrastructure management project focused on making self-hosted infrastructure easier to deploy, operate, and maintain.

It aims to bring infrastructure services and machines together under a unified, modular system — reducing operational complexity while keeping the underlying technologies interoperable.

The goal is simple: infrastructure should work together, adapt to its environment, and require less manual intervention to manage.

## Vision

Modern infrastructure is often fragmented across individual tools, services, machines, and configurations. Managing these components independently creates unnecessary complexity, especially when building and maintaining self-hosted environments.

Yathahi explores a different approach: a unified infrastructure management system where independent components work together through clearly defined interfaces.

Rather than replacing every existing technology, Yathahi aims to coordinate and manage them.

## Architecture

Yathahi is built around three primary architectural components:

### Entry

The entry point into the Yathahi environment. It establishes the initial interaction with the system and connects the entry process to the underlying infrastructure.

### Core

The central management component running on a host. Core coordinates the environment and provides the foundation through which infrastructure components are managed.

### Engine

Independent components responsible for carrying out specific infrastructure operations. Engines extend the system's capabilities while keeping individual responsibilities modular and separated.

Together, Entry, Core, and Engine form the architectural foundation of Yathahi.

## Design Principles

- **Open** — Designed with openness and transparency in mind.
- **Interoperable** — Built to work alongside existing technologies rather than isolate itself from them.
- **Modular** — Independent components with clearly defined responsibilities.
- **Self-managing** — Striving to reduce repetitive manual administration.
- **Secure** — Security considered throughout the architecture and communication between components.

## Goals

Yathahi aims to provide a foundation for:

- Simplified infrastructure deployment and management.
- Unified management of services across machines.
- Modular infrastructure components.
- Automated configuration and environment management.
- Reliable communication between system components.
- Infrastructure that can evolve without depending on a single technology.

## Development Status

Yathahi is under active development.

Its architecture and implementation are evolving as the foundational components are developed and tested. Interfaces, internal structures, and supported capabilities may change during this stage.

Expect incomplete features and breaking changes until the project reaches a more mature release.

## Project Structure

## Project Structure

The project is written in Go and organized around its core architectural components and shared internal libraries:

```text
Yathahi/
├── go.mod
└── source/
    ├── architecture/    # Entry, Core, and Engine components
    └── library/         # Shared UDS communication, package management, and data structures
```

*Note: Because Yathahi is under active development, internal file paths and directory layouts are highly volatile and may change frequently.*


## Contributing

Contributions, discussions, suggestions, and technical feedback are welcome.

The long-term goal is to build an infrastructure ecosystem through collaboration, interoperability, and well-defined technical foundations.

* Read our [Contribution Guidelines](CONTRIBUTING.md) to learn how to fork the repository, sign our one-click Contributor Assignment Agreement, and submit a Pull Request.
* View all project contributors in [CONTRIBUTORS.md](CONTRIBUTORS.md).

## License

Yathahi is licensed under the **Apache License 2.0**. See the [LICENSE](LICENSE) and [NOTICE](NOTICE) files for details.

Copyright 2026 **J Aromal**.

---

<p align="center">
  <strong>Yathahi</strong>
  <br/>
  Unity through interoperability.
</p>
