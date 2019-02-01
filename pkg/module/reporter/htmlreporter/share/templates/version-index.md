+++
title = "{{{shortenId .VersionIdentifier}}}"
versionid = "{{{.VersionIdentifier}}}"
package = "{{{.Package.PackageName}}}"
author = "{{{.Author}}}"
summary = {{ htmlEscape "{{{summary .Message}}}" }}
weight = 1
+++

## {{{.Package.PackageName}}} @ {{{.VersionIdentifier}}}

{{< version-cmpl-header >}}

{{< buildconfig-list >}}