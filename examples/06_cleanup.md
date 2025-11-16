# Cleanup

By default teststeps are skipped if an earlier step has failed. When a `sh` code-block contains `cleanup` (or `cleanup=true`) option, the script is executed without -e flag and even if there are earlier failed teststeps. For example, `` ```sh cleanup``.
    
```sh cleanup
rm -rf file-created-by-test.txt
```
