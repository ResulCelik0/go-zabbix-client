package gozabbix

import (
	"encoding/json"
)

type UserAPI struct {
	*ZabbixClient
}

// ////////////////////// OBJECTS ////////////////////////
// User Object Zabbix ref: https://www.zabbix.com/documentation/current/en/manual/api/reference/user/object
type User struct {
	UserID          string `json:"userid,omitempty"`
	Username        string `json:"username,omitempty"`
	Passwd          string `json:"passwd,omitempty"`
	RoleID          string `json:"roleid,omitempty"`
	AttemptClock    string `json:"attempt_clock,omitempty"`
	AttemptFailed   string `json:"attempt_failed,omitempty"`
	AttemptIp       string `json:"attempt_ip,omitempty"`
	Autologin       string `json:"autologin,omitempty"`
	Autologout      string `json:"autologout,omitempty"`
	Lang            string `json:"lang,omitempty"`
	Name            string `json:"name,omitempty"`
	Refresh         string `json:"refresh,omitempty"`
	RowsPerPage     string `json:"rows_per_page,omitempty"`
	Surname         string `json:"surname,omitempty"`
	Theme           string `json:"theme,omitempty"`
	TsProvisioned   string `json:"ts_provisioned,omitempty"`
	URL             string `json:"url,omitempty"`
	UserDirectoryID string `json:"userdirectoryid,omitempty"`
	Timezone        string `json:"timezone,omitempty"`
}

// Media Object Zabbix ref: https://www.zabbix.com/documentation/current/en/manual/api/reference/user/object
type Media struct {
	MediaTypeID          string `json:"mediatypeid,omitempty"`
	SendTo               string `json:"sendto,omitempty"`
	Active               string `json:"active,omitempty"`
	Severity             string `json:"severity,omitempty"`
	Period               string `json:"period,omitempty"`
	UserDirectoryMediaID string `json:"userdirectory_mediaid,omitempty"`
}

//////////////////////// FUNCTIONS ////////////////////////

// Login Function Zabbix ref: https://www.zabbix.com/documentation/current/en/manual/api/reference/user/login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserData bool   `json:"userData,omitempty"`
}

type LoginResponse struct {
	Type          int    `json:"type"`
	UserIP        string `json:"userip"`
	DebugMode     int    `json:"debug_mode"`
	GuiAccess     string `json:"gui_access"`
	MfaID         int    `json:"mfaid"`
	Deprovisioned bool   `json:"deprovisioned"`
	AuthType      int    `json:"auth_type"`
	SessionID     string `json:"sessionid"`
	Secret        string `json:"secret"`
	User
}

func (r *LoginResponse) Unmarshal(data interface{}) error {
	switch v := data.(type) {
	case map[string]interface{}:
		b, err := json.Marshal(data)
		if err != nil {
			return err
		}
		return json.Unmarshal(b, r)
	case string:
		r.SessionID = v
	}
	return nil
}

func (u *UserAPI) Login(loginReq *LoginRequest, reqId ...int) (token string, loginResp *LoginResponse, err error) {
	if len(reqId) == 0 {
		reqId = append(reqId, 1)
	}
	req := &Request{
		Method: "user.login",
		Params: loginReq,
		Id:     reqId[0],
	}
	resp, err := u.execute(req)
	if err != nil {
		return
	}
	loginResp = &LoginResponse{}
	err = loginResp.Unmarshal(resp.Result)
	if err != nil {
		return
	}
	token = loginResp.SessionID
	return
}
