GOMBZ
=====

This is a [Go][golang] library that provides a serializable data structure
for 3d models and animations. Basically, it provides a set of data structures that represent 3d models, bones
and animations as well as functions to decode and encode them. The serialization
process uses BSON notation and then is run through [zlib][zlib-link].

Additionally, it provides a compiler to take any file that [Assimp][assimp-link]
supports and turns it into a compressed binary file that can be read by this library.


Requirements
------------

This does require `cgo` which means that `gcc` should be in your path
when trying to build using this library.

Software requirements:

* [Mathgl][mgl] - for 3d math
* [BSON][bson-link] - a BSON implementation in Go for binary serialization

The gombz compiler also needs:

* [Assimp][assimp-link] version 3.1 - tested with this version
* [Assimp-go][assimpgo-link] - Go wrappers for Assimp


Compiler Installation (gombzc)
------------------------------

The gombz compiler called `gombzc` can be installed to your `$GOPATH/bin` folder
by using the following command:

```bash
go install github.com/tbogdala/gombz/cmd/gombzc
```

TODO
----

* Documentation
* Better command-line flags to better control ASSIMP
* Consider supporting more than one mesh
* Consider including bitangents? Send feedback or create an issue if you
  feel strongly about this


LICENSE
=======

Gombz is released under the BSD license. See the [LICENSE][license-link] file for more details.


[golang]: https://golang.org/
[license-link]: https://raw.githubusercontent.com/tbogdala/gombz/master/LICENSE
[mgl]: https://github.com/go-gl/mathgl
[bson-link]: http://gopkg.in/mgo.v2/bson
[zlib-link]: https://golang.org/pkg/compress/zlib/
[assimp-link]: http://assimp.sourceforge.net/
[assimpgo-link]: http://www.github.com/tbogdala/assimp-go
