package audit

import (
	iris "github.com/kataras/iris/v12"
	"github.com/markbates/goth"
	cnf "github.com/rzrbld/adminio-api/config"
	log "log"
)

func DefaultAuditLog(user goth.User, ctx iris.Context) {
	ctx.ViewData("", user)
	if cnf.AuditLogEnable {
		log.Print("userNickName: ", user.NickName, "; userID: ", user.UserID, "; method:", ctx.RouteName())
	}
}
