package main

import (
	"fmt"
	"github.com/mqu/go-notify"
	"os"
	"strings"
)

func Index(vs []string, t string) int {
	for i, v := range vs {
		if v == t {
			return i
		}
	}
	return -1
}
func Include(vs []string, t string) bool {
	return Index(vs, t) >= 0
}

const BreakTag = "!break!"

var ValidTags = []string{"status", "url", "file", "artist", "album", "discnumber", "tracknumber", "title", "date", "duration", BreakTag}

type Metadata map[string]string

func Parse(cmusData string) Metadata {
	result := make(Metadata)
	splitedData := strings.Split(cmusData, " ")
	splitedData = append(splitedData, BreakTag)
	var lastFound string
	var valueCollector []string
	for _, part := range splitedData {
		if Include(ValidTags, part) {
			if len(valueCollector) > 0 {
				joinedValue := strings.Join(valueCollector, " ")
				result[lastFound] = joinedValue
				valueCollector = nil
			}
			lastFound = part
			continue
		}
		valueCollector = append(valueCollector, part)
	}
	return result
}

func FormatMessageBody(m Metadata) string {
	notificationBody := ""
	if title, ok := m["title"]; ok == true {
		notificationBody = title
	} else {
		notificationBody = "Unknown"
	}

	if artist, ok := m["artist"]; ok == true {
		notificationBody = fmt.Sprintf("%s by %s", notificationBody, artist)

		if album, ok := m["album"]; ok == true {
			notificationBody = fmt.Sprintf("%s, %s", notificationBody, album)
		}
	}
	return notificationBody
}

type Notifier func(string, string) error

func NotifySend(summary, content string) error {
	notify.Init("cmustify")
	notification := notify.NotificationNew(summary, content, "")
	return notify.NotificationShow(notification)
}

func HandleData(notifier Notifier, cmusData string) {
	metaData := Parse(cmusData)
	notificationBody := FormatMessageBody(metaData)
	notifier("Cmustify - Current song", notificationBody)
}

func printUsage() {
	fmt.Println("You must set cmus to call this script as notifier")
}

func main() {
	if os.Args[1] == "-h" {
		printUsage()
		return
	}

	args := os.Args[1:]
	cmusData := strings.Join(args, " ")

	HandleData(NotifySend, cmusData)
}
