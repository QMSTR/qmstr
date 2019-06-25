---
title: Support for Android Application (apk) builds
draft: false
weight: 10
---

The proposed way to develop and build Android applications is using gradle. Eventhough Android applications are written in Java the gradle plugin for Android is not using the Java plugin, it is actually incompatible and will break the build if the java plugin is loaded.

This is due to the fact that compiling java sources to class files is only one part in the  build life-cycle of android applications. After compiling classes a dex file is created from those classes. This dex file will become part of the apk (a signed jar file). Other steps are obfuscation/optimization of source code via proguard, generating XML files and packing resources into as well as signing the resulting apk.

# Development Plan

Generating QMSTR build graphs for Android applications must be done via instrumenting the build system Gradle.

## QMSTR 0.6: TRL 6

- First working implementation in the qmstr-gradle-plugin.
- Generation of the build graph for android-blockly apk.

## QMSTR 0.6: TRL 7

  - Support natively built parts (NDK) like shared objects built from C/C++ source code.
