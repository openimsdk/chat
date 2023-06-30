# Copyright © 2023 OpenIM open source community. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

protoc --go_out=plugins=grpc:./common --go_opt=module=github.com/OpenIMSDK/chat/pkg/proto/common common/common.proto
protoc --go_out=plugins=grpc:./admin --go_opt=module=github.com/OpenIMSDK/chat/pkg/proto/admin admin/admin.proto
protoc --go_out=plugins=grpc:./chat --go_opt=module=github.com/OpenIMSDK/chat/pkg/proto/chat chat/chat.proto