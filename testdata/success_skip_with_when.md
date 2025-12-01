# Success: skip sh step with when expression

```sh
exit 0
```

```sh when="false"
exit 1
```

```env
VAR=VALUE
```

```sh when='VAR == "VALUE"'
echo "Environment variable VAR is set to VALUE"
```

```sh when='VAR != "VALUE"'
echo "This step should be skipped because VAR is VALUE"
exit 2
```
