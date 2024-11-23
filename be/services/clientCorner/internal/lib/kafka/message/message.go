package message

type Message struct {
	UserID      string `json:"userID"`       // id пользователя порадившего событие
	TypeMessage string `json:"type_message"` // Тип уведомления [change_password / change_email]
	Msg         string `json:"message"`      // Сообщение уведомления
}
