package agent

import (
	"github.com/gin-gonic/gin"

	"github.com/algrv/server/algorave/strudels"
	agentcore "github.com/algrv/server/internal/agent"
	"github.com/algrv/server/internal/attribution"
	"github.com/algrv/server/internal/buffer"
	"github.com/algrv/server/internal/llm"
)

func RegisterRoutes(router *gin.RouterGroup, agentClient *agentcore.Agent, platformLLM llm.LLM, strudelRepo *strudels.Repository, attrService *attribution.Service, sessionBuffer *buffer.SessionBuffer) {
	agentGroup := router.Group("/agent")
	{
		agentGroup.POST("/generate", GenerateHandler(agentClient, platformLLM, strudelRepo, attrService, sessionBuffer))
	}
}
