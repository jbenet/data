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


## data manifest

    data manifest <subcmd>

Generate the data manifest file (`.data-manifest`? `Datamanifest`?), a list of

    <object hash> <object filename>

Hash function? `sha1` for now. Discuss whether to use `sha256`.
This manifest file is all that is needed to reconstruct the dataset.

(Manifest files can support a poor-man's version control (changesets of
filename/filehashes). Basically, a much simpler `git-annex`. Potentially use
git to store Datafile + manifest:
- making data refs = git refs
- repository abstraction great for storing all package versions)

Subcommands:

    data manifest add [<file> | -a, --all]
    data manifest rm [<file> | -m, --missing | -a, --all]
    data manifest hash [<file> | -r, --rehash]
    data manifest check [<file> | -a, --all]


## data upload

    data upload <storage service>

Upload package archive to the chosen storage service. This ensures the final
archive URL is listed as a mirror in the Datafile.

Use a datadex-specific s3 bucket:

    data upload
    data upload datadex
    # PUTs to http://datadex.archives.s3.amazonaws.com/<author>/<dataset>/<ref>/<dataset>-<ref>.tar.gz

Things to do:

- ensure there is a version
- ensure there is a title
- prompt for a description
- prompt for a license


## data publish

    data publish

Upload and register this package to the dataset index (datadex, configurable).
Registered packages require extra definitions in their `Datafile`.

Things to do:

- ensure uploaded
- post datafile to datadex/publish



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
