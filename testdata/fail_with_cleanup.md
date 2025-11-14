# Fail: failing step, skipped step, and (failing) cleanup step

Failing step:

```sh
exit 4
```

Skipped step:

```sh
exit 3
```

Failing cleanup steps:

```sh cleanup=true
exit 2
```

```sh cleanup
exit 1
```
