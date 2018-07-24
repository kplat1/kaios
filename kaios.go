// Copyright (c) 2018, The KaiOS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	//"reflect"

	"github.com/goki/gi"
	"github.com/goki/gi/oswin"
	"github.com/goki/gi/oswin/driver"
	//"github.com/goki/gi/units"
	"github.com/goki/ki"
	//"github.com/goki/ki/kit"

	bolt "github.com/coreos/bbolt"
	"log"
	"time"
)

type LoginRec struct {
	Username string
	Password string
	Points   float64
}

var KaiOSDB *bolt.DB

func LoadLoginTable() []*LoginRec {

	lt := make([]*LoginRec, 0, 100) // 100 is the starting capacity of slice -- increase if you expect more users.

	KaiOSDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("LoginTable"))

		if b != nil {
			b.ForEach(func(k, v []byte) error {
				if v != nil {
					rec := LoginRec{}
					json.Unmarshal(v, &rec) // loads v value as json into rec
					lt = append(lt, &rec)   // adds record to login table

				}
				return nil
			})
		}
		return nil
	})

	return lt
}

func SaveNewLogin(rec *LoginRec) {
	KaiOSDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("LoginTable"))
		jb, err := json.Marshal(rec) // converts rec to json, as bytes jb

		err = b.Put([]byte(rec.Username), jb)
		return err
	})
}

func main() {
	var err error
	KaiOSDB, err = bolt.Open("KaiOS.db", 0600, nil)

	if err != nil {
		log.Fatal(err)
	}
	defer KaiOSDB.Close()

	driver.Main(func(app oswin.App) {
		mainrun()
	})
}

func mainrun() {
	width := 1024
	height := 768

	// turn these on to see a traces of various stages of processing..
	// gi.Update2DTrace = true
	// gi.Render2DTrace = true
	// gi.Layout2DTrace = true
	// ki.SignalTrace = true

	rec := ki.Node{}          // receiver for events
	rec.InitName(&rec, "rec") // this is essential for root objects not owned by other Ki tree nodes

	win := gi.NewWindow2D("gogi-widgets-demo", "KaiOS", width, height, true) // true = pixel sizes

	//icnm := "widget-wedge-down"

	vp := win.WinViewport2D()
	updt := vp.UpdateStart()
	vp.Fill = true

	// style sheet
	var css = ki.Props{
		"button": ki.Props{
			"background-color": gi.Color{255, 240, 240, 255},
		},
		"#combo": ki.Props{
			"background-color": gi.Color{240, 255, 240, 255},
		},
		".hslides": ki.Props{
			"background-color": gi.Color{240, 225, 255, 255},
		},
		"kbd": ki.Props{
			"color": "blue",
		},
	}
	vp.CSS = css

	vlay := vp.AddNewChild(gi.KiT_Frame, "vlay").(*gi.Frame)
	vlay.Lay = gi.LayoutCol
	// vlay.SetProp("background-color", "linear-gradient(to top, red, lighter-80)")
	// vlay.SetProp("background-color", "linear-gradient(to right, red, orange, yellow, green, blue, indigo, violet)")
	// vlay.SetProp("background-color", "linear-gradient(to right, rgba(255,0,0,0), rgba(255,0,0,1))")
	// vlay.SetProp("background-color", "radial-gradient(red, lighter-80)")

	trow := vlay.AddNewChild(gi.KiT_Layout, "trow").(*gi.Layout)
	trow.Lay = gi.LayoutCol
	trow.SetStretchMaxWidth()

	trow.AddNewChild(gi.KiT_Stretch, "str1")
	title := trow.AddNewChild(gi.KiT_Label, "title").(*gi.Label)
	title.Text =
		`<b>KaiOS</b>`
	title.SetProp("align-horiz", gi.AlignCenter)
	title.SetProp("align-vert", gi.AlignTop)
	title.SetProp("font-family", "Times New Roman, serif")
	title.SetProp("font-size", "x-large")
	// title.SetProp("letter-spacing", 2)
	title.SetProp("line-height", 1.5)
	trow.AddNewChild(gi.KiT_Stretch, "str2")

	p1 := trow.AddNewChild(gi.KiT_Label, "p1").(*gi.Label)
	p1.Text = "<b>KaiOS</b>, a <b>customizable, lightweight</b> OS"
	p1.SetProp("align-horiz", gi.AlignCenter)

	trow.AddNewChild(gi.KiT_Space, "spc1")
	buttonStart := trow.AddNewChild(gi.KiT_Button, "buttonStart").(*gi.Button)
	buttonStart.Text = "Load KaiOS v0.000 pre-alpha"
	buttonStart.SetProp("align-horiz", gi.AlignCenter)

	buttonStart.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		//fmt.Printf("Received button signal: %v from button: %v\n", gi.ButtonSignals(sig), send.Name())
		if sig == int64(gi.ButtonClicked) { // note: 3 diff ButtonSig sig's possible -- important to check
			// vp.Win.Quit()
			//gi.PromptDialog(vp, "buttonStart Dialog", "This is a dialog!  Various specific types of dialogs are available.", true, true, nil, nil)
			updt := vp.UpdateStart()

			buttonStartResult := trow.AddNewChild(gi.KiT_Label, "buttonStartResult").(*gi.Label)
			buttonStartResult.Text = "<b>Sign Up:</b>"
			userText := trow.AddNewChild(gi.KiT_TextField, "userText").(*gi.TextField)
			userText.SetText("Username")
			userText.SetProp("width", "20em")
			passwdText := trow.AddNewChild(gi.KiT_TextField, "passwdText").(*gi.TextField)
			passwdText.SetText("Password")
			passwdText.SetProp("width", "20em")

			signUpButton := trow.AddNewChild(gi.KiT_Button, "signUpButton").(*gi.Button)
			signUpButton.Text = "<b>Sign up</b>"

			signUpButton.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
				//fmt.Printf("Received button signal: %v from button: %v\n", gi.ButtonSignals(sig), send.Name())
				if sig == int64(gi.ButtonClicked) { // note: 3 diff ButtonSig sig's possible -- important to check
					// vp.Win.Quit()
					//gi.PromptDialog(vp, "Button1 Dialog", "This is a dialog!  Various specific types of dialogs are available.", true, true, nil, nil)
					updt := vp.UpdateStart()
					usr := userText.Text()
					passwd := passwdText.Text()

					newLoginRec := LoginRec{Username: usr, Password: passwd, Points: 0}
					SaveNewLogin(&newLoginRec)

					vp.UpdateEnd(updt)
				}
			})

			loginButtonStartResult := trow.AddNewChild(gi.KiT_Label, "loginButtonStartResult").(*gi.Label)
			loginButtonStartResult.Text = "<b>Log In:</b>"

			userTextLogIn := trow.AddNewChild(gi.KiT_TextField, "userTextLogIn").(*gi.TextField)
			userTextLogIn.SetText("Username")
			userTextLogIn.SetProp("width", "20em")
			passwdTextLogIn := trow.AddNewChild(gi.KiT_TextField, "passwdTextLogIn").(*gi.TextField)
			passwdTextLogIn.SetText("Password")
			passwdTextLogIn.SetProp("width", "20em")

			loginButton := trow.AddNewChild(gi.KiT_Button, "loginButton").(*gi.Button)
			loginButton.Text = "<b>Log In</b>"

			loginButton.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
				//fmt.Printf("Received button signal: %v from button: %v\n", gi.ButtonSignals(sig), send.Name())
				if sig == int64(gi.ButtonClicked) { // note: 3 diff ButtonSig sig's possible -- important to check
					// vp.Win.Quit()
					//gi.PromptDialog(vp, "Button1 Dialog", "This is a dialog!  Various specific types of dialogs are available.", true, true, nil, nil)
					loginResult := trow.AddNewChild(gi.KiT_Label, "loginResult").(*gi.Label)
					loginResult.Text = "<i>Logging in as guest... user feature is not avaible yet</i><b></b>"
					trow.AddNewChild(gi.KiT_Space, "spc2")
					appsHeader := trow.AddNewChild(gi.KiT_Label, "appsHeader").(*gi.Label)

					appsHeader.Text = "<b>Apps</b>"
					appsHeader.SetProp("font-size", "x-large")
					//trow.AddNewChild(gi.KiT_Space, "spc3")

					dateAndTimeButton := trow.AddNewChild(gi.KiT_Button, "dateAndTimeButton").(*gi.Button)
					dateAndTimeButton.Text = "<b>Date and time</b>"

					dateAndTimeButton.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
						//fmt.Printf("Received button signal: %v from button: %v\n", gi.ButtonSignals(sig), send.Name())
						if sig == int64(gi.ButtonClicked) { // note: 3 diff ButtonSig sig's possible -- important to check
							// vp.Win.Quit()
							//gi.PromptDialog(vp, "Button1 Dialog", "This is a dialog!  Various specific types of dialogs are available.", true, true, nil, nil)
							updt := vp.UpdateStart()

							dateAndTimeResult := trow.AddNewChild(gi.KiT_Label, "dateAndTimeResult").(*gi.Label)
							dateAndTimeResult.Text = fmt.Sprintf(time.Now().Local().Format("2006-01-02 03:04PM"))

							vp.UpdateEnd(updt)
						}
					})

					vp.UpdateEnd(updt)
				}
			})
			vp.UpdateEnd(updt)
		}
	})

	viewlogins := trow.AddNewChild(gi.KiT_Button, "viewlogins").(*gi.Button)
	viewlogins.SetText("View LoginTable")
	viewlogins.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			lt := LoadLoginTable()

			gi.StructTableViewDialog(vp, &lt, true, nil, "Login Table", "", nil, nil, nil)
		}
	})

	addlogin := trow.AddNewChild(gi.KiT_Button, "addlogin").(*gi.Button)
	addlogin.SetText("Add Login")
	addlogin.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			rec := LoginRec{}
			gi.StructViewDialog(vp, &rec, nil, "Enter Login Info", "", recv, func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.DialogAccepted) {
					SaveNewLogin(&rec)
				}
			})
		}
	})

	quit := trow.AddNewChild(gi.KiT_Button, "quit").(*gi.Button)
	quit.SetText("Quit")
	quit.ButtonSig.Connect(rec.This, func(recv, send ki.Ki, sig int64, data interface{}) {
		if sig == int64(gi.ButtonClicked) {
			gi.PromptDialog(vp, "Quit", "Quit: Are You Sure?", true, true, recv, func(recv, send ki.Ki, sig int64, data interface{}) {
				if sig == int64(gi.DialogAccepted) {
					KaiOSDB.Close()
					vp.Win.Quit()
				}
			})
		}
	})

	vp.UpdateEndNoSig(updt)

	win.StartEventLoop()

	// note: never gets here..
	fmt.Printf("ending\n")
}
