package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

// ---------- Data models ----------

type Profile struct {
	Name      string
	Title     string
	Tagline   string
	Location  string
	Email     string
	GitHub    string
	LinkedIn  string
	Repos     int
	Followers int
	Following int
	Bio       []string
}

type SkillGroup struct {
	Category string
	Skills   []string
}

type Project struct {
	Name        string
	Description string
	Language    string
	URL         string
}

type PageData struct {
	Page    string
	Profile Profile
	Skills  []SkillGroup
	Projects []Project
	Year    int
	// Contact form specific
	Submitted bool
	FormError string
	Name      string
	FormEmail string
	Message   string
}

// ---------- Static site data (sourced from github.com/tahmid56) ----------

func profileData() Profile {
	return Profile{
		Name:      "Tahmid Akter",
		Title:     "Software Engineer",
		Tagline:   "Backend Engineer · Systems Thinker · Builder of Scalable Things",
		Location:  "Remote",
		Email:     "tahmidakter56@gmail.com",
		GitHub:    "https://github.com/tahmid56",
		LinkedIn:  "https://linkedin.com/in/tahmid-akter",
		Repos:     71,
		Followers: 19,
		Following: 23,
		Bio: []string{
			"I'm a backend-focused engineer who enjoys building reliable, scalable, and production-grade systems.",
			"I go beyond just writing APIs — I care about how systems behave under load, fail, recover, and scale.",
			"I prefer learning by building rather than by following tutorials, and I'm currently focused on system design, AWS deployments, CI/CD pipelines, and high-performance backend services.",
		},
	}
}

func skillsData() []SkillGroup {
	return []SkillGroup{
		{Category: "Languages", Skills: []string{"Go (Golang)", "Rust", "TypeScript / JavaScript"}},
		{Category: "Backend & Systems", Skills: []string{"REST APIs & Microservices", "Rate limiting", "Graceful shutdown", "Background workers & queues", "Authentication systems"}},
		{Category: "Databases", Skills: []string{"PostgreSQL", "SQLx / SQLc", "Query optimization & migrations"}},
		{Category: "DevOps & Cloud", Skills: []string{"Docker", "GitHub Actions (CI/CD)", "AWS", "Nginx / Reverse Proxy"}},
	}
}

func projectsData() []Project {
	return []Project{
		{
			Name:        "actix-websocket-server",
			Description: "A WebSocket server built in Rust using the Actix framework, exploring real-time, high-concurrency connection handling.",
			Language:    "Rust",
			URL:         "https://github.com/tahmid56/actix-websocket-server",
		},
		{
			Name:        "movie-rest-api",
			Description: "A REST API written in Go for managing and serving movie data, focused on clean architecture and idiomatic Go.",
			Language:    "Go",
			URL:         "https://github.com/tahmid56/movie-rest-api",
		},
		{
			Name:        "dApp-Blog-App",
			Description: "A decentralized blogging application exploring dApp concepts on the web.",
			Language:    "JavaScript",
			URL:         "https://github.com/tahmid56/dApp-Blog-App",
		},
		{
			Name:        "websocket-horizontal-scale",
			Description: "An experiment in scaling WebSocket connections horizontally across multiple server instances.",
			Language:    "Rust",
			URL:         "https://github.com/tahmid56/websocket-horizontal-scale",
		},
	}
}

// ---------- Templating ----------

// Each page gets its own *template.Template built from the shared base
// layout plus that page's content file. Keeping them separate avoids the
// classic html/template pitfall where multiple pages defining a block with
// the same name (e.g. "content") clobber each other in a shared template set.
var pages map[string]*template.Template

func loadTemplates() {
	pages = map[string]*template.Template{
		"home":    template.Must(template.ParseFiles("templates/base.html", "templates/home.html")),
		"about":   template.Must(template.ParseFiles("templates/base.html", "templates/about.html")),
		"contact": template.Must(template.ParseFiles("templates/base.html", "templates/contact.html")),
	}
}

func render(w http.ResponseWriter, page string, data PageData) {
	tmpl, ok := pages[page]
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("template error (%s): %v", page, err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// ---------- Handlers ----------

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	data := PageData{
		Page:     "home",
		Profile:  profileData(),
		Projects: projectsData(),
		Year:     time.Now().Year(),
	}
	render(w, "home", data)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Page:    "about",
		Profile: profileData(),
		Skills:  skillsData(),
		Year:    time.Now().Year(),
	}
	render(w, "about", data)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Page:    "contact",
		Profile: profileData(),
		Year:    time.Now().Year(),
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			data.FormError = "Could not process the form. Please try again."
			render(w, "contact", data)
			return
		}

		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		if name == "" || email == "" || message == "" {
			data.FormError = "Please fill in all fields before sending."
			data.Name = name
			data.FormEmail = email
			data.Message = message
			render(w, "contact", data)
			return
		}

		// No outbound mail service is configured in this demo project, so we
		// just log the submission server-side and confirm receipt to the user.
		log.Printf("New contact message from %s <%s>: %s", name, email, message)

		data.Submitted = true
		render(w, "contact", data)
		return
	}

	render(w, "contact", data)
}

// ---------- Entry point ----------

func main() {
	loadTemplates()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/contact", contactHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("portfolio server listening on :%s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
