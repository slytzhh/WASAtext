package components

type Username struct {
	Username string `json:"username"`
}

type GroupName struct {
	GroupName string `json:"groupname"`
}

type UsernameList struct {
	UsernameList []string `json:"usernamelist"`
}

type UserIdPhoto struct {
	UserId int    `json:"userid"`
	Photo  string `json:"photo"`
}

type UserList struct {
	UserList []User `json:"userlist"`
}

type User struct {
	UserId     int    `json:"userid"`
	Username   string `json:"username"`
	Photo      string `json:"photo"`
	LastAccess string `json:"lastaccess"`
}

type ChatMessId struct {
	ChatId    int `json:"chatid"`
	MessageId int `json:"messageid"`
}

type Chat struct {
	ChatId       int       `json:"chatid"`
	GroupName    string    `json:"groupname"`
	GroupPhoto   string    `json:"groupphoto"`
	TimeCreated  string    `json:"timecreated"`
	IsGroup      bool      `json:"isgroup"`
	UsernameList []string  `json:"usernamelist"`
	MessageList  []Message `json:"messagelist"`
}

type ChatCreation struct {
	UsernameList []string      `json:"usernamelist"`
	GroupName    string        `json:"groupname"`
	GroupPhoto   string        `json:"groupphoto"`
	FirstMessage MessageToSend `json:"firstmessage"`
	ForwardedId  int           `json:"forwardedid"`
}

type ChatPreview struct {
	ChatId      int            `json:"chatid"`
	GroupName   string         `json:"groupname"`
	GroupPhoto  string         `json:"groupphoto"`
	TimeCreated string         `json:"timecreated"`
	LastMessage MessagePreview `json:"lastmessage"`
}

type MessageId struct {
	MessageId int `json:"messageid"`
}

type MessageToSend struct {
	ReplyId int    `json:"replyid"`
	Text    string `json:"text"`
	Photo   string `json:"photo"`
}

type MessagePreview struct {
	MessageId     int    `json:"messageid"`
	ChatId        int    `json:"chatid"`
	UserId        int    `json:"userid"`
	Username      string `json:"username"`
	Text          string `json:"text"`
	Photo         string `json:"photo"`
	TimeStamp     string `json:"timestamp"`
	IsAllReceived bool   `json:"isallreceived"`
	IsAllRead     bool   `json:"isallread"`
}

type Message struct {
	MessageId     int            `json:"messageid"`
	ChatId        int            `json:"chatid"`
	UserId        int            `json:"userid"`
	Username      string         `json:"username"`
	Text          string         `json:"text"`
	Photo         string         `json:"photo"`
	IsForwarded   bool           `json:"isforwarded"`
	TimeStamp     string         `json:"timestamp"`
	IsAllReceived bool           `json:"isallreceived"`
	IsAllRead     bool           `json:"isallread"`
	CommentList   []Comment      `json:"commentlist"`
	ReplyMessage  MessagePreview `json:"replymessage"`
}

type Comment struct {
	MessageId int    `json:"messageid"`
	UserId    int    `json:"userid"`
	Emoji     string `json:"emoji"`
	Username  string `json:"username"`
}

type Photo struct {
	Photo string `json:"photo"`
}

type Emoji struct {
	Emoji string `json:"emoji"`
}
