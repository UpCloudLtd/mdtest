# Skipping sh steps with when option

By using `when` option on `sh` code-blocks, you can conditionally skip execution of the step. Then `when` option expects a [Common Expression Language](https://cel.dev/) (CEL) expression as its value. The CEL expression given in the `when` option is evaluated before executing the script. If the expression evaluates to false, the step is skipped. For example: `` ```sh when="false" ``.

```sh when="false"
echo "This step will be skipped"
exit 1
```

The CEL expression has access to environment variables, such as those defined in `env` code-blocks. Let's define an environment variable and use it in a `when` expression.

```env
CAR=coupe
```

The following step defines `when` expression that checks the value of `car` variable: `` ```sh when='CAR != "coupe" ``. Thus, by default, the following code-block is skipped.

```sh when='CAR != "coupe"'
echo "The car is not a coupe"
exit 1
```

However, if the user changes the value of `car` by using `--env` (e.g., `--env CAR=minivan`) command-line argument to different value, the code-block would be executed.
