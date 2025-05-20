package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/qq2383/tunnel/config"
	"github.com/qq2383/tunnel/https"
	"github.com/qq2383/tunnel/process"
	"github.com/qq2383/tunnel/server"
	"github.com/qq2383/tunnel/user"
)

func main() {
	root, _ := os.Getwd()
	root = filepath.FromSlash(root)

	// logger.New(root, "tunneld", log.LstdFlags | log.Lshortfile | log.LUTC)

	fp := filepath.Join(root, "tunneld.cnf")
	var cnf Config
	config.Add(&cnf)
	err := config.Read(fp)
	if err != nil {
		log.Panicln(err)
	}

	ufp := filepath.Join(root, "user.db")
	user.Load(ufp)

	var wg sync.WaitGroup
	process.New(&wg)

	hs := https.New()
	process.Put("https", hs)

	ts := server.New()
	process.Put("tunnel", ts)

	process.Starts()

	wg.Wait()
	fmt.Println("exit")
}
