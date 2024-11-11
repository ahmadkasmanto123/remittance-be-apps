package repository

import (
	"love-remittance-be-apps/core/query"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
	"strconv"
)

type DomesticRepositoryImpl struct{}

func NewDomesticRepository() interfc.DomesticRepository {
	return &DomesticRepositoryImpl{}
}

func (rep *DomesticRepositoryImpl) GetSourceOfPayment() ([]model.DataSOPayment, *model.ErrorData) {
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

func (rep *DomesticRepositoryImpl) GetDestination(param model.Param, sopId int) ([]model.DataDestination, *model.Pagination, *model.ErrorData) {
	//query get data payment method in database
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

	//mapping
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

func (rep *DomesticRepositoryImpl) PartnerMasterStoreId(masterStoreId int) (*model.MasterStore, *model.ErrorData) {
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
	row.Scan(&masterStore.BankName, &masterStore.BankCode, &masterStore.PartnerId, &masterStore.CountryId)
	return &masterStore, nil
}

func (rep *DomesticRepositoryImpl) DataInitiator(sopId int) (*model.DataSOPayment, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.Row("SELECT\n"+
		"msop.sop_name,\n"+
		"msop.sop_type,\n"+
		"mi.initiator_id,\n"+
		"mi.initiator_adapter_address\n"+
		"FROM\n"+
		"ms_source_of_payment as msop\n"+
		"LEFT JOIN ms_transaction_source_of_payment as mstsop ON msop.sop_id = mstsop.sop_id\n"+
		"LEFT JOIN ms_initiator as mi ON mi.initiator_id = mstsop.initiator_id\n"+
		"LEFT JOIN ms_location_country mslc ON mi.country_id = mslc.country_id\n"+
		"WHERE\n"+
		"mstsop.transaction_sop_flag = TRUE\n"+
		"and mstsop.transaction_sop_id = $1;", sopId)
	if err != nil {
		return nil, err
	}
	//mapping
	var data model.DataSOPayment
	row.Scan(&data.SopName, &data.SopType, &data.InitiatorId, &data.InitiatorAdapAddr)
	return &data, nil
}

func (rep *DomesticRepositoryImpl) DataAccountAll(phonePrefix string, phone string) (*model.DataAccountDetail, *model.ErrorData) {
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
