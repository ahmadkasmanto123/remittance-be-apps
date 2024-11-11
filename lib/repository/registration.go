package repository

import (
	"love-remittance-be-apps/core/query"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
)

type RegistrationRepositoryImpl struct{}

func NewRegistrationRepository() interfc.RegistrationRepository {
	return &RegistrationRepositoryImpl{}
}

func (rep *RegistrationRepositoryImpl) CountAccount(phone string, phonePrefix string) (*model.CountSomething, *model.ErrorData) {
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

func (rep *RegistrationRepositoryImpl) CountEmail(email string) (*model.CountSomething, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.Row("SELECT COUNT(*) \n"+
		"FROM\n"+
		"	ms_store_admin	mssa \n"+
		"WHERE\n"+
		"	mssa.admin_email = $1", email)
	if err != nil {
		return nil, err
	}
	//mapping
	var countSomething model.CountSomething
	row.Scan(&countSomething.Count)
	return &countSomething, nil
}

func (rep *RegistrationRepositoryImpl) CountMsAccount(phone string) (*model.CountSomething, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.Row("SELECT COUNT(*) \n"+
		"FROM\n"+
		"	ms_account msa \n"+
		"WHERE\n"+
		"	msa.account_phone_number = $1 ", phone)
	if err != nil {
		return nil, err
	}
	//mapping
	var countSomething model.CountSomething
	row.Scan(&countSomething.Count)
	return &countSomething, nil
}

func (rep *RegistrationRepositoryImpl) Storecheck(phonePrefix string) (*model.StoreCheck, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.Row("SELECT\n"+
		"	msst.store_id, \n"+
		"	mslct.country_id\n"+
		"FROM\n"+
		"	ms_store msst \n"+
		"	LEFT JOIN ms_store_branch msstbr on msst.branch_id = msstbr.branch_id\n"+
		"	LEFT JOIN ms_initiator msinit on msstbr.initiator_id = msinit.initiator_id\n"+
		"	LEFT JOIN ms_location_country mslct on msinit.country_id = mslct.country_id\n"+
		"WHERE\n"+
		"	msst.store_name = 'LOVE APPS' \n"+
		"	AND msst.flag_active = TRUE\n"+
		"	AND mslct.country_prefix =  $1", phonePrefix)
	if err != nil {
		return nil, err
	}
	//mapping
	var data model.StoreCheck
	row.Scan(&data.StoreId, &data.CountryId)
	return &data, nil
}

func (rep *RegistrationRepositoryImpl) UpdateMsAccount(userName string, phone string) (*int, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.Row("update\n"+
		"ms_account\n"+
		"set\n"+
		"account_name = $1 ,\n"+
		"timestamp_update = now() ,\n"+
		"account_status_id = 2 \n"+
		"where\n"+
		"account_phone_number = $2 returning account_id ", userName, phone)
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

func (rep *RegistrationRepositoryImpl) InsertMsAccount(userName string, phone string) (*int, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.Row("insert\n"+
		"into\n"+
		"ms_account ( account_name,\n"+
		"account_phone_number,\n"+
		"timestamp_insert,\n"+
		"account_status_id )\n"+
		"values ($1,\n"+
		"$2,\n"+
		"now(),\n"+
		"2) RETURNING account_id", userName, phone)
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

func (rep *RegistrationRepositoryImpl) InsertMsStoreAdmin(req model.DefaultRequest[model.CreateCustomer], storeId int, accountId int) (*int, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.Row("insert\n"+
		"into\n"+
		"ms_store_admin (\n"+
		"admin_username,\n"+
		"admin_pass,\n"+
		"store_id,\n"+
		"admin_pin,\n"+
		"admin_counter_pass,\n"+
		"admin_email,\n"+
		"account_id,\n"+
		"device_id )\n"+
		"values\n"+
		"( $1,\n"+
		"$2,\n"+
		"$3,\n"+
		"$4,\n"+
		"0,\n"+
		"$5,\n"+
		"$6,\n"+
		"$7) returning admin_id\n", req.Request.FirstName+" "+req.Request.LastName, req.Request.Pin,
		storeId, req.Request.Pin, req.Request.Email, accountId, "")
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

func (rep *RegistrationRepositoryImpl) DataAccountReg(phone string, phonePrefix string) (*model.DataAccount, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.Row("SELECT\n"+
		"msa.account_id,\n"+
		"msa.account_status_id,\n"+
		"msa.account_name,\n"+
		"mssad.admin_email,\n"+
		"mssad.admin_pin,\n"+
		"mssad.device_id,\n"+
		"mssad.admin_id,\n"+
		"mssad.admin_conf :: json ->> 'device_id' as admin_conf\n"+
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
		"WHERE\n"+
		"msa.account_phone_number = $1\n"+
		"AND mslc.country_prefix = $2;\n", phonePrefix+phone, phonePrefix)
	if err != nil {
		return nil, err
	}
	//mapping
	var dataAccount model.DataAccount
	row.Scan(&dataAccount.AccountId, &dataAccount.AccountStatutId, &dataAccount.AccountName, &dataAccount.AdminEmail, &dataAccount.AdminPin, &dataAccount.DeviceId, &dataAccount.AdminId, &dataAccount.AdminConf)
	return &dataAccount, nil
}

func (rep *RegistrationRepositoryImpl) UpdateMsStoreAdmin(deviceId string, adminId int) *model.ErrorData {
	_, err := query.ExecUpdate("UPDATE ms_store_admin \n"+
		"SET device_id = $1 \n"+
		"WHERE\n"+
		"	\"admin_id\" = $2;", deviceId, adminId)
	if err != nil {
		return err
	}
	return nil
}
