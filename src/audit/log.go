package audit

import (
	iris "github.com/kataras/iris/v12"
	"github.com/markbates/goth"
	cnf "github.com/rzrbld/adminio-api/config"
	log "github.com/sirupsen/logrus"
)

func DefaultAuditLog(user goth.User, ctx iris.Context) {
	ctx.ViewData("", user)
	if cnf.AuditLogEnable {
		log.Infoln("userNickName: ", user.NickName, "; userID: ", user.UserID, "; method:", ctx.RouteName())
	}
}
