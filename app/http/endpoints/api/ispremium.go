package api

import (
	"github.com/TicketsBot/GoPanel/botcontext"
	"github.com/TicketsBot/GoPanel/rpc"
	"github.com/TicketsBot/common/premium"
	"github.com/gin-gonic/gin"
)

func IsPremium(c *gin.Context) {
	if c.Query("key") != config.Conf.bot.premium-lookup-proxy-key {
		c.JSON(403, gin.H{
			"error": "Invalid secret key",
		})

		return
	}

	id := c.Query("id")

	for user, pledge := range app.Pledges {
		if user.Attributes.DiscordId == id {
			for index, rewardId := range config.Conf.Patreon.Tiers {
				if strconv.Itoa(rewardId) == pledge.Relationships.Reward.Data.ID {
					if pledge.Attributes.IsPaused == nil && !pledge.Attributes.DeclinedSince.Valid {
						c.JSON(200, gin.H{
							"premium": true,
							"tier": index,
						})
						return
					}
				}
			}
		}
	}

	c.JSON(200, gin.H{
		"premium": true,
 		"tier": 2,
	})
}
