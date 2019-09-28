package bot

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

const (
	twitchIRCEndpoint = "irc.chat.twitch.tv:80"
)

// Bot holds the state for a Twitch bot talking to a single twitch chatroom
type Bot struct {
	conn        net.Conn
	cmdq        chan string
	accessToken string
	name        string
	channelName string
	channelID   string
	chatroomID  string
}

// NewBot creates a new Bot that is configured and connected but not logged in.
func NewBot(accessToken, name, channelName, channelID, chatroomID string) (*Bot, error) {
	conn, err := net.Dial("tcp", twitchIRCEndpoint)
	if err != nil {
		return nil, err
	}

	bot := &Bot{
		conn:        conn,
		cmdq:        make(chan string, 10),
		accessToken: accessToken,
		name:        name,
		channelName: channelName,
		channelID:   channelID,
		chatroomID:  chatroomID,
	}

	go bot.reader()
	go bot.writer()

	return bot, nil
}

func (b *Bot) reader() {
	var cmdq chan<- string = b.cmdq
	s := bufio.NewScanner(b.conn)

	for s.Scan() {
		line := s.Text()
	}
}

func (b *Bot) writer() {
	var cmdq <-chan string = b.cmdq
	w := bufio.NewWriter(b.conn)
	var err error

	for cmd := range cmdq {
		w.WriteString(cmd)
		if err = w.Flush(); err != nil {
			break
		}
	}

	log.Println("Error while writing command: %v", err)
}

// Login logs in to the IRC chat using the bot's name and access token
func (b *Bot) Login() {
	b.cmdq <- "CAP REQ :twitch.tv/tags twitch.tv/commands twitch.tv/membership"
	b.cmdq <- fmt.Sprintf("PASS oauth:%s", b.accessToken)
	b.cmdq <- fmt.Sprintf("NICK %s", b.name)
	b.cmdq <- fmt.Sprintf("USER %s 8 * :%s", b.name, b.name)

	if b.chatroomID == "" {
		b.cmdq <- fmt.Sprintf("JOIN #%s", b.channelName)
	} else {
		b.cmdq <- fmt.Sprintf("JOIN #chatrooms:%s:%s", b.channelID, b.chatroomID)
	}
}
