# Timeouts

Use the `timeout` option to set a maximum execution time for a step. If the step exceeds this time, it will be terminated and marked as failed. For example: `` ```sh timeout="1s" ``.

```sh timeout=1s
echo "This step will fail after 1 second"
sleep ${SLEEP_TIME:-0}
```

To set a timeout for the entire test suite, use the `--timeout` flag when running the tests. For example: `mdtest --timeout=5s examples/` will terminate the entire test suite if it runs longer than 5 seconds.
