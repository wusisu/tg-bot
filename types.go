package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// FileInfo is to be download info
type FileInfo struct {
	FileID   string
	FileName string
	FileSize int
	Width    int
	Height   int
	MimeType string
	Duration int
}

// FromPhotoSize read data from PhotoSize
func (f FileInfo) FromPhotoSize(ph tgbotapi.PhotoSize) FileInfo {
	f.FileID = ph.FileID
	f.FileSize = ph.FileSize
	f.Width = ph.Width
	f.Height = ph.Height
	return f
}

// FromDocument read data from Document
func (f FileInfo) FromDocument(doc *tgbotapi.Document) FileInfo {
	f.FileID = doc.FileID
	f.FileSize = doc.FileSize
	f.FileName = doc.FileName
	f.MimeType = doc.MimeType
	return f
}

// FromVideo read data from Video
func (f FileInfo) FromVideo(v *tgbotapi.Video) FileInfo {
	f.FileID = v.FileID
	f.FileSize = v.FileSize
	f.Width = v.Width
	f.Height = v.Height
	f.MimeType = v.MimeType
	f.Duration = v.Duration
	return f
}
