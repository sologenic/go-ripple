package data

import "fmt"

func (txm *TransactionWithMetaData) AMM() (*AMM, error) {
	if txm.GetTransactionType() != AMM_DEPOSIT &&
		txm.GetTransactionType() != AMM_WITHDRAW &&
		txm.GetTransactionType() != AMM_CREATE &&
		txm.GetTransactionType() != AMM_VOTE &&
		txm.GetTransactionType() != AMM_BID &&
		txm.GetTransactionType() != PAYMENT {

		return nil, nil
	}
	for _, nodeAffect := range txm.MetaData.AffectedNodes {
		switch {
		case nodeAffect.CreatedNode != nil && nodeAffect.CreatedNode.LedgerEntryType == AMMROOT:
			ammParsed, ok := nodeAffect.CreatedNode.NewFields.(*AMM)
			if ok {
				return ammParsed, nil
			}
		case nodeAffect.ModifiedNode != nil && nodeAffect.ModifiedNode.LedgerEntryType == AMMROOT:
			ammParsed, ok := nodeAffect.ModifiedNode.FinalFields.(*AMM)
			if ok {
				return ammParsed, nil
			}
		}
	}
	return nil, fmt.Errorf("AMM not found")
}
