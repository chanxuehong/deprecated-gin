package response

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"
)

// bitmap = (io.ReaderFrom, stringWriter, http.Flusher, http.Hijacker, http.CloseNotifier)
const (
	ioReaderFromBitmap      = 1 << 4 // io.ReaderFrom
	stringWriterBitmap      = 1 << 3 // stringWriter
	httpFlusherBitmap       = 1 << 2 // http.Flusher
	httpHijackerBitmap      = 1 << 1 // http.Hijacker
	httpCloseNotifierBitmap = 1      // http.CloseNotifier
)

// stringWriter is the interface that wraps the WriteString method.
type stringWriter interface {
	WriteString(s string) (n int, err error)
}

func bitmap(w http.ResponseWriter) int {
	// get bitmap from cache
	typ := reflect.TypeOf(w)
	bitmapCachePtr := (*map[reflect.Type]int)(atomic.LoadPointer(&__bitmapCache))
	if bitmapCachePtr != nil {
		bitmapCache := *bitmapCachePtr
		if n, ok := bitmapCache[typ]; ok {
			return n
		}
	}

	var ok bool
	n := 0x00
	if _, ok = w.(io.ReaderFrom); ok {
		n |= ioReaderFromBitmap
	}
	if _, ok = w.(stringWriter); ok {
		n |= stringWriterBitmap
	}
	if _, ok = w.(http.Flusher); ok {
		n |= httpFlusherBitmap
	}
	if _, ok = w.(http.Hijacker); ok {
		n |= httpHijackerBitmap
	}
	if _, ok = w.(http.CloseNotifier); ok {
		n |= httpCloseNotifierBitmap
	}

	fmt.Println("match", n)

	// save bitmap to cache
	__bitmapCacheLock.Lock()
	var newBitmapCache map[reflect.Type]int
	if bitmapCachePtr != nil {
		bitmapCache := *bitmapCachePtr
		newBitmapCache = make(map[reflect.Type]int, len(bitmapCache)+1)
		for k, v := range bitmapCache {
			newBitmapCache[k] = v
		}
	} else {
		newBitmapCache = make(map[reflect.Type]int, 1)
	}
	newBitmapCache[typ] = n
	atomic.StorePointer(&__bitmapCache, unsafe.Pointer(&newBitmapCache))
	__bitmapCacheLock.Unlock()

	return n
}

var (
	__bitmapCacheLock sync.Mutex
	__bitmapCache     unsafe.Pointer // *map[reflect.Type]bitmap
)
