package yrs

// NOTE: There should be NO space between the comments and the `import "C"` line.
// The -ldl is sometimes necessary to fix linker errors about `dlsym`.

/*
#cgo LDFLAGS: ./lib/libyrs.a -ldl
#include "./include/libyrs.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

// YDoc represents a document in Yrs.
type YDoc struct {
	ptr *C.YDoc
}

// YTransaction represents a transaction in Yrs.
type YTransaction struct {
	ptr *C.YTransaction
}

// YInput represents a data structure used to pass input values into a shared document.
type YInput struct {
	Tag  int8
	Len  uint32
	Data [8]byte
}

// YInput types
const (
	Y_JSON_BOOL    = C.Y_JSON_BOOL
	Y_JSON_NUM     = C.Y_JSON_NUM
	Y_JSON_INT     = C.Y_JSON_INT
	Y_JSON_STR     = C.Y_JSON_STR
	Y_JSON_BUF     = C.Y_JSON_BUF
	Y_JSON_ARR     = C.Y_JSON_ARR
	Y_JSON_MAP     = C.Y_JSON_MAP
	Y_JSON_NULL    = C.Y_JSON_NULL
	Y_JSON_UNDEF   = C.Y_JSON_UNDEF
	Y_ARRAY        = C.Y_ARRAY
	Y_MAP          = C.Y_MAP
	Y_TEXT         = C.Y_TEXT
	Y_XML_ELEM     = C.Y_XML_ELEM
	Y_XML_TEXT     = C.Y_XML_TEXT
	Y_XML_FRAG     = C.Y_XML_FRAG
	Y_DOC          = C.Y_DOC
	Y_WEAK_LINK    = C.Y_WEAK_LINK
	Y_UNDEFINED    = C.Y_UNDEFINED
	Y_TRUE         = C.Y_TRUE
	Y_FALSE        = C.Y_FALSE
	Y_OFFSET_BYTES = C.Y_OFFSET_BYTES
	Y_OFFSET_UTF16 = C.Y_OFFSET_UTF16
)

// NewYInputBool creates a YInput for a boolean value.
func NewYInputBool(value bool) *YInput {
	var flag C.int8_t
	if value {
		flag = C.Y_TRUE
	} else {
		flag = C.Y_FALSE
	}
	return &YInput{
		Tag:  Y_JSON_BOOL,
		Len:  1,
		Data: *(*[8]byte)(unsafe.Pointer(&flag)),
	}
}

// NewYInputInt creates a YInput for an integer value.
func NewYInputInt(value int64) *YInput {
	return &YInput{
		Tag:  Y_JSON_INT,
		Len:  1,
		Data: *(*[8]byte)(unsafe.Pointer(&value)),
	}
}

// NewYInputFloat creates a YInput for a floating-point value.
func NewYInputFloat(value float64) *YInput {
	return &YInput{
		Tag:  Y_JSON_NUM,
		Len:  1,
		Data: *(*[8]byte)(unsafe.Pointer(&value)),
	}
}

// NewYInputString creates a YInput for a string value.
func NewYInputString(value string) *YInput {
	cstr := C.CString(value)
	return &YInput{
		Tag:  Y_JSON_STR,
		Len:  1,
		Data: *(*[8]byte)(unsafe.Pointer(&cstr)),
	}
}

// Free releases the resources associated with a YInput (if necessary).
func (input *YInput) Free() {
	if input.Tag == Y_JSON_STR {
		// Correctly handle the free operation for the CString
		cstr := *(*C.CString)(unsafe.Pointer(&input.Data[0]))
		C.free(unsafe.Pointer(cstr))
	}
}

// toC converts a YInput to C.struct_YInput.
func (input *YInput) toC() C.struct_YInput {
	return C.struct_YInput{
		tag:   C.int8_t(input.Tag),
		len:   C.uint32_t(input.Len),
		value: *(*C.union_YInputContent)(unsafe.Pointer(&input.Data)),
	}
}

// YOutput represents the output from Yrs API methods.
type YOutput struct {
	tag  int8
	len  uint32
	data [8]byte // Handle the union manually data C.union_YOutputContent
}

// GetValueAsString returns the value as a string if the YOutput is of type string.
func (output *YOutput) GetValueAsString() string {
	if output.tag == Y_JSON_STR {
		// Convert the data to a string pointer and dereference it
		strPtr := (*C.char)(unsafe.Pointer(&output.data[0]))
		return C.GoString(strPtr)
	}
	return ""
}

// GetValueAsInt returns the value as an integer if the YOutput is of type int.
func (output *YOutput) GetValueAsInt() int64 {
	if output.tag == Y_JSON_INT {
		// Convert the data to an int64 pointer and dereference it
		return *(*int64)(unsafe.Pointer(&output.data[0]))
	}
	return 0
}

// GetValueAsFloat returns the value as a float64 if the YOutput is of type float.
func (output *YOutput) GetValueAsFloat() float64 {
	if output.tag == Y_JSON_NUM {
		// Convert the data to a float64 pointer and dereference it
		return *(*float64)(unsafe.Pointer(&output.data[0]))
	}
	return 0.0
}

// GetValueAsBool returns the value as a boolean if the YOutput is of type bool.
func (output *YOutput) GetValueAsBool() bool {
	if output.tag == Y_JSON_BOOL {
		// Convert the data to a boolean flag
		flag := *(*C.int8_t)(unsafe.Pointer(&output.data[0]))
		return flag == Y_TRUE
	}
	return false
}

// GetValueAsBytes returns the value as a byte slice if the YOutput is of type buffer.
func (output *YOutput) GetValueAsBytes() []byte {
	if output.tag == Y_JSON_BUF {
		// Convert the data to a byte pointer and use C.GoBytes
		bufPtr := unsafe.Pointer(&output.data[0])
		return C.GoBytes(bufPtr, C.int(output.len))
	}
	return nil
}

// GetValueType returns the type of the value stored in YOutput.
func (output *YOutput) GetValueType() int8 {
	return output.tag
}

// Branch represents a common shared data type in Yrs.
type Branch struct {
	ptr *C.Branch
}

// YMapIter creates an iterator for a YMap.
func (b *Branch) YMapIter(txn *YTransaction) *YMapIter {
	iter := C.ymap_iter(b.ptr, txn.ptr)
	return &YMapIter{ptr: iter}
}

// YArrayIter represents an iterator for YArray.
type YArrayIter struct {
	ptr *C.YArrayIter
}

// YMapIter represents an iterator for YMap.
type YMapIter struct {
	ptr *C.YMapIter
}

// NewYDoc creates a new YDoc instance.
func NewYDoc() *YDoc {
	doc := C.ydoc_new()
	return &YDoc{ptr: doc}
}

// Destroy releases the memory associated with a YDoc.
func (d *YDoc) Destroy() {
	C.ydoc_destroy(d.ptr)
	d.ptr = nil
}

// ReadTransaction starts a new read-only transaction on a document.
func (d *YDoc) ReadTransaction() *YTransaction {
	txn := C.ydoc_read_transaction(d.ptr)
	return &YTransaction{ptr: txn}
}

// WriteTransaction starts a new read-write transaction on a document.
func (d *YDoc) WriteTransaction(origin string) *YTransaction {
	corigin := C.CString(origin)
	defer C.free(unsafe.Pointer(corigin))
	txn := C.ydoc_write_transaction(d.ptr, C.uint32_t(len(origin)), corigin)
	return &YTransaction{ptr: txn}
}

// Commit commits the transaction and releases its resources.
func (t *YTransaction) Commit() {
	C.ytransaction_commit(t.ptr)
	t.ptr = nil
}

// GetYText returns a YText branch by its name.
func (d *YDoc) GetYText(name string) *Branch {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	branch := C.ytext(d.ptr, cname)
	return &Branch{ptr: branch}
}

// GetYArray returns a YArray branch by its name.
func (d *YDoc) GetYArray(name string) *Branch {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	branch := C.yarray(d.ptr, cname)
	return &Branch{ptr: branch}
}

// GetYMap returns a YMap branch by its name.
func (d *YDoc) GetYMap(name string) *Branch {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	branch := C.ymap(d.ptr, cname)
	return &Branch{ptr: branch}
}

// YArrayLen returns the length of an array stored in YArray.
func (b *Branch) YArrayLen() uint32 {
	return uint32(C.yarray_len(b.ptr))
}

// YTextString returns the string content of a YText.
func (b *Branch) YTextString(txn *YTransaction) string {
	cstr := C.ytext_string(b.ptr, txn.ptr)
	defer C.ystring_destroy(cstr)
	return C.GoString(cstr)
}

// YArrayGet returns an element from a YArray by index.
func (b *Branch) YArrayGet(txn *YTransaction, index uint32) *YOutput {
	output := C.yarray_get(b.ptr, txn.ptr, C.uint32_t(index))
	return &YOutput{
		tag:  int8(output.tag),
		len:  uint32(output.len),
		data: output.value,
	}
}

// YArrayInsert inserts elements into a YArray at the specified index.
func (b *Branch) YArrayInsert(txn *YTransaction, index uint32, items []YInput) {
	cItems := make([]C.struct_YInput, len(items))
	for i, item := range items {
		cItems[i] = item.toC()
	}
	C.yarray_insert_range(b.ptr, txn.ptr, C.uint32_t(index), &cItems[0], C.uint32_t(len(items)))
}

// YArrayRemoveRange removes elements from a YArray.
func (b *Branch) YArrayRemoveRange(txn *YTransaction, index, length uint32) {
	C.yarray_remove_range(b.ptr, txn.ptr, C.uint32_t(index), C.uint32_t(length))
}

// YArrayIterNext retrieves the next element from a YArray iterator.
func (iter *YArrayIter) YArrayIterNext() *YOutput {
	output := C.yarray_iter_next(iter.ptr)
	if output == nil {
		return nil
	}
	return &YOutput{
		tag:  int8(output.tag),
		len:  uint32(output.len),
		data: output.value,
	}
}

// YMapInsert inserts a key-value pair into a YMap. Correct the handling when calling value.toC() in C function calls
func (b *Branch) YMapInsert(txn *YTransaction, key string, value *YInput) {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	cValue := value.toC()
	C.ymap_insert(b.ptr, txn.ptr, ckey, &cValue)
}

// YMapRemove removes a key-value pair from a YMap.
func (b *Branch) YMapRemove(txn *YTransaction, key string) bool {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	return C.ymap_remove(b.ptr, txn.ptr, ckey) != 0
}

// YMapIterNext retrieves the next entry from a YMap iterator.
func (iter *YMapIter) YMapIterNext() *YMapEntry {
	entry := C.ymap_iter_next(iter.ptr)
	if entry == nil {
		return nil
	}
	return &YMapEntry{
		Key: C.GoString(entry.key),
		Value: YOutput{
			tag:  int8(entry.value.tag),
			len:  uint32(entry.value.len),
			data: *(*[8]byte)(unsafe.Pointer(&entry.value.value)),
		},
	}
}

// YMapEntry represents a key-value entry in a YMap.
type YMapEntry struct {
	Key   string
	Value YOutput
}

// YMapEntryNext retrieves the next entry from a YMap iterator.
func (iter *YMapIter) YMapEntryNext() *YMapEntry {
	entry := C.ymap_iter_next(iter.ptr)
	if entry == nil {
		return nil
	}
	return &YMapEntry{
		Key: C.GoString(entry.key),
		Value: YOutput{
			tag:  int8(entry.value.tag),
			len:  uint32(entry.value.len),
			data: entry.value.value,
		},
	}
}

// YMapIterDestroy releases the resources associated with a YMap iterator.
func (iter *YMapIter) YMapIterDestroy() {
	C.ymap_iter_destroy(iter.ptr)
	iter.ptr = nil
}

// GetValue returns the value as a string (for example).
func (entry *YMapEntry) GetValueAsString() string {
	if entry.Value.tag == Y_JSON_STR {
		// Convert the data to a string pointer and dereference it
		strPtr := (*C.char)(unsafe.Pointer(&entry.Value.data[0]))
		return C.GoString(strPtr)
	}
	return ""
}

// GetValueAsInt returns the value as an integer (for example).
func (entry *YMapEntry) GetValueAsInt() int64 {
	if entry.Value.tag == Y_JSON_INT {
		// Convert the data to an int64 pointer and dereference it
		return *(*int64)(unsafe.Pointer(&entry.Value.data[0]))
	}
	return 0
}

// YMapGet retrieves a value from a YMap by key.
func (b *Branch) YMapGet(txn *YTransaction, key string) *YOutput {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	output := C.ymap_get(b.ptr, txn.ptr, ckey)
	if output == nil {
		return nil
	}
	return &YOutput{
		tag:  int8(output.tag),
		len:  uint32(output.len),
		data: output.value,
	}
}

// Destroy releases the resources associated with a YMap iterator.
func (iter *YMapIter) Destroy() {
	C.ymap_iter_destroy(iter.ptr)
	iter.ptr = nil
}
