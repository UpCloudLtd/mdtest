# Writing files

When a code-block contains a `filename` option, the content of the code block is written to the file specified as a parameter for the option. For example, `` ```py filename=js_sum.py ``.

```py filename=js_sum.py
def number_or_string(value):
    try:
        num = None
        num = float(value)
        num = int(value)
    finally:
        return num if num is not None else value
```

If the file path contains directories, the directories are created automatically.

```txt filename=a/b/c/d.txt
ABCD
```

If the same filename is used again, the content of the code block is appended to the file.

```py filename=js_sum.py
def js_sum(*args):
    if not args:
        return ""

    result, *values = [number_or_string(arg) for arg in args]
    for value in values:
        try:
            result += value
        except TypeError:
            result = str(result) + str(value)

    return result
```

Files can be appended multiple times.

```py filename=js_sum.py
if __name__ == "__main__":
    import sys

    print(js_sum(*sys.argv[1:]))
```

The filepath is always treated as relative to the workspace of the test case.

```txt filename=/etc/resolv.conf
nameserver 1.1.1.1
```

The created files can be used by later code blocks. For example, we could run the python script and print the text files defined in the above code blocks:

```sh
cat a/b/c/d.txt
cat ./etc/resolv.conf

test $(python3 js_sum.py 1 2 3) = "6"
test $(python3 js_sum.py 1.0 1) = "2.0"
test $(python3 js_sum.py 1 a 2 b 3 c) = "1a2b3c"
```

The filename option overrides the default behavior of the code-block. I.e., `sh` code-blocks are not executed if `filename` option is defined and `env` blocks will not set environment variables if `filename` option is defined.
