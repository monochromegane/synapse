package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/monochromegane/synapse"
	"github.com/olivere/elastic"
)

func main() {
	err := setUp()
	if err != nil {
		panic(err)
	}

	err = run()
	if err != nil {
		panic(err)
	}
}

func setUp() error {
	// setup users
	err := setUpUsers()
	if err != nil {
		return err
	}

	// setup products
	err = setUpProducts()
	if err != nil {
		return err
	}

	return nil
}

func setUpUsers() error {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/")
	if err != nil {
		return err
	}
	defer db.Close()

	dbName := "synapse"
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		return err
	}
	_, err = db.Exec("USE " + dbName)
	if err != nil {
		return err
	}

	tableName := "users"
	_, err = db.Exec("DROP TABLE " + tableName)
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE TABLE users ( id int AUTO_INCREMENT NOT NULL PRIMARY KEY, name VARCHAR(50), birth DATE )")
	if err != nil {
		return err
	}

	now := time.Now().Year()
	loc, _ := time.LoadLocation("Asia/Tokyo")
	for i, name := range names {
		date := time.Date(now-i-10, 1, 1, 0, 0, 0, 0, loc).Format("2006-01-02")
		q := fmt.Sprintf("INSERT INTO %s(name, birth) VALUES('%s', '%s')", tableName, name, date)
		_, err = db.Exec(q)
		if err != nil {
			return err
		}
	}

	return nil
}

type Product struct {
	Name       string    `json:"name"`
	CategoryID int       `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"product":{
			"properties":{
				"name":{
					"type":"text"
				},
				"created_at":{
					"type":"date"
				}
			}
		}
	}
}`

var names = []string{
	"raccoon", "dog", "wild boar", "rabbit", "cow", "horse", "wolf", "hippopotamus", "kangaroo",
	"fox", "giraffe", "bear", "koala", "bat", "gorilla", "rhinoceros", "monkey",
	"deer", "zebra", "jaguar", "polar bear", "skunk", "elephant", "raccoon dog",
	"animal", "reindeer", "rat", "tiger", "cat", "mouse", "buffalo", "hamster", "panda",
	"sheep", "leopard", "pig", "mole", "goat", "lion", "camel", "squirrel", "donkey"}

func setUpProducts() error {
	indexName := "synapse"

	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}

	ctx := context.Background()
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if exists {
		_, err := client.DeleteIndex(indexName).Do(ctx)
		if err != nil {
			return err
		}
	}

	createIndex, err := client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
	if err != nil {
		return err
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
	}

	now := time.Now()
	for i, name := range names {
		product := Product{Name: name, CategoryID: (i % 10) + 1, CreatedAt: now.Add(-(time.Duration(i) * time.Minute))}
		_, err := client.Index().
			Index(indexName).
			Type("product").
			BodyJson(product).
			Do(ctx)
		if err != nil {
			return err
		}
	}
	_, err = client.Flush().Index(indexName).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func run() error {
	config := synapse.Config{
		PluginDir: "example/plugins",
		Matchers: []synapse.ConfigMatcher{
			synapse.ConfigMatcher{
				Name: "users_products",
				Profiles: []synapse.ConfigPlugin{
					synapse.ConfigPlugin{
						Name:    "profiler",
						Version: "0.0.1",
					},
				},
				Associator: synapse.ConfigPlugin{
					Name:    "associator",
					Version: "0.0.1",
				},
				Searcher: synapse.ConfigPlugin{
					Name:    "searcher",
					Version: "0.0.1",
				},
			},
		},
	}

	syn, err := synapse.NewSynapse(config)
	if err != nil {
		return err
	}

	ctx := synapse.Context{"user_id": "1"}
	hits, err := syn.Match("users_products", ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%v\n", hits.IDs)

	return nil
}
