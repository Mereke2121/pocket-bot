package boltdb

import (
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/pkg/errors"

	"github.com/pocket-bot/pkg/repository"
)

type TokenStorage struct {
	db *bolt.DB
}

func NewStorage(db *bolt.DB) *TokenStorage {
	return &TokenStorage{db: db}
}

func (s *TokenStorage) Get(chatId int, bucket repository.Bucket) (string, error) {
	var token string
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		token = string(b.Get(intToBytes(chatId)))
		return nil
	})
	if err != nil {
		return "", errors.Wrap(err, "get token from db")
	}
	if token == "" {
		return "", errors.New("can't find token")
	}
	return token, nil
}

func (s *TokenStorage) Save(chatId int, token string, bucket repository.Bucket) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(intToBytes(chatId), []byte(token))
		if err != nil {
			return errors.Wrap(err, "put token by chat id into the db")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func intToBytes(num int) []byte {
	return []byte(strconv.Itoa(num))
}
