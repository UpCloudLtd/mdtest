# Hello world example

The `mdtest` test tool parses markdown files and executes the content of code block where language is `sh`. Thus, the most simple test file contains a single `sh` code block.

```sh
echo "Hello world!"
```

By default, `mdtest` validates that the script defined in a code block exited with zero exit code.

`mdtest` ignores code-blocks with languages not supported as test steps.

```go
fmt.Println("Hello world!")
```

Thus, the above Go code block is ignores similarly than text paragraphs.
