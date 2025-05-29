package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"bufio"
    "os"
    "strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/ulule/limiter/v3"
    "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
    "github.com/ulule/limiter/v3/drivers/store/memory"
)

type Post struct {
	ID        int
	Title     string
	Content   string
	CreatedAt string
}

type PageData struct {
	Posts      []Post
	NextPage   int
	PrevPage   int
	TotalPages int
}

type DataBaseStruct struct{
	db *sql.DB
}

type AdminCredentials struct{
    Username string
    Password string
}

func readAdminCredentials(filename string) (AdminCredentials, error) {
    file, err := os.Open(filename)
    if err != nil {
        return AdminCredentials{}, err
    }
    defer file.Close()

    credentials := AdminCredentials{}
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        line := scanner.Text()
        if strings.HasPrefix(line, "ADMIN_USERNAME=") {
            credentials.Username = strings.TrimPrefix(line, "ADMIN_USERNAME=")
        } else if strings.HasPrefix(line, "ADMIN_PASSWORD=") {
            credentials.Password = strings.TrimPrefix(line, "ADMIN_PASSWORD=")
        }
    }

    if err := scanner.Err(); err != nil {
        return AdminCredentials{}, err
    }

    return credentials, nil
}

func initDB() (db *sql.DB){
	var err error
	db, err = sql.Open("postgres", "user=meteoriw dbname=godb sslmode=disable password=11111")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func (dbs DataBaseStruct)adminHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		credentials, err := readAdminCredentials("admin_credentials.txt")
        if err != nil {
            http.Error(w, "Ошибка чтения учётных данных", http.StatusInternalServerError)
            return
        }

		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == credentials.Username && password == credentials.Password {
			title := r.FormValue("title")
			content := r.FormValue("content")
			_, err := dbs.db.Exec("INSERT INTO posts (title, content) VALUES ($1, $2)", title, content)
			if err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}

	tmpl := template.Must(template.ParseFiles("htmlPages/admin.html"))
	tmpl.Execute(w, nil)
}

func (dbs DataBaseStruct)indexHandler(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit := 3
	offset := (page - 1) * limit

	rows, err := dbs.db.Query("SELECT id, title, content, created_at FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
		if err != nil {
			log.Fatal(err)
		}
		posts = append(posts, post)
	}

	var totalPosts int
	dbs.db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&totalPosts)
	totalPages := (totalPosts + limit - 1) / limit

	data := PageData{
		Posts:      posts,
		NextPage:   page + 1,
		PrevPage:   page - 1,
		TotalPages: totalPages,
	}

	tmpl := template.Must(template.ParseFiles("htmlPages/index.html"))
	tmpl.Execute(w, data)
}

func (dbs DataBaseStruct)postHandler(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Path[len("/post/"):]
    log.Printf("Запрошена статья с ID: %s", id)

    var post Post
    err := dbs.db.QueryRow("SELECT id, title, content, created_at FROM posts WHERE id = $1", id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Printf("Статья с ID %s не найдена", id)
            http.NotFound(w, r)
        } else {
            log.Printf("Ошибка при запросе статьи: %v", err)
            http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
        }
        return
    }

    tmpl := template.Must(template.ParseFiles("htmlPages/post.html"))
    if err := tmpl.Execute(w, post); err != nil {
        log.Printf("Ошибка при рендеринге шаблона: %v", err)
        http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
    }
}

func main() {
	rate := limiter.Rate{
        Period:    10 * time.Second,                 
        Limit:     100,                          
    }
	store := memory.NewStore()
	instance := limiter.New(store, rate)
	middleware := stdlib.NewMiddleware(instance)
 
	dbs := DataBaseStruct{
		db : initDB(),
	}
	defer dbs.db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", dbs.indexHandler)
	mux.HandleFunc("/post/", dbs.postHandler)
	mux.HandleFunc("/admin", dbs.adminHandler)
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	handler := middleware.Handler(mux)

	log.Println("Сервер запущен на http://localhost:8888")
	log.Fatal(http.ListenAndServe(":8888", handler))
}