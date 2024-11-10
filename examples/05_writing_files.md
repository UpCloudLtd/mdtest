# Writing files

When a code-block contains a `filename` option, the content of the code block is written to the file specified as a parameter for the option. For example, `` ```python filename=uuid4.py ``.

```py filename=uuid4.py
from uuid import uuid4;

print(str(uuid4()))
```

If the file path contains directories, the directories are created automatically.

```txt filename=a/b/c/d.txt
Dd
```

If the same filename is used again, the content of the code block is appended to the file.

```txt filename=a/b/c/d.txt
Dolphin
```

The filepath is always treated as relative to the workspace of the test case.

```txt filename=/etc/resolv.conf
nameserver 1.1.1.1
```

The created files can be used by later code blocks. For example, we could run the python script and print the text files defined in the above code blocks:

```sh
python3 uuid4.py
cat a/b/c/d.txt
cat ./etc/resolv.conf
```

The filename option overrides the default behavior of the code-block. I.e., `sh` code-blocks are not executed if `filename` option is defined and `env` blocks will not set environment variables if `filename` option is defined.
