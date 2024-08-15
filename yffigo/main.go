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
	// 假设你已经有了 YMap 和 YTransaction 的实例
	var yMap *yrs.YMap
	var txn *yrs.YTransaction

	// 迭代 YMap 中的所有键值对
	iter := yMap.YMapIter(txn)
	defer iter.YMapIterDestroy()

	for entry := iter.YMapEntryNext(); entry != nil; entry = iter.YMapEntryNext() {
		valueStr, err := entry.GetValueAsString()
		if err != nil {
			fmt.Printf("Error retrieving string value for key %s: %v\n", entry.Key, err)
			continue
		}
		fmt.Printf("Key: %s, Value: %s\n", entry.Key, valueStr)
	}
}
