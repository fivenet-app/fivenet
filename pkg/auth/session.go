package auth

import (
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func GetSessionInfo(c *gin.Context) (*SessionInfo, error) {
	return getSessionInfo(sessions.DefaultMany(c, SessionName))
}

func getSessionInfo(sess sessions.Session) (*SessionInfo, error) {
	d := sess.Get(keyName)
	if d != nil {
		info, ok := d.(*SessionInfo)
		if !ok {
			// If the data isn't castable, clear the session
			sess.Clear()
			return nil, errors.New("failed to parse session info")
		}
		return info, nil
	}

	return nil, nil
}

func SaveSessionInfo(c *gin.Context, in *SessionInfo) error {
	sess := sessions.DefaultMany(c, SessionName)
	sess.Set(keyName, in)

	return sess.Save()
}
