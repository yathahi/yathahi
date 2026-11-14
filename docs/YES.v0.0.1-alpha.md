# YES — Yathahi Engineering Standards

**Document:** `YES.v0.0.1-alpha.md`  
**Version:** `0.0.1-alpha`  
**Status:** Alpha  
**Project:** Yathahi

---

## 1. Purpose

The Yathahi Engineering Standards (YES) define how Yathahi code is owned, classified, organised, shared, and connected across architectural boundaries.

This document is intended for every person who develops, reviews, contributes to, or maintains Yathahi. A contributor should be able to use it to decide:

- whether code belongs to an architectural component or to the shared library;
- where technology-specific integrations belong;
- when and how functionality should be extracted into a shared package;
- how architectural components are allowed to communicate; and
- how Go packages should be organised without confusing package names with architectural components.

The purpose is not to prescribe a complete directory tree before the implementation requires one. The structure must reflect the actual architecture and the actual responsibilities of the code.

## 2. Foundational model

Yathahi distinguishes three concepts. They must not be treated as interchangeable.

### 2.1 Architectural components

An **architectural component** is an actual, defined part of Yathahi's architecture. It has its own responsibility and owns the implementation that is exclusive to that responsibility.

An architectural component is an architectural identity, not merely a folder name or a convenient grouping of functions. Its existence and name must come from the actual design of Yathahi.

A component may contain multiple source files and Go packages. The number of packages it contains does not determine whether it is one architectural component or several.

### 2.2 Shared library packages

A **shared library package** contains a cohesive capability or implementation that is genuinely needed by more than one architectural component.

A library package is not automatically an architectural component. Its name describes the functionality or domain it provides, not an inner/outer execution role or another architectural identity.

The `library/` hierarchy is the central home for these shared packages. It is expected to grow and to be classified according to the actual responsibility of each package.

### 2.3 APIs between architectural components

An **API** is the defined boundary through which architectural components communicate and request work from one another.

API communication is different from code reuse. A component uses a shared library package to reuse functionality; it uses another component's API to communicate with that component and ask it to perform work within its own responsibility.

These three concepts serve different purposes:

| Concept | What it represents | What belongs there |
|---|---|---|
| Architectural component | An actual part of the architecture | Functionality exclusive to that component |
| Shared library package | A reusable implementation capability | Functionality genuinely shared by multiple components |
| API | A communication boundary | Requests, responses, state exchange, and instructions between components |

## 3. The ownership rule

> **Exclusive functionality belongs to the architectural component that owns it. Functionality genuinely shared by multiple architectural components belongs in an appropriately classified package under `library/`. Architectural components communicate with one another through APIs, not by importing or calling one another's internal implementation.**

This is the governing rule for code placement and architectural boundaries throughout Yathahi.

### 3.1 Keep exclusive functionality with its owner

When functionality belongs to one architectural component and is not shared by another, it remains inside that component.

Do not move code into `library/` merely because it looks generally useful, could possibly be reused in the future, or has a name that sounds generic. There must be a real shared responsibility.

### 3.2 Share the functionality, not the architectural identity

When multiple components need the same capability, extract the shared capability into a library package. Do not turn the architectural components into one component, and do not make one component depend on another component's internal code.

The extraction should be as narrow as the responsibility allows. A component may retain its exclusive workflow while a common operation used by that workflow is moved to the library.

Do not move an entire component, package, or file just because one part of it becomes shared. Move or restructure the specific cohesive functionality that is shared; leave the exclusive parts with their original owner.

### 3.3 Do not create speculative shared packages

A possible future consumer is not sufficient reason to create a shared package today. Begin with the actual owner. When another architectural component genuinely needs the same functionality, evaluate and perform the appropriate extraction.

This keeps the library meaningful instead of allowing it to become a collection of speculative helpers.

## 4. The `library/` hierarchy

`library/` is a central hierarchy for shared implementation. It is not one universal package and is not limited to generic utilities.

Its subdirectories should be classified by the actual responsibility, domain, or technology of the code. For example, shared capabilities may belong under areas such as:

```text
library/
├── core/
│   ├── verifier/
│   ├── configurator/
│   └── analyzer/
└── infrastructure/
    ├── nginx/
    ├── lxc/
    └── powerdns/
```

This is an **illustrative classification only**, not a declaration that these packages or architectural components currently exist. Create a package only when implementation establishes a real need for it.

### 4.1 Package names describe capabilities

A shared package should have a normal functional or domain name. It should not be named after one consumer when the package represents a capability shared across consumers.

For example, a shared verifier capability may be classified as `library/core/verifier`, rather than `library/core/inner-verifier` or `library/core/outer-verifier`. In this path, `core` is a library classification for shared functionality. It is not another architectural component and must not be confused with an architectural component that may happen to use the word *core* in its name.

Names such as `inner-verifier`, `outer-verifier`, `inner-analyzer`, `outer-analyzer`, `inner-core`, `outer-core`, `inner-bootstrap`, and `outer-bootstrap` are illustrative examples in this standard. They do not assert that a component with any particular example name currently exists or must be created. Actual component names must follow the actual Yathahi architecture.

### 4.2 Classify, merge, and split by responsibility

The library is intended to be a carefully classified hierarchy, not a dumping ground.

A package should have a coherent responsibility. Organise packages according to what their code actually does and the domain it belongs to. Merge or split packages when the implementation demonstrates that their responsibilities belong together or need clear separation.

Do not create a separate package for every function, nor keep unrelated responsibilities together simply to reduce the number of directories. Package boundaries should make ownership, purpose, and dependencies easier to understand.

### 4.3 The library is not an architectural layer that owns every capability

A library package provides implementation to its consumers. It does not take ownership of the consumer's architectural workflow or decisions.

For example, a shared configurator may provide common configuration behaviour, while each consuming architectural component retains responsibility for its own workflow and decisions about when to configure something.

## 5. Technology-specific integrations

Technology integrations follow the same ownership rule as every other implementation.

Examples of technologies include NGINX, LXC, and PowerDNS. Their names in this document are examples, not a claim that every integration already exists or belongs in `library/` today.

### 5.1 A single consumer

If an integration is used exclusively by one architectural component, keep that integration inside the owning component.

For example, if only one component needs a particular NGINX operation, that operation belongs with that component. There is no need to create a shared NGINX library package yet.

### 5.2 Multiple consumers

If multiple architectural components need the same technology-specific capability, extract the shared part into an appropriately classified library package, such as `library/infrastructure/nginx` when that is the suitable classification.

Keep any technology-specific functionality that remains exclusive to one component with that component. The existence of a shared NGINX package does not mean every NGINX-related operation must be moved into it.

### 5.3 Keep technology knowledge in its integration

Technology-specific details should be owned by the relevant technology integration. Generic shared workflows should not become scattered collections of technology-specific implementation details.

For example, a generic configurator should coordinate the general configuration process. NGINX-specific templates, paths, or other NGINX-specific knowledge should be provided by the NGINX integration when that integration is shared, or remain with the sole owning component while it is exclusive.

This separation keeps generic functionality from becoming tightly coupled to individual technologies.

## 6. Extracting shared functionality

Extraction is a normal consequence of the ownership model. It should happen when actual implementation reveals a shared responsibility—not merely because reuse might happen someday.

Use the following process:

1. **Identify the owners.** Determine which actual architectural components need the functionality.
2. **Identify the shared responsibility.** Separate the common capability from each component's exclusive workflow, decisions, and environment-specific behaviour.
3. **Choose the library classification.** Place the shared capability in the library area that best represents its real function or domain.
4. **Extract only what is shared.** Move or restructure the smallest cohesive responsibility that can be maintained as a meaningful package or package API. Do not split code into arbitrary fragments just to claim reuse.
5. **Keep exclusive code with its owner.** The original component retains its component-specific orchestration and behaviour.
6. **Update the consumers.** Each component uses the shared package for the common capability while remaining responsible for its own work.
7. **Check the boundaries.** Ensure the library package does not import an architectural component and that the change does not create circular dependencies or direct component-to-component coupling.

If a function combines shared work with component-specific decisions, refactor the responsibilities apart. The shared operation belongs in the library; the component-specific flow remains with its owner.

## 7. Communication between architectural components

Architectural components may depend on one another operationally. They are not necessarily independent systems. They may need to exchange state, coordinate work, report conditions, or request actions.

However, operational cooperation must not become direct source-code dependency.

### 7.1 API communication is mandatory across component boundaries

One architectural component must not import another architectural component's internal Go package or call its internal functions to make that component perform its own work.

Instead, the requesting component communicates through the receiving component's defined API. The receiving component interprets and handles the request within its own responsibility.

For example, the interaction may look like this:

1. One component reports its current state through the API.
2. Another component receives the report and determines what response is required within its own role.
3. It sends an instruction or request through the API.
4. The receiving component performs the requested work in its own environment and reports the result or updated state.

This is an illustrative interaction pattern, not a specification of a particular API protocol.

### 7.2 Do not confuse a request with code reuse

Use the correct mechanism for the intention:

- If a component needs another component to perform an operation as part of that component's responsibility, use the API.
- If both components need the same reusable implementation capability, extract that capability into `library/` and let both use it.

Do not use a library package as a back door into another component's internal workflow. Equally, do not use an API merely to share a local reusable function that belongs in a library package.

### 7.3 Keep the API boundary explicit

The API should define what a component can request, what information it can report, and what response or result is expected. Internal implementation details should remain private to the component that owns them.

YES does not mandate a transport, serialization format, network port, or protocol. Those are separate technical decisions and must not be inferred from this ownership standard.

## 8. Dependency boundaries

The following dependency rules preserve the separation between architecture and implementation:

| Dependency | Rule |
|---|---|
| Architectural component → shared library package | Allowed when the component genuinely needs the capability |
| Shared library package → architectural component | Not allowed |
| Architectural component → another component's internal Go package | Not allowed |
| Architectural component → another component through its API | Required for cross-component communication |
| Shared library package → another shared library package | Allowed when the dependency is justified and does not create a cycle |

### 8.1 No component-to-component imports

Do not import another architectural component to reuse one of its functions, types, or implementation details. Extract genuinely shared functionality into the library instead.

Do not import another architectural component to make it execute work. Use its API.

### 8.2 No reverse ownership dependency

A shared library package must not depend on a consumer's architectural component. The library provides capabilities to components; it does not know or control the components that use it.

### 8.3 Avoid circular dependencies

Package dependencies must remain clear and acyclic. If two packages require each other's internal implementation, reassess their responsibilities and boundaries rather than adding workarounds that preserve the cycle.

## 9. Go package organisation

YES applies the ownership model to Go without confusing Go packages with architectural components.

- A Go package is an implementation and dependency boundary; it is not automatically an architectural component.
- A directory normally represents one Go package. Multiple `.go` files in that directory can belong to the same package.
- A conceptual directory such as `library/` can contain multiple subdirectories, each representing its own Go package. The parent directory does not need to be one Go package.
- An architectural component may contain more than one Go package when its implementation responsibilities justify that separation.
- Do not create a new architectural component merely because a new Go package is needed.
- Do not create a library package merely because a function could technically be imported elsewhere. Its functionality must be genuinely shared across architectural components.

Package structure should follow cohesive responsibility and dependency boundaries. It should not be driven only by file count, line count, naming symmetry, or an imagined future requirement.

## 10. Contributor decision procedure

Before placing new code or moving existing code, answer these questions in order:

### A. Is this part of an actual architectural component?

- **Yes:** Identify the component that owns the responsibility and continue.
- **No, it is a shared capability:** Identify its correct classification under `library/`.
- **Unclear:** Do not invent an architectural component or package name. Resolve the responsibility and ownership first.

### B. Is the functionality exclusive to one architectural component?

- **Yes:** Keep it with that component.
- **No:** Identify the other actual component or components using the same capability and evaluate extraction into `library/`.

### C. Is the code technology-specific?

- **Yes, and exclusive to one component:** Keep it with that component.
- **Yes, and shared by multiple components:** Put the shared technology-specific capability in the appropriate library integration package.
- **No:** Classify it by its actual shared function or domain.

### D. Is one component trying to make another component perform work?

- **Yes:** Use the receiving component's API. Do not import its internal implementation.
- **No; both need the same implementation capability:** Extract the common capability into the library.

### E. Is the proposed package boundary meaningful?

Confirm that the package has a coherent responsibility, a clear owner, and justified dependencies. Merge or split only when the actual responsibilities support doing so.

## 11. Review checklist

A contribution that changes architecture, package ownership, or shared functionality should be reviewed against these questions:

- Is every exclusive implementation still owned by the correct architectural component?
- Has code been placed in `library/` only when it is genuinely shared by multiple architectural components?
- If only part of a component's implementation is shared, has only that cohesive part been extracted?
- Are shared packages named and classified by their actual capability rather than by one consumer?
- Are technology-specific details kept with the relevant integration?
- Does cross-component communication use an explicit API rather than direct imports or internal function calls?
- Does any library package import an architectural component?
- Have circular dependencies or unnecessary component dependencies been introduced?
- Does each package have a clear, cohesive responsibility?
- Are hypothetical examples being kept separate from the actual architecture and implementation?

## 12. What this standard does not assume

To prevent this alpha standard from being mistaken for an implementation inventory, the following limits are explicit:

- All component and technology names used as examples—including Inner Core, Outer Core, Inner Bootstrap, Outer Bootstrap, Inner Verifier, Outer Verifier, Inner Analyzer, Outer Analyzer, NGINX, LXC, and PowerDNS—are illustrative labels only. Their mention does not establish that a named component or package currently exists or must be created.
- YES does not prescribe a complete directory tree or require every example package to be created.
- YES does not define the full list, names, number, or execution environments of Yathahi's architectural components. Those must come from the actual architecture.
- YES does not prescribe an API transport, wire protocol, serialization format, port, or deployment topology.
- YES does not require speculative abstractions or pre-emptive extraction for possible future consumers.
- YES does not require moving an entire package or component when only one responsibility becomes shared.

The standard defines how to make these decisions; the implemented architecture and real responsibilities determine the resulting structure.

## 13. Maintaining this standard

YES is the common engineering reference for contributors and maintainers. Changes to architecture or code organisation should preserve its ownership and communication principles.

If an implementation decision appears to conflict with YES, first determine whether the code has been classified incorrectly, whether a shared responsibility needs extraction, or whether the architectural/API boundary has been crossed. Do not silently introduce a second, conflicting organisation model.

As Yathahi develops, this standard may be revised. A future revision should make intentional changes explicit rather than allowing outdated or contradictory rules to accumulate.

---

## 14. Core principle

**Architectural components own what is exclusive to them. The shared library owns reusable capabilities that are genuinely used by multiple architectural components. APIs connect architectural components without exposing or importing one another's internal implementation.**

This is the foundation of the Yathahi Engineering Standards, version `0.0.1-alpha`.
