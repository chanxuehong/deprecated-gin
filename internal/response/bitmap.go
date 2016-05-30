package response

import (
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
	httpCloseNotifierBitmap = 1 << 0 // http.CloseNotifier
)

// stringWriter is the interface that wraps the WriteString method.
type stringWriter interface {
	WriteString(s string) (n int, err error)
}

func bitmap(v interface{}) int {
	// get bitmap from cache
	typ := reflect.TypeOf(v)
	bitmapCachePtr := (*map[reflect.Type]int)(atomic.LoadPointer(&__bitmapCache))
	if bitmapCachePtr != nil {
		bitmapCache := *bitmapCachePtr
		if n, ok := bitmapCache[typ]; ok {
			return n
		}
	}

	// cache miss
	n := 0
	var ok bool
	if _, ok = v.(io.ReaderFrom); ok {
		n |= ioReaderFromBitmap
	}
	if _, ok = v.(stringWriter); ok {
		n |= stringWriterBitmap
	}
	if _, ok = v.(http.Flusher); ok {
		n |= httpFlusherBitmap
	}
	if _, ok = v.(http.Hijacker); ok {
		n |= httpHijackerBitmap
	}
	if _, ok = v.(http.CloseNotifier); ok {
		n |= httpCloseNotifierBitmap
	}

	// save bitmap to cache
	__bitmapCacheLock.Lock()
	if bitmapCachePtr != nil {
		bitmapCache := *bitmapCachePtr
		newBitmapCache := make(map[reflect.Type]int, len(bitmapCache)+1)
		for k, v := range bitmapCache {
			newBitmapCache[k] = v
		}
		newBitmapCache[typ] = n
		atomic.StorePointer(&__bitmapCache, unsafe.Pointer(&newBitmapCache))
	} else {
		newBitmapCache := make(map[reflect.Type]int, 1)
		newBitmapCache[typ] = n
		atomic.StorePointer(&__bitmapCache, unsafe.Pointer(&newBitmapCache))
	}
	__bitmapCacheLock.Unlock()

	return n
}

var (
	__bitmapCacheLock sync.Mutex
	__bitmapCache     unsafe.Pointer // *map[reflect.Type]bitmap
)
