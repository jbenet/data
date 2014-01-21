# building data

At present (Go 1.2), cross-compiling go does not work with cgo. It seems
(not actually sure, as not very familiar with all the deps) data uses cgo
extensively. While there seems to be a gcc work-around, it would be useful
to test the tool in all the platforms. Thus, for now, all supported archs
will have a vm in this directory. The process, then, is:

1. setup + launch the vm
1. compile + test data in vm
1. place release binary in `/platforms/<arch>/data`
1. `make <arch>-tar` + `make dist` to package bins up
