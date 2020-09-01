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

func (app *App) readUpdates() (err error) {
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

			updateJ, _ := json.Marshal(update)
			log.Debugf("Receive Update %s", updateJ)

			err := app.readMessage(update.Message)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}
	return
}

func (app *App) readMessage(msg *tgbotapi.Message) (err error) {
	err = app.readPhotos(msg)
	if err == nil {
		err = app.readDocument(msg)
	}
	if err == nil {
		err = app.readVideo(msg)
	}
	return
}

func (app *App) readDocument(msg *tgbotapi.Message) (err error) {
	doc := msg.Document
	if doc == nil {
		return nil
	}
	log.Debug("read Document")
	return app.saveFile(new(FileInfo).FromDocument(msg.Document))
}

func (app *App) readPhotos(msg *tgbotapi.Message) (err error) {
	phs := msg.Photo
	if phs == nil || len(*phs) == 0 {
		return nil
	}
	ph := (*phs)[len(*phs)-1]
	// for _, ph := range *phs {

	log.Debugf("read Photo [%d] %s %v", ph.FileSize, ph.FileID, ph)

	// }
	return app.saveFile(new(FileInfo).FromPhotoSize(ph))
}

func (app *App) readVideo(msg *tgbotapi.Message) (err error) {
	v := msg.Video
	if v == nil {
		return nil
	}
	log.Debug("read Video")
	return app.saveFile(new(FileInfo).FromVideo(v))
}

func (app *App) saveFile(fi FileInfo) (err error) {
	has, err := app.db.Where("file_i_d = ?", fi.FileID).Exist(&File{})
	if err != nil {
		return
	}
	if has {
		log.Infof("FileID exists %s", fi.FileID)
		return
	}
	tgFile, err := app.bot.GetFile(tgbotapi.FileConfig{FileID: fi.FileID})
	if err != nil {
		log.Errorf("failed to downlaod [%s]", fi.FileID)
		return
	}
	url := tgFile.Link(viper.GetString("BotToken"))
	log.Debugf("Download File %s", url)
	has, err = app.db.Where("download_u_r_l = ?", url).Exist(&File{})
	if err != nil {
		return
	}
	if has {
		log.Infof("DownloadURL exists %s", url)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	h := md5.New()
	h.Write(b)
	md5sum := h.Sum(nil)
	md5 := hex.EncodeToString(md5sum)
	has, err = app.db.Where("md5 = ?", md5).Exist(&File{})
	if err != nil {
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
		return
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
	f := new(File).FromFileInfo(fi)
	f.Md5 = md5
	f.OutputName = outputName
	f.DownloadURL = url
	affect, err := app.db.Insert(&f)
	if err != nil {
		return
	}
	log.Debugf("%d item insert to db", affect)
	return
}
