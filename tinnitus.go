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
	var prevBlock *GetBitcoinBlockVerboseResult

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

		Logger.Debug("Playing block ", height, " time ", time.Unix(block.Time, 0))

		if err = tinnitus.playBlock(block); err != nil {
			return err
		}

		if prevBlock != nil {
			sleepDuration := time.Unix(block.Time, 0).Sub(time.Unix(prevBlock.Time, 0))
			Logger.Debug("Sleeping ", sleepDuration)
			// time.Sleep(sleepDuration)
			time.Sleep(Config.SleepDuration)
		}

		prevBlock = block
		height++
	}

	return nil
}

func (tinnitus *Tinnitus) Sandbox() error {
	defer RPC.Shutdown()

	if err := tinnitus.createSynthDefs(); err != nil {
		return err
	}

	controls := make(map[string]float32)
	// controls := map[string]float32{
	// 	"freq": float32(20) * 100,
	// 	"gain": float32(0.5),
	// }
	if err := tinnitus.synth("Coinbase", controls); err != nil {
		return err
	}

	return nil
}
