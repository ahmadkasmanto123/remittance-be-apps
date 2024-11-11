package repository

import (
	"fmt"
	"love-remittance-be-apps/core/query"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
)

type InternationalRepositoryImpl struct{}

func NewInternationalRepository() interfc.InternationalRepository {
	return &InternationalRepositoryImpl{}
}

func (rep *InternationalRepositoryImpl) GetSourceOfPayment() ([]model.DataSOPayment, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "SELECT\n" +
		"msop.sop_name,\n" +
		"msop.sop_type,\n" +
		"mstsop.transaction_sop_id as sop_id,\n" +
		"mi.initiator_id\n" +
		"FROM\n" +
		"ms_source_of_payment as msop\n" +
		"LEFT JOIN ms_transaction_source_of_payment as mstsop ON msop.sop_id = mstsop.sop_id\n" +
		"LEFT JOIN ms_initiator as mi ON mi.initiator_id = mstsop.initiator_id\n" +
		"LEFT JOIN ms_location_country mslc ON mi.country_id = mslc.country_id\n" +
		"WHERE\n" +
		"mstsop.transaction_sop_flag = TRUE\n" +
		"AND mslc.country_prefix = $1\n" +
		"ORDER BY msop.sop_type asc;\n"
	rows, err := query.Rows(stringQuery, "62")
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.DataSOPayment
	for rows.Next() {
		var data model.DataSOPayment
		err := rows.Scan(&data.SopName, &data.SopType, &data.SopId, &data.InitiatorId)
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

func (rep *InternationalRepositoryImpl) GetCountry(sopId int) ([]model.CountryData, *model.ErrorData) {
	//query get data payment method in database
	stringQuery := "select\n" +
		"distinct (mpms.country_id) id,\n" +
		"mlc.country_name as name,\n" +
		"mlc.country_currency as currency,\n" +
		"mlc.country_code_alpha as code_2,\n" +
		"mlc.country_code as code_3\n" +
		"from\n" +
		"ms_partner_fee mpf\n" +
		"left join ms_initiator_fee mif on\n" +
		"mif.partner_id = mpf.partner_id\n" +
		"left join ms_transaction_source_of_payment mtsop on\n" +
		"mtsop.initiator_id = mif.initiator_id\n" +
		"left join ms_partner mp on\n" +
		"mp.partner_id = mpf.partner_id\n" +
		"left join ms_partner_master_store mpms on\n" +
		"mpms.partner_id = mp.partner_id\n" +
		"left join ms_location_country mlc on\n" +
		"mlc.country_id = mpms.country_id\n" +
		"where\n" +
		"mtsop.transaction_sop_id = $1\n" +
		"and mpf.flag_active = true\n" +
		"and mif.flag_active = true\n" +
		"and mp.flag_active = true\n" +
		"and mp.type_region = 'international'\n" +
		"order by\n" +
		"mpms.country_id asc;\n"
	rows, err := query.Rows(stringQuery, sopId)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.CountryData
	for rows.Next() {
		var data model.CountryData
		err := rows.Scan(&data.CountryId, &data.CountryName, &data.Currency, &data.Code2, &data.Code3)
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

func (rep *InternationalRepositoryImpl) GetDestination(param model.Param, countryId int, sopiId int) ([]model.BankDest, *model.ErrorData) {
	//query get data payment method in database
	kosong := ""
	var search *string
	if param.Search != nil || param.Search != &kosong {
		search = param.Search
	} else {
		search = &kosong
	}
	valueSerc := *search
	cari := "%" + valueSerc + "%"

	stringQuery := "select\n" +
		"distinct (mpms.partner_master_store_id)master_store_id,\n" +
		"mpms.partner_master_store_name as bank_name,\n" +
		"mpms.partner_master_store_cashout as type_trx\n" +
		"from\n" +
		"ms_partner_master_store mpms\n" +
		"left join ms_partner mp on\n" +
		"mp.partner_id = mpms.partner_id\n" +
		"left join ms_initiator_fee mif on\n" +
		"mif.partner_id = mp.partner_id\n" +
		"left join ms_transaction_source_of_payment mtsop on\n" +
		"mtsop.initiator_id = mif.initiator_id\n" +
		"where\n" +
		"mpms.country_id = $1\n" +
		"and mpms.flag_active = true\n" +
		"and mp.flag_active = true\n" +
		"and mtsop.transaction_sop_id = $2\n" +
		"and mpms.partner_master_store_name ilike $3;\n"
	rows, err := query.Rows(stringQuery, countryId, sopiId, cari)
	// fmt.Println(stringQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	//mapping
	var datas []model.BankDest
	for rows.Next() {
		var data model.BankDest
		err := rows.Scan(&data.MasterStoreId, &data.BankName, &data.TypeDest)
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

func (rep *InternationalRepositoryImpl) GetCredential(masterStoreId int, sopId int) (*model.Credential, *model.ErrorData) {
	queries := "select\n" +
		"mpms.partner_master_store_name as bank_name,\n" +
		"mpms.partner_additional_data_1 as bank_code,\n" +
		"mpms.partner_id ,\n" +
		"mpms.country_id ,\n" +
		"mif.initiator_id,\n" +
		"mi.initiator_user ,\n" +
		"mi.initiator_key,\n" +
		"msa.admin_id,\n" +
		"mlc.country_code,\n" +
		"mi.initiator_adapter_address,\n" +
		"mlc2.city_id\n" +
		"from\n" +
		"ms_partner_master_store mpms\n" +
		"left join ms_initiator_fee mif on\n" +
		"mpms.partner_id = mif.partner_id\n" +
		"left join ms_transaction_source_of_payment mtsop on\n" +
		"mtsop.initiator_id = mif.initiator_id\n" +
		"left join ms_initiator mi on\n" +
		"mi.initiator_id = mtsop.initiator_id\n" +
		"left join ms_store_branch msb on\n" +
		"msb.initiator_id = mtsop.initiator_id\n" +
		"left join ms_store ms on\n" +
		"ms.branch_id = msb.branch_id\n" +
		"left join ms_store_admin msa on\n" +
		"msa.store_id = ms.store_id\n" +
		"left join ms_location_country mlc on\n" +
		"mlc.country_id = mpms.country_id\n" +
		"left join ms_location_province mlp on\n" +
		"mlp.country_id = mlc.country_id\n" +
		"left join ms_location_city mlc2 on\n" +
		"mlc2.province_id = mlp.province_id\n" +
		"where\n" +
		"mpms.partner_master_store_id = $1\n" +
		"and mtsop.transaction_sop_id = $2\n" +
		"and mlp.province_name = $3\n" +
		"order by\n" +
		"msa.admin_id asc\n" +
		"limit 1\n"
	row, err := query.Row(queries, masterStoreId, sopId, "other")
	if err != nil {
		return nil, &model.ErrorData{
			Description: "Db Error Credent " + err.Description,
		}
	}
	//mapping
	var dataAccount model.Credential
	errs := row.Scan(&dataAccount.BankName, &dataAccount.BankCode, &dataAccount.PartnerId,
		&dataAccount.CountryId, &dataAccount.InitiatorId, &dataAccount.InitiatorUser, &dataAccount.InitiatorKey,
		&dataAccount.AdminId, &dataAccount.CountryCode, &dataAccount.UrlAdapter, &dataAccount.CityReceiver)
	if errs != nil {
		return nil, &model.ErrorData{
			Description: errs.Error(),
		}
	}
	return &dataAccount, nil
}

func (rep *InternationalRepositoryImpl) DataAccountAll(phonePrefix string, phone string) (*model.DataAccountDetail, *model.ErrorData) {
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
		"msa.account_selfie_picture,\n"+
		"mssad.admin_id,\n"+
		"msa.account_selfie_identity_picture,\n"+
		"msa.account_signature_picture\n"+
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
		&dataAccount.OccupationName, &dataAccount.Address, &dataAccount.POB, &dataAccount.DOB, &dataAccount.Gender, &dataAccount.PostalCode, &dataAccount.ImgSelf,
		&dataAccount.AdminId, &dataAccount.ImgIdentity, &dataAccount.ImgSign)
	if errs != nil {
		return nil, &model.ErrorData{
			Description: errs.Error(),
		}
	}
	return &dataAccount, nil
}

func (rep *InternationalRepositoryImpl) GetIdentitas(id int) (*model.DataIdNkey, *model.ErrorData) {
	queries := "select\n" +
		"maitp.identity_type_id as id,\n" +
		"maitp.identity_type_partner_key as key\n" +
		"from\n" +
		"ms_account_identity_type_partner maitp\n" +
		"where\n" +
		"maitp.identity_type_partner_id = $1\n"
	row, err := query.Row(queries, id)
	if err != nil {
		return nil, err
	}
	//mapping
	var data model.DataIdNkey
	errs := row.Scan(&data.Id, &data.Key)
	if errs != nil {
		return nil, &model.ErrorData{
			Description: errs.Error(),
		}
	}
	return &data, nil
}

func (rep *InternationalRepositoryImpl) GetRelations(id int) (*model.DataIdNkey, *model.ErrorData) {
	queries := "select\n" +
		"mrp.relationship_id as id ,\n" +
		"mrp.relationship_partner_key as key\n" +
		"from\n" +
		"ms_relationship_partner mrp\n" +
		"where\n" +
		"mrp.relationship_partner_id = $1\n"
	row, err := query.Row(queries, id)
	if err != nil {
		return nil, err
	}
	//mapping
	var data model.DataIdNkey
	errs := row.Scan(&data.Id, &data.Key)
	if errs != nil {
		return nil, &model.ErrorData{
			Description: errs.Error(),
		}
	}
	fmt.Printf("%v", &data.Key)
	return &data, nil
}

func (rep *InternationalRepositoryImpl) GetPurposes(id int) (*model.DataIdNkey, *model.ErrorData) {
	queries := "select\n" +
		"mtpp.transaction_purpose_id as id,\n" +
		"mtpp.transaction_purpose_partner_key as key\n" +
		"from\n" +
		"ms_transaction_purpose_partner mtpp\n" +
		"where\n" +
		"mtpp.transaction_purpose_partner_id = $1\n"
	row, err := query.Row(queries, id)
	if err != nil {
		return nil, err
	}
	//mapping
	var data model.DataIdNkey
	errs := row.Scan(&data.Id, &data.Key)
	if errs != nil {
		return nil, &model.ErrorData{
			Description: errs.Error(),
		}
	}
	return &data, nil
}

func (rep *InternationalRepositoryImpl) GetFunding(id int) (*model.DataIdNkey, *model.ErrorData) {
	queries := "select\n" +
		"mtsofp.transaction_source_of_fund_id as id,\n" +
		"mtsofp.tsofp_key as key\n" +
		"from\n" +
		"ms_transaction_source_of_fund_partner mtsofp\n" +
		"where\n" +
		"mtsofp.tsofp_id = $1\n"
	row, err := query.Row(queries, id)
	if err != nil {
		return nil, err
	}
	//mapping
	var data model.DataIdNkey
	errs := row.Scan(&data.Id, &data.Key)
	if errs != nil {
		return nil, &model.ErrorData{
			Description: errs.Error(),
		}
	}
	return &data, nil
}

func (rep *InternationalRepositoryImpl) GetOccupation(id int) (*model.DataIdNkey, *model.ErrorData) {
	queries := "select\n" +
		"maop.occupation_id as id,\n" +
		"maop.occupation_partner_key as key\n" +
		"from\n" +
		"ms_account_occupation_partner maop\n" +
		"where\n" +
		"maop.occupation_partner_id = $1\n"
	row, err := query.Row(queries, id)
	if err != nil {
		return nil, err
	}
	//mapping
	var data model.DataIdNkey
	errs := row.Scan(&data.Id, &data.Key)
	if errs != nil {
		return nil, &model.ErrorData{
			Description: errs.Error(),
		}
	}
	return &data, nil
}
