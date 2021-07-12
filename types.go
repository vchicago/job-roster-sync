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

type VATUSAReturn struct {
	Controllers []DTOController `json:"data"`
}

type DTOController struct {
	CID        int    `json:"cid"`
	FirstName  string `json:"fname"`
	LastName   string `json:"lname"`
	Email      string `json:"email"`
	RatingId   int    `json:"rating"`
	Membership string `json:"membership"`
}
