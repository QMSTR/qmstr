---
title: "Database commands"
date: 2019-02-27T09:48:15Z
draft: false
---

Quartermaster uses commands to flexibly modify nodes in the database. For example, if you want manually to add nodes in the database and connect them to other nodes or remove a node or the connection between two nodes, that is possible with the above commands.

All commands use the following generic syntax to reference a node:

    qmstrctl <cmd_to_execute> <type_of_node:attribute:value>

, where `type_of_node` can be either `project`, `package` or `file`,

`attribute` can be `name` or `path` and 

`value`is the value of either the `name` or the `path`.