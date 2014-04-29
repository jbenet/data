# data formats


One of the important design goals is format-fluidity: ability to store datasets in various formats and transfer between them. Suppose a graph of formats, datasets should be able to traverse strongly connected components. So, if a dataset is published in XML, I should be able to request it in json.[1]  This is easy for homogeneous datasets, but gets complicated when one dataset includes files in multiple formats, or it has metadata separated out.

This is complicated further when thinking about how datasets get authored/published to the index, and retrieved thereafter. In brief, the idea is to follow github pattern: `<author>/<dataset>`, which reduces namespace problems. This includes versions (tags/branches): `<author>/<dataset>@<tag>`. Note: this handle will be used in projects' Datafiles, to specify dependencies (datasets composed of other datasets[2]), etc.


 Some possibilities:


1. Let formats be branches like any other. `<author>/<dataset>@<format>`. Since version and format are now in the same namespace, would see things like: `foo/bar@1.0-json`, `foo/bar@1.2-xml`. This complicates maintenance: both new versions or new formats require a "row" of "commits" along the formats or versions, respectively.

2. Let formats be dimensions (see [3]). `<author>/<dataset>#format:<fmt>`. Would see things like: `foo/bar@format:json`, `foo/bar@format:xml`  There would be dimensional 'defaults' (as HEAD is default tag) that could be specified in the package description file.

3. Let formats be specified separately. `<author>/<dataset>.<fmt>`. e.g. `foo/bar.json`, `foo/bar.xml`. This seems neat and nice.

4. Punt. let authors choose their formats in the dataset. Would see things like: `foo/bar-json`, `foo/xmlbar`. Would not have format-fluidity :(. Naming wont be held to standard if users control it...


So far, I like 2 and 3 the best. 2 implies building [3] below, or at least a subset of the functionality. Building [3] would also make it easier to convert between formats. Just unclear how likely data across domains would be generalizable to this DIR. Would genomics/proteomics data fit this?



[1] implementation detail to choose where to be in the `index stores one fmt and tool converts locally <--> index stores every format` spectrum. Most likely in between: index stores every format but constructs them lazily)

[2] think of docker images. datasets can be expressed as instructions that construct it (some files from foo/dataset1 + some from bar/dataset2). This implies that a selecting sub-portions of a dataset could be a really useful mechanic.[3]

### selecting

[3] imagine selecting [n-m] rows of a given dataset. Unclear yet how this should work exactly, but i've ideas along a dataset intermediate representation (DIR), where data is expressed as points in a multi-dimensional space, and a dataset is expressed as a subspace, or intervals across some dimensions. This would work well even for tables, allowing one to select slices of a dataset with something like: <author>/<dataset>#<dimension>[:<low>[:<high>]]` e.g.

    lecun/norb#class          # points that have a class
    lecun/norb#class:car      # points that have class `car`
    lecun/norb#set:training   # points in the training set
    lecun/norb#y:0:10         # points where `0 <= y <= 10`

    (and of course, can specify multiple comma-delimited dimensions)

Or:

    lecun/norb#class          # points that have a class
    lecun/norb#class[car]     # points that have class `car`
    lecun/norb#set[training]  # points in the training set
    lecun/norb#y[0, 10]       # points where `0 <= y <= 10`
    lecun/norb#y]0, 10[       # points where `0 < y < 10`

This seems like a really powerful thing to enable. Unclear how to do it well at present. Lots and lots of edge cases.  This can come in later versions but must not close doors to it now. (Another note: i realize this basically is a dumber query string `?param=val`, problem with using a query string is these handles may have to be embedded in URLs :/ though i guess hashes are out in that case...)

