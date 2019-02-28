---
title: "Create command"
date: 2019-02-27T09:48:15Z
draft: false
weight: 1
---

Create command creates a specific node, either a `project`, a `package` or a `file` node.
The command does NOT follow the generic syntax of the database commands to reference nodes.
To specify an attribute with a value, include the corresponding flag. 

Type the following to get more information about the command:

    > qmstrctl create file --help

For example:

    > qmstrctl create file --name debug.o

    > qmstrctl create file --name debug.c
