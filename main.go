package main

import (
    "fmt"
    "net/http"
    "html/template"
 
)

func page(w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles("index.html")
    if err != nil {
        http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
        return
    }
    tmpl.Execute(w, nil)
}

func handleRequest() {
    // Обслуживаем статические файлы (CSS, JS, изображения)
    fs := http.FileServer(http.Dir("./web"))
    http.Handle("/web/", http.StripPrefix("/web/", fs))
    
    // HTML страница
    http.HandleFunc("/", page)
    
    fmt.Println("Сервер запущен на http://localhost:5500")
    http.ListenAndServe(":5500", nil)
}

func main() {
    handleRequest()
}