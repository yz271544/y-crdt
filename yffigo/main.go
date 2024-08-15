/*
Copyright 2024 The y-crdt Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    <http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"yffigo/yrs"
)

func main() {
	doc := yrs.NewYDoc()
	defer doc.Destroy()

	txn := doc.WriteTransaction("")
	defer txn.Commit()

	// 创建一个 YMap 类型
	yMap := doc.GetYMap("example_map")

	// 插入一个字符串键值对到 YMap
	inputStr := yrs.NewYInputString("value1")
	defer inputStr.Free()
	yMap.YMapInsert(txn, "key1", inputStr)

	// 插入一个整数键值对到 YMap
	inputInt := yrs.NewYInputInt(123)
	defer inputInt.Free()
	yMap.YMapInsert(txn, "key2", inputInt)

	// 获取键为 "key1" 的值
	output := yMap.YMapGet(txn, "key1")
	if output != nil {
		fmt.Println("Key1 Value:", output.GetValueAsString())
	}

	// 获取键为 "key2" 的值
	output = yMap.YMapGet(txn, "key2")
	if output != nil {
		fmt.Println("Key2 Value:", output.GetValueAsInt())
	}

	// 迭代 YMap 中的所有键值对
	iter := yMap.YMapIter(txn)
	defer iter.YMapIterDestroy()

	for entry := iter.YMapEntryNext(); entry != nil; entry = iter.YMapEntryNext() {
		fmt.Printf("Key: %s, Value: %v\n", entry.Key, entry.GetValueAsString()) // 这里假设所有值都是字符串
	}
}
