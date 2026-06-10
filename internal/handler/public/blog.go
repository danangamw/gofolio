package public

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type BlogHandler struct {
	tmpl Renderer
}

func NewBlogHandler(tmpl Renderer) *BlogHandler {
	return &BlogHandler{tmpl: tmpl}
}

type blogPost struct {
	Category    string
	Date        string
	Author      string
	Title       string
	Excerpt     string
	Slug        string
	HTMLContent template.HTML // parsed html
}

var blogDatabase = map[string]blogPost{
	"understanding-dependency-injection-in-go": {
		Category: "Go",
		Date:     "June 10, 2026",
		Author:   "Danang",
		Title:    "Understanding Dependency Injection in Go Simply",
		Excerpt:  "How to neatly manage database and third-party dependencies in Go applications without external frameworks.",
		Slug:     "understanding-dependency-injection-in-go",
		HTMLContent: template.HTML(`
			<p>Dependency Injection (DI) is often considered a complex concept because it is associated with large frameworks. However, in the Go programming language, Dependency Injection is actually very simple and does not require additional frameworks (like Wire or Dig) for most applications.</p>
			
			<h2>What is Dependency Injection?</h2>
			<p>Simply put, Dependency Injection means we pass the dependencies (like database connections or third-party clients) required by a function/struct, instead of letting it instantiate or search for them itself from global variables.</p>
			
			<h2>Without DI Approach (Bad)</h2>
			<pre><code>package handler

import "database/sql"

var DB *sql.DB // Global variable

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    // Tightly coupled to a global database variable
    rows, err := DB.Query("SELECT id, username FROM users")
    // ...
}</code></pre>
			<p>This approach makes unit testing difficult because the handler is tightly coupled to a global database variable.</p>

			<h2>With DI Approach (Good)</h2>
			<p>The best way in Go is to create a handler struct that accepts a database interface during initialization:</p>
			<pre><code>type UserHandler struct {
    db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
    return &UserHandler{db: db}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    // Use h.db here
}</code></pre>
			
			<h2>Conclusion</h2>
			<p>By implementing simple Dependency Injection through struct constructors (like <code>NewUserHandler</code>), your Go application code becomes much easier to test, flexible, and exceptionally modular.</p>
		`),
	},
	"integrating-opentelemetry-tracing-in-gorm": {
		Category: "Observability",
		Date:     "June 8, 2026",
		Author:   "Danang",
		Title:    "Integrating OpenTelemetry Tracing in GORM",
		Excerpt:  "A complete guide on recording SQL query performance directly to Grafana Tempo using the OTel GORM plugin.",
		Slug:     "integrating-opentelemetry-tracing-in-gorm",
		HTMLContent: template.HTML(`
			<p>When our application slows down, one of the main suspects is sub-optimal database queries. By setting up distributed tracing using OpenTelemetry (OTel) and GORM, we can track SQL query details, parameters, and duration directly in a visual dashboard like Grafana Tempo.</p>
			
			<h2>Why Use Tracing for Databases?</h2>
			<p>With tracing, every database query executed during an HTTP request will be recorded as a child span. This makes it easy to see the relation between HTTP requests and the SQL queries running under the hood.</p>

			<h2>Integration Steps</h2>
			<p>First, install the GORM tracing plugin:</p>
			<pre><code>go get gorm.io/plugin/opentelemetry/tracing</code></pre>

			<p>Then, after opening the GORM connection, register the plugin:</p>
			<pre><code>import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/plugin/opentelemetry/tracing"
)

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    log.Fatal(err)
}

// Connect GORM with OpenTelemetry
if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
    log.Printf("Failed to register OTel plugin: %v", err)
}</code></pre>

			<h2>Importance of Context</h2>
			<p>For tracing to work, you must always pass the request context (<code>ctx</code>) when executing database queries in GORM. Example:</p>
			<pre><code>// Span context is passed down via WithContext
err := db.WithContext(r.Context()).Where("id = ?", id).First(&user).Error</code></pre>

			<h2>Conclusion</h2>
			<p>Integrating OpenTelemetry with GORM is a crucial step for production-ready applications. It saves hours of debugging time when database performance issues arise.</p>
		`),
	},
	"designing-database-schema-migrations-with-atlas": {
		Category: "Database",
		Date:     "June 5, 2026",
		Author:   "Danang",
		Title:    "Designing Database Schema Migrations with Atlas",
		Excerpt:  "Why declarative migrations with Atlas are safer and more efficient than traditional manual SQL scripts.",
		Slug:     "designing-database-schema-migrations-with-atlas",
		HTMLContent: template.HTML(`
			<p>Managing database schema changes (migrations) often poses a major challenge when working in teams. Using manual approaches like writing raw SQL files is prone to conflicts and human errors. This is where Atlas shines as a modern migration solution.</p>

			<h2>Declarative vs. Imperative Approach</h2>
			<p>Most traditional Go migration libraries (like golang-migrate) use an imperative model: you write explicit SQL commands like <code>CREATE TABLE</code> or <code>ALTER TABLE</code>. If there is a typo, the migration could fail halfway.</p>
			
			<p>Atlas uses a declarative approach: you describe the desired final state of the database (e.g., from GORM structs), and Atlas automatically calculates the safe transition SQL (diff) required.</p>

			<h2>Integrating GORM with Atlas</h2>
			<p>Atlas can read GORM structs directly using a custom loader, comparing them to a local database to generate migration files automatically.</p>
			<pre><code># Command to generate a new migration file
atlas migrate diff migration_name --env local</code></pre>

			<h2>Conclusion</h2>
			<p>By combining declarative safety from Atlas and GORM mapping convenience, developers can change database schemas confidently, with minimal risk, and integrate them smoothly into CI/CD pipelines.</p>
		`),
	},
}

func (h *BlogHandler) List(w http.ResponseWriter, r *http.Request) {
	var blogs []blogPost
	for _, blog := range blogDatabase {
		blogs = append(blogs, blog)
	}

	data := map[string]any{
		"Title":      "Blog — Danang",
		"ActiveMenu": "blog",
		"Blogs":      blogs,
	}

	h.tmpl.Render(w, "blog_list", data)
}

func (h *BlogHandler) Detail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	blog, exists := blogDatabase[slug]
	if !exists {
		h.render404(w, r)
		return
	}

	data := map[string]any{
		"Title":      blog.Title + " — Blog",
		"ActiveMenu": "blog",
		"Blog":       blog,
	}

	h.tmpl.Render(w, "blog_detail", data)
}

func (h *BlogHandler) render404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.tmpl.Render(w, "404", nil)
}
