package main

import (
    "fmt"
	"log"
    "database/sql"
)

type article struct {
    ID    int     `json:"id"`
    Title  string  `json:"title"`
    Authors string `json:"authors"`
}

func (p *article) getArticle(db *sql.DB) error {

	fmt.Println(p.ID);
    return db.QueryRow("SELECT title, authors FROM articles WHERE id=?",
        p.ID).Scan(&p.Title, &p.Authors)
}

func (p *article) updateArticle(db *sql.DB) error {
    _, err :=
        db.Exec("UPDATE articles SET title=?, authors=? WHERE id=?",
            p.Title, p.Authors, p.ID)

    return err
}

func (p *article) deleteArticle(db *sql.DB) error {
    _, err := db.Exec("DELETE FROM articles WHERE id=?", p.ID)

    return err
}

func (p *article) createArticle(db *sql.DB) int64 {

	
    res, err:= db.Exec("INSERT INTO articles(title, authors) VALUES(?, ?)",p.Title, p.Authors)
	lastId, err := res.LastInsertId()
	
    if err != nil {
        log.Fatal(err)
		return 0;
    }
	
	
    return lastId
}

func getArticles(db *sql.DB, start, count int) ([]article, error) {
	
    rows, err := db.Query(
        "SELECT id, title,  authors FROM articles LIMIT ? OFFSET ?",
        count, start)

    if err != nil {
        return nil, err
    }

    defer rows.Close()

    articles := []article{}

    for rows.Next() {
        var p article
        if err := rows.Scan(&p.ID, &p.Title, &p.Authors); err != nil {
            return nil, err
        }
        articles = append(articles, p)
    }

    return articles, nil
}