package organization

import (
	"context"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func genDepartmentID() string {
	r := utils.Md5(strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.FormatUint(rand.Uint64(), 10))
	bi := big.NewInt(0)
	bi.SetString(r[0:8], 16)
	return bi.String()
}

func (o *organizationSvr) GetDepartmentMemberNum(ctx context.Context, parentID string) (map[string]int, error) {
	type Department struct {
		//DepartmentName         string        // 部门名
		DepartmentID           string        // 部门ID
		ParentDepartment       *Department   // 父部门
		ChildrenDepartmentList []*Department // 子部门
	}

	var departmentIDList []string // 涉及的所有部门ID

	var bottomDepartmentList []*Department // 没有子部门

	var traversalDepartment func(parent *Department) error
	traversalDepartment = func(parent *Department) error {
		departmentIDList = append(departmentIDList, parent.DepartmentID)
		departments, err := o.Database.GetParent(ctx, parent.DepartmentID)
		if err != nil {
			return err
		}
		if len(departments) == 0 {
			return nil
		}
		for _, department := range departments {
			departmentIDList = append(departmentIDList, department.DepartmentID)
			children := &Department{
				//DepartmentName:   department.Name,
				DepartmentID:     department.DepartmentID,
				ParentDepartment: parent,
			}
			parent.ChildrenDepartmentList = append(parent.ChildrenDepartmentList, children)
			if err := traversalDepartment(children); err != nil {
				return err
			}
			if len(children.ChildrenDepartmentList) == 0 {
				bottomDepartmentList = append(bottomDepartmentList, children)
			}
		}
		return nil
	}

	root := &Department{
		DepartmentID: parentID,
	}

	if err := traversalDepartment(root); err != nil {
		return nil, err
	}

	members, err := o.Database.FindDepartmentMember(departmentIDList)
	if err != nil {
		return nil, err
	}

	departmentMemberMap := make(map[string][]string) // 部门ID: []用户ID

	for _, member := range members {
		departmentMemberMap[member.DepartmentID] = append(departmentMemberMap[member.DepartmentID], member.UserID)
	}

	departmentChildrenDepartment := make(map[string][]string) // 每个部门下的所有子部门ID

	for i := 0; i < len(bottomDepartmentList); i++ {
		department := bottomDepartmentList[i]
		departmentChildrenDepartment[department.DepartmentID] = []string{}
		parent := department.ParentDepartment
		children := []string{department.DepartmentID}
		for {
			if parent == nil || parent == root {
				break
			}
			children = append(children, parent.DepartmentID)
			departmentChildrenDepartment[parent.DepartmentID] = append(departmentChildrenDepartment[parent.DepartmentID], children...)
			parent = parent.ParentDepartment
		}
	}

	duplicateRemoval := func(arr []string) []string {
		var (
			res   = make([]string, 0, len(arr))
			exist = make(map[string]struct{})
		)
		for _, val := range arr {
			if _, ok := exist[val]; !ok {
				exist[val] = struct{}{}
				res = append(res, val)
			}
		}
		return res
	}

	res := make(map[string]int)

	for departmentID, childrenDepartmentIDList := range departmentChildrenDepartment {
		var userIDList []string
		userIDList = append(userIDList, departmentMemberMap[departmentID]...) // 当前部门成员
		for _, childrenDepartmentID := range childrenDepartmentIDList {
			userIDList = append(userIDList, departmentMemberMap[childrenDepartmentID]...) // 子部门成员
		}
		userIDList = duplicateRemoval(userIDList)
		res[departmentID] = len(userIDList)
	}

	return res, nil
}
