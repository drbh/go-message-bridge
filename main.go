package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nlopes/slack"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"time"

	"encoding/json"
	"io/ioutil"

	"github.com/boltdb/bolt"
	// "github.com/davecgh/go-spew/spew"
	"hash/fnv"
)

var USER_ID string
var BOT_ID string
var BOT_NAME string
var BOT_TOKEN string
var USER_TOKEN string

var lastText string
var DATE_OFFSET int64
var lastGUID string
var seenMessageGuids = []string{}
var seenMessageGuid = map[string]bool{}

type configOptions struct {
	BOT_NAME   string `json:"bot_name"`
	BOT_TOKEN  string `json:"bot_token"`
	USER_TOKEN string `json:"user_token"`
}

func getConfig() {
	plan, _ := ioutil.ReadFile("config.json")
	var data configOptions
	json.Unmarshal([]byte(plan), &data)
	BOT_NAME = data.BOT_NAME
	BOT_TOKEN = data.BOT_TOKEN
	USER_TOKEN = data.USER_TOKEN
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func makeAppleTimestamp() int {
	DATE_OFFSET = 978307200 + 1
	return int((((time.Now().UnixNano()) / int64(time.Millisecond)) / 1000) - DATE_OFFSET)
}

func handleToChannel(handle string) string {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	var v string
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("TEST"))
		if err != nil {
			return err
		}
		v = string(b.Get([]byte(handle)))
		return nil
	})
	defer db.Close()
	return v
}
func setHandleToChannel(handle []byte, channel []byte) {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("TEST"))
		if err != nil {
			return err
		}
		b.Put(handle, channel)
		return nil
	})
	defer db.Close()
}

func getMessages() {
	args := "?mode=ro&_mutex=no&_journal_mode=WAL&_query_only=1&_synchronous=2"

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(usr.HomeDir)

	connectionString := usr.HomeDir + "/Library/Messages/chat.db" + args
	database, err := sql.Open("sqlite3", connectionString)
	defer database.Close()
	// fmt.Println(database)
	if err != nil {
		log.Fatal("Connection Failed ", err)
	}
	// database.Exec("PRAGMA _ignore_check_constraints = 1")
	defer database.Close()

	database.SetMaxOpenConns(1)
	// database.Exec("PRAGMA journal_mode=WAL;")

	rows, qerr := database.Query(`
        SELECT
            guid as id,
            chat_identifier as recipientId,
            service_name as serviceName,
            room_name as roomName,
            display_name as displayName
        FROM chat
        JOIN chat_handle_join ON chat_handle_join.chat_id = chat.ROWID
        JOIN handle ON handle.ROWID = chat_handle_join.handle_id
        ORDER BY handle.rowid DESC
        LIMIT 10;
        `)
	// rows, qerr := database.Query("SELECT name FROM sqlite_master WHERE type='table';")
	defer rows.Close()
	if qerr != nil {
		log.Fatal("Query Failed ", qerr)
	}

	// fmt.Println(rows)

	for rows.Next() {
		var id string
		var recipientId string
		var serviceName string
		var roomName string
		var displayName string
		rows.Scan(
			&id,
			&recipientId,
			&serviceName,
			&roomName,
			&displayName,
		)
		fmt.Println(
			id,
			recipientId,
			serviceName,
			roomName,
			displayName,
		)
	}
}

func runPoller() {
	ticker := time.NewTicker(1000 * time.Millisecond)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				poll()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func poll() {
	args := "?mode=ro&_mutex=no&_journal_mode=WAL&_query_only=1&_synchronous=2"
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(usr.HomeDir)

	connectionString := usr.HomeDir + "/Library/Messages/chat.db" + args
	database, err := sql.Open("sqlite3", connectionString)
	defer database.Close()
	// fmt.Println(database)
	if err != nil {
		log.Fatal("Connection Failed ", err)
	}
	// database.Exec("PRAGMA _ignore_check_constraints = 1")
	defer database.Close()

	database.SetMaxOpenConns(1)
	// database.Exec("PRAGMA journal_mode=WAL;")
	latest_ := makeAppleTimestamp()
	latest := strconv.Itoa(latest_ - 1) // time.Now().UTC().String()

	// fmt.Println(latest, latest_)
	rows, qerr := database.Query(`
			 SELECT
			    guid,
			    id as handle,
			    handle_id,
			    text,
			    date,
			    date_read,
			    is_from_me,
			    cache_roomnames
			FROM message
			LEFT OUTER JOIN handle ON message.handle_id = handle.ROWID
			WHERE date >= ` + latest + `
        `)
	// rows, qerr := database.Query("SELECT name FROM sqlite_master WHERE type='table';")
	defer rows.Close()
	if qerr != nil {
		log.Fatal("Query Failed ", qerr)
	}

	var guid string
	var handle string
	var handle_id string
	var text string
	var date string
	var date_read string
	var is_from_me string
	var cache_roomnames string

	for rows.Next() {
		rows.Scan(
			&guid,
			&handle,
			&handle_id,
			&text,
			&date,
			&date_read,
			&is_from_me,
			&cache_roomnames,
		)
		if seenMessageGuid[guid] {
		} else {
			theirChannel := handleToChannel(handle_id)
			var relative_channel string
			if theirChannel != "" {
				relative_channel = theirChannel
			} else {
				// make a new channel
				userID := handle_id
				channelID, err := user_API.CreateChannel(userID)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
				// take the channel name and handle and add to db
				fmt.Println(channelID.ID)

				// auto add the bot to the channel!
				channelID2, errr := user_API.InviteUserToChannel(channelID.ID, BOT_ID)
				if errr != nil {
					fmt.Printf("%s\n", errr)
				}

				fmt.Println(channelID2)
				fmt.Println(handle_id)
				relative_channel = channelID.ID
				setHandleToChannel([]byte(handle_id), []byte(channelID.ID))
				setHandleToChannel([]byte(channelID.ID+"-handle"), []byte(handle))
				setHandleToChannel([]byte(channelID.ID), []byte(handle_id))
			}

			fmt.Println(relative_channel)
			handleToChannel(handle_id)
			// THE PLACE WE GET THE SINGLE GOOD MESSAGE
			// de dep logic here
			if is_from_me == "1" {
				fmt.Println(hash(text))
				fmt.Println(hash(lastText))
				if hash(text) == hash(lastText) {
					fmt.Println("DONT SEND MESSAGE")
					lastText = "reset"

				} else {
					fmt.Println("SEND MESSAGE")
					global_rtm_2.SendMessage(global_rtm.NewOutgoingMessage(text, relative_channel))
				}

			} else {
				global_rtm.SendMessage(global_rtm_2.NewOutgoingMessage(text, relative_channel))
				// convo is with handle_id
			}

			fmt.Println(guid, handle_id, handle, is_from_me, "\t\t\t\t\t", text, len(seenMessageGuid), len(seenMessageGuids))

			seenMessageGuid[guid] = true
			seenMessageGuids = append(seenMessageGuids, guid)
		}

		if len(seenMessageGuids) > 10 {
			key := seenMessageGuids[0]
			delete(seenMessageGuid, key)
			seenMessageGuids = append(seenMessageGuids[:0], seenMessageGuids[1:]...)
		}
	}
}

var global_rtm *slack.RTM
var global_rtm_2 *slack.RTM
var user_API *slack.Client

func main() {

	getConfig()

	runPoller()

	api := slack.New(
		BOT_TOKEN,
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	global_rtm = rtm
	go rtm.ManageConnection()

	api2 := slack.New(
		USER_TOKEN,
		slack.OptionDebug(false),
		slack.OptionLog(log.New(os.Stdout, "slack-bot-2: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm2 := api2.NewRTM()
	user_API = api2

	users, _ := user_API.GetUsers()

	// lets fetch the users id and bots id on startup
	for i := 0; i < len(users); i++ {
		// spew.Dump(users[i])
		if users[i].IsPrimaryOwner && users[i].IsAdmin {
			USER_ID = users[i].ID
			fmt.Println("USR: ", USER_ID)
		}

		if users[i].Name == BOT_NAME && users[i].IsBot {
			BOT_ID = users[i].ID
			fmt.Println("BOT: ", BOT_ID)
		}

		// fmt.Println(users[i].ID, users[i].Name, users[i].IsBot, users[i].IsAdmin, users[i].IsPrimaryOwner)
	}

	global_rtm_2 = rtm2
	go rtm2.ManageConnection()

	for msg := range rtm2.IncomingEvents {
		fmt.Println("-")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)

		case *slack.MessageEvent:
			aguid := handleToChannel(ev.Channel)
			fmt.Println(aguid)
			handle := handleToChannel(ev.Channel + "-handle")
			fmt.Println(handle)
			if ev.Type == "message" {
				if ev.User == USER_ID {
					fmt.Println("Should send iMessage to")
					lastText = ev.Text
					command := "osascript sendMessage.applescript " + handle + " \"" + ev.Text + "\""
					fmt.Println(command)
					out, _ := exec.Command("sh", "-c", command).Output()

					fmt.Println(out)
				}
			}
			fmt.Printf("Message: %v\n", ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}

	for msg := range rtm.IncomingEvents {
		// fmt.Print("Event Received: ")
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			// Ignore hello

		case *slack.ConnectedEvent:
			fmt.Println("Infos:", ev.Info)
			fmt.Println("Connection counter:", ev.ConnectionCount)
			// Replace C2147483705 with your Channel ID
			// rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "DC7EXT3RT"))

		case *slack.MessageEvent:
			fmt.Printf("Message: %v\n", ev)

		case *slack.PresenceChangeEvent:
			fmt.Printf("Presence Change: %v\n", ev)

		case *slack.LatencyReport:
			fmt.Printf("Current latency: %v\n", ev.Value)

		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			fmt.Printf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
	time.Sleep(5 * time.Hour)

}
