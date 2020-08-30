package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"xorm.io/xorm"
)

// App is main app
type App struct {
	db  *xorm.Engine
	bot *tgbotapi.BotAPI
}

// NewApp will create a new App, which you should connect by yourself
func NewApp() *App {
	return &App{}
}

func (app *App) readUpdates() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := app.bot.GetUpdatesChan(u)

	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message != nil { // ignore any non-Message Updates
			log.Debugf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			// msg.ReplyToMessageID = update.Message.MessageID

			// app.bot.Send(msg)

			phs := update.Message.Photo
			if phs == nil || len(*phs) == 0 {
				updateJ, _ := json.Marshal(update)
				log.Debugf("no photos in %s", updateJ)
				continue
			}
			ph := (*phs)[len(*phs)-1]
			// for _, ph := range *phs {

			log.Debugf("[%d] %s %v", ph.FileSize, ph.FileID, ph)

			// }

			f, err := app.bot.GetFile(tgbotapi.FileConfig{FileID: ph.FileID})
			if err != nil {
				log.Debugf("failed to downlaod [%s]", ph.FileID)
				continue
			}
			url := f.Link(viper.GetString("BotToken"))
			log.Debugf("Download image %s", url)
			app.saveFile(url, ph)

		}
	}
}

func (app *App) saveFile(url string, ph tgbotapi.PhotoSize) {
	has, err := app.db.Where("file_i_d = ?", ph.FileID).Exist(&File{})
	if err != nil {
		log.Error(err)
		return
	}
	if has {
		log.Infof("FileID exists %s", ph.FileID)
		return
	}
	has, err = app.db.Where("download_u_r_l = ?", url).Exist(&File{})
	if err != nil {
		log.Error(err)
		return
	}
	if has {
		log.Infof("DownloadURL exists %s", url)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Panicf("fetch: %v\n", err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Panicf("fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}
	h := md5.New()
	h.Write(b)
	md5sum := h.Sum(nil)
	md5 := hex.EncodeToString(md5sum)
	has, err = app.db.Where("md5 = ?", md5).Exist(&File{})
	if err != nil {
		log.Error(err)
		return
	}
	if has {
		log.Infof("md5 exists %s", md5)
		return
	}

	suffix := filepath.Ext(url)
	nano := fmt.Sprintf("%d", time.Now().UnixNano())
	outputName := nano + suffix
	fn := path.Join(viper.GetString("DownloadDir"), outputName)
	fo, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	// w := bufio.NewWriter(fo)
	// w.Write(b)
	// w.Flush()
	fo.Write(b)
	f := File{}
	f.Md5 = md5
	f.FileID = ph.FileID
	f.FileSize = ph.FileSize
	f.OutputName = outputName
	f.DownloadURL = url
	affect, err := app.db.Insert(&f)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debugf("%d item insert to db", affect)
}
