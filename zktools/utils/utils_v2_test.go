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
	"fmt"
	"reflect"
	"testing"
)

func TestDistinct(t *testing.T) {
	arr := []int{1, 1, 1, 4, 4, 5, 2, 3, 3, 3, 6}
	fmt.Println(Distinct(arr))
}

func TestDeleteAt(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(Delete(arr, 0, 1, -1, -2))
	fmt.Println(Delete(arr))
	fmt.Println(Delete(arr, 1))
}

func TestSliceToMap(t *testing.T) {
	type Item struct {
		ID   string
		Name string
	}
	list := []Item{
		{ID: "111", Name: "111"},
		{ID: "222", Name: "222"},
		{ID: "333", Name: "333"},
	}

	m := SliceToMap(list, func(t Item) string {
		return t.ID
	})

	fmt.Printf("%+v\n", m)

}

func TestIndexOf(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	fmt.Println(IndexOf(3, arr...))

}

func TestSort(t *testing.T) {
	arr := []int{1, 1, 1, 4, 4, 5, 2, 3, 3, 3, 6}
	fmt.Println(Sort(arr, false))
}

func TestBothExist(t *testing.T) {
	arr1 := []int{1, 1, 1, 4, 4, 5, 2, 3, 3, 3, 6}
	arr2 := []int{6, 1, 3}
	arr3 := []int{5, 1, 3, 6}
	fmt.Println(BothExist(arr1, arr2, arr3))
}

func TestCompleteAny(t *testing.T) {
	type Item struct {
		ID    int
		Value string
	}

	ids := []int{1, 2, 3, 4, 5, 6, 7, 8}

	var list []Item

	for _, id := range ids {
		list = append(list, Item{
			ID:    id,
			Value: fmt.Sprintf("%d", id*1000),
		})
	}

	DeleteAt(&list, -1)
	DeleteAt(&ids, -1)

	ok := Complete(ids, Slice(list, func(t Item) int {
		return t.ID
	}))

	fmt.Printf("%+v\n", ok)
}

func TestStructFieldNotNilReplace(t *testing.T) {
	type Req struct {
		GroupID      string `json:"groupID"`
		GroupName    string `json:"groupName"`
		Notification string `json:"notification"`
		Introduction string `json:"introduction"`
		Count        int64  `json:"faceURL"`
		OwnerUserID  string `json:"ownerUserID"`
	}

	tests := []struct {
		name string
		req  Req
		resp Req
		want Req
	}{
		{
			name: "One by one conversion",
			req: Req{
				GroupID:      "groupID",
				GroupName:    "groupName",
				Notification: "notification",
				Introduction: "introduction",
				Count:        123,
				OwnerUserID:  "ownerUserID",
			},
			resp: Req{
				GroupID:      "ID",
				GroupName:    "Name",
				Notification: "notification",
				Introduction: "introduction",
				Count:        456,
				OwnerUserID:  "ownerUserID",
			},
			want: Req{
				GroupID:      "groupID",
				GroupName:    "groupName",
				Notification: "notification",
				Introduction: "introduction",
				Count:        123,
				OwnerUserID:  "ownerUserID",
			},
		},
		{
			name: "Changing the values of some fields",
			req: Req{
				GroupID:      "groupID",
				GroupName:    "groupName",
				Notification: "",
				Introduction: "",
				Count:        123,
				OwnerUserID:  "ownerUserID",
			},
			resp: Req{
				GroupID:      "ID",
				GroupName:    "Name",
				Notification: "notification",
				Introduction: "introduction",
				Count:        456,
				OwnerUserID:  "ownerUserID",
			},
			want: Req{
				GroupID:      "groupID",
				GroupName:    "groupName",
				Notification: "notification",
				Introduction: "introduction",
				Count:        123,
				OwnerUserID:  "ownerUserID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StructFieldNotNilReplace(&tt.resp, &tt.req)
			fmt.Println(tt.resp)
			if !reflect.DeepEqual(tt.want, tt.resp) {
				t.Errorf("%v have a err,%v", tt.name, tt.want)
			}
		})
	}

	type Req11 struct {
		GroupID      *string `json:"groupID"`
		GroupName    *string `json:"groupName"`
		Notification *string `json:"notification"`
		Introduction *string `json:"introduction"`
		Count        *int64  `json:"faceURL"`
		OwnerUserID  *string `json:"ownerUserID"`
	}

	type Req1 struct {
		Re  []*Req
		Re1 *Req
		Re2 Req11
	}
	r := Req{
		GroupID:      "groupID1",
		GroupName:    "groupName2",
		Notification: "1",
		Introduction: "1",
		Count:        123,
		OwnerUserID:  "ownerUserID1",
	}
	tests1 := []struct {
		name string
		req  Req1
		resp Req1
		want Req1
	}{
		{
			name: "name",
			req: Req1{
				Re: []*Req{
					{
						GroupID:      "groupID1",
						GroupName:    "groupName2",
						Notification: "1",
						Introduction: "1",
						Count:        123,
						OwnerUserID:  "ownerUserID1",
					},
					{
						GroupID:      "groupID2",
						GroupName:    "groupName2",
						Notification: "2",
						Introduction: "2",
						Count:        456,
						OwnerUserID:  "ownerUserID2",
					},
				},
				Re1: &r,
				Re2: Req11{
					GroupID:      &r.GroupID,
					GroupName:    &r.GroupName,
					Notification: &r.Notification,
					Introduction: &r.Introduction,
					Count:        &r.Count,
					OwnerUserID:  &r.OwnerUserID,
				},
			},
			resp: Req1{},
			want: Req1{
				Re: []*Req{
					{
						GroupID:      "groupID1",
						GroupName:    "groupName2",
						Notification: "1",
						Introduction: "1",
						Count:        123,
						OwnerUserID:  "ownerUserID1",
					},
					{
						GroupID:      "groupID2",
						GroupName:    "groupName2",
						Notification: "2",
						Introduction: "2",
						Count:        456,
						OwnerUserID:  "ownerUserID2",
					},
				},
				Re1: &r,
				Re2: Req11{
					GroupID:      &r.GroupID,
					GroupName:    &r.GroupName,
					Notification: &r.Notification,
					Introduction: &r.Introduction,
					Count:        &r.Count,
					OwnerUserID:  &r.OwnerUserID,
				},
			},
		},
	}
	for _, tt := range tests1 {
		t.Run(tt.name, func(t *testing.T) {
			StructFieldNotNilReplace(&tt.resp, &tt.req)
			fmt.Println(tt.resp)
			if !reflect.DeepEqual(tt.want, tt.resp) {
				t.Errorf("%v have a err,%v", tt.name, tt.want)
			}
		})
	}

}
