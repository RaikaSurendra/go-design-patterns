# Contributing to Go Design Patterns

Thanks for your interest in contributing! This is an actively maintained collection of Go design and application patterns. Whether you're fixing a typo, improving an existing pattern, or proposing a new one, your help is appreciated.

## Ways to Contribute

- **Improve existing patterns** — better examples, clearer explanations, bug fixes in code snippets
- **Add new patterns** — propose patterns not yet covered (open an issue first to discuss)
- **Fix issues** — check the [open issues](https://github.com/RaikaSurendra/go-design-patterns/issues) for things to work on
- **Review** — read through patterns and report anything unclear or incorrect

## Getting Started

1. Fork this repository
2. Create a feature branch: `git checkout -b <category>/<pattern-name>`
3. Make your changes
4. Commit following the message guidelines below
5. Push to your fork and open a pull request

## Pull Request Guidelines

- Make an individual pull request for each suggestion.
- Choose the corresponding patterns section for your suggestion.
- List, after your addition, should be in lexicographical order.
- Ensure Go code snippets compile and are idiomatic — use `gofmt` style.
- Keep examples minimal and focused. Avoid unnecessary boilerplate.

## Commit Messages Guidelines

- The message should be in imperative form and uncapitalized.
- If possible, please include an explanation in the commit message body.
- Use the form `<pattern-section>/<pattern-name>: <message>`
  - e.g. `creational/singleton: refactor singleton constructor`
  - e.g. `behavioral/visitor: fix interface example`

## Pattern Template

Each pattern should have a single markdown file containing the important part of the implementation, the usage and the explanations for it. This is to ensure that the reader doesn't have to read a bunch of boilerplate to understand what's going on and the code is as simple as possible and not simpler.

Please use the following template for adding new patterns:

```markdown
# <Pattern-Name>
<Pattern description>

## Implementation

```go
// Go implementation here
```

## Usage

```go
// Usage example here
```

## Rules of Thumb
- Bullet points with practical advice
```

## Code Style

- All Go code should follow standard `gofmt` formatting
- Use meaningful variable and function names
- Include comments only where the logic isn't self-evident
- Prefer standard library packages over external dependencies
- Use Go generics where they improve clarity (Go 1.18+)

## Questions?

Open an issue if you're unsure about anything. We're happy to help!
