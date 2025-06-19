package natsutils

import "strconv"

func GenerateConsumerName(accountId uint64, userId int32, connId string) string {
	return "acc_" + strconv.FormatUint(accountId, 10) + "_usr_" + strconv.FormatInt(int64(userId), 10) + "_sess_" + connId
}
