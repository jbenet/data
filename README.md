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
data - dataset package manager

Commands:

    version     Show data version information.
    config      Manage data configuration.
    info        Show dataset information.
    list        List installed datasets.
    get         Download and install dataset.
    manifest    Generate and manipulate dataset manifest.
    pack        Dataset packaging, upload, and download.
    blob        Manage blobs in the blobstore.
    publish     Guided dataset publishing.
    user        Manage users and credentials.

Use "data help <command>" for more information about a command.
```

### data get

```
# author/dataset
> data get jbenet/bar
Downloading jbenet/foo from datadex.
get blob b53ce99 Manifest
get blob 2183ea8 Datafile
get blob 63443e4 data.csv
copy blob 63443e4 data.txt
copy blob 63443e4 data.xsl
get blob b53ce99 Manifest

Installed jbenet/foo@1.0 at datasets/jbenet/foo
```

### data list

```
> data list
jbenet/bar@1.0
```

### data info

```
> data info jbenet/foo
dataset: jbenet/foo@1.0
title: Foo Dataset
description: The first dataset to use data.
license: MIT

# shows the Datafile
> cat datasets/jbenet/bar/Datafile
dataset: foo/bar@1.1
```

### data publish

```
> data publish
==> Guided Data Package Publishing.

==> Step 1/3: Creating the package.
Verifying Datafile fields...
Generating manifest...
data manifest: added Datafile
data manifest: added data.csv
data manifest: added data.txt
data manifest: added data.xsl
data manifest: hashed 2183ea8 Datafile
data manifest: hashed 63443e4 data.csv
data manifest: hashed 63443e4 data.txt
data manifest: hashed 63443e4 data.xsl

==> Step 2/3: Uploading the package contents.
put blob 2183ea8 Datafile - uploading
put blob 63443e4 data.csv - exists
put blob b53ce99 Manifest - uploading

==> Step 3/3: Publishing the package to the index.
data pack: published jbenet/foo@1.0 (b53ce99).
```

Et voila! You can now use `data get foo/bar` to retrieve it!

### data config

```
> data config index.datadex.url http://localhost:8080
> data config index.datadex.url
http://localhost:8080
```

### data user

```
> data user
data user - Manage users and credentials.

Commands:

    add         Register new user with index.
    auth        Authenticate user account.
    pass        Change user password.
    info        Show (or edit) public user information.
    url         Output user profile url.

Use "user help <command>" for more information about a command.

> data user add
Username: juan
Password (6 char min):
Email (for security): juan@benet.ai
juan registered.

> data user auth
Username: juan
Password:
Authenticated as juan.

> data user info
name: ""
email: juan@benet.ai

> data user info jbenet
name: Juan
email: juan@benet.ai
github: jbenet
twitter: '@jbenet'
website: benet.ai

> data user info --edit
Editing user profile. [Current value].
Full Name: [] Juan Batiz-Benet
Website Url: []
Github username: []
Twitter username: []
Profile saved.

> data user info
name: Juan Batiz-Benet
email: juan@benet.ai

> data user pass
Username: juan
Current Password:
New Password (6 char min):
Password changed. You will receive an email notification.

> data user url
http://datadex.io:8080/juan
```

### data manifest (plumbing)

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

### data blob (plumbing)
```
data blob - Manage blobs in the blobstore.

Commands:

    put         Upload blobs to a remote blobstore.
    get         Download blobs from a remote blobstore.
    url         Output Url for blob named by <hash>.

Use "blob help <command>" for more information about a command.
```

```
> cat Manifest
Datafile: 0d0c669b4c2b05402d9cc87298f3d7ce372a4c80
data.csv: 63443e4d74c3a170499fa9cfde5ae2224060b09e
data.txt: 63443e4d74c3a170499fa9cfde5ae2224060b09e
data.xsl: 63443e4d74c3a170499fa9cfde5ae2224060b09e

> data blob put --all
put blob 0d0c669 Datafile
put blob 63443e4 data.csv

> data blob get 63443e4d74c3a170499fa9cfde5ae2224060b09e
data blob get 63443e4d74c3a170499fa9cfde5ae2224060b09e
get blob 63443e4 data.csv
copy blob 63443e4 data.txt
copy blob 63443e4 data.xsl

> data blob url
http://datadex.archives.s3.amazonaws.com/blob/0d0c669b4c2b05402d9cc87298f3d7ce372a4c80
http://datadex.archives.s3.amazonaws.com/blob/63443e4d74c3a170499fa9cfde5ae2224060b09e
```

### data pack (plumbing)

This is probably the most informative command to look at.

```
data pack - Dataset packaging, upload, and download.

Commands:

    make        Create or update package description.
    manifest    Show current package manifest.
    upload      Upload package contents to remote storage.
    download    Download package contents from remote storage.
    publish     Publish package reference to dataset index.
    check       Verify all file checksums match.

Use "pack help <command>" for more information about a command.
```


```
> ls
data.csv  data.txt  data.xsl

> cat data.*
BAR BAR BAR
BAR BAR BAR
BAR BAR BAR

> data pack make # interactive
Verifying Datafile fields...
Enter author name (required): foo
Enter dataset id (required): bar
Enter dataset version (required): 1.1
Enter dataset title (optional): Barrr
Enter description (optional): A bar dataset.
Enter license name (optional): MIT
Generating manifest...
data manifest: hashed 0d0c669 Datafile
data manifest: hashed 63443e4 data.csv
data manifest: hashed 63443e4 data.txt
data manifest: hashed 63443e4 data.xsl

> ls
Datafile  Manifest  data.csv  data.txt  data.xsl

> data pack manifest
Datafile: 0d0c669b4c2b05402d9cc87298f3d7ce372a4c80
data.csv: 63443e4d74c3a170499fa9cfde5ae2224060b09e
data.txt: 63443e4d74c3a170499fa9cfde5ae2224060b09e
data.xsl: 63443e4d74c3a170499fa9cfde5ae2224060b09e

> data pack upload
put blob 0d0c669 Datafile
put blob 63443e4 data.csv
put blob 8a2e6f6 Manifest

> rm data.*

> ls
Datafile  Manifest

> data pack download
get blob 63443e4 data.csv
copy blob 63443e4 data.txt
copy blob 63443e4 data.xsl

> ls
Datafile  Manifest  data.csv  data.txt  data.xsl

> data pack check
data pack: 4 checksums pass

> echo "FOO FOO FOO" > data.csv

> data pack check
data manifest: check 63443e4 data.csv FAIL
data pack: 1/4 checksums failed!

> data pack download
copy blob 63443e4 data.csv

> data pack check
data pack: 4 checksums pass

> data pack publish
data pack: published foo/bar@1.1 (8a2e6f6).
```

## Datafile

data tracks the definition of dataset packages, and dependencies in a
`Datafile` (in the style of `Makefile`, `Vagrantfile`, `Procfile`, and
friends). Both published dataset packages, and regular projects use it.
In a way, your project defines a dataset made up of other datasets, like
`package.json` in `npm`.

```
# Datafile format
# A YAML (inc json) doc with the following keys:

# required
handle: <author>/<name>[.<format>][@<tag>]
title: Dataset Title

# optional functionality
dependencies: [<other dataset handles>]
formats: {<format> : <format url>}

# optional information
description: Text describing dataset.
repository: <repo url>
website: <dataset url>
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
website: http://caltech.edu/~feynman/not-girls/plate-stuff/trial3
```

much more friendly and approachable than this

```
{
  "dataset": "feynman/spinning-plate-measurements",
  "title": "Measurements of Plate Rotation",
  "contributors": [
    "Richard Feynman <feynman@caltech.edu>"
  ],
  "website": "http://caltech.edu/~feynman/not-girls/plate-stuff/trial3"
}
```

It's already hard enough to get anyone to do anything. Don't add more hoops to
jump through than necessary. Each step will cause significant dropoff in
conversion funnels. (Remember, [Apple pays Amazon for 1-click buy](https://www.apple.com/pr/library/2000/09/18Apple-Licenses-Amazon-com-1-Click-Patent-and-Trademark.html)...)

And, since YAML is a superset of json, you can do whatever you want.


## Development

Setup:

1. [install go](http://golang.org/doc/install)
2. Run

    git clone https://github.com/jbenet/data
    cd data
    make install

You'll want to run [datadex](https://github.com/jbenet/datadex) too.

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
