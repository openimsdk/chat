package organization

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/db/database"
	table "github.com/OpenIMSDK/chat/pkg/common/db/table/organization"
	"github.com/OpenIMSDK/chat/pkg/common/dbconn"
	"github.com/OpenIMSDK/chat/pkg/proto/common"
	"github.com/OpenIMSDK/chat/pkg/proto/organization"
	"github.com/OpenIMSDK/chat/pkg/rpclient/openim"
	organizationClient "github.com/OpenIMSDK/chat/pkg/rpclient/organization"
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
		table.Department{},
		table.DepartmentMember{},
		table.OrganizationUser{},
		table.Organization{},
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
	resp := &organization.CreateDepartmentResp{DepartmentInfo: &common.Department{}}
	if req.DepartmentInfo == nil {
		return nil, errs.ErrArgs.Wrap(" req.DepartmentInfo is nil")
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
			return nil, errs.ErrArgs.Wrap("parent department not found")
		} else if err != nil {
			return nil, err
		}
	}
	if err := o.Database.CreateDepartment(ctx, &department); err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) UpdateDepartment(ctx context.Context, req *organization.UpdateDepartmentReq) (*organization.UpdateDepartmentResp, error) {
	resp := &organization.UpdateDepartmentResp{}

	if req.DepartmentInfo == nil {
		return nil, errs.ErrArgs.Wrap(" req.DepartmentInfo is nil")
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
			return nil, errs.ErrArgs.Wrap("parent department not found")
		} else if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (o *organizationSvr) GetOrganizationDepartment(ctx context.Context, req *organization.GetOrganizationDepartmentReq) (*organization.GetOrganizationDepartmentResp, error) {
	resp := &organization.GetOrganizationDepartmentResp{DepartmentList: []*organization.DepartmentInfo{}}

	numMap, err := o.GetDepartmentMemberNum(ctx, "")
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) DeleteDepartment(ctx context.Context, req *organization.DeleteDepartmentReq) (*organization.DeleteDepartmentResp, error) {
	resp := &organization.DeleteDepartmentResp{}
	departmentList, err := o.Database.GetList(ctx, req.DepartmentIDList)
	if err != nil {
		return nil, err
	}
	if len(departmentList) == 0 {
		return nil, errs.ErrArgs.Wrap("parent department not found")
	}
	// 修改删除的子部门的父部门为删除的上级
	for _, department := range departmentList {
		err := o.Database.UpdateParentID(ctx, department.DepartmentID, department.ParentID)
		if err != nil {
			return nil, err
		}
	}
	// 删除部门
	if err := o.Database.DeleteDepartment(ctx, req.DepartmentIDList); err != nil {
		return nil, err
	}
	// 删除职位信息
	if err := o.Database.DeleteDepartmentIDList(ctx, req.DepartmentIDList); err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) GetDepartment(ctx context.Context, req *organization.GetDepartmentReq) (*organization.GetDepartmentResp, error) {
	resp := &organization.GetDepartmentResp{}

	department, err := o.Database.GetDepartment(ctx, req.DepartmentID)
	if err == nil {
		resp.Department = &common.Department{
			DepartmentID:   department.DepartmentID,
			FaceURL:        department.FaceURL,
			Name:           department.Name,
			ParentID:       department.ParentID,
			Order:          department.Order,
			DepartmentType: department.DepartmentType,
			RelatedGroupID: department.RelatedGroupID,
			CreateTime:     department.CreateTime.UnixMilli(),
		}
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.ErrArgs.Wrap("department not found")
	} else {
		return nil, err
	}

	return resp, nil
}

func (o *organizationSvr) CreateOrganizationUser(ctx context.Context, req *organization.CreateOrganizationUserReq) (*organization.CreateOrganizationUserResp, error) {
	resp := &organization.CreateOrganizationUserResp{}
	if req.OrganizationUser == nil {
		return nil, errs.ErrArgs.Wrap(" req.OrganizationUser is nil")
	}
	err := o.Database.CreateOrganizationUser(ctx, &table.OrganizationUser{
		UserID:      req.OrganizationUser.UserID,
		Nickname:    req.OrganizationUser.Nickname,
		EnglishName: req.OrganizationUser.EnglishName,
		FaceURL:     req.OrganizationUser.FaceURL,
		Gender:      req.OrganizationUser.Gender,
		Mobile:      req.OrganizationUser.Mobile,
		Telephone:   req.OrganizationUser.Telephone,
		Birth:       time.UnixMilli(req.OrganizationUser.Birth),
		Email:       req.OrganizationUser.Email,
		Status:      req.OrganizationUser.Status,
		Station:     req.OrganizationUser.Station,
		AreaCode:    req.OrganizationUser.AreaCode,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) UpdateOrganizationUser(ctx context.Context, req *organization.UpdateOrganizationUserReq) (*organization.UpdateOrganizationUserResp, error) {
	resp := &organization.UpdateOrganizationUserResp{}
	if req.OrganizationUser == nil {
		return nil, errs.ErrArgs.Wrap(" req.OrganizationUser is nil")
	}
	err := o.Database.UpdateOrganizationUser(ctx, &table.OrganizationUser{
		UserID:      req.OrganizationUser.UserID,
		Nickname:    req.OrganizationUser.Nickname,
		EnglishName: req.OrganizationUser.EnglishName,
		FaceURL:     req.OrganizationUser.FaceURL,
		Gender:      req.OrganizationUser.Gender,
		Mobile:      req.OrganizationUser.Mobile,
		Telephone:   req.OrganizationUser.Telephone,
		Birth:       time.UnixMilli(req.OrganizationUser.Birth),
		Email:       req.OrganizationUser.Email,
		Status:      req.OrganizationUser.Status,
		Station:     req.OrganizationUser.Station,
		AreaCode:    req.OrganizationUser.AreaCode,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) DeleteOrganizationUser(ctx context.Context, req *organization.DeleteOrganizationUserReq) (*organization.DeleteOrganizationUserResp, error) {
	resp := &organization.DeleteOrganizationUserResp{}
	err := o.Database.DeleteOrganizationUser(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	err = o.Database.DeleteDepartmentMemberByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) CreateDepartmentMember(ctx context.Context, req *organization.CreateDepartmentMemberReq) (*organization.CreateDepartmentMemberResp, error) {
	//todo：departmentid或userId不存在报错
	resp := &organization.CreateDepartmentMemberResp{}
	if req.DepartmentMember == nil {
		return nil, errs.ErrArgs.Wrap("req.DepartmentInfo is nil")
	}
	var terminationTime *time.Time
	if req.DepartmentMember.TerminationTime != constant.NilTimestamp {
		t := time.UnixMilli(req.DepartmentMember.EntryTime)
		terminationTime = &t
	}
	err := o.Database.CreateDepartmentMember(ctx, &table.DepartmentMember{
		UserID:          req.DepartmentMember.UserID,
		DepartmentID:    req.DepartmentMember.DepartmentID,
		Order:           req.DepartmentMember.Order,
		Position:        req.DepartmentMember.Position,
		Leader:          req.DepartmentMember.Leader,
		Status:          req.DepartmentMember.Status,
		EntryTime:       time.UnixMilli(req.DepartmentMember.EntryTime),
		TerminationTime: terminationTime,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) GetUserInDepartment(ctx context.Context, req *organization.GetUserInDepartmentReq) (*organization.GetUserInDepartmentResp, error) {
	resp := &organization.GetUserInDepartmentResp{}
	user, err := o.Database.GetOrganizationUser(ctx, req.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.ErrArgs.Wrap("user not fount")
	} else if err != nil {
		return nil, err
	}
	dms, err := o.Database.GetDepartmentMemberByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	departmentIDList := make([]string, 0, len(dms))
	for _, dm := range dms {
		departmentIDList = append(departmentIDList, dm.DepartmentID)
	}
	departmentList, err := o.Database.GetList(ctx, departmentIDList)
	if err != nil {
		return nil, err
	}
	numMap, err := o.GetDepartmentMemberNum(ctx, "")
	if err != nil {
		return nil, err
	}
	departmentMap := make(map[string]*common.Department)
	for _, department := range departmentList {
		departmentMap[department.DepartmentID] = &common.Department{
			DepartmentID:   department.DepartmentID,
			FaceURL:        department.FaceURL,
			Name:           department.Name,
			ParentID:       department.ParentID,
			Order:          department.Order,
			DepartmentType: department.DepartmentType,
			RelatedGroupID: department.RelatedGroupID,
			CreateTime:     department.CreateTime.UnixMilli(),
			MemberNum:      uint32(numMap[department.DepartmentID]),
		}
	}
	resp.UserInDepartment = &common.UserInDepartment{
		OrganizationUser: &common.OrganizationUser{
			UserID:      user.UserID,
			Nickname:    user.Nickname,
			EnglishName: user.EnglishName,
			FaceURL:     user.FaceURL,
			Gender:      user.Gender,
			Mobile:      user.Mobile,
			Telephone:   user.Telephone,
			Birth:       user.Birth.UnixMilli(),
			Email:       user.Email,
			Order:       user.Order,
			Status:      user.Status,
			CreateTime:  user.CreateTime.Unix(),
			Ex:          "",
			Station:     user.Station,
			AreaCode:    user.AreaCode,
		},
		DepartmentMemberList: make([]*common.DepartmentMember, len(dms)),
	}
	for i, dm := range dms {
		var terminationTime int64
		if dm.TerminationTime == nil {
			terminationTime = constant.NilTimestamp
		} else {
			terminationTime = dm.TerminationTime.UnixMilli()
		}
		resp.UserInDepartment.DepartmentMemberList[i] = &common.DepartmentMember{
			UserID:          dm.UserID,
			DepartmentID:    dm.DepartmentID,
			Order:           dm.Order,
			Position:        dm.Position,
			Leader:          dm.Leader,
			Status:          dm.Status,
			Ex:              "",
			EntryTime:       dm.EntryTime.UnixMilli(),
			TerminationTime: terminationTime,
			Department:      departmentMap[dm.DepartmentID],
		}
	}
	return resp, nil
}

func (o *organizationSvr) DeleteUserInDepartment(ctx context.Context, req *organization.DeleteUserInDepartmentReq) (*organization.DeleteUserInDepartmentResp, error) {
	resp := &organization.DeleteUserInDepartmentResp{}
	err := o.Database.DeleteDepartmentMemberByKey(ctx, req.UserID, req.DepartmentID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) UpdateUserInDepartment(ctx context.Context, req *organization.UpdateUserInDepartmentReq) (*organization.UpdateUserInDepartmentResp, error) {
	resp := &organization.UpdateUserInDepartmentResp{}
	if req.DepartmentMember == nil {
		return nil, errs.ErrArgs.Wrap(" req.DepartmentInfo is nil")
	}
	var (
		entryTime       time.Time
		terminationTime *time.Time
	)
	if req.DepartmentMember.EntryTime != constant.NilTimestamp {
		entryTime = time.UnixMilli(req.DepartmentMember.EntryTime)
	}
	if req.DepartmentMember.TerminationTime != constant.NilTimestamp {
		t := time.UnixMilli(req.DepartmentMember.TerminationTime)
		terminationTime = &t
	}
	err := o.Database.UpdateDepartmentMember(ctx, &table.DepartmentMember{
		UserID:          req.DepartmentMember.UserID,
		DepartmentID:    req.DepartmentMember.DepartmentID,
		Order:           req.DepartmentMember.Order,
		Position:        req.DepartmentMember.Position,
		Leader:          req.DepartmentMember.Leader,
		Status:          req.DepartmentMember.Status,
		EntryTime:       entryTime,
		TerminationTime: terminationTime,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) GetSearchUserList(ctx context.Context, req *organization.GetSearchUserListReq) (*organization.GetSearchUserListResp, error) {
	resp := &organization.GetSearchUserListResp{UserList: []*common.UserInDepartment{}}
	var userIDList []string
	if len(req.DepartmentIDList) > 0 {
		departments, err := o.Database.FindDepartmentMember(ctx, req.DepartmentIDList)
		if err != nil {
			return nil, err
		}
		if len(departments) == 0 {
			return resp, nil
		}
		userIDList = make([]string, len(departments))
		for i, department := range departments {
			userIDList[i] = department.UserID
		}
	}
	total, users, err := o.Database.SearchPage(ctx, req.PositionList, userIDList, req.Text, req.Sorts, req.PageNumber, req.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	findUserIDList := make([]string, len(users))
	for i, user := range users {
		findUserIDList[i] = user.UserID
	}
	departmentMemberList, err := o.Database.FindDepartmentMemberByUserID(ctx, findUserIDList)
	if err != nil {
		return nil, err
	}
	departmentIDList := make([]string, 0, len(departmentMemberList))
	for _, member := range departmentMemberList {
		departmentIDList = append(departmentIDList, member.DepartmentID)
	}
	departmentList, err := o.Database.GetList(ctx, departmentIDList)
	if err != nil {
		return nil, err
	}
	departmentMap := make(map[string]*common.Department)
	for _, department := range departmentList {
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
	departmentMemberMap := make(map[string][]*common.DepartmentMember)
	for _, member := range departmentMemberList {
		var terminationTime int64
		if member.TerminationTime == nil {
			terminationTime = constant.NilTimestamp
		} else {
			terminationTime = member.TerminationTime.UnixMilli()
		}
		departmentMemberMap[member.UserID] = append(departmentMemberMap[member.UserID], &common.DepartmentMember{
			UserID:          member.UserID,
			DepartmentID:    member.DepartmentID,
			Order:           member.Order,
			Position:        member.Position,
			Leader:          member.Leader,
			Status:          member.Status,
			Ex:              "",
			EntryTime:       member.EntryTime.UnixMilli(),
			TerminationTime: terminationTime,
			CreateTime:      member.CreateTime.UnixMilli(),
			Department:      departmentMap[member.DepartmentID],
		})
	}
	for _, user := range users {
		departmentMembers := departmentMemberMap[user.UserID]
		if departmentMembers == nil {
			departmentMembers = []*common.DepartmentMember{}
		}
		resp.UserList = append(resp.UserList, &common.UserInDepartment{
			DepartmentMemberList: departmentMembers,
			OrganizationUser: &common.OrganizationUser{
				UserID:      user.UserID,
				Nickname:    user.Nickname,
				EnglishName: user.EnglishName,
				FaceURL:     user.FaceURL,
				Gender:      user.Gender,
				Mobile:      user.Mobile,
				Telephone:   user.Telephone,
				Birth:       user.Birth.UnixMilli(),
				Email:       user.Email,
				Order:       user.Order,
				Status:      user.Status,
				CreateTime:  user.CreateTime.UnixMilli(),
				Ex:          "",
				Station:     user.Station,
				AreaCode:    user.AreaCode,
			},
		})
	}
	return resp, nil
}

func (o *organizationSvr) SetOrganization(ctx context.Context, req *organization.SetOrganizationReq) (*organization.SetOrganizationResp, error) {
	resp := &organization.SetOrganizationResp{}
	if req.Organization == nil {
		return nil, errs.ErrArgs.Wrap(" req.Organization is nil")
	}
	err := o.Database.SetOrganization(ctx, &table.Organization{
		LogoURL:        req.Organization.LogoURL,
		Name:           req.Organization.Name,
		Homepage:       req.Organization.Homepage,
		RelatedGroupID: req.Organization.RelatedGroupID,
		Introduction:   req.Organization.Introduction,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) GetOrganization(ctx context.Context, req *organization.GetOrganizationReq) (*organization.GetOrganizationResp, error) {
	resp := &organization.GetOrganizationResp{}
	org, err := o.Database.GetOrganization(ctx)
	if err != nil {
		return nil, err
	}
	resp.Organization = &common.Organization{
		LogoURL:        org.LogoURL,
		Name:           org.Name,
		Homepage:       org.Homepage,
		RelatedGroupID: org.RelatedGroupID,
		Introduction:   org.Introduction,
		CreateTime:     org.CreateTime.UnixMilli(),
	}
	return resp, nil
}

func (o *organizationSvr) GetSubDepartment(ctx context.Context, req *organization.GetSubDepartmentReq) (*organization.GetSubDepartmentResp, error) {
	resp := &organization.GetSubDepartmentResp{
		DepartmentMemberList:    []*common.DepartmentMember{},
		DepartmentList:          []*common.Department{},
		DepartmentDirectoryList: []*common.Department{},
	}
	departmentList, err := o.Database.GetParent(ctx, req.DepartmentID)
	if err != nil {
		return nil, err
	}
	memberCountMap, err := o.GetDepartmentMemberNum(ctx, req.DepartmentID)
	if err != nil {
		return nil, errs.ErrArgs.Wrap(" get num ", err.Error())
	}
	for _, department := range departmentList {
		resp.DepartmentList = append(resp.DepartmentList, &common.Department{
			DepartmentID:   department.DepartmentID,
			FaceURL:        department.FaceURL,
			Name:           department.Name,
			ParentID:       department.ParentID,
			Order:          department.Order,
			DepartmentType: department.DepartmentType,
			RelatedGroupID: department.RelatedGroupID,
			CreateTime:     department.CreateTime.UnixMilli(),
			MemberNum:      uint32(memberCountMap[department.DepartmentID]),
		})
	}
	if req.DepartmentID == "" {
		userIDList, err := o.Database.GetNoDepartmentUserIDList(ctx)
		if err != nil {
			return nil, err
		}
		organizationUserList, err := o.Database.GetOrganizationUserList(ctx, userIDList)
		if err != nil {
			return nil, err
		}
		departmentMemberList, err := o.Database.GetUserListInDepartment(ctx, req.DepartmentID, userIDList)
		if err != nil {
			return nil, err
		}
		departmentMemberMap := make(map[string]*table.DepartmentMember)
		for i, member := range departmentMemberList {
			departmentMemberMap[member.UserID] = departmentMemberList[i]
		}
		for _, user := range organizationUserList {
			member := departmentMemberMap[user.UserID]
			resp.DepartmentMemberList = append(resp.DepartmentMemberList, &common.DepartmentMember{
				UserID:          user.UserID,
				DepartmentID:    member.DepartmentID,
				Order:           member.Order,
				Position:        member.Position,
				Leader:          member.Leader,
				Status:          member.Status,
				Ex:              "",
				EntryTime:       constant.NilTimestamp,
				TerminationTime: constant.NilTimestamp,
				CreateTime:      constant.NilTimestamp,
				Department:      nil,
				OrganizationUser: &common.OrganizationUser{
					UserID:      user.UserID,
					Nickname:    user.Nickname,
					EnglishName: user.EnglishName,
					FaceURL:     user.FaceURL,
					Gender:      user.Gender,
					Mobile:      user.Mobile,
					Telephone:   user.Telephone,
					Birth:       user.Birth.UnixMilli(),
					Email:       user.Email,
					Status:      user.Status,
					Order:       user.Order,
					CreateTime:  user.CreateTime.UnixMilli(),
					Ex:          "",
					Station:     user.Station,
					AreaCode:    user.AreaCode,
				},
			})
		}
	} else {
		departmentMemberList, err := o.Database.GetDepartmentMemberByDepartmentID(ctx, req.DepartmentID)
		if err != nil {
			return nil, err
		}
		userIDList := make([]string, len(departmentMemberList))
		for i, member := range departmentMemberList {
			userIDList[i] = member.UserID
		}
		userList, err := o.Database.GetOrganizationUserList(ctx, userIDList)
		if err != nil {
			return nil, err
		}
		userMap := make(map[string]*common.OrganizationUser)
		for _, user := range userList {
			userMap[user.UserID] = &common.OrganizationUser{
				UserID:      user.UserID,
				Nickname:    user.Nickname,
				EnglishName: user.EnglishName,
				FaceURL:     user.FaceURL,
				Gender:      user.Gender,
				Mobile:      user.Mobile,
				Telephone:   user.Telephone,
				Birth:       user.Birth.UnixMilli(),
				Email:       user.Email,
				Order:       user.Order,
				Status:      user.Status,
				CreateTime:  user.CreateTime.UnixMilli(),
				Ex:          "",
				Station:     user.Station,
				AreaCode:    user.AreaCode,
			}
		}
		for _, member := range departmentMemberList {
			var terminationTime int64
			if member.TerminationTime == nil {
				terminationTime = constant.NilTimestamp
			} else {
				terminationTime = member.TerminationTime.UnixMilli()
			}
			resp.DepartmentMemberList = append(resp.DepartmentMemberList, &common.DepartmentMember{
				UserID:           member.UserID,
				DepartmentID:     member.DepartmentID,
				Order:            member.Order,
				Position:         member.Position,
				Leader:           member.Leader,
				Status:           member.Status,
				Ex:               "",
				EntryTime:        member.EntryTime.UnixMilli(),
				TerminationTime:  terminationTime,
				CreateTime:       member.CreateTime.UnixMilli(),
				OrganizationUser: userMap[member.UserID],
			})
		}
	}
	var ds []*common.Department
	if req.DepartmentID != "" {
		departmentID := req.DepartmentID
		for {
			department, err := o.Database.GetDepartment(ctx, departmentID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					break
				}
				return nil, err
			}
			ds = append(ds, &common.Department{
				DepartmentID:   department.DepartmentID,
				FaceURL:        department.FaceURL,
				Name:           department.Name,
				ParentID:       department.ParentID,
				Order:          department.Order,
				DepartmentType: department.DepartmentType,
				RelatedGroupID: department.RelatedGroupID,
				CreateTime:     department.CreateTime.UnixMilli(),
			})
			if department.ParentID == "" {
				break
			}
			departmentID = department.ParentID
		}
	}
	org, err := o.Database.GetOrganization(ctx)
	if err != nil {
		return nil, err
	}
	ds = append(ds, &common.Department{
		DepartmentID: "",
		FaceURL:      org.LogoURL,
		Name:         org.Name,
		CreateTime:   org.CreateTime.UnixMilli(),
	})
	resp.DepartmentDirectoryList = make([]*common.Department, 0, len(ds))
	for i := len(ds) - 1; i >= 0; i-- {
		resp.DepartmentDirectoryList = append(resp.DepartmentDirectoryList, ds[i])
	}
	return resp, nil
}

func (o *organizationSvr) GetSearchDepartmentUser(ctx context.Context, req *organization.GetSearchDepartmentUserReq) (*organization.GetSearchDepartmentUserResp, error) {
	resp := &organization.GetSearchDepartmentUserResp{
		OrganizationUserList: []*organization.GetSearchDepartmentUserOrganizationUser{},
		DepartmentList:       []*common.Department{},
	}
	organizationUserList, err := o.Database.SearchOrganizationUser(ctx, nil, nil, req.Keyword, nil)
	if err != nil {
		return nil, err
	}
	org, err := o.Database.GetOrganization(ctx)
	if err != nil {
		return nil, err
	}
	orgDepartment := &common.Department{
		DepartmentID:   "",
		FaceURL:        org.LogoURL,
		Name:           org.Name,
		ParentID:       "",
		Order:          0,
		DepartmentType: 0,
		RelatedGroupID: org.RelatedGroupID,
		CreateTime:     org.CreateTime.UnixMilli(),
	}
	for _, user := range organizationUserList {
		departmentMemberList, err := o.Database.GetDepartmentMemberByUserID(ctx, user.UserID)
		if err != nil {
			return nil, err
		}
		organizationUser := &common.OrganizationUser{
			UserID:      user.UserID,
			Nickname:    user.Nickname,
			EnglishName: user.EnglishName,
			FaceURL:     user.FaceURL,
			Gender:      user.Gender,
			Mobile:      user.Mobile,
			Telephone:   user.Telephone,
			Birth:       user.Birth.UnixMilli(),
			Email:       user.Email,
			Order:       user.Order,
			Status:      user.Status,
			CreateTime:  user.CreateTime.UnixMilli(),
			Ex:          "",
			Station:     user.Station,
			AreaCode:    user.AreaCode,
		}
		var res []*common.Department
		if len(departmentMemberList) > 0 {
			for _, member := range departmentMemberList {
				if rootDepartment, err := o.Database.GetDepartment(ctx, member.DepartmentID); err == nil {
					var departmentList []*table.Department
					departmentList = append(departmentList, rootDepartment)
					parentID := rootDepartment.ParentID
					for {
						if parentID == "" {
							break
						}
						subDepartment, err := o.Database.GetDepartment(ctx, parentID)
						if err == nil {
							departmentList = append(departmentList, subDepartment)
							parentID = subDepartment.ParentID
						} else if errors.Is(err, gorm.ErrRecordNotFound) {
							break
						} else {
							return nil, err
						}
					}
					res = make([]*common.Department, 0, len(departmentList)+1)
					res = append(res, orgDepartment)
					for i := len(departmentList) - 1; i >= 0; i-- {
						item := departmentList[i]
						res = append(res, &common.Department{
							DepartmentID:   item.DepartmentID,
							FaceURL:        item.FaceURL,
							Name:           item.Name,
							ParentID:       item.ParentID,
							Order:          item.Order,
							DepartmentType: item.DepartmentType,
							RelatedGroupID: item.RelatedGroupID,
							CreateTime:     item.CreateTime.UnixMilli(),
							Position:       member.Position,
						})
					}
				} else if errors.Is(err, gorm.ErrRecordNotFound) {
					res = []*common.Department{orgDepartment}
				} else {
					return nil, err
				}
				if req.IsGetAllDepartment {
					resp.OrganizationUserList = append(resp.OrganizationUserList, &organization.GetSearchDepartmentUserOrganizationUser{
						DepartmentList:   res,
						OrganizationUser: organizationUser,
						Position:         member.Position,
					})
				}
			}
		}
		if !req.IsGetAllDepartment {
			resp.OrganizationUserList = append(resp.OrganizationUserList, &organization.GetSearchDepartmentUserOrganizationUser{
				DepartmentList:   res,
				OrganizationUser: organizationUser,
			})
		}
	}
	return resp, nil
}

func (o *organizationSvr) SortDepartmentList(ctx context.Context, req *organization.SortDepartmentListReq) (*organization.SortDepartmentListResp, error) {
	resp := &organization.SortDepartmentListResp{}
	if req.DepartmentID == req.NextDepartmentID {
		return nil, errs.ErrArgs.Wrap("department id equal")
	}
	idList := append(make([]string, 0, 3), req.DepartmentID)
	if req.ParentID != "" {
		if req.ParentID == req.DepartmentID || req.ParentID == req.NextDepartmentID {
			return nil, errs.ErrArgs.Wrap("parent department id error")
		}
		idList = append(idList, req.ParentID)
	}
	if req.NextDepartmentID != "" {
		idList = append(idList, req.NextDepartmentID)
	}
	departments, err := o.Database.GetList(ctx, idList)
	if err != nil {
		return nil, err
	}
	if len(idList) != len(departments) {
		return nil, errs.ErrArgs.Wrap("department id not found")
	}
	if req.NextDepartmentID == "" { // 添加到最后一个
		order, err := o.Database.GetMaxOrder(ctx, req.ParentID)
		if err != nil {
			return nil, errs.ErrArgs.Wrap(" get max order " + err.Error())
		}
		order++
		if order == 0 {
			order++
		}
		err = o.Database.UpdateDepartment(ctx, &table.Department{
			DepartmentID: req.DepartmentID,
			ParentID:     req.ParentID,
			Order:        order,
		})
		if err != nil {
			return nil, errs.ErrArgs.Wrap("update" + err.Error())
		}
	} else {
		var nextDepartment *table.Department
		for i := 0; i < len(departments); i++ {
			if departments[i].DepartmentID == req.NextDepartmentID {
				nextDepartment = departments[i]
				break
			}
		}
		if nextDepartment == nil {
			return nil, err
		}
		if err := o.Database.UpdateOrderIncrement(ctx, nextDepartment.ParentID, nextDepartment.Order); err != nil {
			return nil, errs.ErrArgs.Wrap(" get max order " + err.Error())
		}
		err = o.Database.UpdateParentIDOrder(ctx, req.DepartmentID, req.ParentID, nextDepartment.Order)
		if err != nil {
			return nil, errs.ErrArgs.Wrap(" update " + err.Error())
		}
	}
	return resp, nil
}

func (o *organizationSvr) SortOrganizationUserList(ctx context.Context, req *organization.SortOrganizationUserListReq) (*organization.SortOrganizationUserListResp, error) {
	resp := &organization.SortOrganizationUserListResp{}
	// TODO 待实现
	return resp, nil
}

func (o *organizationSvr) CreateNewOrganizationMember(ctx context.Context, req *organization.CreateNewOrganizationMemberReq) (*organization.CreateNewOrganizationMemberResp, error) {
	resp := &organization.CreateNewOrganizationMemberResp{}
	if req.OrganizationUser == nil {
		return nil, errs.ErrArgs.Wrap("req.OrganizationUser is nil")
	}
	if req.UserIdentity == nil {
		return nil, errs.ErrArgs.Wrap(" req.UserIdentity is nil")
	}
	if req.UserIdentity.Account == "" {
		return nil, errs.ErrArgs.Wrap("account is empty")
	}

	if len(req.DepartmentMemberList) > 0 {
		departmentIDList := make([]string, 0, len(req.DepartmentMemberList))
		for _, member := range req.DepartmentMemberList {
			//if member.Position == "" {
			//	resp.CommonResp.ErrCode = constant.FormattingError
			//	resp.CommonResp.ErrMsg = "position is empty"
			//	return resp, nil
			//}
			departmentIDList = append(departmentIDList, member.DepartmentID)
		}
		departments, err := o.Database.GetList(ctx, departmentIDList)
		if err != nil {
			return nil, err
		}
		if len(departments) != len(departmentIDList) {
			return nil, errs.ErrArgs.Wrap("department not existence")
		}
	}
	if req.OrganizationUser.UserID == "" {
		req.OrganizationUser.UserID = GenUserID()
	}
	user := &table.OrganizationUser{
		UserID:      req.OrganizationUser.UserID,
		Nickname:    req.OrganizationUser.Nickname,
		EnglishName: req.OrganizationUser.EnglishName,
		FaceURL:     req.OrganizationUser.FaceURL,
		Gender:      req.OrganizationUser.Gender,
		Mobile:      req.OrganizationUser.Mobile,
		Telephone:   req.OrganizationUser.Telephone,
		Birth:       time.UnixMilli(req.OrganizationUser.Birth),
		Email:       req.OrganizationUser.Email,
		Order:       req.OrganizationUser.Order,
		Status:      req.OrganizationUser.Status,
		Station:     req.OrganizationUser.Station,
		AreaCode:    req.OrganizationUser.AreaCode,
	}
	if _, err := o.Database.GetOrganizationUser(ctx, user.UserID); err == nil {
		return nil, errs.ErrArgs.Wrap("has registered")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if user.Nickname == "" && user.EnglishName == "" {
		return nil, errs.ErrArgs.Wrap("nickname and englishName is empty")
	}
	if err := o.Database.CreateOrganizationUser(ctx, user); err != nil {
		return nil, err
	}
	if len(req.DepartmentMemberList) > 0 {
		members := make([]*table.DepartmentMember, 0)
		for _, member := range req.DepartmentMemberList {
			var terminationTime *time.Time
			if member.TerminationTime != constant.NilTimestamp {
				t := time.UnixMilli(member.TerminationTime)
				terminationTime = &t
			}
			members = append(members, &table.DepartmentMember{
				UserID:          user.UserID,
				DepartmentID:    member.DepartmentID,
				Order:           member.Order,
				Position:        member.Position,
				Leader:          member.Leader,
				Status:          member.Status,
				EntryTime:       time.UnixMilli(member.EntryTime),
				TerminationTime: terminationTime,
			})
		}
		if err := o.Database.CreateDepartmentMemberList(ctx, members); err != nil {
			return nil, err
		}
	}
	return resp, nil
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
	resp := &organization.GetUserFullListResp{}
	var (
		total int64
		users []*table.OrganizationUser
		err   error
	)
	if len(req.UserIDList) == 0 {
		if req.ShowNumber == 0 {
			req.ShowNumber = 10
		}
		total, users, err = o.Database.GetPage(ctx, int(req.PageNumber), int(req.ShowNumber))
	} else {
		users, err = o.Database.GetOrganizationUserList(ctx, req.UserIDList)
		total = int64(len(users))
	}
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	for _, user := range users {
		members, err := o.Database.GetDepartmentMemberByUserID(ctx, user.UserID)
		if err != nil {
			return nil, err
		}
		departmentIDList := make([]string, 0, len(members))
		for _, member := range members {
			departmentIDList = append(departmentIDList, member.DepartmentID)
		}
		departments, err := o.Database.GetList(ctx, departmentIDList)
		if err != nil {
			return nil, err
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

func (o *organizationSvr) SearchUsersFullInfo(ctx context.Context, req *organization.SearchUsersFullInfoReq) (*organization.SearchUsersFullInfoResp, error) {
	//TODO implement me
	panic("implement me")
}
