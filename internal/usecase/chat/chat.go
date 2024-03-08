package chat

type ChatPool struct {
	Pool map[string]*Pool
}

func NewChatPool() *ChatPool {
	return &ChatPool{
		Pool: make(map[string]*Pool),
	}
}