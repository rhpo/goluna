# GOLuna 🌙, The Go Port of Luna

**GOLuna** is the Go implementation of [Luna](https://github.com/rhpo/luna), a high-level programming language originally written in TypeScript. It brings Luna's expressive syntax and scripting-oriented design to a Go-based implementation.

<div align="center">

  <img src="/assets/logo-styles/luna-noborders.png" width="140" alt="Luna logo" />

  <h3>Luna</h3>

  <p>
    A high-level programming language designed for expressive syntax,
    readable code, and productive scripting.
  </p>

  <p>
    <a href="https://github.com/rhpo/luna/tree/main/docs"><strong>Original Documentation</strong></a>
    ·
    <a href="https://rhpo.github.io/luna/examples/">Examples</a>
    ·
    <a href="https://github.com/rhpo/luna/issues">Report a Bug</a>
  </p>

  <p>
    <a href="https://github.com/rhpo/luna/blob/main/LICENSE">
      <img src="https://img.shields.io/github/license/rhpo/luna?style=flat-square" alt="MIT License" />
    </a>
    <a href="https://github.com/rhpo/luna">
      <img src="https://img.shields.io/badge/original%20implementation-TypeScript-blue?style=flat-square" alt="Original implementation: TypeScript" />
    </a>
    <img src="https://img.shields.io/badge/Go-implementation-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go implementation" />
  </p>

</div>

---

## Table of Contents

- [About the Project](#about-the-project)
- [Why Go?](#why-go)
- [Language Features](#language-features)
- [Architecture](#architecture)
- [Examples](#examples)
- [Getting Started](#getting-started)
- [Project Structure](#project-structure)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)
- [Author](#author)

## About the Project

Luna is a high-level programming language built around a simple idea: programming should be expressive without sacrificing readability.

The original implementation was written in TypeScript and provided an interpreter, a runtime environment, and language features designed for general-purpose programming and automation.

**GOLuna is the Go port of that implementation.** The project aims to bring Luna's language design to Go while providing a foundation for further development of its interpreter and runtime.

The original language includes concepts such as:

- Procedural, functional, and object-oriented programming styles.
- Functions, anonymous functions, and lambda expressions.
- Variables, constants, and explicit variable actions.
- Reactive values that respond to changes in their dependencies.
- Module imports, aliases, and selective imports.
- Code embedding and direct file execution.
- Asynchronous programming capabilities.

The availability and implementation of individual features may differ between the original TypeScript implementation and GOLuna.

## Why Go?

Rewriting Luna in Go provides an opportunity to explore a different implementation strategy while retaining the language's identity.

- **A different implementation:** Move the language runtime from TypeScript to Go.
- **A native Go ecosystem:** Develop and maintain the implementation using Go's tooling and standard library.
- **A foundation for optimization:** Profile and improve the interpreter and runtime as the project evolves.
- **Maintainability:** Keep the implementation organized around Go's language conventions and development practices.

GOLuna is a language implementation project, not simply a TypeScript-to-Go transpiler.

## Language Features

The original Luna language explores several programming concepts.

| Feature | Description |
|---|---|
| Functions | Named functions with a concise declaration syntax |
| Lambda expressions | Anonymous functions and inline expressions |
| Constants | Immutable bindings using the `const` action |
| Reactivity | Values derived from dependencies that update when those dependencies change |
| Modules | Import modules and selectively import or rename their components |
| Embedding | Include code from other files during preprocessing |
| Direct execution | Execute a file without importing its declarations |
| Asynchronous expressions | Support for asynchronous lambda expressions in the original implementation |

Refer to the implementation and its tests to determine which features are currently supported in GOLuna.

## Architecture

The original Luna implementation follows a language-processing pipeline:

1. **Tokenization** : Convert source code into tokens.
2. **Parsing** : Build an abstract syntax tree (AST) from the token stream.
3. **Preprocessing** : Process supported embedded code and other preprocessing directives.
4. **Evaluation** : Interpret the resulting syntax tree.
5. **Runtime** : Represent values and execute language operations.

<div align="center">

  <img src="https://github.com/rhpo/Luna/assets/69460661/6f390f7f-8a43-41c9-9765-6ec50fbf8b47" width="700" alt="Original Luna architecture diagram" />

  <p><sub>Architecture diagram from the original Luna project.</sub></p>

</div>

The Go implementation may differ internally. Its architecture should be documented according to the current source code.

## Examples

The following snippets illustrate the syntax of the original Luna language. They are reference examples, not a guarantee that every feature is implemented identically in GOLuna.

### Functions

```rust
fn read book {
    print("Reading {book}...")
}

read("Luna Docs")
```

### Default arguments

```rust
fn sum x=(1) y=(1) {
    if isdef x && isdef y: x + y
    else null
}

print(sum(2, 4))
```

Expected output:

```text
6
```

### Constants

```rust
pi: const = 3.141

fn area r {
    pi * r ** 2
}

print(area(2))
```

### Anonymous functions

```rust
sum = lambda x y {
    x + y
}

print(sum(2, 4))
```

### Reactive values

```rust
x = 1

doubled: react<x> = x * 2

print(doubled)

x = 2

print(doubled)
```

In the original implementation, reactive values can be recalculated when their dependencies change.

### Module imports

```rust
use "maths.lnx" as math

use (cos, sin as Sinus, tan, PI as pi) from "math.lnx"
```

### Embedding and direct execution

```rust
embed "file.ln";

tap "exec/this/code.ln";
```

These directives belong to Luna's original language design. Consult the Go implementation for their current support and behavior.

## Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) installed on your system.
- Git, if you intend to clone the repository.

### Clone the repository

```bash
git clone <GOLUNA_REPOSITORY_URL>
cd <GOLUNA_DIRECTORY>
```

### Build and run

The exact commands depend on GOLuna's current entry point and module configuration.

For a standard Go project with a root `main` package, the typical workflow is:

```bash
go build ./...
go run .
```

If the executable is located in a `cmd` directory, use the corresponding package path instead.

### Verify the installation

Run the project's documented test suite:

```bash
go test ./...
```

Replace the generic commands above with the actual build and execution workflow if GOLuna uses a different entry point.

## Project Structure

The implementation should be explored through its source code. Typical components of an interpreter include:

- **Lexer:** Tokenizes source code.
- **Parser:** Converts tokens into syntax trees.
- **AST:** Represents language expressions and statements.
- **Interpreter:** Evaluates statements and expressions.
- **Runtime:** Implements values, variables, functions, and built-in operations.
- **CLI:** Provides the interface for running Luna programs, if included.

These are conceptual interpreter components; their names and organization depend on the actual Go codebase.

## Documentation

- [Original Luna repository](https://github.com/rhpo/luna)
- [Original language documentation](https://github.com/rhpo/luna/tree/main/docs)
- [Original examples](https://rhpo.github.io/luna/examples/)
- [Original API documentation](https://github.com/rhpo/luna/tree/main/api)
- [Original issue tracker](https://github.com/rhpo/luna/issues)

Documentation specific to GOLuna should be added to this repository as the Go implementation evolves.

## Contributing

Contributions, bug reports, and suggestions are welcome.

Before submitting changes:

1. Review the existing source code and project conventions.
2. Run the Go test suite.
3. Add tests for new functionality or bug fixes.
4. Document any changes to language behavior.
5. Submit a pull request describing the change.

For language-design discussions and issues concerning the original implementation, see the [Luna repository](https://github.com/rhpo/luna).

## License

Luna's original implementation is distributed under the MIT License.

See the applicable `LICENSE` file in the GOLuna repository to determine the license governing this Go port and its contributions.

## Author

**Ramy Hadid** : [@rhpo](https://github.com/rhpo)

- [GitHub](https://github.com/rhpo)
- [LinkedIn](https://www.linkedin.com/in/ramy-hadid-15aa70243/)

---

Originally created as Luna, a TypeScript-based programming language, and subsequently rewritten in Go as GOLuna.
