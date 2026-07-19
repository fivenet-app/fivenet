package natsutils

import "strconv"

func GenerateConsumerName(accountId int64, userId int32, connId string) string {
	return "a" + strconv.FormatInt(
		accountId,
		10,
	) + "_u" + strconv.FormatInt(
		int64(userId),
		10,
	) + "_s" + connId
}
