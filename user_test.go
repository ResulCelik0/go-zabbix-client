package gozabbix

import "testing"

func TestLogin(t *testing.T) {
	client, userInfo, err := NewZabbixClient(&Config{
		URL:      "http://localhost:8080/api_jsonrpc.php",
		Username: "Admin",
		Password: "zabbix",
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(userInfo)
	t.Log(client)
}
