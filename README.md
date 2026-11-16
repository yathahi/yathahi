# Yathahi

<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="assets/logos/yathahi-text-logo-white.inkscape.svg">
    <source media="(prefers-color-scheme: light)" srcset="assets/logos/yathahi-text-logo.inkscape.svg">
    <img src="assets/logo-dark.svg" alt="Yathahi Logo" width="180">
  </picture>
</p>

Yathahi is an infrastructure creation, management, orchestration, security, and maintenance project.

It aims to bring machines, services, and infrastructure components together into one unified environment, making it easier to create, operate, coordinate, and maintain infrastructure without requiring extensive manual administration.

Yathahi is designed to be flexible across every scale — from a single computer running multiple services or a small two-device home lab to large, geographically distributed infrastructures.

The goal is simple: infrastructure should work together, adapt to its environment, and require less manual intervention to manage.

## Vision

Infrastructure doesn't have to be large or complicated to benefit from proper management. A personal computer hosting several services, a home laboratory, a small office, or a large distributed environment all have infrastructure needs of their own.

Yathahi explores a unified approach to creating and managing these environments.

Rather than requiring users to manage every service and machine independently, Yathahi aims to coordinate them through a common ecosystem while allowing the underlying technologies and services to remain independent and interoperable.

Its purpose isn't to replace every existing technology, but to make the infrastructure built around those technologies easier to establish, understand, secure, and maintain.

## Architecture

Yathahi is built around four primary architectural components:

### Entry

The initial entry point into Yathahi. Entry is the executable users directly install. It establishes Core on a device when necessary and provides the initial user interaction for collecting information and submitting infrastructure and service requests.

### Core

The central management component on an individual device. Core manages the device's Yathahi environment, including container lifecycles, Engine installation and assignment, authentication, records, and device-level coordination and security.

### Engine

Specialised components responsible for operating particular services and workloads inside their assigned containers. Each Engine understands its own service and handles its configuration, operation, and maintenance without taking responsibility for the wider device or infrastructure.

### Administrator

The infrastructure-wide authority responsible for coordinating participating devices and their Cores. Administrator governs infrastructure-level operations, communication, and network security. It also provides the graphical interface for viewing, configuring, and managing the wider infrastructure.

Together, Entry, Core, Engine, and Administrator form the architectural foundation of Yathahi.

For a detailed explanation, see the [Architecture Document](docs/ARCHITECTURE_v0.0.1-alpha.md).

## Design Principles

- **Open** — Designed with openness and transparency in mind.
- **Interoperable** — Built to work alongside existing technologies rather than isolate itself from them.
- **Modular** — Independent components with clearly defined responsibilities.
- **Flexible** — Intended to work across different environments, configurations, and scales.
- **Self-managing** — Striving to reduce repetitive manual administration and operational complexity.
- **Secure** — Security considered throughout the architecture, from individual workloads to the wider infrastructure.
- **Resilient** — Aiming to preserve the independent operation of underlying infrastructure and limit the impact of failures.

## Goals

Yathahi aims to provide a foundation for:

- Infrastructure creation and automated deployment.
- Hosting and managing services on individual machines or across multiple devices.
- Coordinating infrastructure across different networks and geographical locations.
- Managing specialised services through independent Engines.
- Automated configuration, monitoring, and maintenance.
- Infrastructure-wide security and workload isolation.
- Simplified infrastructure management without requiring users to understand every underlying technology.
- Infrastructure that can evolve without depending on a single technology or service.

## Documentation

The `docs/` directory contains the foundational documents describing Yathahi's current direction.

| Document | Description |
|---|---|
| [Project Vision](docs/PROJECT_VISION_v0.0.1-alpha.md) | The purpose, ambitions, and guiding principles of Yathahi. |
| [Architecture](docs/ARCHITECTURE_v0.0.1-alpha.md) | The architectural components, responsibilities, and relationships. |
| [Specification](docs/SPECIFICATION_v0.0.1-alpha.md) | The intended system behaviour, requirements, and constraints. |
| [Yathahi Engineering Standards (YES)](docs/YES.v0.0.1-alpha.md) | Defines the engineering standards, conventions, and development guidelines followed by Yathahi. |

**Documentation Notice:** Yathahi is in its early alpha stage. Its documentation is highly volatile. Architectural decisions, requirements, and project direction may change, be revised, or be replaced as development progresses. This flexibility is intentional. Refer to the [documentation notice](docs/README.md) for further details.

## Development Status

Yathahi is under active development.

Its architecture, specifications, and implementation are evolving as the foundational components are developed, tested, and refined.

Expect incomplete features, evolving interfaces, architectural changes, and breaking changes until the project reaches a more mature release.

Nothing in the current alpha documentation should be treated as an irreversible architectural or technical commitment.

## Project Structure

The project is written in Go and organised around its architectural components, shared libraries, executable binaries, and service definitions.

```text
Yathahi/
├── bin/                    # Built executables
├── docs/                   # Project documentation
├── source/
│   ├── architecture/       # Entry, Core, and Engine source
│   └── library/            # Shared libraries and utilities
├── systemd/                # Service definitions
├── go.mod
└── README.md
```

*Note: Internal source files, directory layouts, and implementation details are subject to change during development.*

## Contributing

Contributions, discussions, suggestions, and technical feedback are welcome.

The long-term goal is to build an infrastructure ecosystem through collaboration, interoperability, and well-defined technical foundations.

- Read our [Contribution Guidelines](CONTRIBUTING.md) to learn how to fork the repository and submit a Pull Request (no CLA required).
- View all project contributors in [CONTRIBUTORS.md](CONTRIBUTORS.md).
- Read our [Code of Conduct](CODE_OF_CONDUCT.md) before participating.

## License

Yathahi is permanently licensed under the **Apache License 2.0**. See the [LICENSE](LICENSE) and [NOTICE](NOTICE) files for details.

Copyright 2026 **J Aromal**.

---

<p align="center">
  <strong>Yathahi</strong>
  <br/>
  Unity through interoperability.
</p>
