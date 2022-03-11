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

### License
Copyright 2022 ForgetAll

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
