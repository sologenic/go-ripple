package data

import (
	"fmt"
	"sort"
)

type Trade struct {
	LedgerSequence   uint32
	TransactionIndex uint32
	TransactionType  string
	Op               string
	Paid             *Amount
	Got              *Amount
	Giver            Account
	Taker            Account

	FinalFields    *Offer
	PreviousFields *Offer
}

/*
	           {
	                "ModifiedNode": {
	                    "FinalFields": {
	                        "Balance": {
	                            "currency": "434F524500000000000000000000000000000000",
	                            "issuer": "rrrrrrrrrrrrrrrrrrrrBZbvji",
	                            "value": "-5818.434290071079"
	                        },
	                        "Flags": 2228224,
	                        "HighLimit": {
	                            "currency": "434F524500000000000000000000000000000000",
	                            "issuer": "rQ3VRvcp4taJrkp6b2WJW3TPP5zxsFp5Ld",
	                            "value": "500000000"
	                        },
	                        "HighNode": "8",
	                        "LowLimit": {
	                            "currency": "434F524500000000000000000000000000000000",
	                            "issuer": "rcoreNywaoz2ZCQ8Lg2EbSLnGuRBmun6D",
	                            "value": "0"
	                        },
	                        "LowNode": "154"
	                    },
	                    "LedgerEntryType": "RippleState",
	                    "LedgerIndex": "5E7AB54AE88612123DD989C3AEF14B008BD5C38E4066E053C24DB6A87986486B",
	                    "PreviousFields": {
	                        "Balance": {
	                            "currency": "434F524500000000000000000000000000000000",
	                            "issuer": "rrrrrrrrrrrrrrrrrrrrBZbvji",
	                            "value": "-5817.434290071079"
	                        }
	                    },
	                    "PreviousTxnID": "5172D1444AA33ADA8DACBE994008E5777F53DCA043F936F26EBAB4F49F9EFD4D",
	                    "PreviousTxnLgrSeq": 87420168
	                }
	            },

				e.g.:
				abs(modifiednode.finalfields.balance.value) - abs(previousfields.balance.value) = paid amount to user in finalfields.highlimit.issuer

				Does that also work for offer? If so we can get rid of bookdirectory function and solve this a lot easier

				So the offer part looks like:

				{
                "ModifiedNode": {
                    "FinalFields": {
                        "Account": "rhqTdSsJAaEReRsR27YzddqyGoWTNMhEvC",
                        "BookDirectory": "5C8970D155D65DB8FF49B291D7EFFA4A09F9E8A68D9974B25A08B7604E87EA3A",
                        "BookNode": "0",
                        "Flags": 0,
                        "OwnerNode": "2",
                        "Sequence": 71723515,
                        "TakerGets": {
                            "currency": "534F4C4F00000000000000000000000000000000",
                            "issuer": "rsoLo2S1kiGeCcn6hCUXVrCpGMWLrRrLZz",
                            "value": "39.81129921936846"
                        },
                        "TakerPays": "9767400"
                    },
                    "LedgerEntryType": "Offer",
                    "LedgerIndex": "032ABC1F85D88CC888D4872DD7D0AE45CFF79FEDC4EF4C0C6C76800FE00AF56C",
                    "PreviousFields": {
                        "TakerGets": {
                            "currency": "534F4C4F00000000000000000000000000000000",
                            "issuer": "rsoLo2S1kiGeCcn6hCUXVrCpGMWLrRrLZz",
                            "value": "40.759361979"
                        },
                        "TakerPays": "10000000"
                    },
                    "PreviousTxnID": "FD076550A85CC2FBEF394752C45FA128435A059EA768678412255825FAD155F1",
                    "PreviousTxnLgrSeq": 87420135
                }
            },

			And substraction FInalFields - PreviousFields seems to give us the received amount using takerGets

*/
func newRippleState(txm *TransactionWithMetaData, i int) (*Trade, error) {
	// 	_, final, previous, action := txm.MetaData.AffectedNodes[i].AffectedNode()
	// 	v, ok := final.(*RippleState)
	// 	if !ok || action == Created {
	// 		return nil, nil
	// 	}
	// 	p := previous.(*RippleState)
	// 	// if p != nil && p.TakerGets == nil || p.TakerPays == nil {
	// 	// 	// Some "micro" offer consumptions don't change both balances!
	// 	// 	return nil, nil
	// 	// }
	// 	paid := v.HighLimit
	// 	if p.HighLimit != nil {
	// 		var err error
	// 		paid, err = p.HighLimit.Subtract(v.HighLimit)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}
	// 	got := v.LowLimit
	// 	if p.HighLimit != nil {
	// 		var err error
	// 		got, err = p.HighLimit.Subtract(v.HighLimit)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	}

	// 	/*
	// Offer: FinalFields.BookDirectory.String()
	// 	if price, err := priceFromBookDirectory(rplTrade.FinalFields.BookDirectory.String(), rplTrade.Paid.Currency.String(), rplTrade.Got.Currency.String()); err == nil {

	// 	*/
	// 	o:=&Offer{
	// 		FinalFields
	// 	}
	// 	// got, err := p.TakerGets.Subtract(v.TakerGets)
	// 	// if err != nil {
	// 	// 	return nil, err
	// 	// }
	// 	trade := &Trade{
	// 		LedgerSequence:   txm.LedgerSequence,
	// 		TransactionIndex: txm.MetaData.TransactionIndex,
	// 		TransactionType:  txm.GetTransactionType().String(),
	// 		Op:               "Modify",
	// 		Paid:             paid,
	// 		Got:              got,
	// 		Giver:            v.HighLimit.Issuer,
	// 		Taker:            txm.Transaction.GetBase().Account,
	// 		// FinalFields:      v, => We Required this data as *Offer to calculate the price
	// 		// PreviousFields:   p,
	// 	}
	// 	if action == Deleted {
	// 		trade.Op = "Delete"
	// 	}
	// 	return trade, nil
	return nil, nil
}

func newTrade(txm *TransactionWithMetaData, i int) (*Trade, error) {
	_, final, previous, action := txm.MetaData.AffectedNodes[i].AffectedNode()
	v, ok := final.(*Offer)
	if !ok || action == Created {
		return newRippleState(txm, i)
	}
	p := previous.(*Offer)
	if p != nil && p.TakerGets == nil || p.TakerPays == nil {
		// Some "micro" offer consumptions don't change both balances!
		return nil, nil
	}
	paid, err := p.TakerPays.Subtract(v.TakerPays)
	if err != nil {
		return nil, err
	}
	got, err := p.TakerGets.Subtract(v.TakerGets)
	if err != nil {
		return nil, err
	}
	trade := &Trade{
		LedgerSequence:   txm.LedgerSequence,
		TransactionIndex: txm.MetaData.TransactionIndex,
		TransactionType:  txm.GetTransactionType().String(),
		Op:               "Modify",
		Paid:             paid,
		Got:              got,
		Giver:            *v.Account,
		Taker:            txm.Transaction.GetBase().Account,
		FinalFields:      v,
		PreviousFields:   p,
	}
	if action == Deleted {
		trade.Op = "Delete"
	}
	return trade, nil
}

func (t *Trade) Rate() float64 {
	return t.Got.Ratio(*t.Paid).Float()
}

func (t Trade) String() string {
	return fmt.Sprintf("%8d %3d %22.8f %22.8f %-38s %22.8f %-38s %34s %34s %11s %s", t.LedgerSequence, t.TransactionIndex, t.Rate(), t.Paid.Float(), t.Paid.Asset(), t.Got.Float(), t.Got.Asset(), t.Taker, t.Giver, t.TransactionType, t.Op)

}

type TradeSlice []Trade

func NewTradeSlice(txm *TransactionWithMetaData) (TradeSlice, error) {
	var trades TradeSlice
	for i := range txm.MetaData.AffectedNodes {
		trade, err := newTrade(txm, i)
		if err != nil {
			return nil, err
		}
		if trade != nil {
			trades = append(trades, *trade)
		}
	}
	trades.Sort()
	return trades, nil
}

func (s TradeSlice) Filter(account Account) TradeSlice {
	var trades TradeSlice
	for i := range s {
		if s[i].Taker.Equals(account) || s[i].Giver.Equals(account) {
			trades = append(trades, s[i])
		}
	}
	return trades
}

func (s TradeSlice) Sort()         { sort.Sort(s) }
func (s TradeSlice) Len() int      { return len(s) }
func (s TradeSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TradeSlice) Less(i, j int) bool {
	if s[i].LedgerSequence == s[j].LedgerSequence {
		if s[i].TransactionIndex == s[j].TransactionIndex {
			if s[i].Got.Currency.Equals(s[j].Got.Currency) {
				if s[i].Got.Issuer.Equals(s[j].Got.Issuer) {
					if s[i].Paid.Currency.Equals(s[j].Paid.Currency) {
						if s[i].Paid.Issuer.Equals(s[j].Paid.Issuer) {
							return s[i].Rate() > s[j].Rate()
						}
						return s[i].Paid.Issuer.Less(s[j].Paid.Issuer)
					}
					return s[i].Paid.Currency.Less(s[j].Paid.Currency)
				}
				return s[i].Got.Issuer.Less(s[j].Got.Issuer)
			}
			return s[i].Got.Currency.Less(s[j].Got.Currency)
		}
		return s[i].TransactionIndex < s[j].TransactionIndex
	}
	return s[i].LedgerSequence < s[j].LedgerSequence
}
