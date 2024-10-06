# Writing files

When a code-block contains a `filename` option, the content of the code block is written to the file specified as a parameter for the option. For example, `` ```python filename=uuid4.py ``.

```py filename=uuid4.py
from uuid import uuid4;

print(str(uuid4()))
```

The created files can be used by later code blocks. For example, we could run the python script defined in the above code block:

```sh
python3 uuid4.py
```

The filename option overrides the default behavior of the code-block. I.e., `sh` code-blocks are not executed if `filename` option is defined and `env` blocks will not set environment variables if `filename` option is defined.

If the same filename is used again, the content of the code block is appended to the file.