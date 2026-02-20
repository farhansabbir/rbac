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
```bash 

## ⚡ Quick Start
Here is how to wire up the system, create a policy, and enforce it.

```go
package main

import (
	"fmt"
	"[github.com/farhansabbir/rbac/controllers](https://github.com/farhansabbir/rbac/controllers)"
	"[github.com/farhansabbir/rbac/lib](https://github.com/farhansabbir/rbac/lib)"
)

func main() {
	// 1. Initialize the Singleton Controller (Starts Event Loops)
	ctrl := controllers.GetController()
	defer ctrl.Stop() // Graceful shutdown

	// 2. Get Sub-Controllers
	userCtrl := ctrl.GetUserController()
	profCtrl := ctrl.GetProfileController()
	ruleCtrl := ctrl.GetRuleController()

	// 3. Create a Rule: "Allow Read on All Profiles"
	readRule := ruleCtrl.CreateRule(
		"allow-read-profiles",
		"Allows reading any profile",
		lib.ResourceIDAll,      // Target ID: "*"
		lib.VerbRead,           // Action: Read
		lib.ActionAllow,        // Effect: Allow
	)
	// Important: Set the Target Resource Type
	readRule.SetTargetResourceType(lib.ResourceTypeProfile)

	// 4. Create a Profile and Attach Rule
	adminProfile := profCtrl.CreateProfile("AdminProfile", "Read-Only Admin")
	adminProfile.AddRule(readRule)

	// 5. Create a User and Attach Profile
	john := userCtrl.CreateUser("John Doe", "DevOps", "john@example.com")
	john.AddProfile(adminProfile)

	// 6. Initialize the Gatekeeper (Enforcer)
	// Note: In a real app, pass the controller to the Gatekeeper
	gk := lib.NewGatekeeper() 

	// 7. Create a Request Context (The "Ticket")
	// John wants to READ Profile ID 123
	ctx, _ := lib.NewRequestContext(
		john.GetResourceID(),
		lib.ResourceTypeProfile,
		123,
		lib.VerbRead,
		nil,
	)

	// 8. Ask: Is Request Allowed?
	allowed, err := gk.IsRequestAllowed(ctx)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Printf("Request Allowed: %v\n", allowed)
}
```go

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

🧩 Advanced Usage
Bitwise Verbs

You can combine verbs to create complex permissions in a single rule.

```go
// Allow Read AND List AND Execute
complexVerb := lib.VerbRead | lib.VerbList | lib.VerbExecute

rule.UpdateVerb(complexVerb)
Soft Deletion

Deleting a user does not remove them from memory immediately. It sets a DeletedAt timestamp.

```Go
userCtrl.DeleteUser(userID) // User remains in map, but IsActive() returns false
The Gatekeeper automatically rejects requests from inactive users.

🧪 Running Tests
The library includes a comprehensive test suite covering wildcard matching, deny-overrides logic, and bitmask checks.

```Bash
# Run all tests in the lib package
go test ./lib -v
Test Coverage:

TestGatekeeper_DenyOverridesAllow: Ensures security is paramount.

TestGatekeeper_WildcardIDMatching: Validates * vs specific IDs.

TestGatekeeper_VerbBitmaskMatching: Validates bitwise logic.

## 🔮 Roadmap
[ ] Rule Forwarding: Full implementation of ActionAllowAndForwardToNextRule to chain policies.

[ ] Attribute-Based Access Control (ABAC): utilize the Attributes map in RequestContext for finer-grained control (e.g., Owner checks).

[ ] Persistence Layer: Add interfaces for SQL/Redis storage in the Controllers.
