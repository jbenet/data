# data - package manager for datasets

This is the %(arch)s binary distribution of `data` (%(version)s).

Find other distributions, or install instructions at
http://datadex.io/doc/install

Website:        http://datadex.io
Repository:     https://github.com/jbenet/data
Issues:         https://github.com/jbenet/data/issues
Mailing list:   data-discuss@googlegroups.com


## Install

Move the `data` binary included in this archive to a location accessible
in your $PATH E.g.:

    sudo mv data /usr/bin/data

Now, see if data works. Run:

    data version

You should see:

    data version %(version)s

If you run into trouble, check the websites + mailing list above.


## Docs

You can list available commands by running:

    data

You can always get usage instructions with:

    data <command> help

To see a reference of all data commands run:

    data commands help | less
