---
title: "Node labels in Quartermaster"
date: 2019-02-27T09:48:15Z
draft: false
weight: 30
---

Node labels are used to flexibly add data to nodes based on key-value pairs. For example, if a build process consists of multiple stages, a label can be used to mark nodes with the stage during which they are created. Or a target program generated during the build could be marked with a release milestone name. The labelling functionality allows for user-defined concepts to be flexibly represented in the knowledge graph.

## Labels basics

Labels are attached to nodes and consist of key-value pairs. The key can be any [POSIX "fully portable filename"](https://en.wikipedia.org/wiki/Filename) like `build-stage` or `Version_Name`. The value can be any Unicode string, like `"Hello world!"`. Values may be empty, since sometimes it only matters whether or not a label exists.

## Creating and deleting labels

Let's assume a file node exists in the knowledge graph called `file:src/main.c`. We want to mark it as modified in the current release 0.3.0 so that it can be highlighted in the reports.

    > qmstrctl label set file:src/main.c MODIFIED

This command will create the label `MODIFIED` on the node with an empty value. The command will return an error if the node does not exist:

```shell {shelldocexitcode=2}
> qmstrctl label set file:src/nonexistant.c ERROR
Error: Node "file:src/nonexistant.c" not found.
```

The label value is queried using the `get` command:

    > qmstrctl label get file:src/main.c MODIFIED
    ... (empty line)

Note that the command returned the correct empty value, but succeeded. Querying a label that is not set on the node results in an error:

```shell {shelldocexitcode=1}
> qmstrctl label get file:src/main.c some-label
Error: Label "some-label" undefined on node "file:src/main.c".
```

Re-creating an existing label will simply change its value. This adds the version that was missing when the label was initially created:

    > qmstrctl label set file:src/main.c MODIFIED 0.3.0
    > qmstrctl label get file:src/main.c MODIFIED
    0.3.0

A label is removed using the `delete` command:

    > qmstrctl label get file:src/main.c MODIFIED
    0.3.0
    > qmstrctl label delete file:src/main.c MODIFIED

Reading the label after deleting it results in an error:

```shell {shelldocexitcode=1}
> qmstrctl label get file:src/main.c MODIFIED
Error: Label "MODIFIED" undefined on node "file:src/main.c".
```

There are no practical limits to the number of labels or the size of the label values attached to a node.
