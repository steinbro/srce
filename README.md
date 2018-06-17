# srce

An over-simplified git-like revision control system written in Go, inspired by
[RGit](https://github.com/JoelQ/rgit). Not yet capable of tracking its own
source.

Also a word meaning "heart" in various Slavic languages.

```
$ make all
go build -o exe/srce ./bin/cmd/srce
go build -o exe/srce-add ./bin/cmd/srce-add
[...]

$ ./exe/srce init
srce initialized in .srce

$ ./exe/srce add LICENSE

$ ./exe/srce status
M	LICENSE

$ ./exe/srce commit "first commit"

$ ./exe/srce rev-parse HEAD
e160570596bcfc89bc296a0d4118bbe44637cabc

$ ./exe/srce log
commit e160570596bcfc89bc296a0d4118bbe44637cabc
Author: steinbro
Date:   Mon May 21 20:25:57 2018 -0400

	first commit

$ ./exe/srce reflog
e1605705 HEAD@{0}: commit: first commit

```
