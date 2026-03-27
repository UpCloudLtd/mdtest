# Real-time output streaming

Test output can be streamed live using the `--ouput-to-terminal` flag, demonstrating the whole output of an `sh` script block for even more fine-tuned debugging.

Make sure to run the examples using `./bin/mdtest --output-to-terminal --jobs=1 examples/`, otherwise live output will not be visible.

```sh
# Simulate a process that might get stuck
for i in $(seq 1 10); do
  echo "Attempt $i: Checking status..."
  sleep ${SLEEP_TIME:-0.5}
done
echo "Process completed"
```
