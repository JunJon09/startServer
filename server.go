package main

import (
	"fmt" //標入出力を行うため
	"html/template"
	"log" //ロギングの機能を提供するパッケージ
	"net/http" //HTTPクライアントとサーバーの機能の提供
	"os"
	"bufio"
)


type BookList struct {
	Books []string
}

func New(books []string) *BookList {
	return &BookList{Books: books}
}

func fileRead(fileName string) []string {
	var bookList []string
	file, err := os.Open(fileName)
	if os.IsNotExist(err) {
		return nil
	}
	defer file.Close()
	scaner := bufio.NewScanner(file) //データを読み取るた目にスキャナーを作成
	for scaner.Scan() {//1行ずつ読み取るために作成
		bookList = append(bookList, scaner.Text()) //TextをbookListに追加
	}
	return bookList
}

//WebサーバがHTTPリクエストを受けた時に実行される関数
//引数w http.ResponseWriter はクライアントへのレスポンスを書き込むためのインターフェース, r *http.Requestは、クライアントからのリクエストの情報
func helloHandler(w http.ResponseWriter, r *http.Request) {
	hello := []byte("Hello World!!!")
	_, err := w.Write(hello)
	if err != nil {
		log.Fatal(err)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	bookList := fileRead("reading.txt")
	fmt.Println(bookList)
	html, err := template.ParseFiles("view.html")
	if err != nil {
		log.Fatal(err)
	}
	getBooks := New(bookList)
	if err := html.Execute(w, getBooks); err != nil {//view.htmに解析したテンプレートを代入
		log.Fatal(err)
	}
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	formValue := r.FormValue("value")
	file, err := os.OpenFile("reading.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(0600))
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintln(file, formValue)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/view", http.StatusFound)
}

func main() {
	http.HandleFunc("/hello", helloHandler) // /helloのパスで接続がくると関数を呼び出す
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/view/create", createHandler)
	fmt.Println("Server Start Up........")
	log.Fatal(http.ListenAndServe("localhost:8080", nil)) //サーバー起動
}