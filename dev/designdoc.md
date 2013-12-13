## WARNING - WIP

data is in very early development.
This document is too. Track ideas here, and in the [roadmap](roadmap.md)

# data designdoc

    data - a package manager for datasets
    datahub - a centralized dataset hosting service

## Abstract

    data : datasets :: git : source code
    data : datahub :: git : github


## Introduction

### Prerequisites

This document assumes strong familiarity with the following software
engineering concepts and systems:

- data formats
- datasets
- version control: `git, hg`
- central source code repositories: `github, google code`
- package managers: `aptitude, pip, npm, brew`
- package indices: `ubuntu packages, pypi, npm registry, docker index`
- containers: `LXC, docker`

Dataset management is a mess. There are millions of datasets strewn across the
internet, encoded in thousands of formats. [more gripes here]


## Design Goals

data must be

- **format agnostic**: no special treatment of specific formats. ideally, data
  itself does not understand formats.
- **domain agnostic**: no special treatment of specific application domains
  and
  their biases (e.g. machine learning vs genomics vs neuroscience).
- **platform agnostic**: no special treatment of specific platforms/stacks
  (*nix, windows, etc).


- **decentralized**: no requirement on one central package index. There will
  be one (a default), but data should be capable of pointing to other indices.
- **intuitive UX**: to facilitate adoption in a massively competitive/
  entrenched landscape, it is key to craft a highly intuitive user experience,
  with as gradual learning curve as possible.
- **simple to use**: simplicity is key. `data get norb/simple-norb`

- **modular**: to support development and feature exploration, data should be
  modular, isolating functionality. Learn from git.
- **infrastructure**: data is a general infrastructure tool. it aims to solve
  core, wide problems. special cases can be handled by sub-tools /
  applications on top.
- **command line**: data is a unix-style tool.


## datadex - data index

(The name datadex is worse than datahub, but datahub seems to be taken by
a related project. Perhaps collaborate? TODO: talk to datahub people.)

The important power data brings to the table is publishing and downloading
datasets from a repository, public or private. This is achieved by the use
of `datadex`, `data`'s sister tool and a simple website. The plan is to run
one main, global `datadex` (much like most successful package managers out
there) but allow users of `data` to point to whatever `datadex` (repository)
they wish to use.

The datadex is where data finds datasets when you run:

    data get jbenet/foobar

Dataset foobar is looked up at the default data index:
http://datadex.io/jbenet/foobar. Users should be able to point to an
entirely different datadex, or even list a secondary one. This is useful in
case the main datadex is {down, unmaintained, controlled by evil baboons},
and in case a user wishes to run her own private datadex for private datasets.

See more at https://github.com/jbenet/datadex/blob/master/dev/roadmap.md


## data handles

building on the roads paved by git and github, data introduces a standard way
to reference every unique dataset version. This is accomplished with
*data handles*: unique, impure, url-friendly identifiers.

data handle structure:

    <author>/<name>[.<format>][@<ref>]

Where:

- <author> is the datadex username of the author/packager e.g. `feynman`
- <name> is a unique shortname for the dataset e.g. `spinning-plates`
- <format> is an optional format. details TBD, see [`dev/formats`](formats.md)
  defaults to `default` e.g. `json`
- <ref> is an optional reference (hash, version, tag, etc).
  defaults to `latest` e.g. `1.0`

Examples:

    jbenet/cifar-10
    jbenet/cifar-10.matlab
    jbenet/cifar-10@latest
    jbenet/cifar-10@1.0
    jbenet/cifar-10.matlab@0.8.2-rc1

### URL handling

data handles are meant to be embedded in URLs, as in the datadex:

    http://datadex.io/jbenet/cifar-10@1.0

(yes @ symbols get encoded)


## data hashes and refs

data borrows more git concepts: object hashes and references.

**data hashes**: In git, objects are identified by their hashes (sha1); one can
retrieve an object with `git show <object-hash>`. In data, unique datasets --
including different versions of the same dataset -- are identified by the hash
value of their dataset archive (i.e. `hash(tar(dataset directory))`). All
published versions of a dataset are hosted in the datadex, e.g.:

    # these are all different datasets:
    http://datadex.io/jbenet/cifar-10@49be4be15ec96b72323698a710b650ca5a46f9e6
    http://datadex.io/jbenet/cifar-10@e9db19b48ced2631513d2a165e0386686e8a0c8a
    http://datadex.io/jbenet/cifar-10@5b13d6abb15dccabb6aaf8573d5a01cd0d74c86d

These are three different versions of the same named dataset. They may differ
only slightly, or be completely different.

**data references**: It is very useful to reference object versions (their
hashes) via human-friendly and even logical names, like `2.0`. These names are
simply references (pointers, symlinks) to hashes. data is designed to
understand (and de-reference) named references wherever it would normally
expect a hash. Moreover, the term `ref` is used throughout to mean `reference,
hash, version, tag, etc`. e.g.

    # while these could all point to the same dataset:
    http://datadex.io/jbenet/cifar-10 // defaults to @latest
    http://datadex.io/jbenet/cifar-10@latest
    http://datadex.io/jbenet/cifar-10@1.0
    http://datadex.io/jbenet/cifar-10@e9db19b48ced2631513d2a165e0386686e8a0c8a


**default ref**: it is worth noting that often the "default reference" will be
used when a reference is expected but not provided. The "default ref" is
`latest`, and it points to the latest published version.

**tags**: tags are user-specified references, e.g. version numbers like `1.0`.

## data manifest

data uses a manifest of `{ <paths> : <hashes> }` in order to:

- account what files are part of a dataset
- detect data corruption (check hashes match)
- provide minimal version control (manifest changesets)

data functions somewhat like `git-annex`:

- stores (version-controls) the path and object hash in the "repository"
- fetches the large blobs from a storage service

The blobs from all the datasets stored in the same object store. (Blobs from
different datasets are not segregated into separate bundles). This greatly
reduces storage needs, de-duplicating common blobs across datasets. This is
particularly useful for versions of the same dataset, as not all files change
between versions.

This design reduces storage both remotely (the datadex service de-duplicates
across all indexed datasets) and locally (users' computers keep one blob cache
for all installed datasets).
