# Yathahi Roadmap

**Unity through interoperability.**

This roadmap outlines the planned development direction of Yathahi, from its initial functional release towards a unified, interoperable infrastructure ecosystem.

The roadmap is intentionally milestone-oriented. Individual implementation tasks will be tracked separately through GitHub Issues and Pull Requests. Milestones may evolve as development progresses.

---

## Development Principles

Yathahi's development will follow these principles:

- **Interoperability:** Components and services should communicate through clearly defined interfaces.
- **Automation:** Infrastructure provisioning, configuration, and coordination should require minimal manual intervention.
- **Modularity:** Independent Engines should provide specialised infrastructure capabilities.
- **Resource efficiency:** Prefer lightweight deployment approaches, including Linux containers, where appropriate.
- **Reliability:** Operations should be verifiable, with meaningful execution results and failure reporting.
- **Incremental development:** Each milestone should produce demonstrable, functional progress.

---

## Phase 1 — Initial Infrastructure Ecosystem

**Target release:** `v0.1.0`

**Status:** Planned

### Objective

Establish the first functional Yathahi ecosystem capable of automatically provisioning, configuring, and coordinating a miniature web-and-DNS infrastructure within a local network.

The initial release will demonstrate the interaction between Entry, Core, and Engine through a practical end-to-end deployment.

### Scope

**Entry**
- Collect the user's infrastructure requirements.
- Accept the requested service and domain name.
- Transmit the request to Core.

**Core**
- Receive and process infrastructure requests.
- Provision the required Linux containers.
- Install and initialise the required Engines.
- Coordinate Engine execution.
- Collect and report execution results.
- Communicate infrastructure information to the Administrator management interface.

**Web Engine**
- Install and configure the selected web service.
- Launch the web service inside its assigned environment.
- Report its operational status to Core.

**DNS Engine**
- Deploy PowerDNS.
- Configure the required DNS zone and records.
- Obtain the web service's network address through Core.
- Make the web service discoverable through its configured domain.

**Administrator management**
- Receive infrastructure information from Core.
- Register the resulting services.
- Manage the information and connectivity required for service coordination.

### Acceptance Criteria

- A user can initiate a deployment through Entry.
- Core can provision the required containers.
- Web Engine and DNS Engine can be installed and initialised.
- Both Engines can independently execute their assigned responsibilities.
- DNS configuration can automatically reflect the web service's network address.
- Engines can report their execution results to Core.
- Core can communicate the resulting infrastructure state to the Administrator.
- A client on the local network can resolve the configured domain and access the deployed web service.
- Deployment failures are reported meaningfully.
- The deployment process can be repeated.

### Exclusions

The initial release will not target:

- Public internet deployment.
- Authentication and authorisation mechanisms.
- Automated TLS certificate provisioning.
- High availability and clustering.
- Multiple DNS or web-server providers.
- Advanced graphical management interfaces.
- Automated router and DHCP configuration.

The initial environment is intended for controlled local-network testing. Unauthenticated management interfaces must not be exposed to untrusted networks.

---

## Phase 2 — Infrastructure Management and Reliability

**Target release:** `v0.2.0`

**Status:** Proposed

### Objective

Improve the reliability, observability, and management capabilities of the initial ecosystem.

### Direction

- Improve infrastructure and service state tracking.
- Strengthen Engine lifecycle management.
- Improve failure detection and recovery.
- Introduce more comprehensive diagnostics and verification.
- Improve configuration validation.
- Make deployment and configuration operations safely repeatable.
- Improve service discovery and communication between components.

---

## Phase 3 — Security and Network Integration

**Target release:** `v0.3.0`

**Status:** Proposed

### Objective

Introduce the security and networking capabilities required to operate Yathahi in more demanding environments.

### Direction

- Introduce authentication and authorisation.
- Secure communication between components.
- Introduce TLS support.
- Improve network configuration and connectivity management.
- Strengthen isolation between services.
- Improve administrative access and security controls.

---

## Phase 4 — Extensible Infrastructure Ecosystem

**Target release:** `v0.4.0`

**Status:** Proposed

### Objective

Expand Yathahi beyond the initial web-and-DNS infrastructure and establish a more extensible Engine ecosystem.

### Direction

- Define and mature the Engine integration contract.
- Support additional infrastructure services.
- Improve Engine installation and lifecycle management.
- Expand service discovery and interoperability.
- Improve infrastructure configuration and orchestration.
- Establish a foundation for third-party Engine development.

---

## Phase 5 — Ecosystem Maturity

**Target release:** `v1.0.0`

**Status:** Long-term direction

### Objective

Establish Yathahi as a reliable, extensible, and interoperable infrastructure ecosystem suitable for broader real-world adoption.

### Direction

- Mature the architecture and component interfaces.
- Establish stable compatibility and integration standards.
- Improve scalability and resource management.
- Support more complex infrastructure deployments.
- Strengthen security, reliability, and recovery mechanisms.
- Mature documentation and developer tooling.
- Establish a sustainable contributor and Engine ecosystem.

---

## Roadmap Maintenance

This roadmap is a development direction, not a fixed delivery contract.

- Milestones may be revised as technical understanding improves.
- Release targets represent intended outcomes, not guaranteed deadlines.
- Implementation tasks belong in GitHub Issues.
- Completed functionality should be verified before marking a milestone complete.
- Significant changes to project direction should be reflected in this document.

The priority is to build a working, coherent ecosystem rather than prematurely implement every planned capability.
