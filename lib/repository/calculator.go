package repository

import (
	"love-remittance-be-apps/core/query"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
)

type CalculatorRepositoryImp struct{}

func NewCalculatorRepository() interfc.CalculatorRepository {
	return &CalculatorRepositoryImp{}
}

func (rep *CalculatorRepositoryImp) GetPayment() ([]model.CalculotorPayment, *model.ErrorData) {
	//query get data payment method in database
	rows, err := query.Rows("SELECT\n" +
		"CASE\n" +
		"WHEN msop.sop_type = 'VA' THEN\n" +
		"'BANK' else msop.sop_type\n" +
		"END sop_name,\n" +
		"mi.initiator_id\n" +
		"FROM\n" +
		"ms_source_of_payment AS msop\n" +
		"LEFT JOIN ms_transaction_source_of_payment AS mstsop ON msop.sop_id = mstsop.sop_id\n" +
		"LEFT JOIN ms_initiator AS mi ON mi.initiator_id = mstsop.initiator_id\n" +
		"LEFT JOIN ms_location_country mslc ON mi.country_id = mslc.country_id\n" +
		"WHERE\n" +
		"mstsop.transaction_sop_flag = TRUE\n" +
		"AND mslc.country_prefix = '62'\n" +
		"AND msop.sop_name IN ( 'Mitra/Agen LOVE!', 'MANDIRI', 'QRIS' )\n" +
		"ORDER BY\n" +
		"mi.initiator_id ASC;\n")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var payments []model.CalculotorPayment
	for rows.Next() {
		var payment model.CalculotorPayment
		err := rows.Scan(&payment.SopName, &payment.InitiatorId)
		if err == nil {
			payments = append(payments, payment)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return payments, nil

}

func (rep *CalculatorRepositoryImp) Destination(initiator_id int) ([]model.Destination, *model.ErrorData) {
	//query get data payment method in database
	rows, err := query.Rows("SELECT\n"+
		"*\n"+
		"FROM\n"+
		"(\n"+
		"SELECT DISTINCT ON\n"+
		"(  mpms.country_id ) mpms.country_id,\n"+
		"mslc.country_name,\n"+
		"mslc.country_code_alpha\n"+
		"FROM\n"+
		"ms_initiator AS msi\n"+
		"LEFT JOIN ms_initiator_fee AS msif ON msi.initiator_id = msif.initiator_id\n"+
		"left join ms_partner_fee mpf on mpf.initiator_id = msif.initiator_id\n"+
		"LEFT JOIN ms_partner AS msp ON msp.partner_id = mpf.partner_id\n"+
		"LEFT JOIN ms_partner_master_store AS mpms ON mpms.partner_id = msif.partner_id\n"+
		"LEFT JOIN ms_location_country AS mslc ON mpms.country_id = mslc.country_id\n"+
		"WHERE\n"+
		"msi.initiator_id = $1\n"+
		"AND msp.type_region = 'international'\n"+
		"AND mpms.flag_active = TRUE\n"+
		"AND mpms.flag_delete = FALSE\n"+
		"AND msp.flag_active = TRUE\n"+
		"AND msif.flag_active = TRUE\n"+
		"ORDER BY\n"+
		"mpms.country_id ASC\n"+
		"\n"+
		") AS main\n"+
		"ORDER BY\n"+
		"main.country_name ASC;\n", initiator_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//mapping
	var destinations []model.Destination
	for rows.Next() {
		var destination model.Destination
		err := rows.Scan(&destination.CountryId, &destination.CountryName, &destination.CountryCodeAlpha)
		if err == nil {
			destinations = append(destinations, destination)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return destinations, nil
}

func (rep *CalculatorRepositoryImp) Transaction(initiator_id int, country_id int) ([]model.Transaction, *model.ErrorData) {
	//query get data payment method in database
	rows, err := query.Rows("SELECT\n"+
		"*\n"+
		"FROM\n"+
		"(\n"+
		"SELECT DISTINCT ON\n"+
		"(  mpms.partner_master_store_cashout ) mpms.partner_master_store_cashout as transaction_type,\n"+
		"mpms.partner_master_store_id as master_store_id,\n"+
		"mpms.partner_master_store_name as bank_name\n"+
		"FROM\n"+
		"ms_initiator AS msi\n"+
		"LEFT JOIN ms_initiator_fee AS msif ON msi.initiator_id = msif.initiator_id\n"+
		"LEFT JOIN ms_partner AS msp ON msp.partner_id = msif.partner_id\n"+
		"LEFT JOIN ms_partner_master_store AS mpms ON mpms.partner_id = msif.partner_id\n"+
		"LEFT JOIN ms_location_country AS mslc ON mpms.country_id = mslc.country_id\n"+
		"WHERE\n"+
		"msi.initiator_id = $1\n"+
		"AND msp.type_region = 'international'\n"+
		"AND mpms.country_id = $2\n"+
		"AND mpms.flag_active = TRUE\n"+
		"AND mpms.flag_delete = FALSE\n"+
		"AND msp.flag_active = TRUE\n"+
		"AND msif.flag_active = TRUE\n"+
		"ORDER BY\n"+
		"mpms.partner_master_store_cashout ASC\n"+
		"\n"+
		") AS main\n"+
		"ORDER BY\n"+
		"main.bank_name DESC;\n", initiator_id, country_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	//mapping
	var transactions []model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		err := rows.Scan(&transaction.TransactionType, &transaction.MasterStoreId, &transaction.BankName)
		if err == nil {
			transactions = append(transactions, transaction)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return transactions, nil
}

func (rep *CalculatorRepositoryImp) PartnerMasterStoreId(masterStoreId int) (*model.MasterStore, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.Row("select\n"+
		"mpms.partner_master_store_name as bank_name,\n"+
		"mpms.partner_additional_data_1 as bank_code,\n"+
		"mpms.partner_id,\n"+
		"mpms.country_id\n"+
		"from\n"+
		"ms_partner_master_store mpms\n"+
		"where\n"+
		"mpms.partner_master_store_id = $1\n"+
		"and flag_active = true;\n", masterStoreId)
	if err != nil {
		return nil, err
	}
	//mapping
	var masterStore model.MasterStore
	errss := row.Scan(&masterStore.BankName, &masterStore.BankCode, &masterStore.PartnerId, &masterStore.CountryId)
	if errss != nil {
		return nil, &model.ErrorData{
			Description: errss.Error()}
	}
	return &masterStore, nil
}
