# Setting env variables

`env` steps can be used to set env variables. Each non-empty row of the `env` code block is passed to env of future test steps.

```sh
# Test for empty string
test -z "$example_var"
```

The above test will fail because $example_var is not set (or it is empty). Let's defined it in a env block.

```env
example_var=Example value

another_var=Another value
```

The variable is now defined in following test steps.

```sh
# Test for non-empty string
test -n "$example_var"
test -n "$another_var"
```

The environment variables defined in `env` code-block can be overridden by defining environment variable with the same name in the `mdtest` command using `-e`/`--env` parameter, for example `--env TARGET=TEST`.