package manage

import (
	"github.com/TicketsBot/GoPanel/app/http/template"
	"github.com/TicketsBot/GoPanel/config"
	"github.com/TicketsBot/GoPanel/database/table"
	"github.com/TicketsBot/GoPanel/utils"
	"github.com/TicketsBot/GoPanel/utils/discord/objects"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

func TicketListHandler(ctx *gin.Context) {
	store := sessions.Default(ctx)
	if store == nil {
		return
	}
	defer store.Save()

	if utils.IsLoggedIn(store) {
		userIdStr := store.Get("userid").(string)
		userId, err := utils.GetUserId(store)
		if err != nil {
			ctx.String(500, err.Error())
			return
		}

		// Verify the guild exists
		guildIdStr := ctx.Param("id")
		guildId, err := strconv.ParseInt(guildIdStr, 10, 64)
		if err != nil {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 404 Page
			return
		}

		// Get object for selected guild
		var guild objects.Guild
		for _, g := range table.GetGuilds(userIdStr) {
			if g.Id == guildIdStr {
				guild = g
				break
			}
		}

		// Verify the user has permissions to be here
		if !guild.Owner && !table.IsAdmin(guildId, userId) {
			ctx.Redirect(302, config.Conf.Server.BaseUrl) // TODO: 403 Page
			return
		}

		tickets := table.GetOpenTickets(guildId)

		var toFetch []int64
		for _, ticket := range tickets {
			toFetch = append(toFetch, ticket.Owner)

			for _, idStr := range strings.Split(ticket.Members, ",") {
				if memberId, err := strconv.ParseInt(idStr, 10, 64); err == nil {
					toFetch = append(toFetch, memberId)
				}
			}
		}

		nodes := make(map[int64]table.UsernameNode)
		for _, node := range table.GetUserNodes(toFetch) {
			nodes[node.Id] = node
		}

		var ticketsFormatted []map[string]interface{}

		for _, ticket := range tickets {
			var membersFormatted []map[string]interface{}
			for index, memberIdStr := range strings.Split(ticket.Members, ",") {
				if memberId, err := strconv.ParseInt(memberIdStr, 10, 64); err == nil {
					if memberId != 0 {
						var separator string
						if index != len(strings.Split(ticket.Members, ",")) - 1 {
							separator = ", "
						}

						membersFormatted = append(membersFormatted, map[string]interface{}{
							"username": utils.Base64Decode(nodes[memberId].Name),
							"discrim": nodes[memberId].Discriminator,
							"sep": separator,
						})
					}
				}
			}

			ticketsFormatted = append(ticketsFormatted, map[string]interface{}{
				"uuid": ticket.Uuid,
				"ticketId": ticket.TicketId,
				"username": utils.Base64Decode(nodes[ticket.Owner].Name),
				"discrim": nodes[ticket.Owner].Discriminator,
				"members": membersFormatted,
			})
		}

		utils.Respond(ctx, template.TemplateTicketList.Render(map[string]interface{}{
			"name":    store.Get("name").(string),
			"guildId": guildIdStr,
			"csrf": store.Get("csrf").(string),
			"avatar": store.Get("avatar").(string),
			"baseUrl": config.Conf.Server.BaseUrl,
			"tickets": ticketsFormatted,
		}))
	}
}