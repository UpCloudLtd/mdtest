# Real-time output streaming

Output from `sh` blocks is not logged to terminal by default. To print the output to terminal in real-time as the test is running, use the `--output-to-terminal` flag. The real-time output can be useful, for example, when debugging a test that might get stuck or when you want to see the progress of a long-running test.

Real-time terminal output is only available when running tests non-concurrently either by targeting a single test file or by setting `--jobs=1` when running multiple tests. For example, to see the output of below script in real-time, use either `mdtest --output-to-terminal --jobs=1 examples/` or `mdtest --output-to-terminal examples/08_real_time_output.md` commands to run the example(s).

```sh
# Simulate a process that might get stuck
for i in $(seq 1 10); do
  echo "Attempt $i: Checking status..."
  sleep ${SLEEP_TIME:-0.5}
done
echo "Process completed"
```
