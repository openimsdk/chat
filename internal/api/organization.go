package api

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/common/constant"
	"github.com/OpenIMSDK/chat/pkg/common/xlsx"
	"github.com/OpenIMSDK/chat/pkg/common/xlsx/model"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	"github.com/OpenIMSDK/chat/pkg/proto/organization"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func NewOrg(chatConn, orgConn grpc.ClientConnInterface) *Org {
	return &Org{
		organizationClient: organization.NewOrganizationClient(orgConn),
		chatClient:         chat.NewChatClient(chatConn),
	}
}

type Org struct {
	organizationClient organization.OrganizationClient
	chatClient         chat.ChatClient
}

func (o *Org) CreateDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.CreateDepartment, o.organizationClient, c)
}

func (o *Org) UpdateDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.UpdateDepartment, o.organizationClient, c)
}

func (o *Org) DeleteDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.DeleteDepartment, o.organizationClient, c)
}

func (o *Org) GetDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.GetDepartment, o.organizationClient, c)
}

func (o *Org) CreateOrganizationUser(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.CreateOrganizationUser, o.organizationClient, c)
}

func (o *Org) UpdateOrganizationUser(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.UpdateOrganizationUser, o.organizationClient, c)
}
func (o *Org) CreateDepartmentMember(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.CreateDepartmentMember, o.organizationClient, c)
}
func (o *Org) GetUserInDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.GetUserInDepartment, o.organizationClient, c)
}

func (o *Org) UpdateUserInDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.UpdateUserInDepartment, o.organizationClient, c)
}

func (o *Org) DeleteUserInDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.DeleteUserInDepartment, o.organizationClient, c)
}

func (o *Org) GetSearchUserList(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.GetSearchUserList, o.organizationClient, c)
}

func (o *Org) SetOrganization(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.SetOrganization, o.organizationClient, c)
}

func (o *Org) GetOrganization(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.GetOrganization, o.organizationClient, c)
}

func (o *Org) MoveUserDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.MoveUserDepartment, o.organizationClient, c)
}

func (o *Org) GetSubDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.GetSubDepartment, o.organizationClient, c)
}

func (o *Org) GetSearchDepartmentUser(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.GetSearchDepartmentUser, o.organizationClient, c)
}

func (o *Org) GetOrganizationDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.GetOrganizationDepartment, o.organizationClient, c)
}

func (o *Org) SortDepartmentList(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.SortDepartmentList, o.organizationClient, c)
}

func (o *Org) SortOrganizationUserList(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.SortOrganizationUserList, o.organizationClient, c)
}

func (o *Org) CreateNewOrganizationMember(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.CreateNewOrganizationMember, o.organizationClient, c)
}

func (o *Org) DeleteOrganizationUser(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.DeleteOrganizationUser, o.organizationClient, c)
}

func (o *Org) BatchImportTemplate(c *gin.Context) {
	md5Sum := md5.Sum(config.ImportTemplate)
	md5Val := hex.EncodeToString(md5Sum[:])
	if c.GetHeader("If-None-Match") == md5Val {
		c.Status(http.StatusNotModified)
		return
	}
	c.Header("Content-Disposition", "attachment; filename=template.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Length", strconv.Itoa(len(config.ImportTemplate)))
	c.Header("ETag", md5Val)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", config.ImportTemplate)
}

func (o *Org) BatchImport(c *gin.Context) {
	resp, err := o.batchImport(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

func (o *Org) batchImport(c *gin.Context) (*organization.BatchImportResp, error) {
	formFile, err := c.FormFile("data")
	if err != nil {
		return nil, err
	}
	file, err := formFile.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var (
		departments []model.Department
		users       []model.OrganizationUser
	)
	if err := xlsx.ParseAll(file, &departments, &users); err != nil {
		return nil, err
	}
	rpcReq := organization.BatchImportReq{}
	for _, department := range departments {
		var names []string
		for _, name := range strings.Split(department.Parent, "/") {
			name = strings.TrimSpace(name)
			if name != "" {
				names = append(names, name)
			}
		}
		rpcReq.DepartmentList = append(rpcReq.DepartmentList, &organization.BatchImportDepartment{
			DepartmentID:   department.DepartmentID,
			FaceURL:        department.FaceURL,
			Name:           department.Name,
			DepartmentType: department.DepartmentType,
			RelatedGroupID: department.RelatedGroupID,
			ParentDepartmentName: &organization.BatchImportUserDepartmentName{
				HierarchyName: names,
			},
		})
	}
	for _, user := range users {
		var birth int64
		if user.Birth == "" {
			birth = constant.NilTimestamp
		} else {
			var arr []string
			for _, s := range []string{"-", "/"} {
				if index := strings.Index(user.Birth, s); index >= 0 {
					arr = strings.Split(user.Birth, s)
					break
				}
			}
			if len(arr) != 3 {
				return nil, errs.ErrArgs.Wrap(user.Birth + " birth parse error")
			}
			for i, s := range arr[1:] {
				if len(s) == 1 {
					arr[i] = "0" + s
				}
			}
			t, err := time.Parse("2006-01-02", strings.Join(arr, "-"))
			if err != nil {
				return nil, errs.ErrArgs.Wrap(user.Birth + " birth parse error " + err.Error())
			}
			birth = t.UnixMilli()
		}
		var gender int32
		switch strings.ToLower(user.Gender) {
		case "male", "男", "1":
			gender = constant.GenderMale
		case "female", "女", "0":
			gender = constant.GenderFemale
		default:
			gender = constant.GenderUnknown
		}
		if user.Account == "" {
			return nil, errs.ErrArgs.Wrap("account is empty")
		}
		if user.Password == "" {
			return nil, errs.ErrArgs.Wrap("password is empty")
		}
		var list []*organization.BatchImportUserDepartmentNamePosition
		user.Department = strings.TrimSpace(user.Department)
		if user.Department != "" {
			for _, dstr := range strings.Split(user.Department, ";") { // 分为多个部门
				dstr = strings.TrimSpace(dstr)
				temp := strings.Split(dstr, "/") // 部门路径
				var item organization.BatchImportUserDepartmentNamePosition
				for i, name := range temp {
					name = strings.TrimSpace(name)
					if len(temp) == i+1 { // 最后一个
						arr := strings.Split(name, ":")
						switch len(arr) {
						case 1:
							ts := strings.TrimSpace(arr[0])
							if ts != "" {
								item.HierarchyName = append(item.HierarchyName, ts)
							}
						case 2:
							ts := strings.TrimSpace(arr[0])
							if ts != "" {
								item.HierarchyName = append(item.HierarchyName, ts)
							}
							item.Position = strings.TrimSpace(arr[1])
						default:
							return nil, errs.ErrArgs.Wrap("password is empty")
						}
					} else {
						if name != "" { // 非最后一个
							item.HierarchyName = append(item.HierarchyName, name)
						}
					}
				}
				list = append(list, &item)
			}
		}
		pwdSum := md5.Sum([]byte(user.Password))
		rpcReq.UserList = append(rpcReq.UserList, &organization.BatchImportUser{
			UserID:                 user.UserID,
			Nickname:               user.Nickname,
			EnglishName:            user.EnglishName,
			FaceURL:                user.FaceURL,
			Gender:                 gender,
			Mobile:                 user.Mobile,
			Telephone:              user.Telephone,
			Birth:                  birth,
			Email:                  user.Email,
			Account:                user.Account,
			Password:               hex.EncodeToString(pwdSum[:]),
			AreaCode:               user.AreaCode,
			Station:                user.Station,
			UserDepartmentNameList: list,
		})
	}
	return o.organizationClient.BatchImport(c, &rpcReq)
}
