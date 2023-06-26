package organization

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/chat/pkg/common/db/database"
	chat2 "github.com/OpenIMSDK/chat/pkg/common/db/table/chat"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	"github.com/OpenIMSDK/chat/pkg/common/dbconn"
	"github.com/OpenIMSDK/chat/pkg/proto/common"
	"github.com/OpenIMSDK/chat/pkg/proto/organization"
	"github.com/OpenIMSDK/chat/pkg/rpclient/openim"
	organizationClient "github.com/OpenIMSDK/chat/pkg/rpclient/organization"
	"github.com/OpenIMSDK/open_utils/constant"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func Start(discov discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := dbconn.NewGormDB()
	if err != nil {
		return err
	}
	//todo:修改
	tables := []any{
		chat2.Account{},
		chat2.Register{},
		chat2.Attribute{},
		chat2.VerifyCode{},
		chat2.UserLoginRecord{},
	}
	if err := db.AutoMigrate(tables...); err != nil {
		return err
	}
	if err != nil {
		return err
	}
	organization.RegisterOrganizationServer(server, &organizationSvr{
		Database:     database.NewOrganizationDatabase(db),
		Organization: organizationClient.NewOrgClient(discov),
		OpenIM:       openim.NewOpenIMClient(discov),
	})
	return nil
}

type organizationSvr struct {
	Database     database.OrganizationDatabaseInterface
	Organization *organizationClient.OrgClient
	OpenIM       *openim.OpenIMClient
}

func (o *organizationSvr) CreateDepartment(ctx context.Context, req *organization.CreateDepartmentReq) (*organization.CreateDepartmentResp, error) {
	resp := &organization.CreateDepartmentResp{CommonResp: &common.CommonResp{}, DepartmentInfo: &common.Department{}}
	if req.DepartmentInfo == nil {
		resp.CommonResp.ErrCode = constant.ErrArgs.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrArgs.ErrMsg + " req.DepartmentInfo is nil"
		return resp, nil
	}
	department := table.Department{
		DepartmentID:   genDepartmentID(),
		FaceURL:        req.DepartmentInfo.FaceURL,
		Name:           req.DepartmentInfo.Name,
		ParentID:       req.DepartmentInfo.ParentID,
		Order:          req.DepartmentInfo.Order,
		DepartmentType: req.DepartmentInfo.DepartmentType,
		RelatedGroupID: req.DepartmentInfo.RelatedGroupID,
		CreateTime:     time.UnixMilli(req.DepartmentInfo.CreateTime),
	}
	if department.DepartmentID == "" {
		department.DepartmentID = strconv.FormatInt(time.Now().Unix(), 10)
	}
	if req.DepartmentInfo.ParentID != "" {
		_, err := o.Database.GetDepartmentByID(ctx, req.DepartmentInfo.ParentID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.CommonResp.ErrCode = constant.RecordNotFound
			resp.CommonResp.ErrMsg = "parent department not found"
			return resp, nil
		} else if err != nil {
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
			return resp, nil
		}
	}
	if err := o.Database.CreateDepartment(ctx, &department); err != nil {
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
	}
	return resp, nil
}

func (o *organizationSvr) UpdateDepartment(ctx context.Context, req *organization.UpdateDepartmentReq) (*organization.UpdateDepartmentResp, error) {
	resp := &organization.UpdateDepartmentResp{CommonResp: &common.CommonResp{}}

	if req.DepartmentInfo == nil {
		resp.CommonResp.ErrCode = constant.ErrArgs.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrArgs.ErrMsg + " req.DepartmentInfo is nil"
		return resp, nil
	}
	err := o.Database.UpdateDepartment(ctx, &table.Department{
		DepartmentID:   req.DepartmentInfo.DepartmentID,
		FaceURL:        req.DepartmentInfo.FaceURL,
		Name:           req.DepartmentInfo.Name,
		ParentID:       req.DepartmentInfo.ParentID,
		Order:          req.DepartmentInfo.Order,
		DepartmentType: req.DepartmentInfo.DepartmentType,
		RelatedGroupID: req.DepartmentInfo.RelatedGroupID,
	})
	if req.DepartmentInfo.ParentID != "" {
		_, err := o.Database.GetDepartmentByID(ctx, req.DepartmentInfo.ParentID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.CommonResp.ErrCode = constant.RecordNotFound
			resp.CommonResp.ErrMsg = "parent department not found"
			return resp, nil
		} else if err != nil {
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
			return resp, nil
		}
	}
	if err != nil {
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
	}

	return resp, nil
}

func (o *organizationSvr) GetOrganizationDepartment(ctx context.Context, req *organization.GetOrganizationDepartmentReq) (*organization.GetOrganizationDepartmentResp, error) {
	resp := &organization.GetOrganizationDepartmentResp{CommonResp: &common.CommonResp{}, DepartmentList: []*organization.DepartmentInfo{}}

	numMap, err := o.GetDepartmentMemberNum(ctx, "")
	if err != nil {
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
		return resp, nil
	}
	var getSubDepartmentList func(departmentId string, list *[]*organization.DepartmentInfo) error
	getSubDepartmentList = func(departmentId string, list *[]*organization.DepartmentInfo) error {
		departments, err := o.Database.GetParent(ctx, departmentId)
		if err != nil {
			return err
		}
		for _, department := range departments {
			subs := make([]*organization.DepartmentInfo, 0)
			err = getSubDepartmentList(department.DepartmentID, &subs)
			if err != nil {
				return err
			}
			*list = append(*list, &organization.DepartmentInfo{
				Department: &common.Department{
					DepartmentID:   department.DepartmentID,
					FaceURL:        department.FaceURL,
					Name:           department.Name,
					ParentID:       department.ParentID,
					Order:          department.Order,
					DepartmentType: department.DepartmentType,
					RelatedGroupID: department.RelatedGroupID,
					CreateTime:     department.CreateTime.UnixMilli(),
					MemberNum:      uint32(numMap[department.DepartmentID]),
				},
				SubDepartmentList: subs,
			})
		}
		return nil
	}

	if err := getSubDepartmentList("", &resp.DepartmentList); err != nil {
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
		return resp, nil
	}
	return resp, nil
}

func (o *organizationSvr) DeleteDepartment(ctx context.Context, req *organization.DeleteDepartmentReq) (*organization.DeleteDepartmentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) GetDepartment(ctx context.Context, req *organization.GetDepartmentReq) (*organization.GetDepartmentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) CreateOrganizationUser(ctx context.Context, req *organization.CreateOrganizationUserReq) (*organization.CreateOrganizationUserResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) UpdateOrganizationUser(ctx context.Context, req *organization.UpdateOrganizationUserReq) (*organization.UpdateOrganizationUserResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) DeleteOrganizationUser(ctx context.Context, req *organization.DeleteOrganizationUserReq) (*organization.DeleteOrganizationUserResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) CreateDepartmentMember(ctx context.Context, req *organization.CreateDepartmentMemberReq) (*organization.CreateDepartmentMemberResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) GetUserInDepartment(ctx context.Context, req *organization.GetUserInDepartmentReq) (*organization.GetUserInDepartmentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) DeleteUserInDepartment(ctx context.Context, req *organization.DeleteUserInDepartmentReq) (*organization.DeleteUserInDepartmentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) UpdateUserInDepartment(ctx context.Context, req *organization.UpdateUserInDepartmentReq) (*organization.UpdateUserInDepartmentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) GetSearchUserList(ctx context.Context, req *organization.GetSearchUserListReq) (*organization.GetSearchUserListResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) SetOrganization(ctx context.Context, req *organization.SetOrganizationReq) (*organization.SetOrganizationResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) GetOrganization(ctx context.Context, req *organization.GetOrganizationReq) (*organization.GetOrganizationResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) GetSubDepartment(ctx context.Context, req *organization.GetSubDepartmentReq) (*organization.GetSubDepartmentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) GetSearchDepartmentUser(ctx context.Context, req *organization.GetSearchDepartmentUserReq) (*organization.GetSearchDepartmentUserResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) SortDepartmentList(ctx context.Context, req *organization.SortDepartmentListReq) (*organization.SortDepartmentListResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) SortOrganizationUserList(ctx context.Context, req *organization.SortOrganizationUserListReq) (*organization.SortOrganizationUserListResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) CreateNewOrganizationMember(ctx context.Context, req *organization.CreateNewOrganizationMemberReq) (*organization.CreateNewOrganizationMemberResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) GetUserInfo(ctx context.Context, req *organization.GetUserInfoReq) (*organization.GetUserInfoResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) BatchImport(ctx context.Context, req *organization.BatchImportReq) (*organization.BatchImportResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) MoveUserDepartment(ctx context.Context, req *organization.MoveUserDepartmentReq) (*organization.MoveUserDepartmentResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) GetUserFullList(ctx context.Context, req *organization.GetUserFullListReq) (*organization.GetUserFullListResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) SearchUsersFullInfo(ctx context.Context, req *organization.SearchUsersFullInfoReq) (*organization.SearchUsersFullInfoResp, error) {
	//TODO implement me
	panic("implement me")
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
