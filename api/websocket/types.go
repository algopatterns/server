package websocket

type ConnectParams struct {
	SessionID   string `form:"session_id" binding:"required"`
	Token       string `form:"token"`        // JWT token for authenticated users
	InviteToken string `form:"invite"`       // Invite token for joining sessions
	DisplayName string `form:"display_name"` // Optional display name for anonymous users
}
