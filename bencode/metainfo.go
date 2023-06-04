package bencode

/*
There is also a key length or a key files,
but not both or neither.
If length is present then the download represents a single file,
otherwise it represents a set of files which go in a directory
structure.
*/
type metainfo struct {
	announce string
	info     struct {
		name         string
		piece_length int64
		pieces       string
		files        []string
	}
}
