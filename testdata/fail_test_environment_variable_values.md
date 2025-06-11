# Fail: test environment variable values

```env
berry=banana
fruit=apple
```

```sh
test "$berry" = "strawberry"
test "$fruit" = "orange"
```
