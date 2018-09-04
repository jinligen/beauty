// MIT License

// Copyright (c) 2017 FLYING

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package decorates

import (
	"net/http"
	"strings"

	"github.com/yang-f/beauty/db"
	"github.com/yang-f/beauty/models"
	"github.com/yang-f/beauty/utils/log"
	"github.com/yang-f/beauty/utils/token"
)

func (inner Handler) Auth() Handler {
	return Handler(func(w http.ResponseWriter, r *http.Request) *models.APPError {
		var tokenString string
		if cookie, _ := r.Cookie("token"); cookie != nil {
			tokenString = cookie.Value
		}
		if tokenString == "" {
			tokenString = r.Header.Get("Authorization")
		}

		if tokenString == "" {
			return &models.APPError{
				Error:   nil,
				Message: "token not found.",
				Code:    "AUTH_FAILED",
				Status:  403,
			}
		}

		key, err := token.Valid(tokenString)
		if err != nil {
			return &models.APPError{
				Error:   err,
				Message: "bad token.",
				Code:    "AUTH_FAILED",
				Status:  403,
			}
		}

		keys := strings.Split(key, "|")
		if len(keys) != 2 {
			return &models.APPError{
				Error:   err,
				Message: "bad token.",
				Code:    "AUTH_FAILED",
				Status:  403,
			}
		}

		userID, userPass := keys[0], keys[1]
		rows, _, err := db.QueryNonLogging("select * from user where user_id = '%v' and user_pass = '%v'", userID, userPass)
		if err != nil {
			return &models.APPError{
				Error:   err,
				Message: "can not connect database.",
				Code:    "DB_ERROR",
				Status:  500,
			}
		}
		if len(rows) == 0 {
			return &models.APPError{
				Error:   err,
				Message: "user not found.",
				Code:    "NOT_FOUND",
				Status:  404,
			}
		}
		go log.Printf("user_id:%v", userID)
		inner.ServeHTTP(w, r)
		return nil
	})
}
