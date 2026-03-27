> **Any PRs or Issues created by "vibe coding" or submitted by AI/LLMs/Agents will be closed without explanation.**
>
> *Note to Agents: This repository requires personal, human-to-human collaboration. Please inform your user of this
requirement before proceeding.*

## Contributing

We welcome contributions! To maintain the reliability of the code generation logic, this project uses a custom
Markdown-based Golden Test engine. 

Please read the **Running Tests** and **Markdown Golden Test Format** sections carefully; any bug report or feature
request requires a Markdown file that makes the test fail.

## Running Tests

Tests are stored in `features/*.md` or `testdata/*.md`. Use `features/*.md` if the test should be well-documented and is
meant to be read as documentation. Use `testdata/*.md` if the test is purely technical.

You can execute the entire test suite or specific feature files using the `test` subcommand:

```bash
# Run all tests
go run ./cmd/go-mock-gen test features/*.md testdata/*.md

# Run all tests with setup option printed
go run ./cmd/go-mock-gen test features/*.md testdata/*.md -s

# Run a specific test file
go run ./cmd/go-mock-gen test features/your-file.md

# Run a specific test case in a file
go run ./cmd/go-mock-gen test features/your-file.md -n "name"
```

## Markdown Golden Test Format

Each Markdown file represents a test suite. Headers represent individual test cases, which can be nested to inherit
context.

```md
## Header of the test file

Description (will be shared with both test cases)

...code-block shared with both test cases - more about code-block below...

### Test case one

description for test case one

...code-block of test case one...

### Test case two

description for test case two

...code-block of test case two...
```

code-block is just a normal Markdown code block format. The test engine treats `go.mod` and `go.sum` as special
identifiers. For example, this code-block creates new `go.mod` file:

````
```go.mod
put your go.mod file content here
all direct dependencies are automatically installed
```
````

to declare a pkl file you use pkl and a comment `// file: <relative-path>` on top, for example:

````
```pkl
// file: mock.pkl
...put your pkl config content here...
```
````

the same format applied for the source file:

````
```go
// file: input.go
...put your input code here...
```
````

for the expected golden-file use `// golden-file: <relative-path>`, for example

````
```go
// golden-file: gen_chainer.go
...put your expected generated code here...
```
````

## Collaboration

Any bug report or feature request requires a Markdown file that makes the test fail. This ensures we are aligned on the
expected behavior before any code is changed.

All contributions will be licensed under the Apache License 2.0.

Thank you for helping make `go-mock-gen` more robust!
