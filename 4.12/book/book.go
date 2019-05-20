//每个go程序都应该有一个main包
package main

//引入依赖包
import(
    //"fmt"   //格式化输出
     "strconv"
     //"reflect"       //获取变量类型
     "net/http"    //实现http服务
     "encoding/json"   //实现json格式的转换
     "github.com/gorilla/mux"   //实现多路路由设置
     "github.com/jinzhu/gorm"    //实现数据库操作的gorm
     _ "github.com/jinzhu/gorm/dialects/mysql"    //go的MySQL引擎
)

//定义全局变量，存数据库操作结构体
var db *gorm.DB

//定义（构造函数）init初始化函数
func init() {
	var err error
	db, err = gorm.Open("mysql", "root:root@/golang?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&book{})
}

type (
	//数据类型
	book struct {
		gorm.Model
		BookName string `gorm:"type:varchar(50)"`
		Author   string `gorm:"type:varchar(30)"`
		Descs    string `gorm:"type:varchar(300)"`
		Status   int    `gorm:"type:tinyint(1)"`
	}
	//数据模板
	bookForm struct {
		ID        uint
		BookName  string
		Author    string
		Descs     string
		Status    string
	}
	//同意返回模板
	responseForm struct {
		Code     int
		Message  string
		Data     interface{}
	}
)

//定义路由
func main() {
	http.ListenAndServe(":8080", routeInit())
}

func routeInit() *mux.Router {
	//初始化mux路由
	mux :=  mux.NewRouter().StrictSlash(true)
	//绑定路由
	mux.HandleFunc("/book/{id}",func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, OPTIONS")
		}).Methods("OPTIONS")
	mux.HandleFunc("/book", bookList).Methods("GET")
	mux.HandleFunc("/book/{id}", getBook).Methods("GET")
	mux.HandleFunc("/book", createBook).Methods("POST")
	mux.HandleFunc("/book/{id}", deleteBook).Methods("DELETE")
	mux.HandleFunc("/book/{id}", updateBook).Methods("PUT")

	return mux
}

//查询多条数据
func bookList(w http.ResponseWriter, r *http.Request) {
	
	var books []book
	var booksForm []bookForm

	db.Find(&books)

	for _, v := range books {
		status := "发布"
		if v.Status == 0 {
			status = "未发布"
		}
		booksForm = append(booksForm, bookForm{
			ID:          v.ID,
			BookName:    v.BookName,
			Author:      v.Author,
			Descs:       v.Descs,
			Status:      status,
			})
	}
	response(w, responseForm{Code: 200, Data: booksForm})
}

//查询单条
func getBook(w http.ResponseWriter, r *http.Request) {
	var _book book
	vars := mux.Vars(r)
	db.First(&_book, vars["id"])
	if _book.ID == 0 {
		response(w, responseForm{Code: 404, Message:"data is not found!"})
		return
	}
	status := "发布"
	if _book.Status == 0 {
		status = "未发布"
	}
	_bookForm := bookForm {
		ID:          _book.ID,
		BookName:    _book.BookName,
		Author:      _book.Author,
		Descs:       _book.Descs,
		Status:      status,				
	}
	response(w, responseForm{Code: 200, Data: _bookForm})	

}

//添加数据
func createBook(w http.ResponseWriter, r *http.Request) {
	Status, _ := strconv.Atoi(r.PostFormValue("Status"))

    _book := book {
    	BookName:  r.PostFormValue("BookName"),
    	Author:    r.PostFormValue("Author"),
    	Descs:     r.PostFormValue("Descs"),
    	Status:    Status,
    }

    db.Create(&_book)
    response(w, responseForm{Code: 200, Message: "create action is successful!", Data: _book.ID})
}

//删除数据
func deleteBook(w http.ResponseWriter, r *http.Request) {
	var _book book
	vars := mux.Vars(r)
	db.First(&_book, vars["id"])
	if _book.ID == 0 {
		response(w, responseForm{Code: 404, Message:"data is not found!"})
		return
	}
	db.Delete(&_book)
	response(w, responseForm{Code: 200, Message: "delete action is successful!", Data: _book.ID})
}

//修改数据
func updateBook(w http.ResponseWriter, r *http.Request){
    var _book book
    vars := mux.Vars(r)
    db.First(&_book, vars["id"])
    if _book.ID == 0 {
		response(w, responseForm{Code: 404, Message:"data is not found!"})
		return
	}
    Status, _ := strconv.Atoi(r.PostFormValue("Status"))
    _book.BookName = r.PostFormValue("BookName")
    _book.Author   = r.PostFormValue("Author")
    _book.Descs    = r.PostFormValue("Descs")
    _book.Status   = Status
    db.Save(&_book)
    response(w, responseForm{Code: 200, Message: "update action is successful!", Data: _book.ID})
}

func response(w http.ResponseWriter, data ...interface{}) {
	json, _ :=json.Marshal(data[0])
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(json)
}
