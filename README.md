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

$ cat .srce/index
6f53920efa5dee5a54d4f5f6b07f5d7f07c1710f blob LICENSE

$ ./exe/srce commit "first commit"

$ ./exe/srce rev-parse HEAD
e160570596bcfc89bc296a0d4118bbe44637cabc

$ ./exe/srce cat-file -p HEAD
tree 9ddc1cd6eb70ab9f9c0e4537a6b7ca7bc2bc13bd
author steinbro

first commit
```
