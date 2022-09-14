# Multiple commands

A test file can contain multiple code blocks:

```sh
sleep ${SLEEP_TIME:-3}
```

The test progress logging will include information on which steps is currently being executed. E.g., `(Step 1 of 2)`.

```sh
sleep ${SLEEP_TIME:-3}
```

The test case succeeds if all of its steps succeed. I.e., if one of test steps fails the whole test case is marked as failed.
