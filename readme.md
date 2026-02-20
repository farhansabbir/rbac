# Go RBAC Library

A high-performance, thread-safe, and distinct Role-Based Access Control (RBAC) library for Go.

Designed with an **AWS IAM** and **Kubernetes** inspired architecture, this library separates data entities (`User`, `Profile`, `Rule`) from the enforcement logic (`Gatekeeper`) and state management (`Controllers`).



## 🚀 Key Features

* **Thread-Safe Concurrency:** Built-in `sync.RWMutex` protection for all state changes, allowing safe concurrent access in high-load environments.
* **Singleton Controller Pattern:** Centralized state management replaces global variables, ensuring a single source of truth.
* **High-Performance Evaluation:**
    * **Bitwise Verbs:** Permissions (Read, Write, etc.) are evaluated using bitwise operations for $O(1)$ speed.
    * **Indexed Lookups:** Rules are sharded by `ResourceType`, skipping 90% of irrelevant rules during checks.
* **AWS-Style Logic:** Implements **"Explicit Deny Overrides Allow"** logic.
* **Event-Driven:** Integrated non-blocking event loops for auditing and logging state changes.
* **Wildcard Support:** Supports `*` for resource IDs and verbs.

---

## 📦 Installation

```bash
go get [github.com/farhansabbir/rbac](https://github.com/farhansabbir/rbac)
```

## 🏗 Architecture
1. Data Models (lib/)

The data layer is composed of three primary entities. All entities implement the Resource interface.

User: The identity (Principal). Holds a list of Profiles.

Profile: A collection of policies (Rules). Acts as the bridge between Users and Rules.

Rule: The atomic logic unit.

Verbs: Bitmask (VerbRead | VerbList).

Actions: Allow, Deny, or AllowAndForward.

Targets: Defines ResourceType (e.g., URL, Project) and ResourceID.

2. State Management (controllers/)

We avoid global variables by using a Singleton Controller.

Thread Safety: Every map (User store, Rule store) is protected by sync.RWMutex.

Event Loop: A background goroutine listens on buffered channels for events (Creation, Deletion) to handle logging without blocking the API response.

3. The Engine (Gatekeeper)

The IsRequestAllowed method is the heart of the library. It is stateless and relies on the RequestContext.

Evaluation Flow:

Implicit Deny: Default state is False.

Filter: Retrieve only rules matching the requested ResourceType.

Evaluate:

Check ID Match (Exact or Wildcard *).

Check Verb Match (Bitwise &).

Decision Logic:

✅ Allow: Sets allowed = true but continues looping.

❌ Deny: Returns false IMMEDIATELY (stops looping).

Final Result: Returns true only if allowed == true AND no Deny rules were triggered.

## 🔮 Roadmap
[ ] Rule Forwarding: Full implementation of ActionAllowAndForwardToNextRule to chain policies.

[ ] Attribute-Based Access Control (ABAC): utilize the Attributes map in RequestContext for finer-grained control (e.g., Owner checks).

[ ] Persistence Layer: Add interfaces for SQL/Redis storage in the Controllers.
