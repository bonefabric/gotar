# Table of Contents

- [Tar archiver](#tar-archiver)
- [Usage](#usage)
    + [Examples](#examples)
    + [Flags](#flags)
- [Contributing](#contributing)

# Tar archiver
A simple archiver written in go and using only the standard library

Usage
-----

#### Examples

##### Extract files from _files.tar_ archive to _output_ directory
```shell
gotar -x -f files.tar output/
```
##### Create archive with _test.txt_ and _somefiles_ directory to default _archive.tar_ named archive
```shell
gotar -c test.txt somefiles
```
#### Flags

- -x Extract files from archive
- -c Create archive
- -v Verbose output
- -f Archive file name

Contributing
------------

If you find an issue, please report it on the
[issue tracker](https://github.com/bonefabric/gotar/issues/new)