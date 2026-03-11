# Fail: test environment variable values

```env
berry=banana
fruit=apple
berry_fruit="${berry}-${fruit}"
```

```sh
test "$berry" = "strawberry"
test "$fruit" = "orange"
test "$berry_fruit" = "strawberry-orange"
test "$berry_fruit" = "strawberry-orange"
```
