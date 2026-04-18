# Engineering Mandates

## Context & Vision
We are currently porting an existing 2D engine from **raylib-go** to **ebiten**. The goal is to reuse as much of the original codebase as possible, with surgical rewrites and simplifications where necessary to adapt to ebiten's architecture.

## General Principles
- **No Tests:** We do not write tests for this project.
- **Explicit Variable Declaration:** Use `var` instead of `:=` for all variable declarations. The only exception is loop initializers (e.g., `for i := 0; ...`) where `var` is not supported by Go syntax.
- **Value Types Only:** Always use value types. Never use pointers unless absolutely required by an external library interface (e.g., Ebiten).
- **Naming Conventions:** Avoid single-letter parameter names in the public API. Use descriptive names like `value`, `index`, or `text`. Shorter variants are only acceptable for unexported functions or internal implementation details.
- **Minimal Heap Allocations:** Keep heap allocations to an absolute minimum. Avoid any code patterns that cause variables to escape to the heap.

## Architecture
- **State Management:** Follow the **BigFatStruct/LargeArraysOfThings** paradigm. 
    - All serializable engine state must reside in a single large struct (`internal.State` or `internal.Data`).
    - This struct may contain other nested structs.
- **Global Throwaway State:** Any non-serializable or transient state (e.g., input caches, delta timers) should be kept as global variables in the relevant internal package.
- **Engine Encapsulation:** The engine's core state must be hidden within the `internal` package. It should be accessible to all internal packages but never exposed to the end user or game logic.

## API Design
- **Primitive Exports:** The public API must never expose maps, structs, or pointers. 
- **Allowed Types:** Only primitive types (int, float32, bool, byte, etc.), strings, and fixed-size arrays are allowed in the public API.
- **Slices in API:** Avoid slices in the public API unless it can be guaranteed they do not cause heap allocations (which is rare). Prefer fixed-size arrays.
- **Type Aliases:** Custom types based on primitives (e.g., `type Mode byte`) are acceptable as long as they follow the primitive-only rule.
