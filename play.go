package tinnitus

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/scgolang/sc"
	"time"
)

func (tinnitus *Tinnitus) playBlock(block *GetBitcoinBlockVerboseResult) error {
	Logger.Debugf("Transactions number %d, difficulty %f", len(block.Tx), block.Difficulty)
	for _, tx := range block.Tx {

		var err error
		var txHash *chainhash.Hash
		var transaction *btcjson.TxRawResult

		if txHash, err = chainhash.NewHashFromStr(tx.Txid); err != nil {
			return err
		}
		if transaction, err = RPC.GetRawTransactionVerbose(txHash); err != nil {
			return err
		}

		blockTime := time.Unix(block.Time, 0)
		txTime := time.Unix(transaction.Time, 0)

		if block.Time != transaction.Time {
			Logger.Errorf("Block time %s, transaction time %s.", blockTime, txTime)
		}

		for _, vin := range tx.Vin {

			if clippedTx, ok := tinnitus.transactionIDs[vin.Txid]; ok {
				tinnitus.playPreviousTx(clippedTx)
				delete(tinnitus.transactionIDs, vin.Txid)
			}

			if vin.IsCoinBase() {
				tinnitus.playCoinbase(vin)
			}
		}

		controls := make(map[string]float32)

		if err := tinnitus.playTx(tx.Vout, block.Confirmations, &controls); err != nil {
			return err
		}

		tinnitus.transactionIDs[tx.Txid] = &ClippedTransaction{
			Time:     time.Unix(tx.Time, 0),
			Controls: controls,
		}
	}

	return nil
}

func (tinnitus *Tinnitus) playPreviousTx(clippedTx *ClippedTransaction) error {
	if err := tinnitus.synth("PrevVout", clippedTx.Controls); err != nil {
		return err
	}

	return nil
}

func (tinnitus *Tinnitus) playTx(vouts []btcjson.Vout, confirmations int64, controls *map[string]float32) error {
	for _, vout := range vouts {
		*controls = map[string]float32{
			"freq": float32(vout.Value) * 100,
			"gain": float32(confirmations) / 2000000,
		}

		Logger.Debug(*controls, confirmations, len(vouts))
		if err := tinnitus.synth("Vout", *controls); err != nil {
			return err
		}

	}

	return nil
}

func (tinnitus *Tinnitus) playCoinbase(vin btcjson.Vin) error {
	controls := map[string]float32{}

	if err := tinnitus.synth("Coinbase", controls); err != nil {
		return err
	}

	return nil
}

func (tinnitus *Tinnitus) synth(name string, controls map[string]float32) error {
	Logger.Debug("Playing ", name, "...")

	if _, err := SC.group.Synth(name, -1, sc.AddToTail, controls); err != nil {
		return err
	}

	tinnitus.synthID++

	return nil
}

func (tinnitus *Tinnitus) createSynthDefs() error {
	ampEnv := sc.EnvGen{
		Env: sc.EnvTriangle{
			Dur:   sc.C(0.5),
			Level: sc.C(1),
		},
		Gate:       sc.C(1),
		LevelScale: sc.C(0),
		LevelBias:  sc.C(0),
		TimeScale:  sc.C(1),
		Done:       sc.FreeEnclosing,
	}.Rate(sc.KR)

	vout := sc.NewSynthdef("Vout", func(params sc.Params) sc.Ugen {
		freq := params.Add("freq", 440)
		gain := params.Add("gain", 0)
		bus := sc.C(0)
		l := sc.SinOsc{Freq: freq}.Rate(sc.AR)
		r := sc.SinOsc{Freq: freq, Phase: sc.C(0.5)}.Rate(sc.AR)
		pos := sc.SinOsc{Freq: sc.C(0.5)}.Rate(sc.KR)
		sig := sc.Balance2{L: l, R: r, Pos: pos, Level: gain}.Rate(sc.AR)
		return sc.Out{Bus: bus, Channels: sig.Mul(ampEnv)}.Rate(sc.AR)
	})

	if err := SC.SendDef(vout); err != nil {
		return err
	}

	prevVout := sc.NewSynthdef("PrevVout", func(params sc.Params) sc.Ugen {
		freq := params.Add("freq", 440)
		gain := params.Add("gain", 0)
		bus := sc.C(0)
		l := sc.Saw{Freq: freq}.Rate(sc.AR)
		r := sc.Saw{Freq: freq}.Rate(sc.AR)
		pos := sc.SinOsc{Freq: sc.C(0.5)}.Rate(sc.KR)
		sig := sc.Balance2{L: l, R: r, Pos: pos, Level: gain}.Rate(sc.AR)
		return sc.Out{Bus: bus, Channels: sig.Mul(ampEnv)}.Rate(sc.AR)
	})

	if err := SC.SendDef(prevVout); err != nil {
		return err
	}

	coinbase := sc.NewSynthdef("Coinbase", func(params sc.Params) sc.Ugen {
		ampEnv := sc.EnvGen{
			Env: sc.EnvPerc{
				Attack:  sc.C(0.01),
				Release: sc.C(1),
				Level:   sc.C(1),
				Curve:   sc.C(-4),
			},
			Gate:       sc.C(1),
			LevelScale: sc.C(1),
			LevelBias:  sc.C(0),
			TimeScale:  sc.C(1),
			Done:       sc.FreeEnclosing,
		}.Rate(sc.KR)

		return sc.Out{
			Bus:      sc.C(0),
			Channels: sc.PinkNoise{}.Rate(sc.AR).Mul(ampEnv),
		}.Rate(sc.AR)
	})

	return SC.SendDef(coinbase)
}
