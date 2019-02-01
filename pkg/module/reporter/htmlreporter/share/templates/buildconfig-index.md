+++
title = "{{{.BuildConfig}}}"
versionid = "{{{.VersionIdentifier}}}"
package = "{{{.Package.PackageName}}}"
author = "{{{.Author}}}"
summary = {{ htmlEscape "{{{summary .Message}}}" }}
weight = 1
+++

## {{{.Package.PackageName}}} @ {{{.VersionIdentifier}}}

{{< version-cmpl-header >}}

{{< target-list >}}

{{< authors >}}
