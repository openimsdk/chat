#!/usr/bin/env bash
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

echo "start running docker_start_all.sh..."
chmod +x ./*.sh
echo "start running internal bash ====> ./start_all.sh"
./start_all.sh
echo "running ./start_all.sh succeed "
i=1
while ((i == 1))
do
    sleep 5
done
echo "docker  docker is started successfully"