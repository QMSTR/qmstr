---
title: "Connections between nodes"
date: 2019-02-27T09:48:15Z
draft: false
weight: 20
---

Connect command connects two specified nodes. The command follows the generic syntax of the database commands to reference nodes.
The form should be:

    qmstrctl connect <type_of_node:attribute:value> <type_of_node:attribute:value>

for example,

    > qmstrctl connect file:name:debug.o file:name:debug.c
