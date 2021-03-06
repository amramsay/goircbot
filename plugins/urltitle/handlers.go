// Package urltitle implements a plugin to watch web URLs, fetch and display title.
package urltitle

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/StalkR/goircbot/bot"
	"github.com/StalkR/goircbot/lib/url"
	"github.com/fluffle/goirc/client"
)

var (
	linkRE    = regexp.MustCompile(`(?:^|\s)(https?://[^#\s]+)`)
	silenceRE = regexp.MustCompile(`(^|\s)tg(\)|\s|$)`) // Line ignored if matched.
)

func watchLine(b bot.Bot, line *client.Line, ignore map[string]bool) {
	target := line.Args[0]
	if !strings.HasPrefix(target, "#") {
		return
	}
	if _, ignore := ignore[line.Nick]; ignore {
		return
	}
	text := line.Args[1]
	if silenceRE.MatchString(text) {
		return
	}
	match := linkRE.FindStringSubmatch(text)
	if match == nil {
		return
	}
	link := match[1]
	if len(link) > 200 {
		return
	}
	title, err := url.Title(link)
	if err != nil {
		log.Println("urltitle:", err)
		return
	}
	if len(title) > 200 {
		title = title[:200]
	}
	b.Privmsg(target, fmt.Sprintf("%s :: %s", link, title))
}

// Register registers the plugin with a bot.
func Register(b bot.Bot, ignore []string) {
	ignoremap := make(map[string]bool)
	for _, nick := range ignore {
		ignoremap[nick] = true
	}

	b.Conn().HandleFunc("privmsg",
		func(conn *client.Conn, line *client.Line) { watchLine(b, line, ignoremap) })
}
