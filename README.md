# data - package manager for datasets


Imagine installing datasets like this:

    data get jbenet/norb

It's about time we used all we've learned making package managers to fix the
awful data management problem. Read the [designdoc](dev/designdoc.md).


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
