package router

import (
	"git/inspursoft/board/src/apiserver/controller"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSRouter("/sign-in",
				&controller.AuthController{},
				"post:SignInAction"),
			beego.NSRouter("/sign-up",
				&controller.AuthController{},
				"post:SignUpAction"),
			beego.NSRouter("/log-out",
				&controller.AuthController{},
				"get:LogOutAction"),
			beego.NSRouter("/user-exists",
				&controller.AuthController{},
				"get:UserExists"),
			beego.NSRouter("/users/current",
				&controller.AuthController{},
				"get:CurrentUserAction"),
			beego.NSRouter("/users",
				&controller.UserController{},
				"get:GetUsersAction"),
			beego.NSRouter("/users/:id([0-9]+)/password",
				&controller.UserController{},
				"put:ChangePasswordAction"),
			beego.NSRouter("/users/changeaccount",
				&controller.UserController{},
				"put:ChangeUserAccount"),
			beego.NSRouter("/adduser",
				&controller.SystemAdminController{},
				"post:AddUserAction"),
			beego.NSRouter("/users/:id([0-9]+)",
				&controller.SystemAdminController{},
				"get:GetUserAction;put:UpdateUserAction;delete:DeleteUserAction"),
			beego.NSRouter("/users/:id([0-9]+)/systemadmin",
				&controller.SystemAdminController{},
				"put:ToggleSystemAdminAction"),
			beego.NSRouter("/projects",
				&controller.ProjectController{},
				"head:ProjectExists;get:GetProjectsAction;post:CreateProjectAction"),
			beego.NSRouter("/projects/:id([0-9]+)/publicity",
				&controller.ProjectController{},
				"put:ToggleProjectPublicAction"),
			beego.NSRouter("/projects/:id([0-9]+)",
				&controller.ProjectController{},
				"get:GetProjectAction;delete:DeleteProjectAction"),
			beego.NSRouter("/projects/:id([0-9]+)/members",
				&controller.ProjectMemberController{},
				"get:GetProjectMembersAction;post:AddOrUpdateProjectMemberAction"),
			beego.NSRouter("/projects/:projectId([0-9]+)/members/:userId([0-9]+)",
				&controller.ProjectMemberController{},
				"delete:DeleteProjectMemberAction"),
			beego.NSRouter("/images",
				&controller.ImageController{},
				"get:GetImagesAction;delete:DeleteImageAction"),
			beego.NSRouter("/images/:imagename(.*)",
				&controller.ImageController{},
				"get:GetImageDetailAction;delete:DeleteImageTagAction"),
			beego.NSRouter("/images/building",
				&controller.ImageController{},
				"post:BuildImageAction"),
			beego.NSRouter("/images/dockerfile",
				&controller.ImageController{},
				"get:GetImageDockerfileAction"),
			beego.NSRouter("/images/preview",
				&controller.ImageController{},
				"post:DockerfilePreviewAction"),
			beego.NSRouter("/images/configclean",
				&controller.ImageController{},
				"delete:ConfigCleanAction"),
			beego.NSRouter("/search",
				&controller.SearchSourceController{}, "get:Search"),
			beego.NSRouter("/node",
				&controller.NodeController{}, "get:GetNode"),
			beego.NSRouter("/nodes",
				&controller.NodeController{}, "get:NodeList"),
			beego.NSRouter("/node/toggle",
				&controller.NodeController{}, "get:NodeToggle"),
			beego.NSNamespace("/dashboard", beego.NSRouter("/service",
				&controller.DashboardServiceController{},
				"post:GetServiceData"),
				beego.NSRouter("/node",
					&controller.DashboardNodeController{}, "post:GetNodeData"),
				beego.NSRouter("/data",
					&controller.Dashboard{}, "post:GetData"),
				beego.NSRouter("/time",
					&controller.DashboardServiceController{}, "get:GetServerTime"),
			),
			beego.NSRouter("/git/serve",
				&controller.GitRepoController{},
				"post:CreateServeRepo"),
			beego.NSRouter("/git/repo",
				&controller.GitRepoController{},
				"post:InitUserRepo"),
			beego.NSRouter("/git/push",
				&controller.GitRepoController{},
				"post:PushObjects"),
			beego.NSRouter("/git/pull",
				&controller.GitRepoController{},
				"post:PullObjects"),
			beego.NSRouter("/services",
				&controller.ServiceController{},
				"head:ServiceExists;post:CreateServiceConfigAction;get:GetServiceListAction"),
			beego.NSRouter("/services/:id([0-9]+)",
				&controller.ServiceController{},
				"delete:DeleteServiceAction"),
			beego.NSRouter("/services/deployment",
				&controller.ServiceDeployController{},
				"post:DeployServiceAction"),
			beego.NSRouter("/services/config",
				&controller.ServiceConfigController{},
				"post:SetConfigServiceStepAction;get:GetConfigServiceStepAction;delete:DeleteServiceStepAction"),
			beego.NSRouter("/services/status/:service_name([a-z0-9]+)",
				&controller.ServiceController{},
				"get:GetServiceStatusAction"),
			beego.NSRouter("/services/yaml/upload",
				&controller.ConfigFilesController{},
				"post:UploadDeploymentYamlFileAction"),
			beego.NSRouter("/services/yaml/download",
				&controller.ConfigFilesController{},
				"get:DownloadDeploymentYamlFileAction"),
			beego.NSRouter("/services/dockerfile/upload",
				&controller.ConfigFilesController{},
				"post:UploadDockerfileFileAction"),
			beego.NSRouter("/services/dockerfile/download",
				&controller.ConfigFilesController{},
				"get:DownloadDockerfileFileAction"),
			beego.NSRouter("/services/info/:service_name([a-z0-9]+)",
				&controller.ServiceController{},
				"get:GetServiceInfoAction"),
			beego.NSRouter("/services/info",
				&controller.ServiceController{},
				"post:StoreServiceRoute"),
			beego.NSRouter("/services/:id([0-9]+)/test",
				&controller.ServiceController{},
				"post:DeployServiceTestAction"),
			beego.NSRouter("/services/:id([0-9]+)/toggle",
				&controller.ServiceController{},
				"put:ToggleServiceAction"),
			beego.NSRouter("/services/:id([0-9]+)/publicity",
				&controller.ServiceController{},
				"put:ServicePublicityAction"),
			beego.NSRouter("/files/upload",
				&controller.FileUploadController{},
				"post:Upload"),
			beego.NSRouter("/files/list",
				&controller.FileUploadController{},
				"post:ListFiles"),
			beego.NSRouter("/files/remove",
				&controller.FileUploadController{},
				"post:RemoveFile"),
			beego.NSRouter("/jenkins-job/lastbuildnumber",
				&controller.JenkinsJobController{},
				"get:GetLastBuildNumber"),
			beego.NSRouter("/jenkins-job/console",
				&controller.JenkinsJobController{},
				"get:Console"),
			beego.NSRouter("/jenkins-job/stop",
				&controller.JenkinsJobController{},
				"get:Stop"),
		),
	)

	beego.AddNamespace(ns)
	beego.Router("/deploy/:owner_name/:project_name/:service_name", &controller.ServiceShowController{})
	beego.SetStaticPath("/swagger", "swagger")
}
