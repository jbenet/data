```
data

  version     Show data version information.
  config      Manage data configuration.
  info        Show dataset information.
  list        List installed datasets.
  get         Download and install dataset.
  publish     Guided dataset publishing.

  user        Manage users and credentials.
    add         Register new user with index.
    auth        Authenticate user account.
    pass        Change user password.
    info        Show (or edit) public user information.
    url         Output user profile url.

  manifest    Generate and manipulate dataset manifest.
    add <file>      Adds <file> to manifest (does not hash).
    rm <file>       Removes <file> from manifest.
    hash <file>     Hashes <file> and adds checksum to manifest.
    check <file>    Verifies <file> checksum matches manifest.

  pack        Dataset packaging, upload, and download.
    make       Create or update package description.
    manifest   Show current package manifest.
    upload     Upload package to remote storage.
    download   Download package from remote storage.
    checksum   Verify all file checksums match.

  blob        Manage blobs in the blobstore.
    put <hash>    Upload blob named by <hash> to blobstore.
    get <hash>    Download blob named by <hash> from blobstore.
    url <hash>    Output Url for blob named by <hash>.
    check <hash>  Verify blob contents named by <hash> match <hash>.
    show <hash>   Output blob contents named by <hash>.
```
