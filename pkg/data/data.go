package data

import (
	"io/fs"
	"log"
	"os"
	"sync"
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
	FileTrackMap FileTrackMap
}

type FileTrackMap struct {
	FileTrack map[string]time.Time
	Lock sync.Mutex
}

func (n *NodeData) LoadFileSystem(dirwatch string) {
	fileSystem := os.DirFS(dirwatch)
	fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		// log.Println("Evaluating path: ", path)
		if err != nil {
			return err
		}

		if d.IsDir() {
			// log.Println("Encountered DIR")
		} else if d.Type().IsRegular() {
			fileInfo, err := d.Info()
			if err != nil {
				log.Fatal(err)
				return err
			}
			// log.Printf("Loaded tracking for %s", path)
			n.FileTrackMap.Lock.Lock()
			n.FileTrackMap.FileTrack[path] = fileInfo.ModTime()
			n.FileTrackMap.Lock.Unlock()
		}

		return nil
	})
}
