package repository

import (
	"love-remittance-be-apps/core/query"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
)

type UpdateRepositoryImpl struct{}

func NewUpdateRepository() interfc.UpdateRepository {
	return &UpdateRepositoryImpl{}
}

func (rep *UpdateRepositoryImpl) CountAccount(phone string, phonePrefix string) (*model.CountSomething, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.Row("SELECT COUNT(*) \n"+
		"FROM\n"+
		"	ms_account msa\n"+
		"	JOIN ms_store_admin mssad ON msa.account_id = mssad.account_id\n"+
		"	LEFT JOIN ms_store msst ON mssad.store_id = msst.store_id\n"+
		"	LEFT JOIN ms_store_branch mssbr ON msst.branch_id = mssbr.branch_id\n"+
		"	LEFT JOIN ms_initiator msin ON mssbr.initiator_id = msin.initiator_id\n"+
		"	LEFT JOIN ms_location_country mslc ON msin.country_id = mslc.country_id \n"+
		"WHERE\n"+
		"	msa.account_phone_number = $1 \n"+
		"	and mslc.country_prefix = $2 ;", phonePrefix+phone, phonePrefix)
	if err != nil {
		return nil, err
	}
	//mapping
	var countSomething model.CountSomething
	row.Scan(&countSomething.Count)
	return &countSomething, nil
}

func (rep *UpdateRepositoryImpl) UpdateMsAccount(req model.DefaultRequest[model.UpdateProfil], img []string) (*int, *model.ErrorData) {
	// dob, _ := time.Parse(time.RFC3339, req.Request.DOB)
	// expired, _ := time.Parse(time.RFC3339, req.Request.IdentityExpired)
	row, err := query.Row("UPDATE \"public\".\"ms_account\"\n"+
		"SET \"account_status_id\" = $1,\n"+
		"\"city_id\" = $2,\n"+
		"\"account_address\" = $3 ,\n"+
		"\"account_name\" = $4 ,\n"+
		"\"account_identity_type_id\" = $5,\n"+
		"\"account_identity_number\" = $6 ,\n"+
		"\"account_gender\" = $7 ,\n"+
		"\"account_dob\" = $8,\n"+
		"\"account_pob\" = $9 ,\n"+
		"\"occupation_id\" = $10,\n"+
		"\"timestamp_update\" = now(),\n"+
		"\"account_identity_expire_date\" = $11,\n"+
		"\"account_postal_code\" = $12,\n"+
		"\"account_selfie_picture\" = $13,\n"+
		"\"account_signature_picture\" = $14,\n"+
		"\"account_selfie_identity_picture\" = $15\n"+
		"WHERE\n"+
		"\"account_phone_number\" = $16 returning account_id\n", 5, req.Request.City,
		req.Request.Address,
		req.Request.FirstName+" "+req.Request.LastName,
		req.Request.IdentityTypeId, req.Request.IdentityNumber, req.Request.Gender,
		req.Request.DOB, req.Request.POB, req.Request.Occupation, req.Request.IdentityExpired,
		req.Request.PostalCode, img[1], img[0], img[2], req.Request.PhonePrefix+req.Request.Phone)
	if err != nil {
		return nil, err
	}
	//mapping
	var countSomething *int
	errs := row.Scan(&countSomething)
	if errs != nil {
		return nil, &model.ErrorData{
			Description: errs.Error(),
		}
	}
	return countSomething, nil
}

func (rep *UpdateRepositoryImpl) GetDataaaa() ([]model.DataSOPayment, *model.ErrorData) {
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
