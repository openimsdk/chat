package api

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/chat/pkg/common/config"
	"github.com/OpenIMSDK/chat/pkg/proto/chat"
	"github.com/OpenIMSDK/chat/pkg/proto/organization"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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
	//todo:?
}

func (o *Org) CreateOrganizationUser(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.CreateOrganizationUser, o.organizationClient, c)
}

func (o *Org) UpdateOrganizationUser(c *gin.Context) {
	//todo:?
}
func (o *Org) CreateDepartmentMember(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.CreateDepartmentMember, o.organizationClient, c)
}
func (o *Org) GetUserInDepartment(c *gin.Context) {
	//todo:?
	a2r.Call(organization.OrganizationClient.GetUserInDepartment, o.organizationClient, c)
}

func (o *Org) UpdateUserInDepartment(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.UpdateUserInDepartment, o.organizationClient, c)
}

func (o *Org) DeleteUserInDepartment(c *gin.Context) {
	//todo:?
	a2r.Call(organization.OrganizationClient.DeleteUserInDepartment, o.organizationClient, c)
}

func (o *Org) GetSearchUserList(c *gin.Context) {
	//todo:?
}

func (o *Org) SetOrganization(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.SetOrganization, o.organizationClient, c)
}

func (o *Org) GetOrganization(c *gin.Context) {
	a2r.Call(organization.OrganizationClient.GetOrganization, o.organizationClient, c)
}

func (o *Org) MoveUserDepartment(c *gin.Context) {
	//todo:?
}

func (o *Org) GetSubDepartment(c *gin.Context) {
	//todo:?

}

func (o *Org) GetSearchDepartmentUser(c *gin.Context) {
	//todo:?
}

func (o *Org) GetSearchDepartmentUserWithoutToken(c *gin.Context) {
	//todo:?
}

func (o *Org) GetOrganizationDepartment(c *gin.Context) {
	//todo:?
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

func (o *Org) BatchImport(c *gin.Context) {
	//todo:?
}

func (o *Org) BatchImportTemplate(c *gin.Context) {
	//data := []byte(time.Now().Format("2006-01-02 15:04:05") + " hello world!")
	c.Header("Content-Disposition", "attachment; filename=组织架构批量导入模板.xlsx")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Description", "File Transfer")
	//c.Header("Content-Length", strconv.Itoa(len(data)))
	//c.Data(http.StatusOK, "application/octet-stream", data)
	c.File(config.Config.ImportTemplate)
}

func (o *Org) DeleteOrganizationUser(c *gin.Context) {

}