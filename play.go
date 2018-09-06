package tinnitus

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/scgolang/sc"
	"time"
)

func (tinnitus *Tinnitus) playBlock(block *GetBitcoinBlockVerboseResult) error {
	for _, tx := range block.Tx {

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

		if err := tinnitus.playTx(tx.Vout, &controls); err != nil {
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
	return nil
}

func (tinnitus *Tinnitus) playTx(vout []btcjson.Vout, controls *map[string]float32) error {
	for _, vout := range vout {
		*controls = map[string]float32{
			"freq": float32(vout.Value) * 100,
			"gain": float32(0.5),
		}

		if err := tinnitus.synth("Tx", *controls); err != nil {
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
	if _, err := SC.group.Synth("Vout", tinnitus.synthID, sc.AddToTail, controls); err != nil {
		return err
	}

	tinnitus.synthID++

	return nil
}

func (tinnitus *Tinnitus) createSynthDefs() error {
	vout := sc.NewSynthdef("Vout", func(p sc.Params) sc.Ugen {
		freq := p.Add("freq", 440)
		gain := p.Add("gain", 0)
		bus := sc.C(0)
		l := sc.SinOsc{Freq: freq}.Rate(sc.AR)
		r := sc.SinOsc{Freq: freq, Phase: sc.C(0.5)}.Rate(sc.AR)
		pos := sc.SinOsc{Freq: sc.C(0.5)}.Rate(sc.KR)
		sig := sc.Balance2{L: l, R: r, Pos: pos, Level: gain}.Rate(sc.AR)
		return sc.Out{Bus: bus, Channels: sig}.Rate(sc.AR)
	})

	if err := SC.SendDef(vout); err != nil {
		return err
	}

	coinbase := sc.NewSynthdef("Coinbase", func(p sc.Params) sc.Ugen {
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
