package rc

import (
	"love-remittance-be-apps/core/query"
	"love-remittance-be-apps/lib/model"
)

type RCodeData struct {
	RC      string `json:"rc"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type RCData struct{}

func GetResponseMsg(rc string, types int, lang string) (*RCodeData, *model.ErrorData) {
	//query get data payment method in database
	row, err := query.RowRC("select\n"+
		"rc_code,\n"+
		"rc_status,\n"+
		"case\n"+
		"when 'id' = $1 then rc_message_id\n"+
		"else rc_message_en\n"+
		"end message\n"+
		"from\n"+
		"ms_rc_message\n"+
		"where\n"+
		"rc_code = $2\n"+
		"and rc_type_id = $3;\n", lang, rc, types)
	if err != nil {
		return nil, err
	}
	//mapping
	var datas RCodeData
	row.Scan(&datas.RC, &datas.Status, &datas.Message)
	return &datas, nil
}
