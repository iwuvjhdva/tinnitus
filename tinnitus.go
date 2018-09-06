package tinnitus

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"time"
)

func NewTinnitus() *Tinnitus {
	tinnitus := &Tinnitus{
		synthID:        1000,
		transactionIDs: make(map[string]*ClippedTransaction),
	}
	return tinnitus
}

func (tinnitus *Tinnitus) Run() error {
	defer RPC.Shutdown()

	if err := tinnitus.createSynthDefs(); err != nil {
		return err
	}

	height := *Flags.From
	transactionCount := 0
	transactionIDs := make(map[string]bool)

	for running := true; running; running = !ShouldExitGracefully() {
		var err error

		var hash *chainhash.Hash
		var block *GetBitcoinBlockVerboseResult

		if hash, err = RPC.GetBlockHash(height); err != nil {
			return err
		}

		if block, err = RPC.GetBitcoinBlockVerboseTx(hash); err != nil {
			return err
		}

		transactionCount += len(block.Tx)

		for _, tx := range block.Tx {
			transactionIDs[tx.Txid] = true
			completeCount := 0

			for _, vin := range tx.Vin {

				if transactionIDs[vin.Txid] {
					delete(transactionIDs, vin.Txid)
					completeCount++
				}
			}

		}

		if err = tinnitus.playBlock(block); err != nil {
			return err
		}

		time.Sleep(Config.SleepDuration)
		height++
	}

	return nil
}
