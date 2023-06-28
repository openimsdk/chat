# organization rpc测试

## CreateDepartment

postman 输入：
```
{
    "operationID":"abcd",
    "OpUserID":"sdf",
    "DepartmentInfo": {
	"state": null,
	"sizeCache": null,
	"unknownFields": null,
	"departmentID": "",
	"faceURL": "",
	"name": "",
	"parentID": "",
	"order": 0,
	"departmentType": 0,
	"relatedGroupID": "",
	"createTime": 0,
	"memberNum": 0,
	"position": ""
    }
}
```
输出：
```
{
    "errCode": 0,
    "errMsg": "",
    "errDlt": "",
    "data": {
        "commonResp": {
            "errCode": 0,
            "errMsg": ""
        },
        "departmentInfo": {
            "departmentID": "",
            "faceURL": "",
            "name": "",
            "parentID": "",
            "order": 0,
            "departmentType": 0,
            "relatedGroupID": "",
            "createTime": 0,
            "memberNum": 0,
            "position": ""
        }
    }
}
```
## UpdateDepartment

 ```
 
 ```

