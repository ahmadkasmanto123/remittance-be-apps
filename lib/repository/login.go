package repository

import (
	"love-remittance-be-apps/core/query"
	"love-remittance-be-apps/lib/interfc"
	"love-remittance-be-apps/lib/model"
)

type LoginRepositoryImpl struct{}

func NewLoginRepository() interfc.LoginRepository {
	return &LoginRepositoryImpl{}
}

func (rep *LoginRepositoryImpl) CountAccount(phone string, phonePrefix string) (*model.CountSomething, *model.ErrorData) {
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

func (rep *LoginRepositoryImpl) DataAccount(phone string, phonePrefix string) (*model.DataAccount, *model.ErrorData) {
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

func (rep *LoginRepositoryImpl) UpdateAccountStatus(phone string, statusId int) *model.ErrorData {
	_, err := query.ExecUpdate("UPDATE ms_account \n"+
		"SET account_status_id = $2 \n"+
		"WHERE\n"+
		"	\"account_phone_number\" = $1;", phone, statusId)
	if err != nil {
		return err
	}
	return nil
}

func (rep *LoginRepositoryImpl) UpdateMsStoreAdmin(deviceId string, adminId int) *model.ErrorData {
	_, err := query.ExecUpdate("UPDATE ms_store_admin \n"+
		"SET device_id = $1 \n"+
		"WHERE\n"+
		"	\"admin_id\" = $2;", deviceId, adminId)
	if err != nil {
		return err
	}
	return nil
}
