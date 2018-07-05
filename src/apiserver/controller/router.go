package controller

import (
	"github.com/astaxie/beego"
)

func InitRouter() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSRouter("/sign-in",
				&AuthController{},
				"post:SignInAction"),
			beego.NSRouter("/ext-auth",
				&AuthController{},
				"get:ExternalAuthAction"),
			beego.NSRouter("/sign-up",
				&AuthController{},
				"post:SignUpAction"),
			beego.NSRouter("/log-out",
				&AuthController{},
				"get:LogOutAction"),
			beego.NSRouter("/user-exists",
				&AuthController{},
				"get:UserExists"),
			beego.NSRouter("/users/current",
				&AuthController{},
				"get:CurrentUserAction"),
			beego.NSRouter("/systeminfo",
				&AuthController{},
				"get:GetSystemInfo"),
			beego.NSRouter("/reset-password",
				&AuthController{},
				"post:ResetPassword"),
			beego.NSRouter("/users",
				&UserController{},
				"get:GetUsersAction"),
			beego.NSRouter("/users/:id([0-9]+)/password",
				&UserController{},
				"put:ChangePasswordAction"),
			beego.NSRouter("/users/changeaccount",
				&UserController{},
				"put:ChangeUserAccount"),
			beego.NSRouter("/adduser",
				&SystemAdminController{},
				"post:AddUserAction"),
			beego.NSRouter("/users/:id([0-9]+)",
				&SystemAdminController{},
				"get:GetUserAction;put:UpdateUserAction;delete:DeleteUserAction"),
			beego.NSRouter("/users/:id([0-9]+)/systemadmin",
				&SystemAdminController{},
				"put:ToggleSystemAdminAction"),
			beego.NSRouter("/projects",
				&ProjectController{},
				"head:ProjectExists;get:GetProjectsAction;post:CreateProjectAction"),
			beego.NSRouter("/projects/:id([0-9]+)/publicity",
				&ProjectController{},
				"put:ToggleProjectPublicAction"),
			beego.NSRouter("/projects/:id([0-9]+)",
				&ProjectController{},
				"get:GetProjectAction;delete:DeleteProjectAction"),
			beego.NSRouter("/projects/:id([0-9]+)/members",
				&ProjectMemberController{},
				"get:GetProjectMembersAction;post:AddOrUpdateProjectMemberAction"),
			beego.NSRouter("/projects/:projectId([0-9]+)/members/:userId([0-9]+)",
				&ProjectMemberController{},
				"delete:DeleteProjectMemberAction"),
			beego.NSRouter("/images",
				&ImageController{},
				"get:GetImagesAction;delete:DeleteImageAction"),
			beego.NSRouter("/images/:imagename(.*)",
				&ImageController{},
				"get:GetImageDetailAction;delete:DeleteImageTagAction"),
			beego.NSRouter("/images/building",
				&ImageController{},
				"post:BuildImageAction"),
			beego.NSRouter("/images/dockerfilebuilding",
				&ImageController{},
				"post:DockerfileBuildImageAction"),
			beego.NSRouter("/images/dockerfile",
				&ImageController{},
				"get:GetImageDockerfileAction"),
			beego.NSRouter("/images/registry",
				&ImageController{},
				"get:GetImageRegistryAction"),
			beego.NSRouter("/images/preview",
				&ImageController{},
				"post:DockerfilePreviewAction"),
			beego.NSRouter("/images/configclean",
				&ImageController{},
				"delete:ConfigCleanAction"),
			beego.NSRouter("/images/:imagename(.*)/existing",
				&ImageController{},
				"get:CheckImageTagExistingAction"),
			beego.NSRouter("/search",
				&SearchSourceController{}, "get:Search"),
			beego.NSRouter("/node",
				&NodeController{}, "get:GetNode"),
			beego.NSRouter("/nodes",
				&NodeController{}, "get:NodeList"),
			beego.NSRouter("/node/toggle",
				&NodeController{}, "get:NodeToggle"),
			beego.NSRouter("/node/:id([0-9]+)/group",
				&NodeController{},
				"get:GetGroupsOfNodeAction;post:AddNodeToGroupAction;delete:RemoveNodeFromGroupAction"),
			beego.NSRouter("/nodegroup",
				&NodeGroupController{},
				"get:GetNodeGroupsAction;post:AddNodeGroupAction;delete:DeleteNodeGroupAction"),
			beego.NSRouter("/nodegroup/existing",
				&NodeGroupController{},
				"get:CheckNodeGroupNameExistingAction"),
			beego.NSNamespace("/storage",
				beego.NSRouter("/setnfs", &StorageController{}, "post:Storage"),
			),
			beego.NSNamespace("/dashboard", beego.NSRouter("/service",
				&DashboardServiceController{},
				"post:GetServiceData"),
				beego.NSRouter("/node",
					&DashboardNodeController{}, "post:GetNodeData"),
				beego.NSRouter("/data",
					&Dashboard{}, "post:GetData"),
				beego.NSRouter("/time",
					&DashboardServiceController{}, "get:GetServerTime"),
			),
			beego.NSRouter("/services",
				&ServiceController{},
				"post:CreateServiceConfigAction;get:GetServiceListAction"),
			beego.NSRouter("/services/exists",
				&ServiceController{},
				"get:ServiceExists"),
			beego.NSRouter("/services/rollingupdate/image",
				&ServiceRollingUpdateController{},
				"get:GetRollingUpdateServiceImageConfigAction;patch:PatchRollingUpdateServiceImageAction"),
			beego.NSRouter("/services/rollingupdate/nodegroup",
				&ServiceRollingUpdateController{},
				"get:GetRollingUpdateServiceNodeGroupConfigAction;patch:PatchRollingUpdateServiceNodeGroupAction"),
			beego.NSRouter("/services/:id([0-9]+)",
				&ServiceController{},
				"delete:DeleteServiceAction"),
			beego.NSRouter("/services/deployment",
				&ServiceController{},
				"post:DeployServiceAction"),
			beego.NSRouter("/services/config",
				&ServiceConfigController{},
				"post:SetConfigServiceStepAction;get:GetConfigServiceStepAction;delete:DeleteServiceStepAction"),
			beego.NSRouter("/services/reconfig",
				&ServiceConfigController{},
				"get:GetConfigServiceFromDBAction"),
			beego.NSRouter("/services/:id([0-9]+)/status",
				&ServiceController{},
				"get:GetServiceStatusAction"),
			beego.NSRouter("/services/selectservices",
				&ServiceController{},
				"get:GetSelectableServicesAction"),
			beego.NSRouter("/services/yaml/upload",
				&ServiceController{},
				"post:UploadYamlFileAction"),
			beego.NSRouter("/services/yaml/download",
				&ServiceController{},
				"get:DownloadDeploymentYamlFileAction"),
			beego.NSRouter("/images/dockerfile/upload",
				&ImageController{},
				"post:UploadDockerfileFileAction"),
			beego.NSRouter("/images/dockerfile/download",
				&ImageController{},
				"get:DownloadDockerfileFileAction"),
			beego.NSRouter("/services/:id([0-9]+)/info",
				&ServiceController{},
				"get:GetServiceInfoAction"),
			beego.NSRouter("/services/info",
				&ServiceController{},
				"post:StoreServiceRoute"),
			beego.NSRouter("/services/:id([0-9]+)/test",
				&ServiceController{},
				"post:DeployServiceTestAction"),
			beego.NSRouter("/services/:id([0-9]+)/toggle",
				&ServiceController{},
				"put:ToggleServiceAction"),
			beego.NSRouter("/services/:id([0-9]+)/scale",
				&ServiceController{},
				"put:ScaleServiceAction;get:GetScaleStatusAction"),
			beego.NSRouter("/services/:id([0-9]+)/publicity",
				&ServiceController{},
				"put:ServicePublicityAction"),
			beego.NSRouter("/files/upload",
				&FileUploadController{},
				"post:Upload"),
			beego.NSRouter("/files/download",
				&FileUploadController{},
				"head:DownloadProbe;get:Download"),
			beego.NSRouter("/files/list",
				&FileUploadController{},
				"post:ListFiles"),
			beego.NSRouter("/files/remove",
				&FileUploadController{},
				"post:RemoveFile"),
			beego.NSRouter("/jenkins-job/:userID([0-9]+)/:buildNumber([0-9]+)",
				&JenkinsJobCallbackController{},
				"get:BuildNumberCallback"),
			beego.NSRouter("/jenkins-job/console",
				&JenkinsJobController{},
				"get:Console"),
			beego.NSRouter("/jenkins-job/stop",
				&JenkinsJobController{},
				"get:Stop"),
			beego.NSRouter("/email/ping",
				&EmailController{},
				"post:Ping"),
			beego.NSRouter("/email/notification",
				&EmailController{},
				"post:GrafanaNotification"),
			beego.NSRouter("/operations",
				&OperationController{},
				"get:OperationList"),
			beego.NSRouter("/forgot-password",
				&EmailController{},
				"post:ForgotPasswordEmail"),
		),
	)

	beego.AddNamespace(ns)
	beego.Router("/deploy/:owner_name/:project_name/:service_name", &ServiceShowController{})
	beego.SetStaticPath("/swagger", "swagger")
}
