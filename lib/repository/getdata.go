package repository

import (
	"fmt"
	"love-remittance-be-apps/core/query"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
	"strconv"
	"time"
)

type GetDataAllRepositoryImpl struct {
}

func NewGetDataAllRepository() interfc.GetDataAllRepository {
	return &GetDataAllRepositoryImpl{}
}

func (rep *GetDataAllRepositoryImpl) DataAccountAll(phonePrefix string, phone string) (*model.DataAccountDetail, *model.ErrorData) {
	row, err := query.Row("SELECT\n"+
		"msa.account_id,\n"+
		"msa.account_status_id,\n"+
		"mas.account_status_name,\n"+
		"msa.account_identity_type_id,\n"+
		"mait.identity_type_name,\n"+
		"msa.account_identity_number,\n"+
		"msa.account_name,\n"+
		"mssad.admin_email,\n"+
		"mssad.device_id,\n"+
		"mslcy.city_id,\n"+
		"mslcy.city_name,\n"+
		"mao.occupation_id,\n"+
		"mao.occupation_name,\n"+
		"msa.account_address,\n"+
		"msa.account_pob,\n"+
		"msa.account_dob ,\n"+
		"msa.account_gender,\n"+
		"msa.account_postal_code,\n"+
		"msa.account_selfie_picture\n"+
		"FROM\n"+
		"ms_account msa\n"+
		"JOIN ms_store_admin mssad ON msa.account_id = mssad.account_id\n"+
		"LEFT JOIN ms_store msst ON mssad.store_id = msst.store_id\n"+
		"LEFT JOIN ms_store_branch mssbr ON msst.branch_id = mssbr.branch_id\n"+
		"LEFT JOIN ms_initiator msin ON mssbr.initiator_id = msin.initiator_id\n"+
		"LEFT JOIN ms_location_country mslc ON msin.country_id = mslc.country_id\n"+
		"LEFT JOIN ms_location_city mslcy ON mslcy.city_id = msa.city_id\n"+
		"LEFT JOIN ms_location_province mlp ON mlp.province_id = mslcy.province_id\n"+
		"LEFT JOIN ms_account_identity_type msait ON msait.identity_type_id = msa.account_identity_type_id\n"+
		"left join ms_account_status mas on mas.account_status_id = msa.account_status_id\n"+
		"left join ms_account_identity_type mait on mait.identity_type_id = msa.account_identity_type_id\n"+
		"left join ms_account_occupation mao on mao.occupation_id = msa.occupation_id\n"+
		"WHERE\n"+
		"msa.account_phone_number = $1\n"+
		"AND mslc.country_prefix = $2;\n", phonePrefix+phone, phonePrefix)
	if err != nil {
		return nil, err
	}
	//mapping
	var dataAccount model.DataAccountDetail
	errs := row.Scan(&dataAccount.AccountId, &dataAccount.AccountStatutId, &dataAccount.AccountStatusName,
		&dataAccount.IdentityTypeId, &dataAccount.IdentityTypeName, &dataAccount.IdentityNumber, &dataAccount.AccountName,
		&dataAccount.AdminEmail, &dataAccount.DeviceId, &dataAccount.CityId, &dataAccount.CityName, &dataAccount.OccupationId,
		&dataAccount.OccupationName, &dataAccount.Address, &dataAccount.POB, &dataAccount.DOB, &dataAccount.Gender, &dataAccount.PostalCode, &dataAccount.ImgSelf)
	if errs != nil {
		return nil, &model.ErrorData{
			Description: errs.Error(),
		}
	}
	return &dataAccount, nil
}

func (rep *GetDataAllRepositoryImpl) GetProvince(param model.Param) ([]model.DataIdName, *model.Pagination, *model.ErrorData) {
	kosong := ""
	nol := 0
	all := "all"
	asc := "asc"
	var limit string
	if param.Limit != nil || param.Limit != &nol {
		limit = strconv.Itoa(*param.Limit)
	} else {
		limit = all
	}
	var offset *int
	if param.Offset != nil || param.Offset != &nol {
		offset = param.Offset
	} else {
		offset = &nol
	}
	var search *string
	if param.Search != nil || param.Search != &kosong {
		search = param.Search
	} else {
		search = &kosong
	}
	var sort *string
	if param.Sort != nil || param.Sort != &kosong {
		sort = param.Sort
	} else {
		sort = &asc
	}

	valueSerc := *search
	cari := "%" + valueSerc + "%"
	newSort := *sort
	querySort := "mlp.province_name asc"
	if newSort != "asc" {
		querySort = "mlp.province_name desc"
	}
	newOffset := *offset
	queryLimit := " limit all"
	if limit != "all" {
		queryLimit = " limit " + limit
	}
	//query get data payment method in database
	stringQuery := "select\n" +
		"mlp.province_id as id,\n" +
		"mlp.province_name as name,\n" +
		"count(*) over() as full_count\n" +
		"from\n" +
		"ms_location_province mlp\n" +
		"where\n" +
		"mlp.country_id = 77\n" +
		"and province_name ilike $1\n" +
		"order by\n" +
		querySort + queryLimit + " offset $2"
	rows, err := query.Rows(stringQuery, cari, newOffset)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	var total int
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name, &total)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}

	pagination := model.Pagination{
		Limit:  *param.Limit,
		Offset: *offset,
		Search: *search,
		Sort:   *sort,
		Total:  total,
	}

	return datas, &pagination, nil
}

func (rep *GetDataAllRepositoryImpl) GetCity(param model.Param, provinceId int) ([]model.DataIdName, *model.Pagination, *model.ErrorData) {
	kosong := ""
	var search *string
	if param.Search != nil || param.Search != &kosong {
		search = param.Search
	} else {
		search = &kosong
	}

	valueSerc := *search
	cari := "%" + valueSerc + "%"
	//query get data payment method in database
	stringQuery := "select\n" +
		"mlc.city_id as id,\n" +
		"mlc.city_name,\n" +
		"count(*) over() as full_count\n" +
		"from\n" +
		"ms_location_city mlc\n" +
		"where\n" +
		"mlc.province_id = $1\n" +
		"and mlc.city_name ilike $2\n"
	rows, err := query.Rows(stringQuery, provinceId, cari)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	var total int
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name, &total)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}

	pagination := model.Pagination{
		Search: *search,
		Total:  total}

	return datas, &pagination, nil
}

func (rep *GetDataAllRepositoryImpl) GetCountry(param model.Param) ([]model.DataIdName, *model.ErrorData) {
	kosong := ""
	var search *string
	if param.Search != nil || param.Search != &kosong {
		search = param.Search
		fmt.Printf("masuk sini = " + *search)
	} else {
		search = &kosong
		fmt.Printf("masuk sana")
	}

	valueSerc := *search
	cari := "%" + valueSerc + "%"
	fmt.Println("cari = " + cari)
	//query get data payment method in database
	stringQuery := "select\n" +
		"mlc.country_id ,\n" +
		"mlc.country_name\n" +
		// "count(*) over() as full_count\n" +
		"from\n" +
		"ms_location_country mlc\n" +
		"where\n" +
		"mlc.country_name ilike $1\n" +
		"order by mlc.country_name asc\n"
	rows, err := query.Rows(stringQuery, cari)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}

func (rep *GetDataAllRepositoryImpl) GetOccupation() ([]model.DataIdName, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"mao.occupation_id as id,\n" +
		"mao.occupation_name\n" +
		"from\n" +
		"ms_account_occupation mao\n" +
		"where\n" +
		"mao.country_id = $1\n"
	rows, err := query.Rows(stringQuery, 77)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}

func (rep *GetDataAllRepositoryImpl) GetOccupationIntl(masterStoreId int) ([]model.DataIdName, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"maop.occupation_partner_id as id,\n" +
		"maop.occupation_partner_value as name\n" +
		"from\n" +
		"ms_account_occupation_partner maop\n" +
		"left join ms_partner_master_store mpms on\n" +
		"maop.partner_id = mpms.partner_id\n" +
		"where\n" +
		"mpms.partner_master_store_id = $1\n" +
		"order by\n" +
		"maop.occupation_partner_key asc \n"
	rows, err := query.Rows(stringQuery, masterStoreId)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}
func (rep *GetDataAllRepositoryImpl) GetIdentityType() ([]model.DataIdName, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"mait.identity_type_id as id,\n" +
		"mait.identity_type_name\n" +
		"from\n" +
		"ms_account_identity_type mait\n" +
		"order by\n" +
		"mait.identity_type_id asc\n"
	rows, err := query.Rows(stringQuery)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}

func (rep *GetDataAllRepositoryImpl) GetIdentityTypeIntl(masterStoreId int) ([]model.DataIdName, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"maitp.identity_type_partner_id as id,\n" +
		"maitp.identity_type_partner_value as name\n" +
		"from\n" +
		"ms_account_identity_type_partner maitp\n" +
		"left join ms_partner_master_store mpms on\n" +
		"maitp.partner_id = mpms.partner_id\n" +
		"where\n" +
		"mpms.partner_master_store_id = $1\n" +
		"order by maitp.identity_type_partner_key asc\n"
	rows, err := query.Rows(stringQuery, masterStoreId)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}

func (rep *GetDataAllRepositoryImpl) GetSoFund() ([]model.DataIdName, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"mtsof.transaction_source_of_fund_id as id,\n" +
		"mtsof.transaction_source_of_fund_name as name\n" +
		"from\n" +
		"ms_transaction_source_of_fund mtsof\n" +
		"where\n" +
		"mtsof.country_id = 77;\n"
	rows, err := query.Rows(stringQuery)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}

func (rep *GetDataAllRepositoryImpl) GetSoFundIntl(masterStoreId int) ([]model.DataIdName, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"mtsofp.tsofp_id as id,\n" +
		"mtsofp.tsofp_value as name\n" +
		"from\n" +
		"ms_transaction_source_of_fund_partner mtsofp\n" +
		"left join ms_partner_master_store mpms on\n" +
		"mtsofp.partner_id = mpms.partner_id\n" +
		"where\n" +
		"mpms.partner_master_store_id = $1\n" +
		"order by\n" +
		"mtsofp.tsofp_key asc ;\n"
	rows, err := query.Rows(stringQuery, masterStoreId)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}

func (rep *GetDataAllRepositoryImpl) GetPurpose() ([]model.DataIdName, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"mtp.transaction_purpose_id as id,\n" +
		"mtp.transaction_purpose_name as name\n" +
		"from\n" +
		"ms_transaction_purpose mtp\n" +
		"where\n" +
		"mtp.country_id = 77;\n"
	rows, err := query.Rows(stringQuery)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}
func (rep *GetDataAllRepositoryImpl) GetPurposeIntl(masterStoreId int) ([]model.DataIdName, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"mtpp.transaction_purpose_partner_id as id,\n" +
		"mtpp.transaction_purpose_partner_value  as name\n" +
		"from\n" +
		"ms_transaction_purpose_partner mtpp\n" +
		"left join ms_partner_master_store mpms on\n" +
		"mtpp.partner_id = mpms.partner_id\n" +
		"where\n" +
		"mpms.partner_master_store_id = $1\n" +
		"order by\n" +
		"mtpp.transaction_purpose_partner_key  asc ;\n"
	rows, err := query.Rows(stringQuery, masterStoreId)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}
func (rep *GetDataAllRepositoryImpl) GetRelations() ([]model.DataIdName, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"mr.relationship_id as id,\n" +
		"mr.relationship_name as name\n" +
		"from\n" +
		"ms_relationship mr\n" +
		"where\n" +
		"mr.country_id = 77\n" +
		"order by\n" +
		"order_num asc\n"
	rows, err := query.Rows(stringQuery)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}

func (rep *GetDataAllRepositoryImpl) GetRelationsIntl(masterStoreId int) ([]model.DataIdName, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"mrp.relationship_partner_id as id,\n" +
		"mrp.relationship_partner_value as name\n" +
		"from\n" +
		"ms_relationship_partner mrp\n" +
		"left join ms_partner_master_store mpms on\n" +
		"mrp.partner_id = mpms.partner_id\n" +
		"where\n" +
		"mpms.partner_master_store_id = $1\n" +
		"order by\n" +
		"mrp.relationship_partner_key asc \n"
	rows, err := query.Rows(stringQuery, masterStoreId)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataIdName
	for rows.Next() {
		var data model.DataIdName
		err := rows.Scan(&data.Id, &data.Name)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}

func (rep *GetDataAllRepositoryImpl) GetHistory(phone string) ([]model.History, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "SELECT\n" +
		"DISTINCT ON (trx.transaction_id) transaction_id,\n" +
		"trx.transaction_ext_ref AS ext_ref,\n" +
		"trx.transaction_amount AS amount,\n" +
		"trx.timestamp_insert AS time_trx,\n" +
		"mp.type_region,\n" +
		"mpms.partner_master_store_name as bank_name,\n" +
		"mslc.country_currency as receiver_currency,\n" +
		"CASE\n" +
		"WHEN mvs.voucher_status_id = 1 AND (NOW() AT TIME ZONE'utc-7' - INTERVAL '30 minutes' >= trx.timestamp_insert) THEN 1\n" +
		"WHEN mvs.voucher_status_id = 2 AND trx.transaction_status_id = 3 THEN\t2\n" +
		"WHEN mvs.voucher_status_id = 2 AND mp.partner_type = 'BANK' THEN\t3\n" +
		"WHEN mvs.voucher_status_id = 2 AND mp.partner_type = 'WALLET' THEN\t3\n" +
		"WHEN mvs.voucher_status_id = 2 THEN\t4\n" +
		"WHEN mvs.voucher_status_id = 3 THEN\t5\n" +
		"WHEN mvs.voucher_status_id = 4 THEN\t6\n" +
		"ELSE 7\n" +
		"END AS status,\n" +
		"CASE\n" +
		"WHEN mp.partner_type = 'COUNTER' THEN 'CASH'\n" +
		"WHEN mp.partner_type = 'WALLET' THEN 'WALLET'\n" +
		"ELSE 'BANK'\n" +
		"END AS partner_type\n" +
		"FROM\n" +
		"tr_transaction trx\n" +
		"LEFT JOIN ms_account AS msa ON msa.account_id = trx.account_id_sender\n" +
		"LEFT JOIN ms_account AS msarece ON msarece.account_id = trx.account_id_receiver\n" +
		"LEFT JOIN ms_transaction_source_of_payment AS mtsop ON trx.transaction_sop_id = mtsop.transaction_sop_id\n" +
		"LEFT JOIN ms_source_of_payment AS msop ON mtsop.sop_id = msop.sop_id\n" +
		"LEFT JOIN ms_transaction_source_of_fund AS mtsof ON mtsof.transaction_source_of_fund_id = trx.transaction_source_of_fund_id\n" +
		"LEFT JOIN ms_transaction_type AS mttp ON mttp.transaction_type_id = trx.transaction_type_id\n" +
		"LEFT JOIN ms_partner_fee AS mpf ON mpf.partner_fee_id = trx.partner_fee_id\n" +
		"LEFT JOIN ms_partner AS mp ON mp.partner_id = mpf.partner_id\n" +
		"LEFT JOIN ms_partner_master_store AS mpms ON mpms.partner_additional_data_1 = trx.transaction_ad_bank_code AND mpms.partner_id = mp.partner_id\n" +
		"LEFT JOIN ms_transaction_status AS mts ON mts.transaction_status_id = trx.transaction_status_id\n" +
		"LEFT JOIN tr_voucher AS trv ON trv.voucher_id = trx.voucher_id\n" +
		"LEFT JOIN ms_voucher_status AS mvs ON mvs.voucher_status_id = trv.voucher_status_id\n" +
		"LEFT JOIN ms_location_country as mslc ON mslc.country_id = mpms.country_id\n" +
		"WHERE\n" +
		"msa.account_phone_number = $1\n" +
		"ORDER BY\n" +
		"trx.transaction_id desc limit 25;\n"
	rows, err := query.Rows(stringQuery, phone)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.History
	for rows.Next() {
		var data model.History
		var trxTime time.Time
		senderCurrency := "IDR"
		err := rows.Scan(&data.TransactionId, &data.ExtRef, &data.TransactionAmount,
			&trxTime, &data.Region, &data.BankName, &data.ReceiverCurrency, &data.Status, &data.TypeTrx)
		if err == nil {
			datas = append(datas, model.History{
				TransactionId:     data.TransactionId,
				ExtRef:            data.ExtRef,
				TransactionAmount: data.TransactionAmount,
				TimeTrx:           trxTime.Format("2006-01-02"),
				Region:            data.Region,
				BankName:          data.BankName,
				ReceiverCurrency:  data.ReceiverCurrency,
				SenderCurrency:    senderCurrency,
				Status:            data.Status,
				TypeTrx:           data.TypeTrx,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	return datas, nil
}

func (rep *GetDataAllRepositoryImpl) GetTransaction(param model.Param, sopId int) ([]model.DataDestination, *model.Pagination, *model.ErrorData) {
	// query get data payment method in database
	kosong := ""
	nol := 0
	all := "all"
	asc := "asc"
	var limit string
	if param.Limit != nil || param.Limit != &nol {
		limit = strconv.Itoa(*param.Limit)
	} else {
		limit = all
	}
	var offset *int
	if param.Offset != nil || param.Offset != &nol {
		offset = param.Offset
	} else {
		offset = &nol
	}
	var search *string
	if param.Search != nil || param.Search != &kosong {
		search = param.Search
	} else {
		search = &kosong
	}
	var sort *string
	if param.Sort != nil || param.Sort != &kosong {
		sort = param.Sort
	} else {
		sort = &asc
	}

	valueSerc := *search
	cari := "%" + valueSerc + "%"
	newSort := *sort
	querySort := "mpms.partner_master_store_name asc"
	if newSort != "asc" {
		querySort = "mpms.partner_master_store_name desc"
	}
	newOffset := *offset
	queryLimit := " limit all"
	if limit != "all" {
		queryLimit = " limit " + limit
	}
	rows, err := query.Rows("select distinct (mpms.partner_master_store_name) as bank_name,\n"+
		"mpms.partner_master_store_cashout as type,\n"+
		"mpms.partner_master_store_id as master_store_id ,\n"+
		"count(*) over() as full_count,\n"+
		"mpms.partner_additional_data_1 as bank_code\n"+
		"from\n"+
		"ms_transaction_source_of_payment as mtsop\n"+
		"left join ms_partner_fee as mspf on\n"+
		"mspf.initiator_id = mtsop.initiator_id\n"+
		"left join ms_partner_master_store as mpms on\n"+
		"mpms.partner_id = mspf.partner_id\n"+
		"left join ms_partner as mp on\n"+
		"mp.partner_id = mspf.partner_id\n"+
		"where\n"+
		"mtsop.transaction_sop_id = $1\n"+
		"and mspf.flag_active = true\n"+
		"and mspf.flag_delete = false\n"+
		"and mp.flag_active = true\n"+
		"and mp.type_region = 'domestic'\n"+
		"and mpms.flag_active = true\n"+
		"and mpms.partner_master_store_name ilike $2\n"+
		"order by\n"+
		querySort+
		queryLimit+" offset $3;", sopId, cari, newOffset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	// mapping
	var datas []model.DataDestination
	var total int
	for rows.Next() {
		var data model.DataDestination
		err := rows.Scan(&data.BankName, &data.Types, &data.MasterStoreId, &total, &data.BankCode)
		if err == nil {
			datas = append(datas, data)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, nil, &model.ErrorData{
			Title:       "Error Row",
			Description: err.Error(),
		}
	}
	pagination := model.Pagination{
		Limit:  *param.Limit,
		Offset: *offset,
		Search: *search,
		Sort:   *sort,
		Total:  total,
	}
	return datas, &pagination, nil
}

func (rep *GetDataAllRepositoryImpl) GetHistoryDetail(transactionId int) (*model.History, *model.ErrorData) {
	row, err := query.Row("SELECT\n"+
		"DISTINCT ON (trx.transaction_id) transaction_id,\n"+
		"trx.transaction_ext_ref AS ext_ref,\n"+
		"trx.transaction_amount AS amount,\n"+
		"trx.timestamp_insert AS time_trx,\n"+
		"mp.type_region,\n"+
		"mpms.partner_master_store_name as bank_name,\n"+
		"mslc.country_currency as receiver_currency,\n"+
		"CASE\n"+
		"WHEN mvs.voucher_status_id = 1 AND (NOW() AT TIME ZONE'utc-7' - INTERVAL '30 minutes' >= trx.timestamp_insert) THEN 1\n"+
		"WHEN mvs.voucher_status_id = 2 AND trx.transaction_status_id = 3 THEN\t2\n"+
		"WHEN mvs.voucher_status_id = 2 AND mp.partner_type = 'BANK' THEN\t3\n"+
		"WHEN mvs.voucher_status_id = 2 AND mp.partner_type = 'WALLET' THEN\t3\n"+
		"WHEN mvs.voucher_status_id = 2 THEN\t4\n"+
		"WHEN mvs.voucher_status_id = 3 THEN\t5\n"+
		"WHEN mvs.voucher_status_id = 4 THEN\t6\n"+
		"ELSE 7\n"+
		"END AS status,\n"+
		"CASE\n"+
		"WHEN mp.partner_type = 'COUNTER' THEN 'CASH'\n"+
		"WHEN mp.partner_type = 'WALLET' THEN 'WALLET'\n"+
		"ELSE 'BANK'\n"+
		"END AS partner_type,\n"+
		"mvs.voucher_status_id AS type_id,\n"+
		"mi.initiator_adapter_address as url_adapter,\n"+
		"msarece.account_name as receive_name,\n"+
		"msarece.account_phone_number as receive_phone,\n"+
		"mlc2.country_name as receive_country\n"+
		"FROM\n"+
		"tr_transaction trx\n"+
		"LEFT JOIN ms_account AS msa ON msa.account_id = trx.account_id_sender\n"+
		"LEFT JOIN ms_account AS msarece ON msarece.account_id = trx.account_id_receiver\n"+
		"LEFT JOIN ms_transaction_source_of_payment AS mtsop ON trx.transaction_sop_id = mtsop.transaction_sop_id\n"+
		"LEFT JOIN ms_source_of_payment AS msop ON mtsop.sop_id = msop.sop_id\n"+
		"LEFT JOIN ms_transaction_source_of_fund AS mtsof ON mtsof.transaction_source_of_fund_id = trx.transaction_source_of_fund_id\n"+
		"LEFT JOIN ms_transaction_type AS mttp ON mttp.transaction_type_id = trx.transaction_type_id\n"+
		"LEFT JOIN ms_partner_fee AS mpf ON mpf.partner_fee_id = trx.partner_fee_id\n"+
		"LEFT JOIN ms_partner AS mp ON mp.partner_id = mpf.partner_id\n"+
		"LEFT JOIN ms_partner_master_store AS mpms ON mpms.partner_additional_data_1 = trx.transaction_ad_bank_code AND mpms.partner_id = mp.partner_id\n"+
		"LEFT JOIN ms_transaction_status AS mts ON mts.transaction_status_id = trx.transaction_status_id\n"+
		"LEFT JOIN tr_voucher AS trv ON trv.voucher_id = trx.voucher_id\n"+
		"LEFT JOIN ms_voucher_status AS mvs ON mvs.voucher_status_id = trv.voucher_status_id\n"+
		"LEFT JOIN ms_location_country as mslc ON mslc.country_id = mpms.country_id\n"+
		"left join ms_initiator mi on mi.initiator_id = mtsop.initiator_id\n"+
		"left join ms_location_city mlc on mlc.city_id = msarece.city_id\n"+
		"left join ms_location_province mlp on mlp.province_id = mlc.province_id\n"+
		"left join ms_location_country mlc2 on mlc2.country_id = mlp.country_id\n"+
		"WHERE\n"+
		"trx.transaction_id = $1\n", transactionId)
	if err != nil {
		return nil, err
	}
	//mapping
	var data model.History
	var trxTime time.Time
	senderCurrency := "IDR"
	errs := row.Scan(&data.TransactionId, &data.ExtRef, &data.TransactionAmount,
		&trxTime, &data.Region, &data.BankName, &data.ReceiverCurrency, &data.Status,
		&data.TypeTrx, &data.TypeId, &data.UrlAdapter, &data.ReceiverName, &data.ReceiverPhone, &data.ReceiverCountry)
	if errs != nil {
		return nil, &model.ErrorData{
			Description: errs.Error(),
		}
	}
	return &model.History{
		TransactionId:     data.TransactionId,
		ExtRef:            data.ExtRef,
		TransactionAmount: data.TransactionAmount,
		TimeTrx:           trxTime.Format("2006-01-02"),
		Region:            data.Region,
		BankName:          data.BankName,
		ReceiverCurrency:  data.ReceiverCurrency,
		SenderCurrency:    senderCurrency,
		Status:            data.Status,
		TypeTrx:           data.TypeTrx,
		TypeId:            data.TypeId,
		UrlAdapter:        data.UrlAdapter,
		ReceiverName:      data.ReceiverName,
		ReceiverPhone:     data.ReceiverPhone,
		ReceiverCountry:   data.ReceiverCountry,
	}, nil
}
