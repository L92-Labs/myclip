// cliphelper drives MyClip's Linux smoke test (not shipped):
//
//	cliphelper write TEXT [hint] [secs]  own the clipboard with TEXT (and the
//	                                     KDE password-manager hint) for secs
//	cliphelper read                      print the clipboard text
//	cliphelper dump DB                   print the stored history, newest first
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/transfer"

	"github.com/L92-Labs/myclip/internal/store"
)

func main() {
	switch os.Args[1] {
	case "dump":
		s, err := store.Open(context.Background(), os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		items, err := s.Load(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		for _, it := range items {
			fmt.Printf("%d\t%d\t%q\t%s\n", it.ID, it.NumberOfCopies, it.Text(), it.Application)
		}
		return
	}
	mygo.App.OnWindowAllClosed(func() {})
	mygo.App.WhenReady(func() {
		switch os.Args[1] {
		case "write":
			reps := []transfer.Representation{transfer.Bytes(transfer.Text, []byte(os.Args[2]))}
			if len(os.Args) > 3 && os.Args[3] == "hint" {
				reps = append(reps, transfer.Bytes("x-kde-passwordManagerHint", []byte("secret")))
			}
			if err := mygo.Clipboard.Write(transfer.New(transfer.NewItem(reps...))); err != nil {
				log.Fatal(err)
			}
			secs := 3
			if len(os.Args) > 4 {
				secs, _ = strconv.Atoi(os.Args[4])
			}
			go func() { time.Sleep(time.Duration(secs) * time.Second); mygo.App.Quit() }()
		case "read":
			fmt.Println(mygo.Clipboard.ReadText())
			fmt.Println("formats:", mygo.Clipboard.AvailableFormats())
			mygo.App.Quit()
		}
	})
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}
