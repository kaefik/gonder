package api

import (
	"github.com/supme/gonder/models"
	"encoding/json"
	"fmt"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"database/sql"
)

func users(req request) (js []byte, err error) {
	type UserList struct {
		Id   int64  `json:"recid"`
		UnitId int64 `json:"unitid"`
		GroupsId []int64 `json:"groupsid"`
		Name string `json:"name"`
		Password string `json:"password"`
	}

	type Users struct{
		Total int64 `json:"total"`
		Records []UserList `json:"records"`
	}

	if auth.IsAdmin() {
	switch req.Cmd {
		case "get":
			var (
				sl = Users{}
				id, unitid, grpid int64
				name string
				partWhere, where string
				partParams, params []interface{}
			)
			sl.Records = []UserList{}

			where = " WHERE 1=1 "
			partWhere, partParams, err = createSqlPart(req, where, params, map[string]string{"recid":"id", "name":"name"}, true)
			if err != nil {
				return nil, err
			}
			query, err := models.Db.Query("SELECT `id`, `auth_unit_id`, `name` FROM `auth_user`" + partWhere, partParams...)
			if err != nil {
				return nil, err
			}
			defer query.Close()
			for query.Next() {
				err = query.Scan(&id, &unitid, &name)
				groupq, err := models.Db.Query("SELECT `group_id` FROM `auth_user_group` WHERE `auth_user_id`=?", id)
				if err != nil {
					return nil, err
				}
				defer groupq.Close()
				groupsid := []int64{}
				for groupq.Next() {
					groupq.Scan(&grpid)
					groupsid = append(groupsid, grpid)
				}
				sl.Records = append(sl.Records, UserList{
					Id:   id,
					UnitId: unitid,
					GroupsId: groupsid,
					Name: name,
				})
			}

			partWhere, partParams, err = createSqlPart(req, where, params, map[string]string{"recid":"id", "name":"name"}, false)
			if err != nil {
				apilog.Print(err)
			}
			err = models.Db.QueryRow("SELECT COUNT(*) FROM `auth_user` " + partWhere, partParams...).Scan(&sl.Total)
			if err != nil {
				return nil, err
			}
			js, err = json.Marshal(sl)
			return js, err

		case "save":
			// Change user
			if req.Record.Password != "" {
				hash := sha256.New()
				hash.Write([]byte(fmt.Sprint(req.Record.Password)))
				md := hash.Sum(nil)
				shaPassword := hex.EncodeToString(md)
				_, err = models.Db.Exec("UPDATE `auth_user` SET `password`=?, auth_unit_id=? WHERE `id`=?", shaPassword, req.Record.Unit.Id, req.Record.Id)
				if err != nil {
					return nil, err
				}
			} else {
				_, err = models.Db.Exec("UPDATE `auth_user` SET auth_unit_id=? WHERE `id`=?", req.Record.Unit.Id, req.Record.Id)
				if err != nil {
					return nil, err
				}
			}
			_, err = models.Db.Exec("DELETE FROM `auth_user_group` WHERE `auth_user_id`=?", req.Record.Id)
			if err != nil {
				return nil, err
			}
			for _, g := range req.Record.Group {
				_, err = models.Db.Exec("INSERT INTO `auth_user_group` (`auth_user_id`, `group_id`) VALUES (?, ?)", req.Record.Id, g.Id)
				if err != nil {
					return nil, err
				}
			}

		case "add":
			// Add user
			row := models.Db.QueryRow("SELECT 1 FROM `auth_user` WHERE `name`=?", req.Record.Name).Scan()
			if row == sql.ErrNoRows {
				fmt.Print(req.Record)
				if req.Record.Password == "" {
					return nil, errors.New("Password can not be empty")
				}
				if req.Record.Name == "" {
					return nil, errors.New("Name can not be empty")
				}
				hash := sha256.New()
				hash.Write([]byte(fmt.Sprint(req.Record.Password)))
				md := hash.Sum(nil)
				shaPassword := hex.EncodeToString(md)
				s, err := models.Db.Exec("INSERT INTO `auth_user`(`auth_unit_id`, `name`, `password`) VALUES (?, ?, ?)", req.Record.Unit.Id, req.Record.Name, shaPassword)
				if err != nil {
					return nil, err
				}
				req.Record.Id, err = s.LastInsertId()
				if err != nil {
					return nil, err
				}
				for _, g := range req.Record.Group {
					_, err = models.Db.Exec("INSERT INTO `auth_user_group` (`auth_user_id`, `group_id`) VALUES (?, ?)", req.Record.Id, g.Id)
					if err != nil {
						return nil, err
					}
				}
			} else {
				return nil, errors.New("User name already exist")
			}

		default:
			err = errors.New("Command not found")

		}
	} else {
		return nil, errors.New("Access denied")
	}

	return js, err

}