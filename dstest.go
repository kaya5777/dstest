package main

import (
	"fmt"
	"net/http"
	"appengine"
	"appengine/datastore"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"appengine/urlfetch"
	"html/template"
)

type User struct {
	FirstName string
	FamilyName string
	BirthYear int
	BirthMonth int
	BirthDay int
	NickName string
	Address string
	Password string
}

func init() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/input", inHandler)
	http.HandleFunc("/input2", inHandler2)
	http.HandleFunc("/set", setHandler)
}

func inHandler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>")
	data := url.Values{"FirstName":{"学"}, "FamilyName":{"栢本"}, "BirthYear":{strconv.Itoa(1982)}, "BirthMonth":{strconv.Itoa(7)}, "BirthDay":{strconv.Itoa(7)}, }
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	//	resp,err := client.Post("http://localhost:8080/set", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	resp,err := client.Post("http://learned-now-845.appspot.com/set", "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		body,_ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		fmt.Fprintf(w, "%s<br>", string(body))
	}
	fmt.Fprintf(w, "</body></html>")
}

func inHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<html>
	<head>

	<script type="text/javascript">
	<!--

	function check(){

		var flag = 0;


		// 設定開始（必須にする項目を設定してください）

		if(document.form1.field1.value == ""){ // 「お名前」の入力をチェック

		flag = 1;

		}
		else if(document.form1.field2.value == ""){ // 「パスワード」の入力をチェック

		flag = 1;

		}
		else if(document.form1.field3.value == ""){ // 「コメント」の入力をチェック

		flag = 1;

		}

		// 設定終了


		if(flag){

			window.alert('必須項目に未入力がありました'); // 入力漏れがあれば警告ダイアログを表示
			return false; // 送信を中止

		}
		else{

			return true; // 送信を実行

		}

	}

	// -->

	</script>

	</head>
	<body>

	<form method="POST" action="/set" enctype="multipart/form-data" name="form1">
	<p>姓：<br><input type="text" name="FamilyName" size="40"> （必須）</p>
	<p>名：<br><input type="text" name="FirstName" size="40"> （必須）</p>
	<p>誕生年：<br><input type="text" name="BirthYear" size="40"> （必須）</p>
	<p>誕生月：<br><input type="text" name="BirthMonth" size="40"> （必須）</p>
	<p>誕生日：<br><input type="text" name="BirthDay" size="40"> （必須）</p>
	<p>パスワード：<br><input type="password" name="Password" size="40"> （必須）</p>
	<p><input type="submit" value="送信"></p>

	</form>

	</body>
</html>
	`
	fmt.Fprint(w, html)
}

func setHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>")
	c := appengine.NewContext(r)

	//	usr := User{
	//		FirstName: "学",
	//		FamilyName: "栢本",
	//		BirthYear: 1982,
	//		BirthMonth: 7,
	//		BirthDay: 7,
	//		NickName: "kaya",
	//		Address: "mkayamoto@gmail.com",
	//		Password: "hogehoge",
	//	}

	year, _ := strconv.Atoi(r.FormValue("BirthYear"))
	month, _ := strconv.Atoi(r.FormValue("BirthMonth"))
	day, _ := strconv.Atoi(r.FormValue("BirthDay"))

	usr := User{
		FirstName: r.FormValue("FirstName"),
		FamilyName: r.FormValue("FamilyName"),
		BirthYear: year,
		BirthMonth: month,
		BirthDay: day,
		NickName: r.FormValue("NickName"),
		Address: r.FormValue("Address"),
		Password: r.FormValue("Password"),
	}

	if usr.FamilyName == "" && usr.FirstName == "" {
		fmt.Fprintf(w, "<b>データが空っぽだよ</b>")
	} else {
		key := datastore.NewIncompleteKey(c, "user", nil)
		key, err := datastore.Put(c, key, &usr)
		if err != nil {
			fmt.Fprintf(w, "<b>%s</b>", err)
		} else {
			fmt.Fprintf(w, "<b>set success</b>%v", key)
		}
	}
	fmt.Fprintf(w, "</body></html>")
}

func handler(w http.ResponseWriter, r *http.Request) {
	/*
	fmt.Fprintf(w, "<html><body>")

	c := appengine.NewContext(r)
	q := datastore.NewQuery("user")
	var usr []User
	_, err := q.GetAll(c, &usr)
	if err != nil {
		fmt.Fprintf(w, "<b>%s</b>", err)
	} else {
		for _, u := range usr {
			fmt.Fprintf(w, "%s<br>", u.FirstName)
			fmt.Fprintf(w, "%s<br>", u.FamilyName)
			fmt.Fprintf(w, "%d/%d/%d<br>", u.BirthYear, u.BirthMonth, u.BirthDay)
			fmt.Fprintf(w, "%s<br>", u.NickName)
			fmt.Fprintf(w, "%s<br>", u.Address)
			fmt.Fprintf(w, "%s<br>", u.Password)
			fmt.Fprintf(w, "===============<br>")
		}
	}
	fmt.Fprintf(w, "</body></html>")
	*/
	c := appengine.NewContext(r)
	q := datastore.NewQuery("user")
	var usr []User
	q.GetAll(c, &usr)
	t,err := template.ParseFiles("./list.html")
	if err != nil {
		fmt.Fprintln(w, err)
	} else {
		t.Execute(w, usr)
	}
}
