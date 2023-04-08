package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/padhitheboss/key-value-db/pkg/model"
)

func CommandHandler(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Command string `json:"command"`
	}{}
	json.NewDecoder(r.Body).Decode(&req)
	command := req.Command
	commandAttrib := strings.Split(command, " ")
	for i := 0; i < len(commandAttrib); i++ {
		commandAttrib[i] = strings.Trim(commandAttrib[i], " ")
		if commandAttrib[i] == "" {
			commandAttrib = append(commandAttrib[:i], commandAttrib[i+1:]...)
		}
	}
	switch strings.ToUpper(commandAttrib[0]) {
	case "SET":
		if len(commandAttrib) > 6 || len(commandAttrib) < 3 {
			json.NewEncoder(w).Encode(model.Response{Error: "invalid number of arguments."})
			return
		}
		key := commandAttrib[1]
		value := commandAttrib[2]
		var expiry int64 = 0
		var err error
		var condition string
		if len(commandAttrib) > 3 && strings.ToUpper(commandAttrib[3]) == "EX" {
			expiry, err = strconv.ParseInt(commandAttrib[4], 10, 64)
			if err != nil {
				json.NewEncoder(w).Encode(model.Response{Error: "invalid expiry time."})
				return
			}
			if len(commandAttrib) == 6 && (strings.ToUpper(commandAttrib[5]) == "XX" || strings.ToUpper(commandAttrib[5]) == "NX") {
				condition = strings.ToUpper(commandAttrib[5])
			}
		} else if len(commandAttrib) > 3 && (strings.ToUpper(commandAttrib[3]) == "XX" || strings.ToUpper(commandAttrib[3]) == "NX") {
			condition = strings.ToUpper(commandAttrib[3])
		}
		var expTime time.Time
		if expiry != 0 {
			expTime = time.Now().Add(time.Duration(time.Second * time.Duration(expiry)))
		}
		status, err := model.DB.Set(key, value, expTime, condition)
		if err != nil {
			json.NewEncoder(w).Encode(model.Response{Error: fmt.Sprint(err)})
			return
		}
		json.NewEncoder(w).Encode(model.Response{Status: status})
	case "GET":
		if len(commandAttrib) != 2 {
			json.NewEncoder(w).Encode(model.Response{Error: "invalid number of arguments."})
			return
		}
		value, err := model.DB.Get(commandAttrib[1])
		if err != nil {
			json.NewEncoder(w).Encode(model.Response{Error: fmt.Sprint(err)})
			return
		}
		json.NewEncoder(w).Encode(model.Response{Value: value})
	case "QPUSH":
		if len(commandAttrib) < 3 {
			json.NewEncoder(w).Encode(model.Response{Error: "invalid number of arguments."})
			return
		}
		status, err := model.DB.QPush(commandAttrib[1], commandAttrib[2:])
		if err != nil {
			json.NewEncoder(w).Encode(model.Response{Error: fmt.Sprint(err)})
			return
		}
		json.NewEncoder(w).Encode(model.Response{Status: status})
	case "QPOP":
		if len(commandAttrib) != 2 {
			json.NewEncoder(w).Encode(model.Response{Error: "invalid number of arguments."})
			return
		}
		value, err := model.DB.QPop(commandAttrib[1])
		if err != nil {
			json.NewEncoder(w).Encode(model.Response{Error: fmt.Sprint(err)})
			return
		}
		json.NewEncoder(w).Encode(model.Response{Value: value})
	case "BQPOP":
		if len(commandAttrib) != 3 {
			json.NewEncoder(w).Encode(model.Response{Error: "invalid number of arguments."})
			return
		}
		_, err := strconv.ParseFloat(commandAttrib[2], 64)
		if err != nil {
			json.NewEncoder(w).Encode(model.Response{Error: "invalid timeout value."})
		}
		value, err := model.DB.QPop(commandAttrib[1])
		if err != nil {
			json.NewEncoder(w).Encode(model.Response{Error: fmt.Sprint(err)})
			return
		}
		json.NewEncoder(w).Encode(model.Response{Value: value})
	default:
		json.NewEncoder(w).Encode(model.Response{Error: "invalid commands."})
	}
}
