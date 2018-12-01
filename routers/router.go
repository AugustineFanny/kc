package routers

import (
	"kuangchi_backend/controllers"
	"kuangchi_backend/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/http"
)

func init() {
	ns := beego.NewNamespace("/kc",
		beego.NSRouter("/register-captcha", &controllers.AuthController{}, "post:RegisterCaptcha"),
		beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
		beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
		beego.NSRouter("/logout", &controllers.AuthController{}, "get:Logout"),
		beego.NSRouter("/reset-password", &controllers.AuthController{}, "post:ResetPassword"),
		beego.NSNamespace("public",
			beego.NSRouter("/currencies", &controllers.PublicController{}, "get:Currencies"),
			beego.NSRouter("/currencies/detail", &controllers.PublicController{}, "get:CurrenciesDetail"),
			beego.NSRouter("/exchanges", &controllers.PublicController{}, "get:Exchanges"),
			beego.NSRouter("/exchange/:exchange", &controllers.PublicController{}, "get:Exchange"),
			beego.NSRouter("/:currency/nodes", &controllers.PublicController{}, "get:Nodes"),
			beego.NSRouter("/IUU/exchange", &controllers.PublicController{}, "get:IUUExchange"),

			//beego.NSRouter("/btcdo-script", &controllers.PublicController{}, "get:GetBtcdoScript;post:BtcdoScript"),
			//beego.NSRouter("/btcdo-script-logs", &controllers.PublicController{}, "get:GetBtcdoScriptLogs"),
			//beego.NSRouter("/btcdo-script-start", &controllers.PublicController{}, "post:BtcdoScriptStart"),
			//beego.NSRouter("/btcdo-script-stop", &controllers.PublicController{}, "post:BtcdoScriptStop"),
		),
		beego.NSNamespace("self",
			beego.NSRouter("/safe", &controllers.SelfController{}, "get:Safe"),
			beego.NSRouter("/info", &controllers.SelfController{}, "get:Info"),
			beego.NSRouter("/avatar", &controllers.SelfController{}, "post:SetAvatar"),
			beego.NSRouter("/change-password", &controllers.SelfController{}, "post:ChangePassword"),
			beego.NSRouter("/name-auth", &controllers.SelfController{}, "get:RealName;post:NameAuth"),
			beego.NSRouter("/two-factor", &controllers.SelfController{}, "get:TfUrl;post:TfInit"),
			beego.NSRouter("/two-factor/d", &controllers.SelfController{}, "post:TfDelete"),
			beego.NSRouter("/send-bind-captcha", &controllers.SelfController{}, "get:SendBindCaptcha"),
			beego.NSRouter("/fund-password", &controllers.SelfController{}, "post:FundPassword"),
			beego.NSRouter("/bind", &controllers.SelfController{}, "post:Bind"),
			beego.NSRouter("/verify-code", &controllers.SelfController{}, "get:SendVerifyCode"),
			beego.NSRouter("/messages", &controllers.SelfController{}, "get:Messages;post:MessagesAllRead"),
			beego.NSRouter("/msg-num", &controllers.SelfController{}, "get:GetMsgNum"),
			beego.NSRouter("/message/:id", &controllers.SelfController{}, "get:MessageRead"),
			beego.NSRouter("/fishing-code", &controllers.SelfController{}, "get:GetFishingCode;post:SetFishingCode"),
			beego.NSRouter("/kyc", &controllers.SelfController{}, "get:GetKyc;post:Kyc"),
			beego.NSRouter("/invitelink", &controllers.SelfController{}, "get:GetInvitelink"),
			beego.NSRouter("/invite-num", &controllers.SelfController{}, "get:GetInviteNum"),
			beego.NSRouter("/invites", &controllers.SelfController{}, "get:GetInvites"),
		),
		beego.NSNamespace("user",
			beego.NSRouter("/:username", &controllers.UserController{}, "get:Info"),
		),
		beego.NSNamespace("captcha",
			beego.NSRouter("/code", &controllers.CaptchaController{}, "get:Code"),
			beego.NSRouter("/send-captcha", &controllers.CaptchaController{}, "get:SendCaptcha"),
		),
		beego.NSNamespace("wallet",
			beego.NSRouter("/currency/:currency", &controllers.WalletController{}, "get:Address"),
			//beego.NSRouter("/withdraw/new", &controllers.WalletController{}, "post:NewWithdraw"),
			beego.NSRouter("/finance", &controllers.WalletController{}, "get:Finance"),
			beego.NSRouter("/transfer-out", &controllers.WalletController{}, "post:TransferOut"),
			beego.NSRouter("/transfers", &controllers.WalletController{}, "get:Transfers"),
			beego.NSRouter("/usable/:currency", &controllers.WalletController{}, "get:Usable"),
			beego.NSRouter("/:currency/locked", &controllers.WalletController{}, "get:GetLocked;post:Locked"),
			beego.NSRouter("/:currency/mining", &controllers.WalletController{}, "get:GetMining"),
			beego.NSRouter("/:currency/mining-stat", &controllers.WalletController{}, "get:GetMiningStat"),
			beego.NSRouter("/locked/:id/transfer", &controllers.WalletController{}, "post:TransferLocked"),
			beego.NSRouter("/instations", &controllers.WalletController{}, "get:Instations"),
			beego.NSRouter("/inlocked", &controllers.WalletController{}, "get:Inlocked"),
			beego.NSRouter("/orders", &controllers.WalletController{}, "get:Orders;post:CreateOrder"),
			beego.NSRouter("/order/:order", &controllers.WalletController{}, "get:GetOrder;delete:DeleteOrder"),
			beego.NSRouter("/hashrate", &controllers.WalletController{}, "get:Hashrate"),
		),
		//后台管理API
		beego.NSRouter("/admin-login", &controllers.AdminController{}, "post:Login"),
		//beego.NSRouter("/simulations", &controllers.AdminController{}, "get:Simulations;post:CreateSimulation"),
		//beego.NSRouter("/simulation/users", &controllers.AdminController{}, "get:GetSUsers;post:CreateSUer;delete:DeleteAllSUsers"),
		//beego.NSRouter("/simulation/batch/users", &controllers.AdminController{}, "post:BatchSUsers"),
		//beego.NSRouter("/simulation/currency", &controllers.AdminController{}, "get:GetSCurrency;post:ChangeSCurrency"),
		//beego.NSRouter("/simulation/user/:uid", &controllers.AdminController{}, "post:ChangeSUser;delete:DeleteSUser"),
		//beego.NSRouter("/simulation/activities", &controllers.AdminController{}, "get:GetSActivities;post:CreateSActivity;delete:DeleteAllSActivities"),
		//beego.NSRouter("/simulation/batch/activities", &controllers.AdminController{}, "post:BatchSActivities"),
		//beego.NSRouter("/simulation/activity/:id", &controllers.AdminController{}, "post:ChangeSActivity;delete:DeleteSActivity"),
		//beego.NSRouter("/simulation/calculation", &controllers.AdminController{}, "get:Calculation"),
		//beego.NSRouter("/simulation/stat", &controllers.AdminController{}, "get:SimulationStat"),
		beego.NSNamespace("admin",
			beego.NSRouter("/change-password", &controllers.AdminController{}, "post:ChangePassword"),
			beego.NSRouter("/set-mobile", &controllers.AdminController{}, "post:SetMobile"),
			beego.NSRouter("/logout", &controllers.AdminController{}, "get:Logout"),
			beego.NSRouter("/users", &controllers.AdminController{}, "get:Users"),
			beego.NSRouter("/user/thaw", &controllers.AdminController{}, "post:Thaw"),
			beego.NSRouter("/real-names", &controllers.AdminController{}, "get:RealNames"),
			beego.NSRouter("/name-auth", &controllers.AdminController{}, "post:NameAuth"),
			beego.NSRouter("/kycs", &controllers.AdminController{}, "get:Kycs"),
			beego.NSRouter("/kyc-auth", &controllers.AdminController{}, "post:KycAuth"),
			beego.NSRouter("/transfers", &controllers.AdminController{}, "get:Transfers"),
			beego.NSRouter("/transfer-auth", &controllers.AdminController{}, "post:TransferAuth"),
			beego.NSRouter("/wallets", &controllers.AdminController{}, "get:Wallets"),
			beego.NSRouter("/addresses", &controllers.AdminController{}, "get:Addresses"),
			beego.NSRouter("/address/pool", &controllers.AdminController{}, "get:AddressPool"),
			beego.NSRouter("/statistics", &controllers.AdminController{}, "get:Statistics"),
			beego.NSRouter("/invites", &controllers.AdminController{}, "get:Invites"),
			beego.NSRouter("/child/amounts", &controllers.AdminController{}, "get:ChildAmounts"),
			beego.NSRouter("/groups", &controllers.AdminController{}, "get:Groups"),
			beego.NSRouter("/fund-changes", &controllers.AdminController{}, "get:FundChanges"),
			beego.NSRouter("/profit-date", &controllers.AdminController{}, "get:ProfitDate"),
			beego.NSRouter("/profit-month", &controllers.AdminController{}, "get:ProfitMonth"),
			beego.NSRouter("/predistribution", &controllers.AdminController{}, "get:Predistribution"),
			beego.NSRouter("/subscriptions", &controllers.AdminController{}, "get:GetSubscriptions"),
			beego.NSRouter("/order/:order/submissions", &controllers.AdminController{}, "get:GetSubmissions"),
			beego.NSRouter("/order", &controllers.AdminController{}, "post:ConfirmOrder"),
			beego.NSNamespace("super",
				beego.NSRouter("/user/:id", &controllers.AdminController{}, "post:ChangeUser"),
				beego.NSRouter("/admins", &controllers.AdminController{}, "get:Admins;post:CreateAdmin"),
				beego.NSRouter("/admin/:id", &controllers.AdminController{}, "post:ChangeAdmin"),
				beego.NSRouter("/currencies", &controllers.AdminController{}, "get:Currencies;post:CreateCurrency"),
				beego.NSRouter("/currency/:id", &controllers.AdminController{}, "post:ChangeCurrency"),
				beego.NSRouter("/groups", &controllers.AdminController{}, "post:CreateGroup"),
				beego.NSRouter("/logs", &controllers.AdminController{}, "get:Logs"),
				beego.NSRouter("/captcha", &controllers.AdminController{}, "get:SendCaptcha"),
				beego.NSRouter("/super-password", &controllers.AdminController{}, "post:SuperPassword"),
				beego.NSRouter("/user/:id/recharge", &controllers.AdminController{}, "post:UserRecharge"),
				beego.NSRouter("/batch-sms", &controllers.AdminController{}, "post:BatchSms"),
				beego.NSRouter("/batch-email", &controllers.AdminController{}, "post:BatchEmail"),
				beego.NSRouter("/captcha-to-user", &controllers.AdminController{}, "get:BackSendCaptchaToUser"),
			),
		),
		//后台管理页面 vue构建
		beego.NSGet("/admin-front/*", func(ctx *context.Context) {
			length := len("/kc/admin-front")
			name := ctx.Request.URL.Path[length:]
			http.ServeFile(ctx.ResponseWriter, ctx.Request, "front/dist"+name)
		}),
	)
	beego.AddNamespace(ns)

	if beego.BConfig.RunMode == "dev" {
		//访问上传文件 生产环境由nginx代理
		beego.SetStaticPath("uphp/gcexserver", utils.MediaPath)
	}
}
