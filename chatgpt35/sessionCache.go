package gpt35

import (
	"encoding/json"
	"github.com/patrickmn/go-cache"
	"time"
)

type SessionService struct {
	cache *cache.Cache
}

type SessionMeta struct {
	Msg []Message `json:"msg,omitempty"`
}

type CacheInterface interface {
	GetMsg(sessionId string) []Message
	SetMsg(sessionId string, msg []Message)
	Clear(sessionId string)
}

var sessionServices *SessionService

func getLength(strPool []Message) int {
	var total int
	for _, v := range strPool {
		bytes, _ := json.Marshal(v)
		total += len(string(bytes))
	}
	return total
}

func (s *SessionService) GetMsg(sessionId string) (msg []Message) {
	sessionContext, ok := s.cache.Get(sessionId)
	if !ok {
		return nil
	}
	sessionMeta := sessionContext.(*SessionMeta)
	return sessionMeta.Msg
}

func (s *SessionService) SetMsg(sessionId string, msg []Message) {
	maxLength := 4096
	maxCacheTime := time.Hour * 12

	for getLength(msg) > maxLength {
		msg = append(msg[:1], msg[2:]...)
	}

	sessionContext, ok := s.cache.Get(sessionId)
	if !ok {
		sessionMeta := &SessionMeta{Msg: msg}
		s.cache.Set(sessionId, sessionMeta, maxCacheTime)
		return
	}
	sessionMeta := sessionContext.(*SessionMeta)
	sessionMeta.Msg = msg
	s.cache.Set(sessionId, sessionMeta, maxCacheTime)
}

func (s *SessionService) Clear(sessionId string) {
	s.cache.Delete(sessionId)
}

func GetSessionCache() CacheInterface {
	if sessionServices == nil {
		sessionServices = &SessionService{cache: cache.New(time.Hour*12, time.Hour*1)}
	}
	return sessionServices
}
