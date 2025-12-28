package generate

import (
	"github.com/algorave/server/algorave/anonsessions"
	"github.com/algorave/server/internal/agent"
	"github.com/algorave/server/internal/auth"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, agentClient *agent.Agent, strudelRepo StrudelGetter, sessionMgr *anonsessions.Manager) {
	router.POST("/generate", auth.OptionalAuthMiddleware(), Handler(agentClient, strudelRepo, sessionMgr))
}
