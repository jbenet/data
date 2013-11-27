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




