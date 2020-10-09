package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.q123123.net/ligmar/boot"
	"gitlab.q123123.net/ligmar/console-back/models"

	"github.com/gin-gonic/gin"
)

func (h *Handler) userIdentity(c *gin.Context) (user models.User, err error) {
	var userID int
	userID, err = handlerMiddlewareTelegramAuth(c.Request.URL.Query())
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid params")
		return
	}

	user, err = h.services.Authentication.Check(userID)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "user not found")
	}

	return
}

const (
	hash      = "hash"
	authDate  = "authDate"
	firstName = "first_name"
	id        = "id"
	lastName  = "last_name"
	photoUrl  = "photo_url"
	username  = "username"
)

func handlerMiddlewareTelegramAuth(data url.Values) (userID int, err error) {
	var dataCheckStr []string

	if v := data.Get(hash); v == "" {
		err = errors.New("try to auth without params: " + hash)
		return
	}

	if v := data.Get(authDate); v == "" {
		err = errors.New("try to auth without params: " + authDate)
		return
	}
	// 60 sec for auth
	date, err := strconv.Atoi(authDate)
	if err != nil || date+60 < int(time.Now().Unix()) {
		err = errors.New("try to auth with bad params: " + authDate)
		return
	}
	dataCheckStr = append(dataCheckStr, authDate+"="+data.Get(authDate))

	if v := data.Get(firstName); v == "" {
		err = errors.New("try to auth without params: " + firstName)
		return
	}
	dataCheckStr = append(dataCheckStr, firstName+"="+data.Get(firstName))

	if v := data.Get(id); v == "" {
		err = errors.New("try to auth without params: " + id)
		return
	}
	userID, err = strconv.Atoi(data.Get(id))
	if err != nil {
		err = errors.New("try to auth with bad params: " + id)
		return
	}
	dataCheckStr = append(dataCheckStr, id+"="+data.Get(id))

	if v := data.Get(lastName); v == "" {
		err = errors.New("try to auth without params: " + lastName)
		return
	}
	dataCheckStr = append(dataCheckStr, lastName+"="+data.Get(lastName))

	if v := data.Get(photoUrl); v == "" {
		err = errors.New("try to auth without params: " + photoUrl)
		return
	}
	dataCheckStr = append(dataCheckStr, photoUrl+"="+data.Get(photoUrl))

	if v := data.Get(username); v != "" {
		err = errors.New("try to auth without params: " + username)
		return
	}
	dataCheckStr = append(dataCheckStr, username+"="+data.Get(username))

	sort.Strings(dataCheckStr)
	imploded := strings.Join(dataCheckStr, "\n")
	sha256hash := sha256.New()
	if _, err = io.WriteString(sha256hash, boot.TelegramSupportBotAPI); err != nil {
		return
	}
	hmacHash := hmac.New(sha256.New, sha256hash.Sum(nil))
	if _, err = io.WriteString(hmacHash, imploded); err != nil {
		return
	}
	ss := hex.EncodeToString(hmacHash.Sum(nil))

	if hash != ss {
		boot.Log.Warn(data)
		err = errors.New("hash not equal")
		return
	}
	return
}
