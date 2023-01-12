package server

import (
	"context"
	"net/http"
	"strconv"

	"github.com/zhashkevych/go-pocket-sdk"
	"go.uber.org/zap"

	"github.com/pocket-bot/pkg/repository"
	"github.com/pocket-bot/pkg/repository/boltdb"
)

type Server struct {
	httpServer *http.Server

	storage      repository.TokenStorage
	pocketClient *pocket.Client
	logger       *zap.Logger
	tlgURL       string
}

func NewHttpServer(storage *boltdb.TokenStorage, pocketClient *pocket.Client, logger *zap.Logger, tlgURL string) *Server {
	return &Server{
		storage:      storage,
		pocketClient: pocketClient,
		logger:       logger,
		tlgURL:       tlgURL,
	}
}

func (s *Server) Start(address string) error {
	s.httpServer = &http.Server{
		Addr:    address,
		Handler: s,
	}
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	chatId, err := strconv.Atoi(r.URL.Query().Get("chat_id"))
	if err != nil {
		s.logger.Error("get chat id from url", zap.Error(err))
	}

	token, err := s.storage.Get(chatId, repository.RequestTokens)
	if err != nil {
		s.logger.Error("get request token from db by chat id", zap.Error(err))
	}

	authResp, err := s.pocketClient.Authorize(context.Background(), token)
	if err != nil {
		s.logger.Error("authorize user by token")
	}

	err = s.storage.Save(chatId, authResp.AccessToken, repository.AccessTokens)
	if err != nil {
		s.logger.Error("save access token by chat id", zap.Error(err))
	}

	w.Header().Set("Location", s.tlgURL)
	w.WriteHeader(http.StatusMovedPermanently)
}
