// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"hash/crc32"
	"math/rand"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// copy a by b  b->a
func CopyStructFields(a interface{}, b interface{}, fields ...string) (err error) {
	return copier.Copy(a, b)
}

func Wrap1(err error) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine())
}

func Wrap2[T any](a T, err error) (T, error) {
	if err != nil {
		return a, errors.Wrap(err, "==> "+printCallerNameAndLine())
	}
	return a, nil
}

func Wrap3[T any, V any](a T, b V, err error) (T, V, error) {
	if err != nil {
		return a, b, errors.Wrap(err, "==> "+printCallerNameAndLine())
	}
	return a, b, nil
}

func Wrap(err error, message string) error {
	return errors.Wrap(err, "==> "+printCallerNameAndLine()+message)
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, "==> "+printCallerNameAndLine()+message)
}

func printCallerNameAndLine() string {
	pc, _, line, _ := runtime.Caller(2)
	return runtime.FuncForPC(pc).Name() + "()@" + strconv.Itoa(line) + ": "
}

func GetSelfFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return cleanUpFuncName(runtime.FuncForPC(pc).Name())
}

func GetFuncName(skips ...int) string {
	skip := 1
	if len(skips) > 0 {
		skip = skips[0] + 1
	}
	pc, _, _, _ := runtime.Caller(skip)
	return cleanUpFuncName(runtime.FuncForPC(pc).Name())
}

func cleanUpFuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		return ""
	}
	return funcName[end+1:]
}

// Get the intersection of two slices
func Intersect(slice1, slice2 []int64) []int64 {
	m := make(map[int64]bool)
	n := make([]int64, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func Difference(slice1, slice2 []int64) []int64 {
	m := make(map[int64]bool)
	n := make([]int64, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}

// Get the intersection of two slices
func IntersectString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// Get the diff of two slices
func DifferenceString(slice1, slice2 []string) []string {
	m := make(map[string]bool)
	n := make([]string, 0)
	inter := IntersectString(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}

	for _, v := range slice2 {
		if !m[v] {
			n = append(n, v)
		}
	}
	return n
}

func OperationIDGenerator() string {
	return strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10)
}

func GetHashCode(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

func GenConversationIDForSingle(sendID, recvID string) string {
	l := []string{sendID, recvID}
	sort.Strings(l)
	return "si_" + strings.Join(l, "_")
}

func GenConversationUniqueKeyForGroup(groupID string) string {
	return groupID
}

func GenGroupConversationID(groupID string) string {
	return "sg_" + groupID
}

func GenConversationUniqueKeyForSingle(sendID, recvID string) string {
	l := []string{sendID, recvID}
	sort.Strings(l)
	return strings.Join(l, "_")
}

func GetNotificationConversationIDByConversationID(conversationID string) string {
	l := strings.Split(conversationID, "_")
	if len(l) > 1 {
		l[0] = "n"
		return strings.Join(l, "_")
	}
	return ""
}

func GetSelfNotificationConversationID(userID string) string {
	return "n_" + userID + "_" + userID
}

func GetSeqsBeginEnd(seqs []int64) (int64, int64) {
	if len(seqs) == 0 {
		return 0, 0
	}
	return seqs[0], seqs[len(seqs)-1]
}
