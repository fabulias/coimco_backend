package model

import (
	"log"
	"testing"
)

func TestInsertOneClient(t *testing.T) {
	cli := &Client{id: 0, name: "nameTest", phone: "+9101010"}
	log.Println(cli)
	_, flag := InsertCustomers(cli)
	if !flag {
		t.Logf("Error insert 1 Client")
	}
}
