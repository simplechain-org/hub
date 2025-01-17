package courier

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/simplechain-org/crosshub/fabric/courier/contractlib"
	"github.com/simplechain-org/crosshub/fabric/courier/utils"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
)

type FieldName = string

const (
	PK             FieldName = "PK"
	CrossIdIndex   FieldName = "CrossID"
	StatusField    FieldName = "Status"
	TimestampField FieldName = "Timestamp"
)

type DB interface {
	Save(txList []*CrossTx) error
	Updates(idList []string, updaters []func(c *CrossTx)) error
	One(fieldName string, value interface{}) *CrossTx
	Set(bucketName string, key interface{}, value interface{}) error
	Get(bucketName string, key interface{}, to interface{})
	Query(pageSize int, startPage int, orderBy []FieldName, reverse bool, filter ...q.Matcher) []*CrossTx
}

type Store struct {
	db storm.Node
}

func OpenStormDB(dataDir string) (*storm.DB, error) {
	var workDir = os.TempDir()

	if dataDir != "" {
		workDir = dataDir
	}
	os.MkdirAll(workDir, os.ModePerm)
	return storm.Open(filepath.Join(workDir, "rootdb"))
}

func NewStore(root *storm.DB) (*Store, error) {
	s := &Store{}
	s.db = root.From("mychannel").WithBatch(true)
	return s, nil
}

func (s *Store) Set(bucketName string, key interface{}, value interface{}) error {
	return s.db.Set(bucketName, key, value)
}

func (s *Store) Get(bucketName string, key interface{}, to interface{}) {
	if err := s.db.Get(bucketName, key, to); err != nil {
		utils.Logger.Warn("[courier.Store] Get ", "bucketName", bucketName, "key", key, "err", err)
	}
}

func (s *Store) Save(txList []*CrossTx) error {
	utils.Logger.Debug("[courier.RemoteStore] to save cross txs", "len(txList)", len(txList))

	withTransaction, err := s.db.Begin(true)
	if err != nil {
		return fmt.Errorf("db begin err: %w", err)
	}
	defer withTransaction.Rollback()

	for _, newTx := range txList {
		var oldTx CrossTx
		err = withTransaction.One(CrossIdIndex, newTx.CrossID, &oldTx)

		if err == storm.ErrNotFound {
			utils.Logger.Debug("[courier.RemoteStore] save new cross tx", "crossID", newTx.CrossID, "status", newTx.GetStatus(), "blockNumber", newTx.BlockNumber)
			if err = withTransaction.Save(newTx); err != nil {
				return fmt.Errorf("db save err: %w", err)
			}
		} else if oldTx.IContract == nil {
			utils.Logger.Warn("[courier.Store] parse old crossTx failed", "crossID", oldTx.CrossID)
		} else if newTx.IsFinished() {
			utils.Logger.Debug("[courier.Store] receive Finished crossTx ", "crossID", newTx.CrossID, "txId", newTx.TxID)
			// update old status, discard new
			oldTx.UpdateStatus(contractlib.Completed)
			if err = withTransaction.Update(&oldTx); err != nil {
				return fmt.Errorf("db update err: %w", err)
			}
			utils.Logger.Info("[courier.RemoteStore] update Finished to Completed, cross chain transaction completed", "crossID", newTx.CrossID, "txId", newTx.TxID)
		} else {
			utils.Logger.Warn("[courier.RemoteStore] duplicate crossTx", "crossID", newTx.CrossID, "old.status", oldTx.GetStatus(), "new.status", newTx.GetStatus())
			continue
		}
	}

	return withTransaction.Commit()
}

func (s *Store) One(fieldName string, value interface{}) *CrossTx {
	to := CrossTx{}
	if err := s.db.One(fieldName, value, &to); err != nil {
		return nil
	}

	return &to
}

func (s *Store) Updates(idList []string, updaters []func(c *CrossTx)) error {
	if len(idList) != len(updaters) {
		return fmt.Errorf("invalid update params")
	}

	utils.Logger.Debug("[courier.RemoteStore] update list", "idList", idList)

	withTransaction, err := s.db.Begin(true)
	if err != nil {
		return fmt.Errorf("db begin err: %w", err)
	}
	defer withTransaction.Rollback()

	for i, id := range idList {
		var c CrossTx
		if err = withTransaction.One(CrossIdIndex, id, &c); err != nil {
			return fmt.Errorf("db query err: %w", err)
		}

		updaters[i](&c)

		if err = withTransaction.Update(&c); err != nil {
			return fmt.Errorf("db update err: %w", err)
		}
	}

	utils.Logger.Debug("[courier.RemoteStore] update list", "successes", len(idList))

	return withTransaction.Commit()
}

func (s *Store) Query(pageSize int, startPage int, orderBy []FieldName, reverse bool, filter ...q.Matcher) (crossTxs []*CrossTx) {
	if pageSize > 0 && startPage <= 0 {
		return nil
	}

	query := s.db.Select(filter...)
	if len(orderBy) > 0 {
		query.OrderBy(orderBy...)
	}
	if reverse {
		query.Reverse()
	}
	if pageSize > 0 {
		query.Limit(pageSize).Skip(pageSize * (startPage - 1))
	}
	query.Find(&crossTxs)

	return crossTxs
}
