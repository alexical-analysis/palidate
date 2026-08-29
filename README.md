# Palidate Reflect Workshop

Palidate is a fun little struct "validation" library that uses struct tags to create validation rules.
The rules should work as follows:
- "required" means that the field can not have it's zero value
- "max" is either the maximum valid value for an int type (uint, uint8-64, int, int8-64) or the max lenght for a string.
- "min" is either the minimum valid value for an int type (uint, uint8-64, int, int8-64) or the min lenght for a string.
- "regex" matches the provided pattern against a string and is valid only if a match is found. The regex 
    must be enclosed by single quotes. 2 single quotes in a row become an escaped quote.

Currently, these rules are not working as the implementation has been removed.
By adding the implementation back in you'll learn now to use Go's `reflect` package to inspect structs,
read struct tags, and validate field values.

## Workshop instructions

Implement the missing functionality in **`palidate.go`**. The other files provide
supporting types, parsing behavior, error handling, and tests.
You can leave these files unchanged.

Validation your implementation gradually by running the 3 tests suites provided in the `palidate_test.go`
test file.

```bash
go test -run '^TestStructRequired$'
```
will validate that your code correctly handles the `min` and `max` rules.

```bash
go test -run '^TestStructRegex$'
```
will validate that your code correctly handles the `required` rule.

```bash
go test -run '^TestStructMaxMin$'
```
will validate that your code correctly handles the `regex` rule.

```bash
go test ./...
```
will run all the validation tests including a new test suite that will validate all your validation
rules work together.

The tests compare error messages exactly, so rely on the pre-defined error types at the top of the 
`palidate.go` file

For reference, see the Go `reflect` package documentation:

<https://pkg.go.dev/reflect>
