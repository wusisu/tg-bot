package main

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// FileInfo is to be download info
type FileInfo struct {
	FileID   string
	FileName string
	FileSize int
	MimeType string
}

// FromPhotoSize read data from PhotoSize
func (f FileInfo) FromPhotoSize(ph tgbotapi.PhotoSize) FileInfo {
	f.FileID = ph.FileID
	f.FileSize = ph.FileSize
	return f
}

// FromDocument read data from Document
func (f FileInfo) FromDocument(doc tgbotapi.Document) FileInfo {
	f.FileID = doc.FileID
	f.FileSize = doc.FileSize
	return f
}
