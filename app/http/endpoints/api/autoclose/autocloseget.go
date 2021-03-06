package api

import (
	"github.com/TicketsBot/GoPanel/database"
	"github.com/gin-gonic/gin"
)

func GetAutoClose(ctx *gin.Context) {
	guildId := ctx.Keys["guildid"].(uint64)

	settings, err := database.Client.AutoClose.Get(guildId); if err != nil {
		ctx.AbortWithStatusJSON(500, gin.H{
			"success": false,
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, settings)
}
