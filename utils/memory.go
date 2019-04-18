package utils

import (
	"os/exec"
	"runtime"
)

// FlushMem will try to reduce the memory usage of the container it is running in
// run this after a build
func FlushMem() {
	defer func() {
		runtime.GC()
		runtime.GC()
	}()
	if runtime.GOOS != "linux" {
		return
	}
	log.Println("Flushing memory.")
	// it's ok if these fail
	// flush memory buffers
	err := exec.Command("sync").Run()
	if err != nil {
		log.Printf("flushMem error (sync): %v", err)
	}
	// clear page cache
	err = exec.Command("bash", "-c", "echo 1 > /proc/sys/vm/drop_caches").Run()
	if err != nil {
		log.Printf("flushMem error (page cache): %v", err)
	}
}
