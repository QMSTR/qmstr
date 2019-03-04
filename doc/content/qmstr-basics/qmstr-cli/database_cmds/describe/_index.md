---
title: "Describe command"
date: 2019-02-27T09:48:15Z
draft: false
weight: 3
---

Describe command prints a description of the referenced node and traverses the tree 
to print a description of the nodes connected to it. 

Describe, can be very useful not only for developers (for debugging purposes) but also for user 
as they can review the nodes, the attributes and the connection between the nodes that are stored 
in the database.

The command follows the generic syntax of the database commands to reference nodes:

    qmstrctl describe <type_of_node:attribute:value> 

For example:

    > qmstrctl describe file:name:debug.o

    > qmstrctl describe file:name:debug.o --less

the --less flag, will traverse the tree and will ignore the information that has been collected 
from the analyzers (info nodes and data nodes). If --less flag is not provided then the `describe` 
command will print out all the nodes in the database connected to our node.
