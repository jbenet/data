# data designdoc

    data - a package manager for datasets
    datahub - a centralized dataset hosting service

## Abstract

data : datasets :: git : source code
data : datahub :: git : github


## Introduction

### Prerequisites

This document assumes strong familiarity with the following software engineering
concepts and systems:

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

- format agnostic: no special treatment of specific formats. ideally, data itself
  does not understand formats.
- domain agnostic: no special treatment of specific application domains and their
  biases (e.g. machine learning vs genomics vs neuroscience).
- platform agnostic: no special treatment of specific platforms/stacks (*nix,
  windows, etc).


- decentralized: no requirement on one central package index. There will be one (a
  default), but data should be capable of pointing to other indices.
- intuitive UX: to facilitate adoption in a massively competitive/entrenched
  landscape, it is key to craft a highly intuitive user experience, with as gradual
  learning curve as possible.
- simple to use: simplicity is key. `data get norb/simple-norb`

- modular: to support development and feature exploration, data should be modular,
  isolating functionality. Learn from git.
- infrastructure: data is a general infrastructure tool. it aims to solve core,
  wide problems. special cases can be handled by sub-tools / applications on top.
- command line: data is a unix-style tool.


