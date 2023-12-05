// Copyright © 2023 OpenIM open source community. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package emoticon_pack

import (
	"github.com/OpenIMSDK/tools/errs"
)

// Check UploadEmoticon
func (e *AddEmoticonReq) Check() error {
	if e.OwnerId == "" {
		return errs.ErrArgs.Wrap("emoticon ID is empty")
	}
	if e.ImageData == "" {
		return errs.ErrArgs.Wrap("image URL is empty")
	}
	return nil
}

// Check RemoveEmoticon
func (e *RemoveEmoticonReq) Check() error {
	if e.EmoticonId == "" {
		return errs.ErrArgs.Wrap("emoticon ID is empty")
	}
	return nil
}

// Check GetEmoticon
//func (e *GetEmoticonReq) Check() error {
//	if e.UserId == "" {
//		return errs.ErrArgs.Wrap("User ID is empty")
//	}
//	return nil
//}
