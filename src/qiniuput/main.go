package main

import (
	"flag"
	"log"
	"net/url"
	"os"
)

import (
	"github.com/qiniu/api/resumable/io"
)

func main() {
	token := flag.String("t", "", "upload token")
	file := flag.String("f", "", "file path")
	key := flag.String("k", "", "file key")
	custom := flag.String("x", "", "custom args")
	flag.Parse()
	if token == nil || file == nil || key == nil {
		log.Fatalln("invalid args")
		return
	}

	f, err := os.Open(*file)
	if err != nil {
		log.Fatalln("file not exist")
		return
	}
	stat, err := f.Stat()
	if err != nil || stat.IsDir() {
		log.Fatalln("invalid file")
		return
	}

	blockNotify := func(blkIdx int, blkSize int, ret *io.BlkputRet) {
		log.Println("size", stat.Size(), "block id", blkIdx, "offset", ret.Offset)
	}

	params := map[string]string{}
	extra := &io.PutExtra{
		ChunkSize: 8192,
		Notify:    blockNotify,
		Params:    params,
	}
	if custom != nil && *custom != "" {
		values, err := url.ParseQuery(*custom)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		for k, v := range values {
			params["x:"+k] = v[0]
		}
		extra.Params = params
	}
	var ret io.PutRet
	err = io.PutFile(nil, &ret, *token, *key, *file, extra)
	if err != nil {
		log.Fatalln(err)
	}
}
