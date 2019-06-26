---
title: Towards self-contained modules
draft: false
weight: 10
---

The fact that modules need to be started by the master server leads requires all modules to be present in the PATH for qmstr master server to start them. This leads to several short-commings:

    - Long cycles of reassembling the whole image when changing something in one module.
    - All dependencies of a module need to be installed in the master image. 
    - It is difficult to debug modules.

To overcome this modules should be self-contained and able to run outside the qmstr master server. Starting the modules should be done by qmstrctl since analysis and reporting phase are triggered via qmstrctl anyway. Modules could run on the host qmstrctl runs on or in a dedicated container that -- with the help of qmstrctl -- is started in the master server's container network.

For this to work master and module need to have a conversation. The protocol for this conversation uses gRPC's bi-directional streaming as follows:

module              master
|     register      |
|------------------>|
|                   |
|     config        |
|<------------------|
|                   |
|     ready         |
|------------------>|
|                   |
|     start         |
|<------------------|
|                   |
|     send result 1 |
|------------------>|
...
|     send result n |
|------------------>|
|     end           |
|------------------>|
