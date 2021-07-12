/*
   ZAU Job - Users Sync
   Copyright (C) 2021  Daniel A. Hawton <daniel@hawton.org>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/dhawton/log4g"
	"github.com/joho/godotenv"
	"github.com/vchicago/common/utils"
	"github.com/vchicago/job-user-sync/db"
	dbTypes "github.com/vchicago/types/database"
	"gorm.io/gorm"
)

var log = log4g.Category("main")

func main() {
	intro := figure.NewFigure("ZAU Job - Users", "", false).Slicify()
	for i := 0; i < len(intro); i++ {
		log.Info(intro[i])
	}

	log.Info("-------------------------")
	log.Info("Checking for .env, loading if exists")
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			log.Error("Could not load .env: %s", err.Error())
		}
	}

	if utils.Getenv("APP_ENV", "prod") != "dev" && utils.Getenv("APP_DEBUG", "false") != "true" {
		log.Info("Setting log level to info")
		log4g.SetLogLevel(log4g.INFO)
	} else {
		log.Info("Setting log level to debug")
		log4g.SetLogLevel(log4g.DEBUG)
	}

	log.Info("Building database connection")
	db.Connect(utils.Getenv("DB_USERNAME", "root"), utils.Getenv("DB_PASSWORD", "secret12345"), utils.Getenv("DB_HOSTNAME", "localhost"), utils.Getenv("DB_PORT", "3306"), utils.Getenv("DB_DATABASE", "zau"))

	log.Info("Getting ZAU Roster from VATUSA")
	start := time.Now()
	response, err := http.Get(fmt.Sprintf("https://api.vatusa.net/v2/facility/ZAU/roster/both?apikey=%s", utils.Getenv("VATUSA_API_KEY", "")))
	if err != nil {
		log.Fatal("Error querying VATUSA API: %s", err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Fatal("Error reading response body: %s", err.Error())
	}

	log.Info("Request completed in %s", time.Now().Sub(start))
	log.Info("Unmarshalling data.")

	data := &VATUSAReturn{}
	if err := json.Unmarshal([]byte(responseData), &data); err != nil {
		log.Fatal("Could not unmarshal data from VATUSA: %s", err.Error())
	}

	log.Info("Processing data")
	start = time.Now()
	for i := 0; i < len(data.Controllers); i++ {
		controller := data.Controllers[i]
		rec := dbTypes.User{}
		isNew := false

		if err := db.DB.Where(&dbTypes.User{CID: uint(controller.CID)}).First(&rec).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Error("Error looking up user %d, %s", controller.CID, err.Error())
				return
			} else {
				isNew = true
			}
		}

		rec.FirstName = controller.FirstName
		rec.LastName = controller.LastName
		rec.Email = controller.Email
		if controller.Membership == "visit" {
			rec.ControllerType = "visitor"
		} else {
			rec.ControllerType = controller.Membership
		}
		rec.RatingID = controller.RatingId

		if isNew {
			rec.Status = "active"
		}

		if err := db.DB.Save(&rec).Error; err != nil {
			log.Error("Error saving controller, %d to database: %s", controller.CID, err.Error())
		}
	}

	log.Info("Job completed in %s", time.Now().Sub(start))
}
