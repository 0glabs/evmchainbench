package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

type Store struct {
	TxStoreDir string
	PrepareTxCache types.Transactions
}

func NewStore(txStoreDir string) *Store {
	return &Store{
		TxStoreDir: txStoreDir,
		PrepareTxCache: types.Transactions{},
	}
}

func (s *Store) AddPrepareTx(tx *types.Transaction) {
	s.PrepareTxCache = append(s.PrepareTxCache, tx)
}

func (s *Store) PersistPrepareTxs() error {
	return persistTxs(s.prepareFilePath(), s.PrepareTxCache)
}

func (s *Store) PersistTxsMap(txsMap map[int]types.Transactions) error {
	for index, txs := range txsMap {
		err := persistTxs(s.txsFilePath(index), txs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) LoadPrepareTxs() (types.Transactions, error) {
	return loadTxs(s.prepareFilePath())
}

func (s *Store) LoadTxsMap() (map[int]types.Transactions, error) {
	txsMap := make(map[int]types.Transactions)

	pattern := fmt.Sprintf("%s/transactions-*.rlp", s.TxStoreDir)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return txsMap, err
	}

	if len(matches) == 0 {
		return txsMap, fmt.Errorf("No transaction files are found")
	} else {
		for index, match := range matches {
			txs, err := loadTxs(match)
			if err != nil {
				return txsMap, err
			}

			txsMap[index] = txs
		}
	}

	return txsMap, nil
}

func (s *Store) prepareFilePath() string {
	return filepath.Join(s.TxStoreDir, "prepare.rlp")
}

func (s *Store) txsFilePath(index int) string {
	return filepath.Join(s.TxStoreDir, fmt.Sprintf("transactions-%d.rlp", index))
}

func persistTxs(path string, txs types.Transactions) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = rlp.Encode(file, txs)
	if err != nil {
		return err
	}

	return nil
}

func loadTxs(path string) (types.Transactions, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var txs types.Transactions

	err = rlp.Decode(file, &txs)
	if err != nil {
		return nil, err
	}

	return txs, nil
}
