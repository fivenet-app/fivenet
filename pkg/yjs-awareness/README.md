# Go Yjs Awareness Protocol Module

This Go module provides encoding, decoding, and modification of Yjs awareness update packets. It allows you to create awareness update bytes from state data, decode updates to read their contents, and modify updates (e.g., to filter or alter fields) programmatically.

## Installation

Copy the `awareness.go`, `encoder.go`, and `decoder.go` files into a Go module (for example, a folder named `awareness`). Then import the package in your code:

```go
import "path/to/awareness"
```

## Usage

### Encoding an Awareness Update

Create an Awareness instance, set some client state, and encode an update for those clients:

```go
aw := awareness.NewAwareness()
aw.SetState(1, map[string]interface{}{"name": "Alice", "cursor": map[string]int{"x": 100, "y": 200}})
aw.SetState(2, map[string]interface{}{"name": "Bob"})  // another client

updateBytes, err := aw.EncodeUpdate([]uint64{1, 2})
if err != nil {
    panic(err)
}
// updateBytes now contains a Yjs-compatible awareness update for clients 1 and 2.
```

### Decoding an Awareness Update

Given an awareness update []byte, you can decode it to inspect the states:

```go
states, err := awareness.DecodeAwarenessUpdate(updateBytes)
if err != nil {
    panic(err)
}
for clientID, state := range states {
    fmt.Printf("Client %d state: %v\n", clientID, state)
}
```

This will output each client's awareness state (as a nested Go interface{} structure). If a state is nil, that client was marked offline in the update. If a value is awareness.Undefined, it represents an undefined value from Yjs.
Modifying an Awareness Update

Use ModifyAwarenessUpdate to transform the states within an update. For example, suppose we want to sanitize user names in the update (e.g., mask them):
```go
modifyFn := func(state interface{}) interface{} {
    // If state is a map and has a "name" field, mask it
    if m, ok := state.(map[string]interface{}); ok {
        if _, hasName := m["name"]; hasName {
            newState := make(map[string]interface{})
            for k, v := range m {
                if k == "name" {
                    newState[k] = "***"
                } else {
                    newState[k] = v
                }
            }
            return newState
        }
    }
    return state // no change for other states
}

sanitizedUpdate, err := awareness.ModifyAwarenessUpdate(updateBytes, modifyFn)
if err != nil {
    panic(err)
}
// sanitizedUpdate now contains the same awareness update with all names replaced by "***".
```

This updated byte sequence can then be applied or forwarded just like a normal awareness update.
