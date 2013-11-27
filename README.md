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

## Development

Setup:

1. [install go](http://golang.org/doc/install)
2. run `go build`

Build and install:

    make
    make install
