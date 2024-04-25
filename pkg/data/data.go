package data

import (
	"io/fs"
	"log"
	"os"
	"time"
)

// The data package is used to maintain "in-memory" snapshots
// of the disk.
//
// In the first mode of data sychrnonisation, the data package
// is meant to store a snapshot of the file system that is being tracked.
// These snaps are used to compare the current state of FS on the disk with
// its snapshot. WIth this, if any change is encountered, it generates a diff
// and sends it to the comms server, which must create a response and share
// the change with all its clients.
//
// The disk is to be tracked indefinitely.

type NodeData struct {
	// This node has the fs attached to it
	FileSystem fs.FS
	BaseDir    string

	// map of file path to mod time
	FileTrack map[string]time.Time
}

func (n *NodeData) LoadFileSystem(dirwatch string) {
	fileSystem := os.DirFS(dirwatch)
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		// log.Println("Evaluating path: ", path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			// add dir to file track
			// log.Println("Encountered DIR")
		} else if d.Type().IsRegular() {
			fileInfo, err := d.Info()
			if err != nil {
				log.Fatal(err)
				return err
			}
			// log.Printf("Loaded tracking for %s", path)
			n.FileTrack[path] = fileInfo.ModTime()
			// buff, err := fs.ReadFile(fileSystem, path)
			// if err != nil {
			// 	log.Fatalf("Failed to read file")
			// }
			// log.Println(string(buff))
		}

		return nil
	})
}
