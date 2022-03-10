![work flow](https://github.com/ForgetAll/xcli/actions/workflows/build.yml/badge.svg)

use cobra to build self command line tool

### Usage
```shell
go get -u github.com/ForgetAll/xcli
```

uuid:
    * xcli uuid: generate uuid
    * xcli uuid -t/--trim: generate uuid without '-' char

rand:
    * xcli rand: generate a rand number, from 0 to 1<<63 - 1
    * xcli rand --from/-f valueFrom --to/-t valueTo: generate a rand in [from, to)
    * xcli rand --real/-r: use crypto/rand instead of math/rand to generate random nunber

### Contribution
Welcome to discuss or commit pull request!

Before commit, please make sure `lint` and `test` which defined in Makefile can pass.
Commit message should observe [rule](https://github.com/woai3c/Front-end-articles/blob/master/git%20commit%20style.md).