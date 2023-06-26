package organization

import (
	common2 "Open-IM-Organization/internal/rpc/common"
	"Open-IM-Organization/pkg/common/config"
	"Open-IM-Organization/pkg/common/constant"
	"Open-IM-Organization/pkg/common/db"
	"Open-IM-Organization/pkg/common/db/mysql_model/organization"
	"Open-IM-Organization/pkg/proto/chat"
	"Open-IM-Organization/pkg/proto/common"
	rpc "Open-IM-Organization/pkg/proto/organization"
	"context"
	log "github.com/OpenIMSDK/open_log"
	utils "github.com/OpenIMSDK/open_utils"
)

func (s *organizationServer) GetUserFullList(ctx context.Context, req *rpc.GetUserFullListReq) (*rpc.GetUserFullListResp, error) {
	resp := &rpc.GetUserFullListResp{CommonResp: &common.CommonResp{}}
	defer func(name string) {
		log.NewInfo(req.OperationID, name, " rpc req ", req.String())
		if resp.CommonResp.ErrCode == 0 {
			log.NewInfo(req.OperationID, name, " rpc resp ", resp.String())
		} else {
			log.NewError(req.OperationID, name, " rpc resp ", resp.String())
		}
	}(utils.GetSelfFuncName())
	var (
		total int64
		users []organization.OrganizationUser
		err   error
	)
	if len(req.UserIDList) == 0 {
		if req.ShowNumber == 0 {
			req.ShowNumber = 10
		}
		total, users, err = db.DB.MysqlDB.OrganizationUser.GetPage(int(req.PageNumber), int(req.ShowNumber))
	} else {
		users, err = db.DB.MysqlDB.OrganizationUser.GetList(req.UserIDList)
		total = int64(len(users))
	}
	if err != nil {
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
		return resp, nil
	}
	resp.Total = int32(total)
	for _, user := range users {
		members, err := db.DB.MysqlDB.DepartmentMember.GetUser(user.UserID)
		if err != nil {
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
			return resp, nil
		}
		departmentIDList := make([]string, 0, len(members))
		for _, member := range members {
			departmentIDList = append(departmentIDList, member.DepartmentID)
		}
		departments, err := db.DB.MysqlDB.Department.GetList(departmentIDList)
		if err != nil {
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
			return resp, nil
		}
		departmentMap := make(map[string]*common.Department)
		for _, department := range departments {
			departmentMap[department.DepartmentID] = &common.Department{
				DepartmentID:   department.DepartmentID,
				FaceURL:        department.FaceURL,
				Name:           department.Name,
				ParentID:       department.ParentID,
				Order:          department.Order,
				DepartmentType: department.DepartmentType,
				RelatedGroupID: department.RelatedGroupID,
				CreateTime:     department.CreateTime.UnixMilli(),
			}
		}
		var departmentMemberList []*common.DepartmentMember
		for _, member := range members {
			department := departmentMap[member.DepartmentID]
			if department == nil {
				continue
			}
			var terminationTime int64
			if member.TerminationTime == nil {
				terminationTime = constant.NilTimestamp
			} else {
				terminationTime = member.TerminationTime.UnixMilli()
			}
			departmentMemberList = append(departmentMemberList, &common.DepartmentMember{
				UserID:          member.UserID,
				DepartmentID:    member.DepartmentID,
				Order:           member.Order,
				Position:        member.Position,
				Leader:          member.Leader,
				Status:          member.Status,
				EntryTime:       member.EntryTime.UnixMilli(),
				TerminationTime: terminationTime,
				CreateTime:      member.CreateTime.UnixMilli(),
				Department:      department,
			})
		}
		resp.OrganizationUserList = append(resp.OrganizationUserList, &common.OrganizationUser{
			UserID:               user.UserID,
			Nickname:             user.Nickname,
			EnglishName:          user.EnglishName,
			FaceURL:              user.FaceURL,
			Gender:               user.Gender,
			Mobile:               user.Mobile,
			Telephone:            user.Telephone,
			Birth:                user.Birth.UnixMilli(),
			Email:                user.Email,
			Order:                user.Order,
			Status:               user.Status,
			CreateTime:           user.CreateTime.UnixMilli(),
			DepartmentMemberList: departmentMemberList,
			Station:              user.Station,
			AreaCode:             user.AreaCode,
		})
	}
	return resp, nil
}

func (s *organizationServer) SearchUsersFullInfo(ctx context.Context, req *rpc.SearchUsersFullInfoReq) (*rpc.SearchUsersFullInfoResp, error) {
	resp := &rpc.SearchUsersFullInfoResp{CommonResp: &common.CommonResp{}}
	if req.Operation == nil {
		resp.CommonResp.ErrCode = constant.ErrArgs.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrArgs.ErrMsg + " req.Operation is nil"
		return resp, nil
	}
	if req.Pagination == nil {
		resp.CommonResp.ErrCode = constant.ErrArgs.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrArgs.ErrMsg + " req.Pagination is nil"
		return resp, nil
	}
	if req.UserPublicInfo == nil {
		resp.CommonResp.ErrCode = constant.ErrArgs.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrArgs.ErrMsg + " req.UserPublicInfo is nil"
		return resp, nil
	}
	if req.Pagination.ShowNumber == 0 {
		req.Pagination.ShowNumber = 10
	}
	var userIDList []string
	if req.Content != "" {
		// 搜索账号匹配的用户
		etcdConn := common2.GetDefaultConn(config.Config.RpcRegisterName.OpenImChatName, req.Operation.OperationID, resp.CommonResp)
		if etcdConn == nil {
			resp.CommonResp.ErrCode = constant.ErrServer.ErrCode
			resp.CommonResp.ErrMsg = "get base server etcd conn is empty"
			return resp, nil
		}
		rpcResp, err := chat.NewChatClient(etcdConn).GetAccountUser(ctx, &chat.GetAccountUserReq{
			Operation:   req.Operation,
			AccountList: []string{req.Content},
		})
		if err != nil {
			resp.CommonResp.ErrCode = constant.ErrServer.ErrCode
			resp.CommonResp.ErrMsg = "get base server GetAccountUser error: " + err.Error()
			return resp, nil
		}
		if rpcResp.CommonResp.ErrCode != constant.NoError {
			resp.CommonResp.ErrCode = rpcResp.CommonResp.ErrCode
			resp.CommonResp.ErrMsg = "get base server GetAccountUser errMsg: " + rpcResp.CommonResp.ErrMsg
			return resp, nil
		}
		for _, userID := range rpcResp.Accounts {
			userIDList = append(userIDList, userID)
		}
		//ds, err := db.DB.MysqlDB.Department.Search(req.Content)
		//if err != nil {
		//	resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		//	resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
		//	return resp, nil
		//}
		//if len(ds) > 0 {
		//	idList := make([]string, 0, len(ds))
		//	for _, d := range ds {
		//		idList = append(idList, d.DepartmentID)
		//	}
		//	dms, err := db.DB.MysqlDB.DepartmentMember.GetDepartmentList(idList)
		//	if err != nil {
		//		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		//		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
		//		return resp, nil
		//	}
		//	for _, dm := range dms {
		//		userIDList = append(userIDList, dm.UserID)
		//	}
		//}
		//dms, err := db.DB.MysqlDB.DepartmentMember.Search(req.Content)
		//if err != nil {
		//	resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		//	resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
		//	return resp, nil
		//}
		//for _, dm := range dms {
		//	userIDList = append(userIDList, dm.UserID)
		//}
	}
	if req.UserPublicInfo.UserID != "" {
		userIDList = append(userIDList, req.UserPublicInfo.UserID)
	}
	total, users, err := db.DB.MysqlDB.OrganizationUser.SearchV2(req.Content, userIDList, int(req.Pagination.PageNumber), int(req.Pagination.ShowNumber))
	if err != nil {
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
		return resp, nil
	}
	resp.Total = int32(total)
	for _, user := range users {
		members, err := db.DB.MysqlDB.DepartmentMember.GetUser(user.UserID)
		if err != nil {
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
			return resp, nil
		}
		departmentIDList := make([]string, 0, len(members))
		for _, member := range members {
			departmentIDList = append(departmentIDList, member.DepartmentID)
		}
		departments, err := db.DB.MysqlDB.Department.GetList(departmentIDList)
		if err != nil {
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
			return resp, nil
		}
		departmentMap := make(map[string]*common.Department)
		for _, department := range departments {
			departmentMap[department.DepartmentID] = &common.Department{
				DepartmentID:   department.DepartmentID,
				FaceURL:        department.FaceURL,
				Name:           department.Name,
				ParentID:       department.ParentID,
				Order:          department.Order,
				DepartmentType: department.DepartmentType,
				RelatedGroupID: department.RelatedGroupID,
				CreateTime:     department.CreateTime.UnixMilli(),
			}
		}
		var departmentMemberList []*common.DepartmentMember
		for _, member := range members {
			department := departmentMap[member.DepartmentID]
			if department == nil {
				continue
			}
			var terminationTime int64
			if member.TerminationTime == nil {
				terminationTime = constant.NilTimestamp
			} else {
				terminationTime = member.TerminationTime.UnixMilli()
			}
			departmentMemberList = append(departmentMemberList, &common.DepartmentMember{
				UserID:          member.UserID,
				DepartmentID:    member.DepartmentID,
				Order:           member.Order,
				Position:        member.Position,
				Leader:          member.Leader,
				Status:          member.Status,
				EntryTime:       member.EntryTime.UnixMilli(),
				TerminationTime: terminationTime,
				CreateTime:      member.CreateTime.UnixMilli(),
				Department:      department,
			})
		}
		resp.OrganizationUserList = append(resp.OrganizationUserList, &common.OrganizationUser{
			UserID:               user.UserID,
			Nickname:             user.Nickname,
			EnglishName:          user.EnglishName,
			FaceURL:              user.FaceURL,
			Gender:               user.Gender,
			Mobile:               user.Mobile,
			Telephone:            user.Telephone,
			Birth:                user.Birth.UnixMilli(),
			Email:                user.Email,
			Order:                user.Order,
			Status:               user.Status,
			CreateTime:           user.CreateTime.UnixMilli(),
			DepartmentMemberList: departmentMemberList,
			Station:              user.Station,
			AreaCode:             user.AreaCode,
		})
	}
	return resp, nil
}
