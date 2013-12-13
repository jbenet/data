# data - package manager for datasets

Imagine installing datasets like this:

    data get jbenet/norb

It's about time we used all we've learned making package managers to fix the
awful data management problem. Read the [designdoc](dev/designdoc.md) and
the [roadmap](dev/roadmap.md).

#### Table of Contents

- [Usage](#usage)
- [Datafile](#datafile)
- [Development](#development)
- [About](#about)

## Usage

```
data is a dataset package manager.

Usage:

    data <command> [arguments]

Commands:

    get       Download and install dataset.
    list      List installed datasets.
    info      Show dataset information.
    help      Show usage information.
    version   Show data version information.
    upload    Upload dataset to storage service.
    manifest  Generate dataset manifest.
```

### data get

```
# author/dataset
> data get foo/bar
Downloading archive at http://datadex.io/foo/bar/archive/master.tar.gz
foo/bar@1.1 downloaded
foo/bar@1.1 installed

# url
> data get http://datadex.io/foo/bar/archive/master.tar.gz
Downloading archive at http://datadex.io/foo/bar/archive/master.tar.gz
foo/bar@1.1 downloaded
foo/bar@1.1 installed
```

### data list

```
> data list
    foo/bar              @1.1
```

### data info

```
> data info foo/bar
dataset: foo/bar@1.1

# shows the Datafile
> cat datasets/foo/bar/Datafile
dataset: foo/bar@1.1
```

### data upload

```
> data upload
Uploading objects to datadex storage service...(12/123) 54%
```

### data manifest

```
> data manifest add filename
data manifest: added filename

> data manifest hash filename
data manifest: hashed 61a66fd filename

> cat .data-manifest
filename: 61a66fda64e397a82d9f0c8b7b3f7ba6bca79b12

> data manifest rm filename
data manifest: removed filename
```

## Datafile

data tracks the definition of dataset packages, and dependencies in a
`Datafile` (in the style of `Makefile`, `Vagrantfile`, `Procfile`, and
friends). Both published dataset packages, and regular projects use it.
In a way, your project defines a dataset made up of other datasets, like
`package.json` in `npm`.

```
Datafile format

A YAML (inc json) doc with the following keys:

(required:)
handle: <author>/<name>[.<format>][@<tag>]
title: Dataset Title

(optional functionality:)
dependencies: [<other dataset handles>]
formats: {<format> : <format url>}

(optional information:)
description: Text describing dataset.
repository: <repo url>
homepage: <dataset url>
license: <license url>
contributors: ["Author Name [<email>] [(url)]>", ...]
sources: [<source urls>]
```
May be outdated. See [datafile.go](datafile.go).

### why yaml?

YAML is much more readable than json. One of `data`'s [design goals
](https://github.com/jbenet/data/blob/master/dev/designdoc.md#design-goals)
is an Intuitive UX. Since the target users are scientists in various domains,
any extra syntax, parse errors, and other annoyances could cease to provide
the ease of use `data` aims for. I've always found this

```
dataset: feynman/spinning-plate-measurements
title: Measurements of Plate Rotation
contributors:
  - Richard Feynman <feynman@caltech.edu>
homepage: http://caltech.edu/~feynman/not-girls/plate-stuff/trial3
```

much more friendly and approachable than this

```
{
  "dataset": "feynman/spinning-plate-measurements",
  "title": "Measurements of Plate Rotation",
  "contributors": [
    "Richard Feynman <feynman@caltech.edu>"
  ],
  "homepage": "http://caltech.edu/~feynman/not-girls/plate-stuff/trial3"
}
```

It's already hard enough to get anyone to do anything. Don't add more hoops to
jump through than necessary. Each step will cause significant dropoff in
conversion funnels. (Remember, [Apple pays Amazon for 1-click buy](https://www.apple.com/pr/library/2000/09/18Apple-Licenses-Amazon-com-1-Click-Patent-and-Trademark.html)...)

And, since YAML is a superset of json, you can do whatever you want.


## Development

Setup:

1. [install go](http://golang.org/doc/install)
2. run `go build`

Build and install:

    make
    make install


## About

This project started because data management is a massive problem in science*.
It should be **trivial** to (a) find, (b) download, (c) track, (d) manage,
(e) re-format, (f) publish, (g) cite, and (h) collaborate on datasets. Data
management is a problem in other domains (engineering, civics, etc), and `data`
seeks to be general enough to be used with any kind of dataset, but the target
use case is saving scientists' time.

Many people agree we direly need the
"[GitHub for Science](http://static.benet.ai/t/github-for-science.md)";
scientific collaboration problems are large and numerous.
It is not entirely clear how, and in which order, to tackle these
challenges, or even how to drive adoption of solutions across fields. I think
simple and powerful tools can solve large problems neatly. Perhaps the best
way to tackle scientific collaboration is by decoupling interconnected
problems, and building simple tools to solve them. Over time, reliable
infrastructure can be built with these. git, github, and arxiv are great
examples to follow.

`data` is an attempt to solve the fairly self-contained issue of downloading,
publishing, and managing datasets. Let's take what computer scientists have
learned about version control and distributed collaboration on source code,
and apply it to the data management problem. Let's build new data tools and
infrastructure with the software engineering and systems design principles
that made git, apt, npm, and github successful.

### Acknowledgements

`data` is released under the MIT License.

Authored by [@jbenet](https://github.com/jbenet). Feel free to contact me
at <juan@benet.ai>, but please post
[issues](https://github.com/jbenet/data/issues) on github first.

Special thanks to
[@colah](https://github.com/colah) (original idea and
[data.py](https://github.com/colah/data)),
[@damodei](https://github.com/damodei), and
[@davidad](https://github.com/davidad),
who provided valuable thoughts + discussion on this problem.
