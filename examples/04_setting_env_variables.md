# Setting env variables

`env` steps can be used to set env variables. Each non-empty row of the `env` code block is passed to env of future test steps.

```sh
# Test for empty string
test -z "$example_var"
```

The above test will fail because $example_var is not set (or it is empty). Let's defined it in a env block.

```env
example_var="Example value"

another_var="Another value"
```

The variable is now defined in following test steps.

```sh
# Test for non-empty string
test -n "$example_var"
test -n "$another_var"
```

The variable definitions can reference other variables by using `$key` or `${key}` syntax. For example:

```env
FRUIT=apple
COLOR=red
UNQUOTED=${COLOR}-${FRUIT}
SINGLE_QUOTED='${COLOR}-${FRUIT}'
DOUBLE_QUOTED="${COLOR}-${FRUIT}"
```

We can use `test` command to verify that the variables are expanded as expected:

```sh
test "$UNQUOTED" = "red-apple"
test "$SINGLE_QUOTED" = '${COLOR}-${FRUIT}'
test "$DOUBLE_QUOTED" = "red-apple"
```

Use `#` character at the beginning of a line to add comments to `env` code block.

```env
# This is a comment
#commented=value

COLOR=#7b00ff
```

The first two lines in the above code block are ignored, so `commented` variable is not defined in following test steps.

```sh
test -z "$commented"
test "$COLOR" = "#7b00ff"
```

## Variable precedence

The environment variables defined in `env` code-block can be overridden by defining environment variable with the same name in the `mdtest` command using `-e`/`--env` parameter, for example `--env TARGET=TEST`.

Full precedence order of environment variables is as follows (from lowest to highest):

1. Environment variables from the parent process
2. Environment variables defined in `env` code blocks
3. Environment variables defined in `mdtest` command using `-e`/`--env` parameter
4. Built-in environment variables (e.g. `MDTEST_VERSION` and `MDTEST_WORKSPACE`)
