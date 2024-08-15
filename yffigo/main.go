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
	"log"
	"yffigo/yrs" // 根据实际路径调整导入路径
)

func main() {
	// 初始化一个YDoc
	doc := yrs.NewYDoc()
	if doc == nil {
		log.Fatal("Failed to create YDoc")
	}
	defer doc.Destroy() // 确保在退出前销毁YDoc

	doc.GetDocID()

	// 获取或创建一个YMap
	ymap := doc.GetYMap("exampleMap")
	if ymap == nil {
		log.Fatal("Failed to create or get YMap")
	}

	// 向YMap添加一些数据
	err := ymap.Set("key1", "value1")
	if err != nil {
		log.Fatalf("Failed to set value in YMap: %v", err)
	}

	// 添加更多数据
	err = ymap.Set("key2", "value2")
	if err != nil {
		log.Fatalf("Failed to set value in YMap: %v", err)
	}

	// 使用迭代器遍历YMap
	iter := ymap.NewIterator()
	for iter.HasNext() {
		key, value := iter.Next()
		fmt.Printf("Key: %s, Value: %s\n", key, value)
	}

	// 假设有一个保存或同步的方法
	err = doc.Save()
	if err != nil {
		log.Fatalf("Failed to save document: %v", err)
	}
}
