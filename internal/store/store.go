package store

import (
	"fmt"
	"time"

	"github.com/VladPetriv/tg_scanner/config"
	"github.com/VladPetriv/tg_scanner/internal/store/pg"
	"github.com/VladPetriv/tg_scanner/logger"
)

type Store struct {
	Pg     *pg.DB
	Logger *logger.Logger

	Channel ChannelRepo
	Message MessageRepo
	Replie  ReplieRepo
	User    UserRepo
}

func New(cfg config.Config, log *logger.Logger) (*Store, error) {
	pgDB, err := pg.Dial(cfg)
	if err != nil {
		return nil, fmt.Errorf("pg.Dial() failed: %w", err)
	}

	if pgDB != nil {
		log.Info("Running migrations...")
		err := runMigrations(&cfg)
		if err != nil {
			return nil, fmt.Errorf("run migrations error: %w", err)
		}
	}

	var store Store
	store.Logger = log
	if pgDB != nil {
		store.Pg = pgDB
		store.Channel = pg.NewChannelRepo(pgDB)
		store.Message = pg.NewMessageRepo(pgDB)
		store.Replie = pg.NewReplieRepo(pgDB)
		store.User = pg.NewUserRepo(pgDB)

		go store.KeepAliveDB(cfg)
	}

	return &store, nil
}

func (s *Store) KeepAliveDB(cfg config.Config) {
	var err error
	for {
		time.Sleep(time.Second * 5)

		lostConnection := false
		if s.Pg == nil {
			lostConnection = true
		} else if _, err := s.Pg.Exec("SELECT 1;"); err != nil {
			lostConnection = true
		}

		if !lostConnection {
			continue
		}

		s.Logger.Debug("[store.KeepAliveDB] Lost db connection. Restoring...")

		s.Pg, err = pg.Dial(cfg)
		if err != nil {
			s.Logger.Error(err)

			continue
		}

		s.Logger.Debug("[store.KeepAliveDB] DB reconnected")
	}
}
