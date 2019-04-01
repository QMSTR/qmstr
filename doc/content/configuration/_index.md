---
title: "Configuration"
date: 2019-01-17T15:26:15Z
draft: false
weight: 3
---

`qmstr.yaml` is QMSTR's project level configuration file. It holds
information about the analysis and reporting phase of a QMSTR process.

Let's assume we want to run the JSON-C library as the project under analysis.
The `qmstr.yaml` file for this project would look like this:

``` 
package:
  name: "json-c"
  metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":50051"
    dbaddress: "localhost:9080"
    dbworkers: 4
  analysis:
    - analyzer: pkg-analyzer
      name: "Package Analyzer"
      config:
        targetdir: "/buildroot/jsonc/.libs"
        targets: "libjson-c\\.(?:so)|(?:dylib)\\.\\d\\.\\d.\\d"
    - analyzer: spdx-analyzer
      name: "Simple SPDX Analyzer"
      trustlevel: 300
    - analyzer: scancode-analyzer
      name: "Scancode Analyzer"
      trustlevel: 400
      config:
        workdir: "/buildroot"
        resultfile: "/buildroot/scancode.json"
        # cached: "true"
    - analyzer: git-analyzer
      name: "Git Analyzer"
      config:
        workdir: "/buildroot/jsonc"
    - analyzer: test-analyzer
      name: "Simple CI Test Analyzer"
      config:
        tests: "TestPackageNode"
    - analyzer: pyqmstr-spdx-analyzer
      name: "SPDX Analyzer"
      trustlevel: 300
      config:
        spdxfile: "/buildroot/SPDX.tag"
  reporting:
    - reporter: test-reporter
      name: "Test Reporter"
      config:
        siteprovider: "Endocode"
    - reporter: qmstr-reporter-html
      name: "HTML Reporter"
      config:
        siteprovider: "Endocode"
        baseurl: "http://qmstr.org/packages/"

```

Depending on the project you want to run QMSTR with, you should change the package name to your project's name.

``` 
package:
  name: project_name
```

The metadata and server sections, are the default data and there is no need to modify them.

```
metadata:
    Vendor: "Endocode"
    OcFossLiaison: "Mirko Boehm"
    OcComplianceContact: "foss@endocode.com"
  server:
    rpcaddress: ":50051"
    dbaddress: "localhost:9080"
    dbworkers: 4
```

## Analyzers

The next step is to define and configure the analyzers to be used on our project.
Analyzers run after the build process has finished successfully.
All analyzers include the above sections:

``` 
- analyzer:
   name:
```

In the `analyzer` section we provide the module name as it is defined in QMSTR.
In the `name` section we provide a string with the full name of the analyzer.

Most of the analyzers include an extra section **`config`**. 
The config section diverses for each analyzer, regarding the dependencies.

### Missing Pieces Analyzer

The `missing pieces analyzer` adds extra build information into the database. 
It should be the first analyzer to run in the analysis phase, so it should
be the first analyzer to be configured in the `qmstr.yaml`, if you want to use it,
as it is the only analyzer which extends the build graph with additional filenodes.
The rest of the analyzers extend the build graph with additional infonodes.
To use the missing pieces analyzer, please provide the directory to your input file
in the `config` section. 
The input file should be a yaml file containing the extra build
information to be added in the database.
The configuration would look like this:

```
analysis:
    - analyzer: missing_pieces-analyzer
      name: "Missing pieces Analyzer"
      config:
        inputfile: "/buildroot/mpa.yml"
```

where mpa.yaml is our input yaml file.

### Package Analyzer

`Package analyzer` concatenates the package node to the rest of the graph. 
This is achieved by connecting the package node to the desired target file nodes.

To use the package analyzer, please fill in the path to the target files and 
the target names in the `config` section, in your `qmstr.yaml` file:

```
        targetdir: target_directory
        targets: list_of_targets
```

To add many targets, you can seperate them with the `;` delimiter. 
For example:

```
        targetdir: "/buildroot/my_project/"
        targets: "a;a.so;b;b.a"
```

Targets can also include regular expressions.

### SPDX Analyzer

`SPDX analyzer` scans input files and package manifests for project metadata 
in SPDX format.

### Scancode Analyzer

`Scancode analyzer` uses the [scancode-toolkit](https://github.com/nexB/scancode-toolkit) to
scan a codebase for licenses and copyrights. 
Provide the working directory in order to use it.
Scancode analyzer may be time consuming. To save time you can provide a result file, 
to save the scancode result. Next time you run
your program with QMSTR you can provide the following line in your qmstr.yaml file:
```
cached: "true"
```

That way, the analyzer will use directly the scan result file from the previous build and 
it will skip the phase where it scans all the files. For the first build this option can not 
be used, as we would not have a scan result file. 

### Python SPDX Analyzer

QMSTR is programed primarily in Golang, but it also includes some Python and Java programming. 
Our aim is for users to be able to include their own analyzers in the QMSTR process, which may 
be written in different programming languages. With the use of [gRPC](https://grpc.io) 
and [protocol buffers](https://developers.google.com/protocol-buffers/docs/overview) 
we can write client applications in different languages and connect them to the QMSTR server 
applications. 
`Python SPDX analyzer` is an example of this use case. It is written in Python and searches for 
SPDX licenses. 
To use this analyzer you have to include in the `config` section the directory 
of the file with all the information to be parsed by the analyzer. 


### Git Analyzer

The `git-analyzer` extracts the current revision and other metadata from the repository 
of the project under analysis.
Provide the path to the repository of the file in the `config` section, 
in order to use the analyzer. 

### Test Analyzer

`Test analyzer` runs tests to confirm that the build graph is valid. 


As you may have already noticed, in the qmstr.yaml you can define a `trustlevel`. 
Trust level does not influence anything yet, but for future purposes it will be a tool
to choose certain results between the analyzers.


## Reporters

Any functionality that processes the information in the knowledge graph and performs 
some action, is considered a report. A reporter may create output files, 
or submit messages to an IRC channel, or interact with an issue tracker.
Reporters have to be configured after the declaration of the analyzers on the section
`package/reporting:`.

### Test Reporter 

`Test reporter` runs tests on the build graph to reassure the graph contains the expected results 
after the analysis phase and that we can proceed to the reporting phase. 

### HTML Reporter 

`HTML reporter` produces a site with the QMSTR reports of the project.
In the `config` section please provide the URL which will host the site.
