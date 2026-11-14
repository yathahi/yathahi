# Yathahi — Specification

**Document version:** 0.0.1-alpha  
**Document type:** System specification  
**Status:** Initial specification; requirements may be refined as the project develops.

## 1. Purpose

This document defines the intended capabilities, responsibilities, behavioural requirements and boundaries of Yathahi. It translates the goals described in the Project Vision and the component responsibilities described in the Architecture into requirements that can guide implementation and later be verified through tests and demonstrations.

This specification describes what Yathahi is expected to do. It does not claim that every requirement is already implemented, and it does not prescribe every internal implementation choice.

The specification applies to Yathahi whether it runs on one computer hosting several services, a small home laboratory, a small office, or a much larger infrastructure distributed across devices and locations.

## 2. Requirement language

The terms below indicate the strength of a requirement:

- **MUST** means the requirement is necessary for the intended behaviour described by this specification.
- **SHOULD** means the requirement is expected in the normal design, but a justified exception may be possible.
- **MAY** means the capability is optional or may be introduced as the project develops.

A requirement in this document describes intended system behaviour. It must not be interpreted as evidence that the behaviour is already present in the current implementation.

## 3. System-wide requirements

### 3.1 Scale independence

- **SYS-001:** Yathahi MUST be usable on a single device, including a computer that runs Yathahi and hosts multiple services.
- **SYS-002:** Yathahi MUST NOT require a second device, a separate server, or a distributed network for a single-device deployment to be meaningful.
- **SYS-003:** Yathahi MUST support small multi-device environments, including home laboratories and small offices.
- **SYS-004:** Yathahi's design MUST allow an environment to expand across devices and locations without changing the fundamental responsibilities of Entry, Core, Engine and Administrator.
- **SYS-005:** A small deployment MUST NOT be required to adopt unnecessary complexity solely because Yathahi can also support large deployments.

### 3.2 Infrastructure and service management

- **SYS-006:** Yathahi MUST allow users to request the creation and management of services without requiring them to perform every underlying infrastructure task manually.
- **SYS-007:** Yathahi MUST keep the responsibilities of collecting user requirements, managing device resources, and operating specialised services separated between the appropriate components.
- **SYS-008:** Yathahi MUST maintain a distinction between infrastructure management and the independent operation of the services deployed through it.
- **SYS-009:** The underlying systems and services established with Yathahi SHOULD continue to operate independently if Yathahi's management components become unavailable or are removed, to the extent supported by those underlying systems.
- **SYS-010:** The absence of Yathahi MUST NOT, by itself, be designed to make an otherwise independently operating service stop merely because Yathahi is no longer managing it.

## 4. Entry requirements

Entry is the starting point directly installed and invoked by the user. It is responsible for initial Core establishment when necessary and for collecting and submitting user requests. It is not the infrastructure-wide graphical administration interface.

### 4.1 Distribution and first use

- **ENT-001:** Entry MUST be the only Yathahi executable that a user is expected to install manually as the starting point for using Yathahi on a device. Other required Yathahi components, including Core and specialised Engines, are established or installed through Yathahi rather than requiring the user to install each one separately.
- **ENT-002:** Entry MUST be usable regardless of which type of service the user intends to run.
- **ENT-003:** On use, Entry MUST determine whether a functioning Core is already available on the device.
- **ENT-004:** If Core is not available, Entry MUST be capable of securely establishing Core on that device, including the setup needed to enable and start it.
- **ENT-005:** The initial Core establishment process MUST minimise the technical work required from the user. The user should not need to separately understand or manually perform Core's internal installation and setup tasks.
- **ENT-006:** Entry MUST NOT treat the absence of Core as a reason to abandon the user's request. It must be able to collect the necessary initial information and pass the request to Core when Core becomes available.

### 4.2 Request collection and interaction

- **ENT-007:** Entry MUST provide an interactive way for the user to state what service or infrastructure capability they want to initiate or manage.
- **ENT-008:** Entry MUST collect the information required to express the user's request and the intended service configuration, to the extent that the information is known at that stage.
- **ENT-009:** Entry MUST pass the collected request and relevant information to Core for handling.
- **ENT-010:** Once Core is available, Entry MUST act as a user-facing means of submitting requests to Core for creating or managing services.
- **ENT-011:** Entry MUST NOT take over Core's container-management responsibilities or perform an Engine's service-specific work itself.
- **ENT-012:** Entry MAY provide simple interaction or request status, but it MUST NOT be treated as a substitute for the Administrator's full graphical administration interface.

## 5. Core requirements

Core is the central management component on an individual device. It connects infrastructure-level direction from Administrator with the specialised work carried out by Engines.

### 5.1 Container lifecycle

- **CORE-001:** Core MUST manage the lifecycle of containers that are managed by Yathahi on its device.
- **CORE-002:** Core MUST be able to create a managed container when a requested service requires one.
- **CORE-003:** Core MUST be able to remove a container that belongs to Yathahi's managed environment when its removal is authorised and appropriate.
- **CORE-004:** Core MUST keep container lifecycle management separate from the service-specific operation performed inside a container.
- **CORE-005:** Core MUST NOT enter managed containers or operate the service internals inside them as part of its normal service-management responsibility.

### 5.2 Engine installation and assignment

- **CORE-006:** Core MUST manage the availability of the Engines required by the services requested on its device.
- **CORE-007:** If a required Engine is not available, Core MUST be able to arrange for that Engine to be installed before assigning the work to it.
- **CORE-008:** Core MUST assign an Engine to the container or containers in which that Engine is authorised to operate.
- **CORE-009:** Core MUST provide the assigned Engine with the request information, configuration data and instructions needed for the Engine to perform its specialised work.
- **CORE-010:** Core MUST coordinate Engines through a common communication method that does not require Core to understand the service-specific procedures of each Engine.
- **CORE-011:** Core MUST NOT perform the specialised installation, configuration or operation of a service that belongs to an Engine's responsibility.

### 5.3 Device records, state and logs

- **CORE-012:** Core MUST maintain an accurate record of the Yathahi-managed containers that exist on its device.
- **CORE-013:** Core MUST maintain records of the Engines available on the device and the containers assigned to each Engine.
- **CORE-014:** Core MUST maintain a record of the services each assigned Engine is responsible for managing.
- **CORE-015:** Core MUST maintain relevant operational state and records for its managed containers, Engines and services.
- **CORE-016:** Core MUST manage or maintain access to the operational logs needed to understand the state and behaviour of its managed environment, without taking over the service-specific work of an Engine.
- **CORE-017:** Core SHOULD be able to identify when the observed state of a managed container, Engine or service differs unexpectedly from its recorded or expected state.

### 5.4 Authentication, authority and device security

- **CORE-018:** Core MUST identify and authenticate Administrators before accepting instructions that require Administrator authority.
- **CORE-019:** Core MUST determine whether an authenticated instruction has the authority required for the requested action.
- **CORE-020:** Core MUST reject or safely handle unauthenticated, invalid or unauthorised privileged instructions.
- **CORE-021:** Core MUST protect the integrity of its own device-level management responsibilities and the records it maintains.
- **CORE-022:** Core MUST monitor for suspicious conditions within the device-level environment for which it is responsible.
- **CORE-023:** When Core identifies a condition that may threaten the wider infrastructure, it MUST report the concern to Administrator through an authorised communication path.
- **CORE-024:** Core MUST remain responsible for device-level management and protection even when Administrator provides infrastructure-wide direction.

## 6. Engine requirements

An Engine is a specialised component responsible for a particular service or workload. An Engine understands how to operate its own service; it does not manage the device's container lifecycle or the wider infrastructure.

### 6.1 Service-specific operation

- **ENG-001:** Each Engine MUST have a defined service or workload responsibility.
- **ENG-002:** An Engine MUST perform the service-specific work required to install, configure, start, operate, monitor and maintain the service within its assigned container or containers, according to its purpose and capabilities.
- **ENG-003:** An Engine MUST interpret the common instructions supplied through Yathahi according to its own service-specific knowledge.
- **ENG-004:** An Engine MUST use the request and configuration information provided to it by Core when establishing or changing its assigned service.
- **ENG-005:** An Engine MUST report relevant service state, completion results, failures and unexpected conditions to Core.

### 6.2 Workload boundaries

- **ENG-006:** An Engine MUST operate within the container or containers assigned to it.
- **ENG-007:** An Engine MUST NOT automatically receive access to containers or workload internals assigned to other Engines.
- **ENG-008:** An Engine MUST NOT create, allocate or delete containers as an independent responsibility. Container lifecycle management belongs to Core.
- **ENG-009:** An Engine MUST NOT assume responsibility for installing, assigning or managing other Engines.
- **ENG-010:** An Engine MUST protect the service and workload within its own area of responsibility and MUST report relevant suspicious activity or compromise to Core.

### 6.3 Instructions and wider coordination

- **ENG-011:** An Engine MUST be able to receive and interpret authorised instructions from Core through the common communication method.
- **ENG-012:** Where required by the architecture and properly authorised, an Engine MAY receive an immediate instruction from Administrator.
- **ENG-013:** An immediate Administrator instruction MUST NOT grant an Engine general authority over Core's device-management responsibilities or another Engine's assigned workload.

## 7. Administrator requirements

Administrator is the logical authority responsible for the infrastructure as a whole. It coordinates participating devices and Cores and is responsible for infrastructure-wide and network-level security.

### 7.1 Infrastructure coordination

- **ADM-001:** Administrator MUST be able to coordinate the Cores of participating devices as parts of one logical infrastructure.
- **ADM-002:** Administrator MUST maintain or obtain the information needed to understand participating devices, their reported state and their relationships within the infrastructure.
- **ADM-003:** Administrator MUST be able to coordinate infrastructure-level changes that involve one or more participating devices.
- **ADM-004:** Administrator MUST NOT need to understand the service-specific procedures performed by individual Engines in order to coordinate the infrastructure.
- **ADM-005:** Administrator MUST be able to coordinate environments in which some devices provide services while others participate without hosting services.

### 7.2 Network-level security and isolation

- **ADM-006:** Administrator MUST be responsible for infrastructure-wide and network-level security.
- **ADM-007:** Administrator MUST govern authorised communication between participating devices according to the infrastructure's security requirements.
- **ADM-008:** Administrator MUST be able to assess security reports received from Cores and coordinate an infrastructure-wide response where necessary.
- **ADM-009:** Administrator MUST be able to isolate a device from the wider infrastructure when that device is considered compromised or presents an unacceptable risk.
- **ADM-010:** Isolation of one compromised device SHOULD limit the ability of that device to affect other participating devices and services.

### 7.3 Multiple Administrators

- **ADM-011:** Yathahi MUST support multiple Administrators for an infrastructure.
- **ADM-012:** Administrators belonging to the same infrastructure MUST have equal rights and capabilities under the current architecture.
- **ADM-013:** The infrastructure's continued governance MUST NOT depend permanently on one specific person or one physical device acting as its sole Administrator.

### 7.4 Administrator graphical interface

- **ADM-014:** Yathahi MUST provide its full graphical infrastructure administration and customisation experience through the Administrator interface, rather than making Entry the full administration application.
- **ADM-015:** The Administrator graphical interface SHOULD allow authorised users to understand the participating devices, services, their states and their relationships.
- **ADM-016:** The Administrator graphical interface SHOULD provide a way to configure and customise the infrastructure and its services through the authority and coordination provided by Administrator.
- **ADM-017:** Where a proposed change may cause a conflict, failure or security concern, the interface SHOULD make the concern understandable before the change is applied, where that can be determined.

## 8. Communication and coordination requirements

This section defines the required communication capabilities without choosing a final protocol, message format or transport.

- **COM-001:** Entry MUST have an authorised way to submit user requests and collected information to Core.
- **COM-002:** Administrator and Core MUST be able to exchange infrastructure instructions, device state, authentication-related information and relevant reports through an authorised communication method.
- **COM-003:** Core and Engines MUST share a common communication method for instructions, status, results and reports.
- **COM-004:** The common Core–Engine communication method MUST allow new specialised Engines to be introduced without requiring Core to understand each Engine's service-specific operating procedures.
- **COM-005:** Communications that request privileged operations MUST be authenticated and authorised by the component responsible for the requested operation.
- **COM-006:** Communication failures, invalid requests and rejected instructions MUST be handled without silently treating an unsuccessful operation as successful.
- **COM-007:** Yathahi MUST support coordination between services and devices when a requested outcome depends on more than one service or device.
- **COM-008:** Cross-service or cross-device coordination MUST respect the responsibilities and security boundaries of the involved Core, Engines and Administrator.
- **COM-009:** An Engine MUST NOT bypass another Engine's workload boundary merely because a service-to-service interaction is required.

The exact routing of cross-service requests—including when a request is handled through Core, coordinated by Administrator, or delivered to another relevant service—is not fixed by this initial specification and must be defined in a later communication specification.

## 9. Security and trust requirements

Yathahi's security responsibilities are layered across infrastructure, device and workload boundaries.

- **SEC-001:** Administrator MUST provide infrastructure-wide and network-level protection.
- **SEC-002:** Core MUST provide device-level protection and protect its own management responsibilities and records.
- **SEC-003:** Each Engine MUST protect the service and workload assigned to it.
- **SEC-004:** Each component MUST report security concerns to the component responsible for the next relevant level of coordination: Engine to Core, and Core to Administrator.
- **SEC-005:** A compromise of one Engine MUST NOT automatically grant access to other Engines' assigned containers.
- **SEC-006:** A compromise of one device MUST NOT automatically grant that device unrestricted authority over the wider infrastructure.
- **SEC-007:** Administrator MUST be able to coordinate containment of a compromised device at the infrastructure level.
- **SEC-008:** Core MUST authenticate and authorise Administrator instructions before carrying out privileged device-level actions.
- **SEC-009:** The initial establishment of Core by Entry MUST be performed securely. The distribution and setup process MUST provide a trustworthy basis for installing and starting Core.
- **SEC-010:** Yathahi's additional security capabilities SHOULD complement, rather than unnecessarily replace or disable, the security controls provided by the underlying operating systems, container environment and services.

## 10. Service lifecycle and operational requirements

- **OPS-001:** A service request MUST be represented with enough information for Core to determine the required local container and Engine arrangements, or to identify what additional information is needed.
- **OPS-002:** Core MUST coordinate the creation or selection of the required managed container environment and the availability and assignment of the appropriate Engine.
- **OPS-003:** The assigned Engine MUST perform the service-specific setup and operation using the information and instructions it receives.
- **OPS-004:** Core MUST maintain records of the resulting container, Engine assignment and service responsibility.
- **OPS-005:** Relevant service and workload status MUST be reportable from Engine to Core.
- **OPS-006:** Relevant device and managed-environment status MUST be reportable from Core to Administrator.
- **OPS-007:** A service, Engine or container failure MUST NOT be recorded or presented as a successful operation merely because a request was accepted.
- **OPS-008:** Yathahi SHOULD identify unexpected changes or failures in managed resources and make the relevant condition available to the responsible component.
- **OPS-009:** Corrective action and recovery MUST respect the authority and workload boundaries defined in this specification.
- **OPS-010:** Yathahi SHOULD automate safe, well-understood maintenance and recovery actions while preserving appropriate human visibility and control for consequential actions.

## 11. Availability and independent operation

- **AVL-001:** The loss or removal of Entry after Core has been established MUST NOT, by itself, stop Core or the services managed through it.
- **AVL-002:** The loss or removal of Yathahi's management components MUST NOT be designed to disable independently operating services solely because their Yathahi management layer is unavailable.
- **AVL-003:** When Yathahi's management functions are unavailable, users may lose the ability to submit new Yathahi-managed requests, coordinate infrastructure changes, receive Yathahi monitoring and use Yathahi-specific security and maintenance capabilities.
- **AVL-004:** The exact behaviour of pending operations, disconnected Cores, unavailable Administrators and recovery after reconnection MUST be defined in a later availability and recovery specification.

## 12. Boundaries and exclusions

The following boundaries are part of the intended specification:

- Entry is the user's starting point and request-collection tool; it is not the full Administrator graphical interface.
- Core manages device-level containers, Engine installation and assignment, records, authentication and device-level protection; it does not perform service-specific operations inside containers.
- An Engine operates and protects its assigned service or workload; it does not manage the device's container lifecycle or the wider infrastructure.
- Administrator coordinates the wider infrastructure and network-level security; it does not replace Core's device-management role or an Engine's service-specific knowledge.
- The four components are the architectural components. The Administrator graphical interface is an interface for Administrator, not an additional fifth architectural component.
- This specification does not select a programming language, container technology, network transport, wire protocol, message encoding, database, authentication algorithm or specific service implementation.

## 13. Items requiring more detailed definition

The following matters are intentionally left open for more focused specifications or technical designs. Their absence here is not permission to implement them inconsistently; it means the relevant rules have not yet been fully defined.

1. **Entry and Core establishment:** How Core's authenticity and integrity are verified, how installation failures are handled, and how a first-use request is retained until Core is ready.
2. **Communication contracts:** Exact message structures, request identifiers, response formats, error codes, version compatibility and transport mechanisms for Entry–Core, Core–Engine and Administrator–Core communication.
3. **Authentication and trust:** Identity representation, credential lifecycle, trust establishment, revocation, authority delegation and the detailed rules for multiple Administrators.
4. **Cross-service coordination:** The exact path for service-to-service requests, including when Core coordinates locally and when Administrator coordinates across devices.
5. **Container and Engine lifecycle:** Detailed states, assignment changes, replacement, upgrade, removal and handling of partially completed operations.
6. **Logs and records:** Which events must be recorded, log format, retention, access, integrity and how records are reconciled after interruption.
7. **Failure and recovery:** Behaviour when Entry, Core, an Engine, a container or Administrator becomes unavailable; retry rules; recovery; and safe handling of incomplete operations.
8. **Security response:** Detection criteria, containment procedures, isolation scope, review and recovery after a device or workload is considered compromised.
9. **Administrator interface:** Detailed GUI behaviour, visual configuration model, permissions, confirmation flows and presentation of warnings.
10. **Scale and resource policy:** Resource allocation, limits, scheduling and performance expectations for deployments of different sizes.

These subjects may become separate protocol, security, interface, lifecycle or technical-design documents as the requirements become sufficiently clear.

## 14. Verification and evolution

Requirements in this document should become the basis for implementation tests, integration tests and demonstrations. Where practical, work should be traceable to the requirement identifiers used above. A requirement should not be considered satisfied solely because a related component or feature exists; its expected behaviour and boundaries must also be demonstrated.

This document records the specification at version **0.0.1-alpha**. It is an initial statement of intended behaviour, not a claim that all details are final. Requirements may be clarified, expanded, replaced or removed as Yathahi develops. Changes should be reflected in the document's version history and reviewed against the Project Vision and Architecture so that the three documents remain consistent.
