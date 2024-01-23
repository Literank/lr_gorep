# lr_gorep

Example project: grep-like CLI app implemented in Go.

## Install

```bash
make install
```

## Usage

```txt
Options:
  -c    count
        Only a count of selected lines is written to standard output.
  -i    ignore-case
        Perform case insensitive matching. By default, it is case sensitive.
  -n    line-number
        Each output line is preceded by its relative line number in the file, starting at line 1. This option is ignored if -count is specified.
  -r    recursive
        Recursively search subdirectories listed.
  -v    invert-match
        Selected lines are those not matching any of the specified patterns.
```

### Search

```bash
gorep pattern *txt

cat *py | gorep pattern
```

### Recursive Search

```bash
gorep -r pattern .
```

### Search Multiple Files

```bash
gorep pattern a.txt b.py c.cpp
```

### Show Line Numbers

```bash
gorep -n pattern *txt
```

See [project tutorial](https://literank.com/project/10/intro) here.
