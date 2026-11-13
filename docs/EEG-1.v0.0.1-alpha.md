# EEG 1 — Ekathvam Engineering Guide: Infrastructure Architecture

**Document ID:** EEG 1  
**Full name:** Ekathvam Engineering Guide 1 — Infrastructure Architecture  
**Filename:** `EEG-1.v0.0.1-alpha.md`  
**Version:** `0.0.1-alpha`  
**Status:** Alpha  
**Companion standard:** EES — Ekathvam Engineering Standards, version `0.0.1-alpha`  
**Scope:** The organisation, ownership, reuse, and communication of Ekathvam's infrastructure implementation.

---

> **Already comfortable with programming basics?** You can skip the beginner primer in **Part 0** and go directly to [Part I — The Ekathvam Infrastructure Engineering Model](#part-i--the-ekathvam-infrastructure-engineering-model). Part I is where the actual Ekathvam engineering guide begins.

> **New to programming or software architecture?** Start at Part 0. It introduces only the basic ideas needed to understand the rest of this guide. You do not need to learn all of computer science before continuing.

---

# How to use this guide

This guide is deliberately arranged as a learning path. Each part introduces the ideas needed by the next part. Read it from the beginning if the concepts are new to you; use the table of contents to return to a topic when needed.

The guide uses short sections, small examples, tables, and simple diagrams. A section should explain one main idea before moving on to another.

## Contents

- [Part 0 — A small beginner's primer](#part-0--a-small-beginners-primer)
- [Part I — The Ekathvam Infrastructure Engineering Model](#part-i--the-ekathvam-infrastructure-engineering-model)
  - [1. What this guide explains](#1-what-this-guide-explains)
  - [2. Understanding the repository's structural distinction](#2-understanding-the-repositorys-structural-distinction)
  - [3. Architectural components and ownership](#3-architectural-components-and-ownership)
- [Part II — Exclusive and shared implementation](#part-ii--exclusive-and-shared-implementation)
  - [4. Understanding exclusive functionality](#4-understanding-exclusive-functionality)
  - [5. Understanding the shared library](#5-understanding-the-shared-library)
  - [6. When functionality changes from exclusive to shared](#6-when-functionality-changes-from-exclusive-to-shared)
  - [7. Technology integrations follow the same ownership rule](#7-technology-integrations-follow-the-same-ownership-rule)
- [Part III — Communication between architectural components](#part-iii--communication-between-architectural-components)
  - [8. Reusing functionality is not the same as communicating](#8-reusing-functionality-is-not-the-same-as-communicating)
  - [9. API-based communication](#9-api-based-communication)
- [Part IV — Package organisation and dependency decisions](#part-iv--package-organisation-and-dependency-decisions)
  - [10. Deciding where new code belongs](#10-deciding-where-new-code-belongs)
  - [11. Go packages and architectural boundaries](#11-go-packages-and-architectural-boundaries)
  - [12. When to create, split, merge, or extract a package](#12-when-to-create-split-merge-or-extract-a-package)
- [Part V — A complete hypothetical implementation walkthrough](#part-v--a-complete-hypothetical-implementation-walkthrough)
  - [13. The scenario: two components and a growing requirement](#13-the-scenario-two-components-and-a-growing-requirement)
- [Part VI — Mistakes, contribution, and maintenance](#part-vi--mistakes-contribution-and-maintenance)
  - [14. Common mistakes and their corrections](#14-common-mistakes-and-their-corrections)
  - [15. Contributor workflow](#15-contributor-workflow)
  - [16. Maintainer responsibilities and review](#16-maintainer-responsibilities-and-review)
- [Reference section](#reference-section)
  - [Quick decision reference](#a-quick-decision-reference)
  - [Glossary](#c-glossary)

## Important note about examples

**Every sample component name, package name, directory tree, file name, and code fragment in this guide is illustrative unless it is explicitly identified as an existing Ekathvam implementation.** Examples use fictional names such as `component-a` and `component-b` so that readers do not mistake a teaching example for the project's actual architecture.

The examples are intended to teach the decisions behind a structure, not to prescribe names that contributors must create. The actual repository must be understood from its real architectural components and its actual packages.

---

# Part 0 — A small beginner's primer

This part is a short bridge for readers who are new to programming. It is not a complete programming course. It explains the few basic terms that appear repeatedly in the infrastructure guide.

## 0.1 What is software?

A computer can perform tasks when it receives instructions it can understand. A collection of instructions that makes a computer perform useful work is called a **program**. A larger collection of related programs and supporting pieces may form a **software system**.

For example, software might receive a request, check some information, prepare a configuration file, and start a service. Each of those activities may be implemented by different parts of the software.

## 0.2 What is source code?

**Source code** is the human-readable form of the instructions developers write. A programming language gives those instructions a defined structure and meaning.

Ekathvam uses Go for its implementation. This guide occasionally shows small Go examples, but you do not need to understand every Go feature to understand the architectural principles.

## 0.3 What is a file?

A **file** is a place where information is stored. A source-code file contains part of a program's implementation.

A project may contain many source files. A file is not automatically a component, and a file does not have to represent an entire responsibility by itself.

## 0.4 What is a function?

A **function** is a named piece of code that performs an operation. Other code can call that function when it needs that operation.

For example, imagine a function that checks whether a configuration value is valid:

```go
func IsValidPort(port int) bool {
    return port > 0 && port <= 65535
}
```

The example checks a value and returns a result. It does not decide which architectural component owns the larger workflow that uses the result.

A function is one small unit of implementation. It is not automatically a package, library, or architectural component.

## 0.5 What is a directory?

A **directory** (also called a folder) groups files and other directories. It helps people navigate a project and helps a repository present its organisation clearly.

A directory name alone does not determine the architectural meaning of the code inside it. To understand that meaning, look at the responsibility, ownership, and dependencies of the implementation.

## 0.6 What is a Go package?

A **Go package** is a group of Go source files that are compiled together as one package. In normal Go project organisation, files in the same directory belong to the same package when they use the same package declaration.

For example, these two files may belong to one Go package:

```text
example-package/
├── validation.go
└── rules.go
```

The directory is an organisational boundary, and the package is a Go language boundary. A package may contain several files. A single file is not normally treated as a separate package simply because it has a different filename.

## 0.7 What is a component?

A **component** is a recognisable part of a larger system that has a responsibility of its own. It performs work as part of the wider system.

A component is a concept about the system's design and responsibilities. It is not simply another word for a Go package.

A component may be implemented using one Go package or several packages. The number of files or packages does not, by itself, decide how many architectural components exist.

## 0.8 What is a library?

A **library** provides implementation that other parts of a program can use. A library package may offer a reusable operation without owning the complete workflow of the code that uses it.

For this guide, the important distinction is that Ekathvam's shared `library/` hierarchy contains reusable implementation capabilities. It is not a second name for an architectural component.

## 0.9 What is an API?

An **API** is a defined way for one part of a system to request information or work from another part. It describes the operations that can be requested and the information exchanged.

Think of it like a service counter: you make a request through the agreed counter, the responsible service handles it, and you receive a response. You do not walk into the service's private workspace and operate its internal machinery yourself.

In software, the API is the agreed boundary. The mechanism carrying the request—such as HTTP or another transport—is a separate technical choice.

## 0.10 Two different reasons for interaction

This difference will appear throughout the guide:

- **Reuse:** Two parts need the same implementation capability. The shared capability belongs in an appropriate library package.
- **Communication:** One component needs another component to perform work that belongs to the receiving component. They communicate through the receiving component's API.

These are not interchangeable. Part III explains the distinction in detail.

## 0.11 A quick vocabulary map

| Term | Simple meaning |
|---|---|
| Function | A piece of code that performs an operation |
| File | A place where source code is stored |
| Directory | A folder that organises files and other folders |
| Go package | Go source files grouped as one package |
| Architectural component | An actual part of the system with its own responsibility |
| Shared library package | A reusable implementation capability |
| API | A defined boundary for communication and requests |
| Dependency | Something a piece of code relies on |
| Extraction | Separating a shared responsibility from its former location and placing it in the appropriate shared package |

You now have the basic vocabulary needed for the main guide.

---

# Part I — The Ekathvam Infrastructure Engineering Model

This is the beginning of the actual Ekathvam infrastructure guide. It explains the project's engineering model and the rules contributors use to organise infrastructure implementation.

## 1. What this guide explains

Ekathvam's infrastructure implementation contains different kinds of code. Some code belongs exclusively to one architectural component. Other code is genuinely shared by multiple components. Components may also need to communicate with one another to coordinate their work.

Those situations require different solutions. Without clear boundaries, a project can become difficult to understand: code may be duplicated, one component may become dependent on another component's private implementation, or a shared package may accidentally take responsibility for an entire architectural workflow.

This guide teaches how to distinguish those situations and how to respond to each one.

### 1.1 The three central ideas

The infrastructure engineering model rests on three distinct ideas:

1. **Architectural ownership:** An actual architectural component owns the implementation that is exclusive to its responsibility.
2. **Shared implementation:** Functionality genuinely needed by multiple architectural components is placed in the appropriate package within the central `library/` hierarchy.
3. **Architectural communication:** When one component needs another component to do work that belongs to it, the components communicate through an API. They do not call one another's internal Go functions or import one another's internal implementation.

These ideas work together, but each answers a different question:

| Question | Correct concept |
|---|---|
| Who owns this exclusive work? | Architectural component |
| Where should genuinely shared functionality live? | Appropriate package under `library/` |
| How does one component ask another to perform its own work? | API communication |

### 1.2 Relationship with EES

**EES is the authoritative engineering standard. EEG 1 is its explanatory and educational guide for infrastructure architecture.**

EES states the engineering requirements. This guide explains the meaning of those requirements, why the boundaries matter, and how to apply them in practical situations.

The examples in EEG 1 must not be treated as permission to contradict EES. If an example seems inconsistent with the standard, follow EES and review the example for clarification.

### 1.3 What this guide does not prescribe

This guide does not create a list of architectural components or require every illustrative package to exist. The actual architecture determines the real components and their names.

It also does not prescribe a transport protocol, network port, data format, or deployment topology for APIs. Those decisions are separate from the ownership principles explained here.

---

## 2. Understanding the repository's structural distinction

The most important structural distinction is between the architectural components and the shared library hierarchy.

### 2.1 The central `library/` hierarchy

Ekathvam uses one central, standalone `library/` hierarchy for shared implementation packages.

It is not a separate library folder inside every component. It is the common location where functionality that has become genuinely shared is classified and maintained.

### 2.2 Architectural components live outside `library/`

Actual architectural component directories are organised outside the central `library/` hierarchy. Each component retains the implementation that is exclusive to its own responsibility.

This arrangement makes the distinction visible:

- Outside `library/`: actual architectural components and their exclusive implementation.
- Inside `library/`: packages that provide shared implementation capabilities.

### 2.3 Not every repository item is an architectural component

A repository may also contain supporting material such as a README, licence, documentation, build metadata, or other project files. Such material helps maintain or explain the repository; it is not automatically an architectural component or a shared library package.

The important structural rule applies to implementation ownership: exclusive component implementation belongs with its architectural owner, while genuinely shared implementation belongs in the central library.

### 2.4 An illustrative structural picture

The following tree uses deliberately fictional component labels:

```text
example-project/
├── component-a/          # Fictional architectural component
├── component-b/          # Fictional architectural component
├── library/              # Central shared implementation hierarchy
│   ├── core/
│   │   └── shared-capability/
│   └── infrastructure/
│       └── shared-integration/
├── README.md             # Example supporting repository file
└── docs/                 # Example supporting documentation
```

**This is not Ekathvam's actual directory tree.** It illustrates only the relationship between architectural components, the single central library hierarchy, and ordinary repository-supporting material.

Do not create `component-a/`, `component-b/`, or the example library packages because they appear here. Use the actual architecture and real implementation responsibilities to determine the repository's contents.

---

## 3. Architectural components and ownership

### 3.1 What makes something an architectural component?

An architectural component is an actual part of the system's architecture. It has a defined role, a responsibility, and implementation that belongs to that role.

Its identity comes from the architecture—not merely from a directory name, a package name, or the fact that some code has been grouped together.

A component can contain multiple functions, files, and Go packages. It may coordinate several activities while retaining ownership of its own work.

### 3.2 A component owns its exclusive implementation

When an operation is part of one component's responsibility and is not a shared capability, that implementation belongs to that component.

For example, imagine a fictional component that prepares a particular runtime environment. Its environment-specific preparation steps belong to that component because those steps are part of its own work.

The fact that another component might someday perform a similar operation is not enough to move the code into the library.

### 3.3 Component identity is not the same as a functional package name

Consider two different names:

- `component-a`: a fictional architectural component that owns a defined part of the system.
- `validation`: a fictional package name that describes a reusable validation capability.

The first identifies an architectural owner. The second describes what a package provides.

A name such as `verifier`, `configurator`, or `analyzer` does not automatically identify a complete architectural component. It may simply name a shared package that contains common implementation used by more than one component.

### 3.4 Do not invent architectural components to organise ordinary code

A new package does not automatically require a new architectural component. A collection of related functions does not automatically become a component either.

Create or identify an architectural component only when the actual system architecture has a distinct responsibility that warrants it. Use ordinary packages to organise implementation within the existing ownership model.

### 3.5 A component can use shared capabilities without giving away its responsibility

A component may use several shared library packages and still remain responsible for its complete workflow.

For example, a component might:

1. decide that a particular operation is needed;
2. use a shared validation package to check input;
3. use a shared technology integration to obtain a reusable capability; and
4. perform its own component-specific orchestration.

The library supplies reusable implementation. The architectural component retains its decisions and exclusive workflow.

---

# Part II — Exclusive and shared implementation

## 4. Understanding exclusive functionality

### 4.1 The meaning of exclusive

In this engineering model, **exclusive functionality** is implementation that belongs to one architectural component and is not a shared implementation capability used by other architectural components.

"Exclusive" describes ownership. It does not mean that the code is hidden from every developer, permanently fixed in one file, or impossible to reuse in the future.

### 4.2 A simple example

Imagine a fictional project with one component responsible for preparing a particular environment. It contains a function that performs a step specific to that environment.

```text
example-project/
└── component-a/
    └── environment_setup.go
```

The file and directory names are illustrative. The important point is that the operation belongs only to Component A's responsibility, so it stays with Component A.

### 4.3 Why exclusive code should stay with its owner

Keeping exclusive implementation with its owner makes it easier to answer:

- Which component is responsible for this behaviour?
- Where should a change to that behaviour be made?
- Which component should test and review it?
- Which other components should not depend on its internal details?

If exclusive code is scattered into unrelated locations, the real owner becomes harder to identify.

### 4.4 Exclusive does not mean "never reusable"

A function may begin as exclusive and later become a shared capability. That is a normal possibility.

The code should be classified according to the real needs of the current implementation. Do not create an abstraction just to prepare for an imagined future. If another component genuinely needs the same capability later, evaluate extraction at that time.

### 4.5 What if only one small part is shared?

Do not move an entire component or an entire package merely because one function becomes shared.

Separate the shared responsibility from the remaining exclusive behaviour. Extract the shared part and leave the component-specific work with its owner.

This function-level or responsibility-level decision is central to the Ekathvam model.

---

## 5. Understanding the shared library

### 5.1 Why the library exists

Imagine two architectural components that both need to validate the same kind of input. If each component maintains a separate copy of the same validation capability, the implementations can drift apart.

A correction may then need to be made in multiple places. One copy might be updated while another is forgotten.

When the responsibility is genuinely shared, a central package gives both components one maintained implementation to use.

### 5.2 What belongs in `library/`?

The central `library/` hierarchy contains reusable implementation capabilities used by multiple architectural components.

A package inside it should have a clear, cohesive responsibility. Its name should describe the capability, function, domain, or technology it provides—not the identity of one particular consumer.

### 5.3 The library is not one giant Go package

`library/` is the central parent directory for shared packages. Its subdirectories can represent separate Go packages.

For example, the following is an illustrative classification:

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

This tree is **hypothetical**. It does not state that these packages currently exist or that every one must be created.

In this example:

- `core/` is a classification inside the library, not another architectural component.
- `verifier/`, `configurator/`, and `analyzer/` are functional package names, not automatically complete architectural components.
- `infrastructure/` groups technology-related shared packages in this illustration.

### 5.4 Why shared package names do not need inner/outer prefixes

Suppose two actual architectural components each have their own verification responsibility. Their component-specific verification workflow remains with each component.

The common validation capability they both use can live in a shared package with a normal functional name, such as the illustrative `library/core/verifier`.

The shared package does not become an "inner verifier" or an "outer verifier" merely because one of those components uses it. It provides the common capability to both.

### 5.5 The library must not become a dumping ground

Not every function that looks generally useful belongs in the library.

A package belongs there because its implementation is genuinely shared—not because:

- its function name sounds generic;
- it might be useful one day;
- it is convenient to move code away from a large component; or
- a developer wants every package in one place.

Each shared package should have a clear responsibility, justified consumers, and understandable dependencies.

### 5.6 The library does not own the components' workflows

A shared package provides a capability. It does not take over the decisions or complete workflow of the components that use it.

For example, a shared validation package may determine whether a value is valid. Each architectural component still decides when to validate that value, what to do with the result, and how that result affects its own workflow.

---

## 6. When functionality changes from exclusive to shared

This section explains one of the most important changes a contributor may need to make during development.

### 6.1 Start with the actual owner

Imagine that Component A needs a validation operation. It is not currently used by any other architectural component.

It is reasonable to implement that operation within Component A's own implementation. There is no need to create a shared package before a second consumer exists.

### 6.2 A second component now needs the same capability

Later, Component B needs the same validation capability.

Pause and identify exactly what Component B needs. Is it the same reusable validation operation, or does Component B actually need Component A to perform a task that belongs to Component A?

Those are different situations. This section covers the first one; Part III explains the second.

### 6.3 Identify the common responsibility

Suppose both components independently need to check the same input rule. The common rule is the shared responsibility.

The rest of Component A's workflow may still be exclusive to Component A. Component B's workflow may also be different.

Do not merge the two workflows simply because one operation is common to both.

### 6.4 Extract only the shared capability

A possible illustrative result is:

**Before the second consumer exists:**

```text
example-project/
└── component-a/
    └── validation.go
```

**After the common validation capability is extracted:**

```text
example-project/
├── component-a/
│   └── component_workflow.go
├── component-b/
│   └── component_workflow.go
└── library/
    └── core/
        └── validation/
            └── rules.go
```

These names and paths are fictional. They illustrate the change in ownership, not an Ekathvam directory prescription.

The shared validation rule is now in the library. Each component retains its own workflow and uses the common package when needed.

### 6.5 What happens to the original component?

Component A does not disappear, lose its identity, or hand its entire responsibility to the library.

It continues to own its exclusive work. It simply uses the shared package for the capability that is no longer exclusive.

Component B uses the same shared package independently.

### 6.6 Extraction is not necessarily a whole-file move

A source file may contain both shared and exclusive functions. The correct action is not automatically to move the whole file.

Separate the responsibilities in a way that gives the shared code a coherent home. The shared implementation goes into the library; the exclusive implementation stays with its architectural owner.

Sometimes that means moving a function. Sometimes it means refactoring a group of related functions into a package. The goal is correct ownership and a cohesive package boundary—not the smallest possible number of lines.

### 6.7 A practical extraction checklist

Before extracting, ask:

1. Are at least two actual architectural components using or genuinely requiring this capability?
2. Is the functionality itself common, rather than merely having a similar name?
3. Can the shared responsibility be separated from each component's exclusive workflow?
4. Is there an appropriate existing library package, or is a new cohesive package justified?
5. Can both consumers use the shared package without importing one another?
6. Does the extracted package avoid depending on an architectural component?

If the functionality is not genuinely shared, keep it with its existing owner.

---

## 7. Technology integrations follow the same ownership rule

Technology-specific implementation is not automatically shared. It follows exactly the same ownership principle as every other kind of implementation.

Examples in this section use NGINX, LXC, and PowerDNS only to explain the rule. Their appearance here does not mean that every such integration currently exists or should be created in a particular location.

### 7.1 A technology integration used by one component

Imagine only Component A needs a particular NGINX operation.

That implementation is exclusive to Component A at this stage. Keep it with Component A.

```text
example-project/
└── component-a/
    └── nginx_support/
        └── setup.go
```

This is an illustrative example. There is no need to create a central NGINX package merely because NGINX is a known technology.

### 7.2 The same technology capability becomes shared

Later, Component B also needs the same NGINX capability.

Identify the specific NGINX functionality that both components need. Extract that shared technology-specific capability into the appropriate library package.

```text
example-project/
├── component-a/
│   └── component_specific_setup.go
├── component-b/
│   └── component_specific_setup.go
└── library/
    └── infrastructure/
        └── nginx/
            └── shared_capability.go
```

Again, the paths and filenames are hypothetical. The point is that the shared NGINX implementation is centralised, while exclusive component-specific work remains with its owner.

### 7.3 Do not move every NGINX-related function automatically

The fact that an NGINX library package now exists does not mean all NGINX-related code must move into it.

Some operations may still belong exclusively to Component A. Other operations may belong exclusively to Component B. Only genuinely shared NGINX functionality belongs in the shared package.

The same rule applies to LXC, PowerDNS, SQL systems, web servers, and any other technology.

### 7.4 Keep technology-specific knowledge with the integration

When a technology integration is shared, it should own the details specific to that technology.

For example, an NGINX integration may know about NGINX-specific configuration templates, paths, options, or metadata. A generic configuration capability should not become filled with hardcoded NGINX-only details.

This separation allows a generic capability to coordinate a general process while a technology integration supplies the details needed for its own technology.

### 7.5 A generic configurator and a technology integration

Consider a hypothetical shared configurator and a hypothetical shared NGINX integration:

- The **configurator** handles the general configuration workflow.
- The **NGINX integration** provides NGINX-specific knowledge, such as the relevant template or configuration metadata.
- The component using them decides when and why that configuration should happen as part of its own workflow.

The NGINX integration does not need to import the configurator just to provide its technology-specific capability. The generic configurator should not contain scattered hardcoded assumptions about every technology.

A suitable interface or provider contract can be used when the implementation needs one. In Go, a type can satisfy an interface without explicitly declaring that it implements it. Do not create an interface for every package automatically; use one when it gives the consumer a clear and useful boundary.

### 7.6 Technology name does not determine ownership

The question is not simply, "Is this NGINX code?" The useful questions are:

- Which component currently owns this operation?
- Is the operation exclusive to that component?
- Do multiple components genuinely need the same capability?
- Which part is technology-specific, and which part is generic?

Answer those questions before deciding where the implementation belongs.

---

# Part III — Communication between architectural components

## 8. Reusing functionality is not the same as communicating

This distinction must be understood before implementing any interaction between architectural components.

### 8.1 Situation A: Both components need the same function

Imagine Component A and Component B both need the same reusable operation, such as checking whether an input value follows a shared rule.

Neither component needs the other to perform its work. Each needs the same capability.

**Correct approach:** Put the common functionality in the appropriate shared library package. Both components use that package.

```text
component-a/ ─────┐
                  ├────> library/shared-validation/
component-b/ ─────┘
```

The arrows represent use of a shared implementation package, not API requests between components.

### 8.2 Situation B: One component needs another component to do its own work

Now imagine Component A owns the responsibility for preparing a particular environment. Component B needs that environment prepared.

Component B is not asking for a reusable helper function so that it can perform the preparation itself. It is asking Component A—the owner—to perform the operation within Component A's responsibility.

**Correct approach:** Component B sends a request through Component A's API. Component A handles the request using its own internal implementation.

```text
component-b  ───── API request ─────>  component-a
component-b  <──── API response ─────  component-a
```

The components do not import each other's internal Go packages or call each other's internal functions.

### 8.3 The key question

Ask:

**Do I need to reuse the capability, or do I need the other component to perform its responsibility?**

- Reuse the capability: use or create the appropriate shared library package.
- Ask the other component to perform its responsibility: communicate through its API.

This question prevents a common architectural mistake: using direct function calls where an API boundary is required, or using an API when the real need is simply shared implementation.

---

## 9. API-based communication

### 9.1 What an API provides

An API defines the operations a component makes available to other parts of the system. It establishes what can be requested, what information must be supplied, and what response can be expected.

The component receiving the request remains responsible for handling it.

### 9.2 A simple hypothetical sequence

Imagine Component A needs Component B to check and prepare something that belongs to B's responsibility.

1. Component A identifies the operation it needs.
2. Component A sends a request using Component B's defined API.
3. Component B receives the request at its API boundary.
4. Component B validates and handles the request according to its own rules.
5. Component B performs its internal work.
6. Component B returns a result or reports an updated state through the API.

Component A does not need to know which private functions Component B used internally.

### 9.3 State reporting and instructions

Components may need to report their state, communicate a condition, receive an instruction, or return a result.

For example, a component might report:

> "The service is not in the required state."

Another component may decide that an action is needed and send an instruction through the responsible component's API.

The receiving component then handles the instruction within its own environment and responsibilities.

This is an illustrative interaction. It does not define a specific Ekathvam message format or protocol.

### 9.4 Why direct internal function calls are not the solution

Suppose Component A imports Component B's internal Go package and calls a function to make B perform its work.

That creates a direct source-code dependency on B's implementation. A now needs to know how B is internally organised. A change inside B can force changes in A even when B's external responsibility has not changed.

It also blurs ownership: A can begin reaching into B's implementation instead of communicating through the boundary B controls.

The API approach keeps the boundary explicit. B can change its internal implementation while preserving the API contract that other components use.

### 9.5 API does not mean one specific transport

The word **API** describes the agreed communication interface. It does not, by itself, mean HTTP, REST, TCP, JSON, Protocol Buffers, or any other specific technology.

An API may be implemented using a suitable transport and data format. The actual choice must come from the system's technical requirements and architecture. This guide does not choose one on Ekathvam's behalf.

### 9.6 Keep internal implementation private

A component should expose the operations and information needed at its API boundary—not its entire internal implementation.

Other components should not need to know:

- the names of its internal functions;
- how its internal packages are arranged;
- which private helper functions it calls; or
- the sequence of private implementation steps it uses.

They should understand the API contract and the result they can expect.

---

# Part IV — Package organisation and dependency decisions

## 10. Deciding where new code belongs

This chapter turns the ownership principles into a repeatable decision process.

### 10.1 Start with the responsibility, not the filename

Before creating a file or package, describe in one sentence what the code is responsible for.

For example:

- "This operation prepares a particular component's runtime environment."
- "This package provides a validation capability used by multiple components."
- "This operation receives requests and performs work owned by another component."

If the responsibility is unclear, the correct directory will probably be unclear too.

### 10.2 Identify the actual owner

If the functionality is part of an actual architectural component's responsibility, identify that component.

Do not guess based on an example name in this guide. Use the actual architecture and current implementation.

### 10.3 Test whether the functionality is exclusive

Ask whether another actual architectural component needs the same implementation capability today.

- If not, keep it with its owner.
- If yes, investigate whether the common responsibility should be extracted into the library.

Do not count a possible future use as an existing shared requirement.

### 10.4 Decide whether the requirement is reuse or communication

If another component needs the capability itself, share the implementation through the library.

If another component needs the owning component to perform work, use the owning component's API.

### 10.5 Choose a package based on its real responsibility

A package should be named and classified by what it provides. Avoid names that merely describe which component happened to need the code first when the package is now shared.

For example, a shared validation capability should have a functional classification, not a name that makes it appear to belong exclusively to one consumer.

### 10.6 Check the dependency direction

Before finalising the placement, check that:

- architectural components may use the shared library packages they need;
- shared library packages do not import architectural components;
- architectural components do not import one another's internal implementation; and
- cross-component work is performed through APIs.

### 10.7 A compact decision path

```text
What responsibility does this code perform?
        |
        v
Which actual architectural component owns that responsibility?
        |
        v
Is the same implementation capability genuinely needed
by multiple architectural components?
        |
        +-- No --> Keep the exclusive implementation with its owner.
        |
        +-- Yes -> Extract only the shared capability into
                   the appropriate package under library/.

Separate question:
Does one component need another component to perform
work that belongs to the receiving component?
        |
        +-- Yes --> Use the receiving component's API.
```

This diagram is a thinking aid, not a replacement for understanding the responsibility. In a real implementation, answer each question carefully; do not force code into a category just to complete a diagram. A request for another component to perform its own work is an API interaction, not a reason to extract that component's private workflow into the library.

---

## 11. Go packages and architectural boundaries

Go's package system is the implementation mechanism. Ekathvam's ownership model determines how that mechanism should be used.

### 11.1 A directory normally represents a Go package

A Go package is normally formed by Go files in one directory. Several files may belong to the same package and work together.

Do not create a separate directory and package for every function. Group related implementation around a cohesive responsibility.

### 11.2 A component can contain multiple packages

An architectural component may need multiple Go packages because it contains distinct implementation responsibilities that benefit from separate package boundaries.

That does not mean each package becomes a new architectural component. The component remains the owner of its exclusive implementation.

### 11.3 The library parent is not automatically a package

The central `library/` directory is a hierarchy. Its subdirectories can contain individual Go packages.

For example, `library/core/validation` and `library/infrastructure/nginx` would be separate packages in an illustrative Go repository. The parent `library/` does not need to be one giant package.

### 11.4 Avoid circular dependencies

A circular dependency appears when packages depend on one another in a cycle. It is often a sign that responsibilities or boundaries need to be reconsidered.

Do not solve a circular dependency by making an architectural component import another component's implementation. Reassess which responsibility is shared, which component owns it, and whether an appropriate library or API boundary is missing.

### 11.5 Interfaces are tools, not decorations

A Go interface can describe behaviour that a consumer needs. A type can satisfy that interface by implementing its required methods; it does not need to declare that relationship explicitly.

Interfaces can be useful when they create a clear boundary between a generic capability and a technology-specific implementation. They are not required for every package, function, or dependency.

Create an interface when it improves the design—for example, by allowing a generic configurator to use a technology-specific template provider without taking a hard dependency on that technology's internal implementation.

### 11.6 Do not confuse an API with a Go interface

These terms are related in some designs but do not mean the same thing:

- A **Go interface** is a language feature that describes a set of methods a type can provide.
- An **API** is a defined boundary through which a component offers operations or exchanges information.

A Go interface may be part of an API or an internal package design. It does not automatically provide the cross-component communication boundary required by this architecture.

---

## 12. When to create, split, merge, or extract a package

A good package boundary makes the responsibility easier to understand and maintain. It is not a contest to create the most packages or the fewest packages.

### 12.1 When should a package be created?

Create a package when there is a coherent responsibility that benefits from its own package boundary and ownership is clear.

For a shared library package, there must also be a genuine shared capability used by multiple architectural components.

### 12.2 When should a package be split?

Consider splitting a package when it contains responsibilities that are genuinely different and should evolve or be depended on separately.

For example, a package that mixes generic configuration logic with unrelated technology-specific operations may have unclear responsibilities. Separate the concerns when doing so makes ownership and dependencies clearer.

Do not split a package merely because it has many lines or files. Size alone does not prove that the package has the wrong boundary.

### 12.3 When should packages be merged?

Consider merging packages when they represent one cohesive responsibility, are tightly coupled for a good reason, and separating them adds navigation or dependency complexity without a meaningful boundary.

Do not merge packages simply to reduce directory count if their responsibilities are genuinely different.

### 12.4 When should functionality be extracted?

Extract functionality when a capability that was previously exclusive becomes genuinely needed by multiple architectural components.

Extraction should preserve the distinction between shared implementation and each component's own workflow.

### 12.5 A practical package-boundary test

Before creating, splitting, merging, or extracting, ask:

- Can I explain the package's responsibility in one or two sentences?
- Is its owner clear?
- Are its dependencies appropriate?
- Does it contain unrelated implementation?
- Is the proposed sharing real, or only an imagined future possibility?
- Would this boundary make a contributor's job easier to understand?

The correct structure is the one that expresses real responsibilities clearly—not the one with the most elaborate hierarchy.

---

# Part V — A complete hypothetical implementation walkthrough

This part follows one fictional project scenario from beginning to end. It brings the earlier principles together so a reader can see how decisions evolve during real development.

## 13. The scenario: two components and a growing requirement

Imagine a fictional infrastructure project with two architectural components:

- **Component A** owns the preparation of one environment.
- **Component B** owns a different responsibility and has its own workflow.

These are teaching labels only. They are not Ekathvam component names.

At the beginning, each component has its own exclusive implementation. The project also has a central `library/` hierarchy, but it does not yet contain every capability that might someday be shared.

### 13.1 First implementation: Component A has exclusive work

Component A needs a function to prepare its environment. No other component uses this operation.

The implementation belongs to Component A.

```text
example-project/
├── component-a/
│   └── prepare_environment.go
├── component-b/
└── library/
```

There is no need to move the operation into the library. It is exclusive to Component A.

### 13.2 A reusable validation requirement appears

Component A also needs to validate a configuration value. At first, only Component A needs that validation capability, so it is reasonable for the implementation to begin with Component A.

Later, Component B needs to apply the same validation rule to its own configuration.

The two components need the same capability. Component B is not asking Component A to do Component A's work; both need to perform the common validation as part of their separate workflows.

### 13.3 Extract the common validation capability

The shared validation responsibility is extracted into the library:

```text
example-project/
├── component-a/
│   └── prepare_environment.go
├── component-b/
│   └── component_workflow.go
└── library/
    └── core/
        └── validation/
            └── rules.go
```

Component A and Component B both use the shared validation package. Each component still decides when to validate and what to do with the result.

### 13.4 A technology integration appears

Now imagine Component A needs an NGINX-specific operation, while Component B does not.

That NGINX implementation stays with Component A because it has only one actual consumer.

```text
component-a/
└── nginx_support/
    └── setup.go
```

Do not create a library integration just because the technology is widely used in the software industry.

### 13.5 The technology capability becomes shared

Later, Component B also needs the same NGINX capability. The common technology-specific implementation is extracted into an appropriately classified library package.

```text
example-project/
├── component-a/
│   └── component_specific_work.go
├── component-b/
│   └── component_specific_work.go
└── library/
    ├── core/
    │   └── validation/
    │       └── rules.go
    └── infrastructure/
        └── nginx/
            └── shared_operation.go
```

Only the shared NGINX responsibility moves. Any NGINX-related workflow that is still exclusive to Component A remains with Component A.

### 13.6 Component B needs Component A to do something

A different requirement appears. Component B needs Component A to prepare the environment that A owns.

This is not code reuse. Component B should not import a function from Component A and call it directly.

Instead, Component B sends a request through Component A's API. Component A receives the request, performs the work internally, and returns the result through the API.

The interaction looks like this:

```text
Component B  ───── request through API ─────>  Component A
Component B  <──── result through API ───────  Component A
```

The components cooperate, but neither takes control of the other's private implementation.

### 13.7 What the finished example teaches

The example contains three different implementation decisions:

| Requirement | Correct location or mechanism | Reason |
|---|---|---|
| Prepare Component A's own environment | Component A | Exclusive responsibility |
| Validate a capability needed by both components | Shared library package | Genuine shared implementation |
| Provide an NGINX capability needed by both components | Shared technology integration | Genuine shared technology-specific implementation |
| Ask Component A to prepare its own environment | Component A's API | Component A owns and performs the requested work |

### 13.8 The entire lesson in one sentence

**Keep each component's exclusive work with its owner, move genuinely shared capabilities into the central library, and use APIs whenever one component needs another to perform its own responsibility.**

---

# Part VI — Mistakes, contribution, and maintenance

## 14. Common mistakes and their corrections

### 14.1 Mistake: Treating every package as an architectural component

**Why it is wrong:** A Go package is an implementation boundary. It does not automatically represent an independent architectural responsibility.

**Correct approach:** Identify the actual architectural owner first. Use packages to organise implementation within the ownership model.

### 14.2 Mistake: Treating every generic-looking function as shared

**Why it is wrong:** A function may look generally useful while still belonging to one component's exclusive responsibility.

**Correct approach:** Share functionality only when multiple actual architectural components genuinely need the capability.

### 14.3 Mistake: Duplicating a genuinely shared capability

**Why it is wrong:** Multiple copies can drift apart and require repeated maintenance.

**Correct approach:** Extract the common responsibility into one appropriately classified library package.

### 14.4 Mistake: Moving an entire package because one function is shared

**Why it is wrong:** Exclusive implementation can lose its natural owner, and the shared package can inherit unrelated responsibilities.

**Correct approach:** Extract only the shared, cohesive responsibility. Leave exclusive code with its component.

### 14.5 Mistake: Making one component import another component's internals

**Why it is wrong:** It creates direct implementation coupling and crosses the architectural ownership boundary.

**Correct approach:** If the capability is shared, extract it into the library. If the other component must perform its own work, use its API.

### 14.6 Mistake: Using an API merely to call a shared helper

**Why it is wrong:** It turns ordinary code reuse into unnecessary component communication.

**Correct approach:** Put genuinely shared implementation in the library and let each component use it.

### 14.7 Mistake: Using a shared library as a back door into another component

**Why it is wrong:** A library package should provide reusable capabilities, not expose or take over another component's private workflow.

**Correct approach:** Keep component-specific decisions and orchestration with the component. Use the API for requests to that component.

### 14.8 Mistake: Putting every technology integration in the library immediately

**Why it is wrong:** A technology package may have only one actual consumer.

**Correct approach:** Keep an exclusive integration with its owning component. Extract it when the same capability becomes genuinely shared.

### 14.9 Mistake: Creating abstractions for hypothetical future consumers

**Why it is wrong:** Speculative packages and interfaces add complexity before there is a real requirement.

**Correct approach:** Implement according to current responsibilities. Refactor when implementation reveals a genuine shared capability or a meaningful boundary.

### 14.10 Mistake: Making the library a miscellaneous folder

**Why it is wrong:** Unrelated implementation becomes difficult to discover, review, and maintain.

**Correct approach:** Classify shared packages by their actual responsibility and maintain clear package boundaries.

### 14.11 Mistake: Assuming an API means HTTP or JSON

**Why it is wrong:** An API is a contract, not a specific transport or serialization format.

**Correct approach:** Keep the contract concept separate from the technical mechanism selected to carry communication.

### 14.12 Mistake: Assuming example names must exist in Ekathvam

**Why it is wrong:** Teaching examples are not an implementation inventory.

**Correct approach:** Identify the real components and packages from the actual repository and architecture.

---

## 15. Contributor workflow

A contributor should be able to use this process whenever they add, move, or reorganise infrastructure code.

### 15.1 Before implementation

1. Understand the task and the responsibility being implemented.
2. Identify the actual architectural component that owns that responsibility, if it is component-specific.
3. Check whether the required capability is already available in the shared library.
4. Determine whether the requirement is code reuse or a request for another component to perform its own work.
5. Avoid creating speculative packages or abstractions.

### 15.2 During implementation

1. Keep exclusive functionality with its architectural owner.
2. Use existing shared library packages where they provide the required capability.
3. When genuine sharing is discovered, separate the shared responsibility from exclusive workflows.
4. Use an API for cross-component requests and communication.
5. Keep package dependencies consistent with the ownership boundaries.

### 15.3 When a package boundary needs to change

1. Explain the current responsibility of the package.
2. Explain why that responsibility needs to be split, merged, or extracted.
3. Identify which functionality is exclusive and which is shared.
4. Confirm that the proposed change does not make a library package depend on a component.
5. Confirm that no component imports another component's internal implementation.
6. Review the resulting structure from the perspective of a new contributor: can they understand where the code belongs and why?

### 15.4 Before submitting a contribution

A contribution should make the architecture easier to understand, not merely move code to a different location. The author should be able to explain the ownership and dependency decisions in ordinary language.

---

## 16. Maintainer responsibilities and review

Maintainers preserve the engineering model as the project grows. Their task is not to demand a fixed folder for every possible future feature. It is to ensure that actual implementation follows clear ownership and dependency boundaries.

### 16.1 Review architectural ownership

Check that exclusive functionality remains with the correct architectural component and that no package has quietly taken over another component's responsibility.

### 16.2 Review shared-library decisions

Check that library packages represent real shared capabilities, are classified by responsibility, and do not become dumping grounds for unrelated functions.

### 16.3 Review extraction quality

When functionality is extracted, check that only the shared responsibility moved and that exclusive workflows remain with their owners.

### 16.4 Review component communication

Check that components communicate through APIs when requesting another component's work. Reject direct imports or calls into another component's internal implementation.

### 16.5 Review dependencies

Check for reverse dependencies from the library into architectural components, circular package dependencies, and unnecessary coupling.

### 16.6 Preserve clarity for future contributors

A good change should help a new contributor answer:

- What is this package for?
- Who owns this behaviour?
- Why is this code in the library or in a component?
- How does this component communicate with another one?

If those questions become harder to answer after a change, the structure deserves another review.

---

# Reference section

## A. Quick decision reference

Use this table as a reminder after reading the explanations above.

| Situation | Decision |
|---|---|
| Functionality belongs to one component and is not shared | Keep it with that component |
| Two or more components genuinely need the same implementation capability | Extract the shared capability into the appropriate library package |
| Only one function in a package becomes shared | Extract that responsibility; do not move unrelated exclusive code |
| A technology integration is used by only one component | Keep it with that component |
| A technology capability becomes shared by multiple components | Extract the shared technology-specific capability into the appropriate library package |
| One component needs another component to perform work it owns | Communicate through the receiving component's API |
| Two components need the same reusable function | Use a shared library package, not a component-to-component call |
| A shared package needs to import an architectural component | Reassess the ownership and dependency design; the library must not depend on its consumer |
| A new package is proposed only for possible future reuse | Do not create it speculatively |

## B. The ownership model at a glance

```text
                 ACTUAL ARCHITECTURAL COMPONENTS

             Component A                 Component B
           exclusive work              exclusive work
                 |                           |
                 | uses                      | uses
                 |                           |
                 +-----------+---------------+
                             |
                             v
                   CENTRAL library/
                shared implementation packages


       Component A  <------ API ------>  Component B

       API communication is separate from library reuse.
```

This is a conceptual illustration. It is not a map of actual Ekathvam components or a protocol diagram.

## C. Glossary

**API:** A defined boundary through which a component exposes operations or exchanges information with another part of the system.

**Architectural component:** An actual part of the system architecture with its own responsibility and ownership of exclusive implementation.

**Architectural ownership:** The responsibility for maintaining and controlling the implementation that belongs to a particular architectural component.

**Capability:** An operation or set of related operations that provides a useful function to its consumers.

**Component-specific workflow:** The decisions and sequence of work that belong to one architectural component.

**Dependency:** A relationship in which one package or component relies on another for functionality or information.

**Exclusive functionality:** Implementation that belongs to one architectural component and is not a shared implementation capability.

**Extraction:** Separating a genuinely shared responsibility from its original location and placing it in an appropriate shared library package.

**Go package:** A group of Go source files compiled together as a package.

**Library:** In this guide, the central `library/` hierarchy containing reusable implementation packages.

**Shared functionality:** A capability genuinely needed by multiple architectural components and maintained in an appropriate library package.

**Technology integration:** Implementation that understands or interacts with a particular technology, such as a web server, container system, or DNS service.

## D. Questions to ask when uncertain

When a placement decision is not obvious, do not guess from the name of a function or directory. Work through these questions:

1. What exact responsibility does this code perform?
2. Which actual architectural component owns that responsibility?
3. Is the functionality exclusive, or do multiple components genuinely need the same capability?
4. If it is shared, what is the smallest cohesive responsibility that should be extracted?
5. Is the requirement actually a request for another component to perform its own work?
6. Does the proposed design preserve API boundaries and avoid direct component-to-component imports?
7. Can the resulting package and ownership be explained clearly to a new contributor?

If the answers remain unclear, ask the responsible maintainer to resolve the ownership before adding another package or dependency.

## E. Final principle

Ekathvam's infrastructure engineering model can be remembered through three rules:

1. **Exclusive work stays with its architectural owner.**
2. **Genuinely shared implementation lives in the appropriate package under the central `library/` hierarchy.**
3. **Architectural components communicate through APIs; they do not call one another's internal implementation.**

Everything else in this guide explains how to recognise and apply those three rules consistently.

---

**End of EEG 1 — Ekathvam Engineering Guide: Infrastructure Architecture, version `0.0.1-alpha`.**
