# data changelog

## v0.1.1 2014-02-

- publish guide messages
- default dataset id to cwd basename
- changed Manifest -> .data/Manifest filename
- data get: install path is handle
- data get: no littering if not found
- data blob: creates dir(path)
- data config flexibility
- semver support


## v0.1.0 2014-01-21

First preview (alpha)

- release builds
- data commands (for reference)
- data pack make -- Datafile defaults
- datadex api suffix
- data blob put -- verify hash
- data blob {hash, check}
- datadex interop
- data config: env var, --edit
- s3 token based auth for uploading
- s3 anonymous downloading

## v0.0.5 2014-01-09

Publishing + downloading packages.

- data pack publish
- data publish
- data get (using pack)
- data user {add, auth, pass, info, url}
- data config

## v0.0.4 2014-01-03

Manifest manipulation and packaging.

- data manifest {add, rm, hash, check}
- data pack {make, manifest, upload, download, check}

## v0.0.3 2013-12-13

Uploading datasets.

- data manifest (list + hash files)
- data blob (blobs to storage service)


## v0.0.2 2013-11-24

Downloading datasets.

- data get (downloads + installs a dataset)

## v0.0.1 2013-11-22

Initial version.

- command dispatch
- datafile format (yml + structure)
- datafile parsing (loading/dumping)
- data version
- data help (just usage for now)
- data list (show installed datasets)
- data info (loads/dumps dataset's Datafile)
