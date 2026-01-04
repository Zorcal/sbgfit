# Agent Guidelines

**IMPORTANT**: These guidelines apply to ALL Go development work in this project. Agents MUST reference this document when:

- Writing or modifying Go code
- Creating or updating tests
- Completing any Go-related tasks
- Before marking any Go work as complete

## Go Style Guidelines

### Variable Scope and Declaration

- **DO** use short variable declarations with assignment (`if myVar := ...`) when there's no else branch to minimize scope
- **DO** prefer minimal variable scope over broader declarations
- **Examples:**

  ```go
  // Good - minimal scope, no else branch
  if result := processData(); result != nil {
      return result.Value
  }

  // Good - minimal scope for error checking
  if err := validate(input); err != nil {
      return fmt.Errorf("validation failed: %w", err)
  }

  // Avoid - unnecessarily broad scope
  result := processData()
  if result != nil {
      return result.Value
  }
  ```

## Post-Completion Requirements

### Linting

- **MUST** run `golangci-lint run` after completing all TODO list items
- **MUST** fix any linting issues before marking task as complete
- **DO** run linting from the backend directory

## Go Testing Best Practices

This document outlines testing guidelines for agents working on Go code, based on the [Go Testing Guidelines](https://go.dev/wiki/TestComments).

### Test Comment Guidelines

#### 1. Subtest Naming

- **DO** choose human-readable names for subtests
- **DO** ensure names remain useful after escaping
- **DO** include input details in test logs or failure messages
- **Example:**
  ```go
  t.Run("empty input", func(t *testing.T) { ... })
  t.Run("valid user data", func(t *testing.T) { ... })
  ```

#### 2. Comparison Strategies

- **DO** compare full structures instead of individual fields
- **DO** use `cmp.Diff` or `cmp.Equal` for complex comparisons
- **DO** prefer comparing semantic meaning over exact output
- **Example:**
  ```go
  if diff := cmp.Diff(got, want); diff != "" {
      t.Errorf("YourFunc() diff mismatch (-got +want):\n%s", diff)
  }
  ```

#### 3. Error Reporting Best Practices

- **DO** include function name in failure messages
- **DO** print actual (got) value before expected (want) value
- **DO** format like: `YourFunc(%v) = %v, want %v`
- **DO** include input details in error messages
- **Example:**
  ```go
  if got != want {
      t.Errorf("Sum(%d, %d) = %d, want %d", a, b, got, want)
  }
  ```

#### 4. Test Design

- **DO** use table-driven tests for similar test cases
- **DO** use separate test functions for significantly different logic
- **DO** separate success and error cases into different test functions using `TestXyz` and `TestXyz_error` naming convention
- **DO NOT** use `wantErr` boolean flags inside table tests - split into separate functions instead
- **DO NOT** add extra whitespace between table declaration and for loop
- **DO** keep tests running even after initial failures
- **DO** prefer `t.Error` over `t.Fatal` to continue testing
- **DO** use `t.Fatal` when the test cannot meaningfully continue (e.g., failed setup)
- **Examples:**

  **Success cases:**

  ```go
  func TestCalculate(t *testing.T) {
      tests := []struct {
          name string
          in int
          want  int
      }{
          {"zero", 0, 0},
          {"positive", 5, 25},
          {"negative", -3, 9},
      }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              got := Calculate(tt.in)
              if got != tt.want {
                  t.Errorf("Calculate(%d) = %d, want %d", tt.in, got, tt.want)
              }
          })
      }
  }
  ```

  **Error cases:**

  ```go
  func TestCalculate_error(t *testing.T) {
      tests := []struct {
          name string
          in int
      }{
          {"too large", math.MaxInt},
          {"invalid range", -1000},
      }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              _, err := Calculate(tt.in)
              if err == nil {
                  t.Errorf("Calculate(%d) error = nil, want error", tt.in)
              }
          })
      }
  }
  ```

#### 5. Error Semantics

- **DO NOT** use strict string comparisons for error types
- **DO** focus on testing error presence/absence
- **DO** use programmatic error structures when possible
- **DO** check against substrings when no sentinel error is available
- **Example:**

  ```go
  // Good - check for error presence
  if err == nil {
      t.Errorf("YourFunc() error = nil, want error")
  }

  // Good - check error type
  var target *MyError
  if !errors.As(err, &target) {
      t.Errorf("YourFunc() error = %v, want MyError", err)
  }

  // Good - substring check when no sentinel error available
  if err == nil {
      t.Errorf("YourFunc() error = nil, want error")
  }
  if wantSubstr := "connection failed"; !strings.Contains(err.Error(), wantSubstr) {
      t.Errorf("YourFunc() error = %v, want error containing %q", err, wantSubstr)
  }

  // Avoid - exact string comparison
  if err.Error() != "exact error message" { ... }
  ```

### Summary

The overarching principle is to write clear, informative, and robust tests that provide meaningful feedback. Focus on:

- Native Go constructs over external libraries
- Clear, descriptive error messages
- Comprehensive structure comparisons
- Human-readable test organization
- Semantic correctness over exact matches
