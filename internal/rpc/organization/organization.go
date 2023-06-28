package organization

import (
	common2 "Open-IM-Organization/internal/rpc/common"
	"Open-IM-Organization/pkg/common/config"
	"Open-IM-Organization/pkg/proto/chat"
	"context"
	"errors"
	"strconv"
	"time"

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
	log "github.com/OpenIMSDK/open_log"
	utils "github.com/OpenIMSDK/open_utils"
	"google.golang.org/grpc"
	"gorm.io/gorm"
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
	resp := &organization.UpdateDepartmentResp{CommonResp: &common.CommonResp{}}

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
	resp := &organization.GetOrganizationDepartmentResp{CommonResp: &common.CommonResp{}, DepartmentList: []*organization.DepartmentInfo{}}

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
	resp := &organization.DeleteDepartmentResp{CommonResp: &common.CommonResp{}}
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
	resp := &organization.GetDepartmentResp{CommonResp: &common.CommonResp{}}

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
	resp := &organization.CreateOrganizationUserResp{CommonResp: &common.CommonResp{}}
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
	resp := &organization.UpdateOrganizationUserResp{CommonResp: &common.CommonResp{}}
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
	resp := &organization.DeleteOrganizationUserResp{CommonResp: &common.CommonResp{}}
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
	resp := &organization.CreateDepartmentMemberResp{CommonResp: &common.CommonResp{}}
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
	resp := &organization.GetUserInDepartmentResp{CommonResp: &common.CommonResp{}}
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
	resp := &organization.DeleteUserInDepartmentResp{CommonResp: &common.CommonResp{}}
	err := o.Database.DeleteDepartmentMemberByKey(ctx, req.UserID, req.DepartmentID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *organizationSvr) UpdateUserInDepartment(ctx context.Context, req *organization.UpdateUserInDepartmentReq) (*organization.UpdateUserInDepartmentResp, error) {
	resp := &organization.UpdateUserInDepartmentResp{CommonResp: &common.CommonResp{}}
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
	resp := &organization.GetSearchUserListResp{CommonResp: &common.CommonResp{}, UserList: []*common.UserInDepartment{}}
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
	resp := &organization.SetOrganizationResp{CommonResp: &common.CommonResp{}}
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
	resp := &organization.GetOrganizationResp{CommonResp: &common.CommonResp{}}
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
		CommonResp:              &common.CommonResp{},
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
			return resp, nil
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
		CommonResp:           &common.CommonResp{},
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
	resp := &organization.GetUserInfoResp{CommonResp: &common.CommonResp{}}

	user, err := o.Database.GetOrganizationUser(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.CommonResp.ErrCode = constant.ErrRecordNotFound.ErrCode
			resp.CommonResp.ErrMsg = "user id not found"
		} else {
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
		}
		return resp, nil
	}

	memberList, err := o.Database.GetDepartmentMemberByUserID(ctx, req.UserID)
	if err != nil {
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
		return resp, nil
	}

	departmentMap := make(map[string]*common.Department)
	if len(memberList) > 0 {
		departmentIDList := make([]string, 0, len(memberList))
		for _, member := range memberList {
			departmentIDList = append(departmentIDList, member.DepartmentID)
		}
		departmentList, err := o.Database.GetList(ctx, departmentIDList)
		if err != nil {
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + err.Error()
			return resp, nil
		}
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
	}

	resp.User = &common.OrganizationUser{
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
		DepartmentMemberList: make([]*common.DepartmentMember, 0, len(memberList)),
		Station:              user.Station,
		AreaCode:             user.AreaCode,
	}

	for _, member := range memberList {
		var terminationTime int64
		if member.TerminationTime == nil {
			terminationTime = constant.NilTimestamp
		} else {
			terminationTime = member.TerminationTime.UnixMilli()
		}
		resp.User.DepartmentMemberList = append(resp.User.DepartmentMemberList, &common.DepartmentMember{
			UserID:          member.UserID,
			DepartmentID:    member.DepartmentID,
			Order:           member.Order,
			Position:        member.Position,
			Leader:          member.Leader,
			Status:          member.Status,
			EntryTime:       member.EntryTime.UnixMilli(),
			TerminationTime: terminationTime,
			CreateTime:      member.CreateTime.UnixMilli(),
			Department:      departmentMap[member.DepartmentID],
		})
	}

	return resp, nil
}

func (o *organizationSvr) BatchImport(ctx context.Context, req *organization.BatchImportReq) (*organization.BatchImportResp, error) {
	resp := &organization.BatchImportResp{CommonResp: &common.CommonResp{}}

	createDepartment := func(department *table.Department) error {
		department.DepartmentID = genDepartmentID()
		return o.Database.CreateDepartment(ctx, department)
	}

	if len(req.DepartmentList) > 0 {
		for _, department := range req.DepartmentList {
			var parentID string
			for _, name := range department.ParentDepartmentName.HierarchyName {
				if name == "" {
					return nil, errs.ErrArgs.Wrap("department name is empty")
				}
				d, err := o.Database.GetDepartmentByName(ctx, name, parentID)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					d = &table.Department{
						Name:     name,
						ParentID: parentID,
					}
					if err = createDepartment(d); err != nil {
						return resp, nil
					}
					return nil, err
				} else if err != nil {
					return nil, err
				}
				parentID = d.DepartmentID
			}
		}
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "import department count:", len(req.DepartmentList))

	if len(req.UserList) > 0 {
		var (
			reqs   = make([]*organization.CreateNewOrganizationMemberReq, 0, len(req.UserList))
			idList = make([]string, 0, len(req.UserList))
		)

		for _, user := range req.UserList {
			departmentIDList := make([]string, 0, len(user.UserDepartmentNameList))
			positions := make([]string, 0, len(user.UserDepartmentNameList))

			for _, nameList := range user.UserDepartmentNameList {
				var parentID string
				for _, name := range nameList.HierarchyName {
					if name == "" {
						return nil, errs.ErrArgs.Wrap("department name is empty")
					}
					d, err := o.Database.GetDepartmentByName(ctx, name, parentID)
					if errors.Is(err, gorm.ErrRecordNotFound) {
						d = &table.Department{
							Name:     name,
							ParentID: parentID,
						}
						if err = createDepartment(d); err != nil {
							return nil, errs.ErrArgs.Wrap("create not found department")
						}
						return nil, err
					} else if err != nil {
						return nil, err
					}
					parentID = d.DepartmentID
				}
				departmentIDList = append(departmentIDList, parentID)
				positions = append(positions, nameList.Position)
			}

			departmentMemberList := make([]*common.DepartmentMember, 0, len(user.UserDepartmentNameList))

			for i, departmentID := range departmentIDList {
				log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "user add department", departmentID, "position", positions[i])
				departmentMemberList = append(departmentMemberList, &common.DepartmentMember{
					UserID:       user.UserID,
					DepartmentID: departmentID,
					Position:     positions[i],
				})
			}

			idList = append(idList, user.UserID)

			reqs = append(reqs, &organization.CreateNewOrganizationMemberReq{
				OperationID: req.OperationID,
				OpUserID:    req.OpUserID,
				OrganizationUser: &common.OrganizationUser{
					UserID:      user.UserID,
					Nickname:    user.Nickname,
					EnglishName: user.EnglishName,
					FaceURL:     user.FaceURL,
					Gender:      user.Gender,
					Mobile:      user.Mobile,
					Telephone:   user.Telephone,
					Birth:       user.Birth,
					Email:       user.Email,
					Order:       user.Order,
					Status:      user.Status,
					Station:     user.Station,
					AreaCode:    user.AreaCode,
				},
				DepartmentMemberList: departmentMemberList,
				UserIdentity: &common.UserIdentity{
					Account:  user.Account,
					Password: user.Password,
				},
			})
		}

		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "user parse end", len(req.UserList))

		users, err := o.Database.GetOrganizationUserList(ctx, idList)
		if err != nil {
			return nil, err
		}

		if len(users) != 0 {
			return nil, err
		}

		for i := 0; i < len(reqs); i++ {
			rpcResp, err := o.CreateNewOrganizationMember(context.Background(), reqs[i])
			if err != nil {
				return nil, err
			}
			if rpcResp.CommonResp.ErrCode != constant.NoError {
				return nil, err
			}
		}
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "import user count:", len(req.DepartmentList))
	return resp, nil
}

func (o *organizationSvr) MoveUserDepartment(ctx context.Context, req *organization.MoveUserDepartmentReq) (*organization.MoveUserDepartmentResp, error) {
	resp := &organization.MoveUserDepartmentResp{CommonResp: &common.CommonResp{}}

	if len(req.MoveUserDepartmentList) == 0 {
		return nil, errs.ErrArgs.Wrap("move user department list is empty")
	}

	tx, err := o.Database.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}

	for _, d := range req.MoveUserDepartmentList {
		if _, err := o.Database.GetOrganizationUser(ctx, d.UserID); err != nil {
			return resp, err
		}
		if _, err := o.Database.GetDepartment(ctx, d.DepartmentID); err != nil {
			return resp, err
		}
		if d.CurrentDepartmentID != "" {
			if err := o.Database.DeleteDepartmentMemberByKey(ctx, d.UserID, d.CurrentDepartmentID); err != nil {
				return resp, err
			}
		}
		var terminationTime *time.Time
		if d.TerminationTime != constant.NilTimestamp {
			t := time.UnixMilli(d.TerminationTime)
			terminationTime = &t
		}

		_, err := o.Database.GetDepartmentMemberByKey(ctx, d.UserID, d.DepartmentID)
		if err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrArgs.Wrap("record")
		}

		m := table.DepartmentMember{
			UserID:          d.UserID,
			DepartmentID:    d.DepartmentID,
			Order:           d.Order,
			Position:        d.Position,
			Leader:          d.Leader,
			Status:          d.Status,
			EntryTime:       time.UnixMilli(d.EntryTime),
			TerminationTime: terminationTime,
		}
		if err := o.Database.CreateDepartmentMember(ctx, &m); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return resp, err
	}

	return resp, nil
}

func (o *organizationSvr) GetUserFullList(ctx context.Context, req *organization.GetUserFullListReq) (*organization.GetUserFullListResp, error) {
	//TODO implement me
	panic("implement me")
}

func (o *organizationSvr) SearchUsersFullInfo(ctx context.Context, req *organization.SearchUsersFullInfoReq) (*organization.SearchUsersFullInfoResp, error) {
	resp := &organization.SearchUsersFullInfoResp{CommonResp: &common.CommonResp{}}
	if req.Operation == nil || req.Pagination == nil || req.UserPublicInfo == nil {
		return nil, errs.ErrArgs.Wrap("req.Operation or req.Pagination or req.UserPublicInfo is nil")
	}
	if req.Pagination.ShowNumber == 0 {
		req.Pagination.ShowNumber = 10
	}
	var userIDList []string
	if req.Content != "" {
		// Search for users with matching accounts
		etcdConn := common2.GetDefaultConn(config.Config.RpcRegisterName.OpenImChatName, req.Operation.OperationID, resp.CommonResp)
		if etcdConn == nil {
			return nil, errs.ErrArgs.Wrap("get base server etcd conn is empty")
		}
		rpcResp, err := chat.NewChatClient(etcdConn).GetAccountUser(ctx, &chat.GetAccountUserReq{
			Operation:   req.Operation,
			AccountList: []string{req.Content},
		})
		if err != nil {
			return nil, errs.ErrArgs.Wrap("get base server GetAccountUser error")
		}
		if rpcResp.CommonResp.ErrCode != constant.NoError {
			return nil, errs.ErrArgs.Wrap("get base server GetAccountUser errMsg")
		}
		userIDList = rpcResp.Accounts
	}
	if req.UserPublicInfo.UserID != "" {
		userIDList = append(userIDList, req.UserPublicInfo.UserID)
	}
	total, users, err := o.Database.SearchV2(ctx, req.Content, userIDList, int(req.Pagination.PageNumber), int(req.Pagination.ShowNumber))
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	for _, user := range users {
		members, err := o.Database.GetDepartmentMemberByUserID(ctx, user.UserID)

		if err != nil {
			return resp, err
		}
		departmentIDList := make([]string, 0, len(members))
		for _, member := range members {
			departmentIDList = append(departmentIDList, member.DepartmentID)
		}
		departments, err := o.Database.GetList(ctx, departmentIDList)
		if err != nil {
			return resp, err
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

// func (o *organizationSvr) SearchUsersFullInfo(ctx context.Context, req *organization.SearchUsersFullInfoReq) (*organization.SearchUsersFullInfoResp, error) {
// 	resp := &organization.SearchUsersFullInfoResp{CommonResp: &common.CommonResp{}}
// 	if err := validateSearchUsersFullInfoReq(req); err != nil {
// 		return nil, err
// 	}
// 	userIDList, err := getUserIDList(ctx, o, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	total, _, err := searchUsers(ctx, resp, o, req, userIDList)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp.Total = int32(total)
// 	return nil, err
// }

// func validateSearchUsersFullInfoReq(req *organization.SearchUsersFullInfoReq) error {
// 	if req.Operation == nil || req.Pagination == nil || req.UserPublicInfo == nil {
// 		return errs.ErrArgs.Wrap("operation or pagination or user public info is nil")
// 	}
// 	if req.Pagination.ShowNumber == 0 {
// 		req.Pagination.ShowNumber = 10
// 	}
// 	return nil
// }

// func getUserIDList(ctx context.Context, o *organizationSvr, req *organization.SearchUsersFullInfoReq) ([]string, error) {
// 	var userIDList []string
// 	if req.Content != "" {
// 		etcdConn := common2.GetDefaultConn(config.Config.RpcRegisterName.OpenImChatName, req.Operation.OperationID, nil)
// 		if etcdConn == nil {
// 			return nil, errs.ErrArgs.Wrap("get base server conn error")
// 		}
// 		rpcResp, err := chat.NewChatClient(etcdConn).GetAccountUser(ctx, &chat.GetAccountUserReq{
// 			Operation:   req.Operation,
// 			AccountList: []string{req.Content},
// 		})
// 		if err != nil {
// 			return nil, errs.ErrArgs.Wrap("get base server GetAccountUser error")
// 		}
// 		if rpcResp.CommonResp.ErrCode != constant.NoError {
// 			return nil, errs.ErrArgs.Wrap("get base server GetAccountUser error")
// 		}
// 		for _, userID := range rpcResp.Accounts {
// 			userIDList = append(userIDList, userID)
// 		}
// 	}
// 	if req.UserPublicInfo.UserID != "" {
// 		userIDList = append(userIDList, req.UserPublicInfo.UserID)
// 	}
// 	return userIDList, nil
// }

// func searchUsers(ctx context.Context, resp *organization.SearchUsersFullInfoResp, o *organizationSvr, req *organization.SearchUsersFullInfoReq, userIDList []string) (int, []*table.OrganizationUser, error) {
// 	total, users, err := o.Database.SearchV2(ctx, req.Content, userIDList, int(req.Pagination.PageNumber), int(req.Pagination.ShowNumber))
// 	if err != nil {
// 		return 0, nil, errs.ErrArgs.Wrap("search users error")
// 	}
// 	var orgUsers []*table.OrganizationUser
// 	for _, user := range users {
// 		orgUser, err := getOrgUser(ctx, resp, o, user)
// 		if err != nil {
// 			return 0, nil, err
// 		}
// 		orgUsers = append(orgUsers, orgUser)
// 	}
// 	return int(total), orgUsers, nil
// }

// func getOrgUser(ctx context.Context, resp *organization.SearchUsersFullInfoResp, o *organizationSvr, user *table.OrganizationUser) (*organization.SearchUsersFullInfoResp, error) {
// 	members, err := o.Database.GetDepartmentMemberByUserID(ctx, user.UserID)
// 	if err != nil {
// 		return nil, errs.ErrArgs.Wrap("get user department member error")
// 	}
// 	departmentIDList := make([]string, 0, len(members))
// 	for _, member := range members {
// 		departmentIDList = append(departmentIDList, member.DepartmentID)
// 	}
// 	departments, err := o.Database.GetList(ctx, departmentIDList)
// 	if err != nil {
// 		return nil, errs.ErrArgs.Wrap("get department list error")
// 	}
// 	departmentMap := make(map[string]*common.Department)
// 	for _, department := range departments {
// 		departmentMap[department.DepartmentID] = &common.Department{
// 			DepartmentID:   department.DepartmentID,
// 			FaceURL:        department.FaceURL,
// 			Name:           department.Name,
// 			ParentID:       department.ParentID,
// 			Order:          department.Order,
// 			DepartmentType: department.DepartmentType,
// 			RelatedGroupID: department.RelatedGroupID,
// 			CreateTime:     department.CreateTime.UnixMilli(),
// 		}
// 	}
// 	var departmentMemberList []*common.DepartmentMember
// 	for _, member := range members {
// 		department := departmentMap[member.DepartmentID]
// 		if department == nil {
// 			continue
// 		}
// 		var terminationTime int64
// 		if member.TerminationTime == nil {
// 			terminationTime = constant.NilTimestamp
// 		} else {
// 			terminationTime = member.TerminationTime.UnixMilli()
// 		}
// 		departmentMemberList = append(departmentMemberList, &common.DepartmentMember{
// 			UserID:          member.UserID,
// 			DepartmentID:    member.DepartmentID,
// 			Order:           member.Order,
// 			Position:        member.Position,
// 			Leader:          member.Leader,
// 			Status:          member.Status,
// 			EntryTime:       member.EntryTime.UnixMilli(),
// 			TerminationTime: terminationTime,
// 			CreateTime:      member.CreateTime.UnixMilli(),
// 			Department:      department,
// 		})
// 	}
// 	resp.OrganizationUserList = append(resp.OrganizationUserList,
// 		&common.OrganizationUser{
// 			UserID:               user.UserID,
// 			Nickname:             user.Nickname,
// 			EnglishName:          user.EnglishName,
// 			FaceURL:              user.FaceURL,
// 			Gender:               user.Gender,
// 			Mobile:               user.Mobile,
// 			Telephone:            user.Telephone,
// 			Birth:                user.Birth.UnixMilli(),
// 			Email:                user.Email,
// 			Order:                user.Order,
// 			Status:               user.Status,
// 			CreateTime:           user.CreateTime.UnixMilli(),
// 			DepartmentMemberList: departmentMemberList,
// 			Station:              user.Station,
// 			AreaCode:             user.AreaCode,
// 		})
// 	return resp, nil
// }