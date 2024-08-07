// Ortelius v11 Domain Microservice that handles creating and retrieving Domains
package main

import (
	"context"
	"encoding/json"

	_ "github.com/ortelius/scec-validate-provenance/docs"

	driver "github.com/arangodb/go-driver/v2/arangodb"
	"github.com/arangodb/go-driver/v2/arangodb/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/ortelius/scec-commons/database"
	"github.com/ortelius/scec-commons/model"
)

var logger = database.InitLogger()
var dbconn = database.InitializeDatabase()

// GetDomains godoc
// @Summary Get a List of Domains
// @Description Get a list of domains for the user.
// @Tags domain
// @Accept */*
// @Produce json
// @Success 200
// @Router /msapi/domain [get]
func GetDomains(c *fiber.Ctx) error {

	var cursor driver.Cursor       // db cursor for rows
	var err error                  // for error handling
	var ctx = context.Background() // use default database context

	// query all the domains in the collection
	aql := `FOR domain in evidence
			RETURN domain`

	// execute the query with no parameters
	if cursor, err = dbconn.Database.Query(ctx, aql, nil); err != nil {
		logger.Sugar().Errorf("Failed to run query: %v", err) // log error
	}

	defer cursor.Close() // close the cursor when returning from this function

	var domains []model.Domain // define a list of domains to be returned

	for cursor.HasMore() { // loop thru all of the documents

		var domain model.Domain      // fetched domain
		var meta driver.DocumentMeta // data about the fetch

		// fetch a document from the cursor
		if meta, err = cursor.ReadDocument(ctx, &domain); err != nil {
			logger.Sugar().Errorf("Failed to read document: %v", err)
		}
		domains = append(domains, domain)                                    // add the domain to the list
		logger.Sugar().Infof("Got doc with key '%s' from query\n", meta.Key) // log the key
	}

	return c.JSON(domains) // return the list of domains in JSON format
}

// GetDomain godoc
// @Summary Get a Domain
// @Description Get a domain based on the _key or name.
// @Tags domain
// @Accept */*
// @Produce json
// @Success 200
// @Router /msapi/domain/:key [get]
func GetDomain(c *fiber.Ctx) error {

	var cursor driver.Cursor       // db cursor for rows
	var err error                  // for error handling
	var ctx = context.Background() // use default database context

	key := c.Params("key")                // key from URL
	parameters := map[string]interface{}{ // parameters
		"key": key,
	}

	// query the domains that match the key or name
	aql := `FOR domain in books
			FILTER (domain.name == @key or domain._key == @key)
			RETURN domain`

	// run the query with patameters
	if cursor, err = dbconn.Database.Query(ctx, aql, &driver.QueryOptions{BindVars: parameters}); err != nil {
		logger.Sugar().Errorf("Failed to run query: %v", err)
	}

	defer cursor.Close() // close the cursor when returning from this function

	var domain model.Domain // define a domain to be returned

	if cursor.HasMore() { // domain found
		var meta driver.DocumentMeta // data about the fetch

		if meta, err = cursor.ReadDocument(ctx, &domain); err != nil {
			logger.Sugar().Errorf("Failed to read document: %v", err)
		}
		logger.Sugar().Infof("Got doc with key '%s' from query\n", meta.Key)

	} else { // not found so get from NFT Storage
		if jsonStr, exists := database.MakeJSON(key); exists {
			if err := json.Unmarshal([]byte(jsonStr), &domain); err != nil { // convert the JSON string from LTF into the object
				logger.Sugar().Errorf("Failed to unmarshal from LTS: %v", err)
			}
		}
	}

	return c.JSON(domain) // return the domain in JSON format
}

// NewDomain godoc
// @Summary Create a Domain
// @Description Create a new Domain and persist it
// @Tags domain
// @Accept application/json
// @Produce json
// @Success 200
// @Router /msapi/domain [post]
func NewDomain(c *fiber.Ctx) error {

	var err error                  // for error handling
	var meta driver.DocumentMeta   // data about the document
	var ctx = context.Background() // use default database context
	domain := new(model.Domain)    // define a domain to be returned

	if err = c.BodyParser(domain); err != nil { // parse the JSON into the domain object
		return c.Status(503).Send([]byte(err.Error()))
	}

	cid, dbStr := database.MakeNFT(&domain) // normalize the object into NFTs and JSON string for db persistence

	logger.Sugar().Infof("%s=%s\n", cid, dbStr) // log the new nft
	var resp driver.CollectionDocumentCreateResponse
	// add the domain to the database.  Ignore if it already exists since it will be identical
	if resp, err = dbconn.Collections["usergroups"].CreateDocument(ctx, domain); err != nil && !shared.IsConflict(err) {
		logger.Sugar().Errorf("Failed to create document: %v", err)
	}
	meta = resp.DocumentMeta
	logger.Sugar().Infof("Created document in collection '%s' in db '%s' key='%s'\n", dbconn.Collections["usergroups"].Name(), dbconn.Database.Name(), meta.Key)

	return c.JSON(domain) // return the domain object in JSON format.  This includes the new _key
}

// setupRoutes defines maps the routes to the functions
func setupRoutes(app *fiber.App) {

	app.Get("/swagger/*", swagger.HandlerDefault) // handle displaying the swagger
	app.Get("/msapi/domain", GetDomains)          // list of domains
	app.Get("/msapi/domain/:key", GetDomain)      // single domain based on name or key
	app.Post("/msapi/domain", NewDomain)          // save a single domain
}

// @title Ortelius v11 Domain Microservice
// @version 11.0.0
// @description RestAPI for the Domain Object
// @description ![Release](https://img.shields.io/github/v/release/ortelius/scec-validate-provenance?sort=semver)
// @description ![license](https://img.shields.io/github/license/ortelius/.github)
// @description
// @description ![Build](https://img.shields.io/github/actions/workflow/status/ortelius/scec-validate-provenance/build-push-chart.yml)
// @description [![MegaLinter](https://github.com/ortelius/scec-validate-provenance/workflows/MegaLinter/badge.svg?branch=main)](https://github.com/ortelius/scec-validate-provenance/actions?query=workflow%3AMegaLinter+branch%3Amain)
// @description ![CodeQL](https://github.com/ortelius/scec-validate-provenance/workflows/CodeQL/badge.svg)
// @description [![OpenSSF-Scorecard](https://api.securityscorecards.dev/projects/github.com/ortelius/scec-validate-provenance/badge)](https://api.securityscorecards.dev/projects/github.com/ortelius/scec-validate-provenance)
// @description
// @description ![Discord](https://img.shields.io/discord/722468819091849316)

// @termsOfService http://swagger.io/terms/
// @contact.name Ortelius Google Group
// @contact.email ortelius-dev@googlegroups.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /msapi/domain
func main() {
	port := ":" + database.GetEnvDefault("MS_POST", "8080")
	app := fiber.New()                       // create a new fiber application
	setupRoutes(app)                         // define the routes for this microservice
	if err := app.Listen(port); err != nil { // start listening for incoming connections
		logger.Sugar().Fatalf("Failed get the microservice running: %v", err)
	}
}
