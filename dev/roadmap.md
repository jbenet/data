# data roadmap

This document briefly outlines desired features to implement.


## command dispatch

Need to implement the skeleton of the project: command parsing/dispatch.

## data list

    data list

List the datasets in the current project

## data config

    data config user.name = 'jbenet'
    data config --global user.name = 'jbenet'

Allow the configuration of `data`, using (`git` like) config files.
Consider using a `~/.dataconfig` global config file.
Consider using a `data/config` (or `.dataconfig`) local config file.

## data update

    data update

Download and install newer version.
Also, check whether data is up-to-date on every run (inc option to silence).

## data get

    data get <author>/<dataset>
    data get http://datadex.io/<author>/<dataset>

Download and install packages from the dataset index (datadex, configurable).
No arguments looks into the directory's `Datafile` (configurable)
Allow installation of packages using `<author>/<dataset>` ref-naming.
Allow installation of packages using `https?://.../<author>/<dataset>` urls.
Use a `--save` flag to store into a `Datafile`.
Installed datasets go into the `data/` directory (configurable) of the project.
Should download compressed files, and use array of mirrors.

## data put

    data put <author>/<dataset>

Upload and register this package to the dataset index (datadex, configurable).
Registered packages require extra definitions in their `Datafile`.

## data format

    data format <author>/<dataset> <desired format>
    data put <author>/<dataset>.<format>
    ref: <author>/<dataset>.<format>

Convert a dataset from one format to another.
Allow datasets to have multiple formats.
Formats should be convertible -- `f : f(dataset.fmt1) -> dataset.fmt2`
Formats should be defined/enabled per-dataset (in their Datafile).

## data tag

    data tag
    data get <author>/<dataset>@<tag>
    data put <author>/<dataset>@[<src tag>:]<tag>
    ref: <author>/<dataset>@<tag>

List the available (named) tags.
Allow referencing of datasets using specific tags.
Unnamed tags are version hashes.
Named tags are aliases to version hashes.
Put tags to create aliases.

## data slice

    ref: <author>/<dataset>#<slice>

See [`dev/formats`](formats.md).
