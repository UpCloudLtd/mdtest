# Non-zero exit codes

By default, test step will fail if script returns with non-zero exit code. To test scripts that return failing exit codes specify `exit_code` option on the scripts start line. For example: `` ```sh exit_code=3 ``.

```sh exit_code=3
exit 3
```

This test case will thus succeed, as script exits with the exit code specified in the `exit_code` option.