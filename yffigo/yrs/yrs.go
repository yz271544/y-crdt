package yrs

// NOTE: There should be NO space between the comments and the `import "C"` line.
// The -ldl is sometimes necessary to fix linker errors about `dlsym`.

/*
#cgo LDFLAGS: -L ./lib -lyrs -lm -ldl
#include "./include/libyrs.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// YDoc wraps a Yrs document pointer
type YDoc struct {
	ptr *C.YDoc
}

// NewYDoc creates a new Yrs document
func NewYDoc() *YDoc {
	return &YDoc{
		ptr: C.ydoc_new(),
	}
}

// Destroy releases resources held by the YDoc
func (d *YDoc) Destroy() {
	if d.ptr != nil {
		C.ydoc_destroy(d.ptr)
		d.ptr = nil
	}
}

// GetDocID returns the unique ID of the YDoc
func (d *YDoc) GetDocID() uint64 {
	return uint64(C.ydoc_id(d.ptr))
}

// BeginTransaction starts a read-write transaction on the YDoc
func (d *YDoc) BeginTransaction() *YTransaction {
	return &YTransaction{
		ptr: C.ydoc_write_transaction(d.ptr, 0, nil),
	}
}

// YTransaction wraps a Yrs transaction pointer
type YTransaction struct {
	ptr *C.YTransaction
}

// Commit commits the transaction and cleans up resources
func (t *YTransaction) Commit() {
	if t.ptr != nil {
		C.ytransaction_commit(t.ptr)
		t.ptr = nil
	}
}

// Example usage of a complex function
// ApplyUpdate applies an update to the document within a transaction
func (t *YTransaction) ApplyUpdate(update []byte) error {
	cUpdate := (*C.char)(unsafe.Pointer(&update[0]))
	length := C.uint32_t(len(update))

	result := C.ytransaction_apply(t.ptr, cUpdate, length)
	if result != 0 {
		return errors.New("failed to apply update")
	}
	return nil
}

// YOutput wraps the C YOutput structure
type YOutput struct {
	ptr *C.YOutput
}

// NewYOutput creates a new YOutput for a given C YOutput pointer
func NewYOutput(ptr *C.YOutput) *YOutput {
	if ptr == nil {
		return nil
	}
	return &YOutput{ptr: ptr}
}

// GetValueAsString tries to convert the YOutput value to a Go string
func (yo *YOutput) GetValueAsString() (string, error) {
	if yo.ptr.tag == C.Y_JSON_STR {
		// 解析 ptr.value 作为字符串
		strPtr := *(**C.char)(unsafe.Pointer(&yo.ptr.value[0]))
		return C.GoString(strPtr), nil
	}
	return "", errors.New("value is not a string")
}

// GetValueAsInt tries to convert the YOutput value to a Go int
func (yo *YOutput) GetValueAsInt() (int64, error) {
	if yo.ptr.tag == C.Y_JSON_INT {
		// 解析 ptr.value 作为 int64
		intVal := *(*int64)(unsafe.Pointer(&yo.ptr.value[0]))
		return int64(intVal), nil
	}
	return 0, errors.New("value is not an integer")
}

// GetValueAsBool tries to convert the YOutput value to a Go bool
func (yo *YOutput) GetValueAsBool() (bool, error) {
	if yo.ptr.tag == C.Y_JSON_BOOL {
		// 解析 ptr.value 作为 bool
		boolVal := *(*C.int8_t)(unsafe.Pointer(&yo.ptr.value[0]))
		return boolVal == C.Y_TRUE, nil
	}
	return false, errors.New("value is not a boolean")
}

// Free cleans up any allocated C resources
func (yo *YOutput) Free() {
	if yo.ptr != nil {
		C.free(unsafe.Pointer(yo.ptr))
		yo.ptr = nil
	}
}

// YMap wraps a Yrs map pointer
type YMap struct {
	ptr *C.Branch // Assuming Branch is already defined somewhere
}

// YMapIter wraps a Yrs map iterator
type YMapIter struct {
	ptr *C.YMapIter
}

// YMapEntry wraps a single key-value entry from a YMap
type YMapEntry struct {
	Key   string
	Value *YOutput
}

// NewYMap creates a new YMap wrapper
func NewYMap(ptr *C.Branch) *YMap {
	return &YMap{ptr: ptr}
}

// YMapIter starts an iterator over the map entries
func (m *YMap) YMapIter(txn *YTransaction) *YMapIter {
	return &YMapIter{
		ptr: C.ymap_iter(m.ptr, txn.ptr),
	}
}

// YMapIterDestroy cleans up the iterator
func (iter *YMapIter) YMapIterDestroy() {
	if iter.ptr != nil {
		C.ymap_iter_destroy(iter.ptr)
		iter.ptr = nil
	}
}

// YMapEntryNext returns the next entry in the map or nil if done
func (iter *YMapIter) YMapEntryNext() *YMapEntry {
	cEntry := C.ymap_iter_next(iter.ptr)
	if cEntry == nil {
		return nil
	}
	goEntry := &YMapEntry{
		Key:   C.GoString(cEntry.key),
		Value: NewYOutput(&cEntry.value), // Assuming NewYOutput is correctly defined
	}
	return goEntry
}

func (entry *YMapEntry) GetValueAsString() (string, error) {
	if entry.Value.ptr.tag == C.Y_JSON_STR {
		// 将value数组的第一个元素地址转换为 *C.char
		strPtr := *(**C.char)(unsafe.Pointer(&entry.Value.ptr.value[0]))
		return C.GoString(strPtr), nil
	}
	return "", errors.New("value is not a string")
}
