# Here-Documents

To use this flag, set Dockerfile version to `labs` channel. This feature is available
since `docker/dockerfile:1.3.0-labs` release.

```dockerfile
# syntax=docker/dockerfile:1.3-labs
```

Here-documents allow redirection of subsequent Dockerfile lines to the input of `RUN` or `COPY` commands.
If such command contains a [here-document](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_07_04)
Dockerfile will consider the next lines until the line only containing a here-doc delimiter as part of the same command.

## Example: Running a multi-line script

```dockerfile
# syntax=docker/dockerfile:1.3-labs
FROM debian
RUN <<EOT bash
  apt-get update
  apt-get install -y vim
EOT
```

If the command only contains a here-document, its contents is evaluated with the default shell.

```dockerfile
# syntax=docker/dockerfile:1.3-labs
FROM debian
RUN <<EOT
  mkdir -p foo/bar
EOT
```

Alternatively, shebang header can be used to define an interpreter.

```dockerfile
# syntax=docker/dockerfile:1.3-labs
FROM python:3.6
RUN <<EOT
#!/usr/bin/env python
print("hello world")
EOT
```

More complex examples may use multiple here-documents.

```dockerfile
# syntax=docker/dockerfile:1.3-labs
FROM alpine
RUN <<FILE1 cat > file1 && <<FILE2 cat > file2
I am
first
FILE1
I am
second
FILE2
```

## Example: Creating inline files

In `COPY` commands source parameters can be replaced with here-doc indicators.
Regular here-doc [variable expansion and tab stripping rules](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_07_04) apply.

```dockerfile
# syntax=docker/dockerfile:1.3-labs
FROM alpine
ARG FOO=bar
COPY <<-EOT /app/foo
  hello ${FOO}
EOT
```

```dockerfile
# syntax=docker/dockerfile:1.3-labs
FROM alpine
COPY <<-"EOT" /app/script.sh
  echo hello ${FOO}
EOT
RUN FOO=abc ash /app/script.sh
```
