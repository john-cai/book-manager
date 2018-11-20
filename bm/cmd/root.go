package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/john-cai/book-manager/models"
	"github.com/john-cai/book-manager/responder"
	"github.com/labstack/gommon/log"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "Book Manager is a nifty system to manage books\n",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("This is Book Manager")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Book Manager",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Book Manager CLI v0.1 -- HEAD")
	},
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add books or collections",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("must specify either book or collection")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "book":
			err := AddBook(isbn, title, author, description, publishedYear, genres)
			if err != nil {
				if errResp, ok := err.(responder.ErrorResponse); ok {
					for _, e := range errResp.Errors {
						fmt.Printf("problem with %v: %v\n", e.Field, e.Message)
					}
					return
				}
				fmt.Println("Something went horribly wrong and I'm so sorry")
				return
			}
			fmt.Printf("%s successfully added to books\n", title)
		case "collection":
			collection, err := AddCollection(collectionName, collectionDescription)
			if err != nil {
				if errResp, ok := err.(responder.ErrorResponse); ok {
					for _, e := range errResp.Errors {
						fmt.Printf("problem with %v: %v\n", e.Field, e.Message)
					}
					return
				}
				fmt.Println("Something went horribly wrong and I'm so sorry")
				return
			}
			fmt.Printf("collection %s successfully added to collections with id %d\n", collectionName, collection.ID)
		default:
			log.Error("unrecognized command")
		}
	},
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View books and collections",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("must specify either book(s) or collection(s)")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "books":
			books, err := ViewBooks(title, author, description, publishedYear, genres)
			if err != nil {
				if errResp, ok := err.(responder.ErrorResponse); ok {
					for _, e := range errResp.Errors {
						fmt.Printf("problem with %v: %v\n", e.Field, e.Message)
					}
					return
				}
				log.Error("Something went horribly wrong and I'm so sorry")
				return
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ISBN", "Title", "Author", "Description", "Published"})
			for _, book := range books {
				table.Append([]string{
					book.ISBN,
					book.Title,
					book.Author,
					book.Description,
					book.PublishedAt.Format("2006"),
				})
			}
			table.Render()
			return
		case "book":
			//TODO: implement me
		case "collections":
			//TODO: implement me
			collections, err := ViewCollections()
			if err != nil {
				if errResp, ok := err.(responder.ErrorResponse); ok {
					for _, e := range errResp.Errors {
						fmt.Printf("problem with %v: %v\n", e.Field, e.Message)
					}
					return
				}
				log.Error("Something went horribly wrong and I'm so sorry")
				return
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"ID", "Name", "Description", "Books"})
			for _, collection := range collections {
				table.Append([]string{
					strconv.Itoa(collection.ID),
					collection.Name,
					collection.Description,
					strconv.Itoa(len(collection.Books)),
				})
			}
			table.Render()
			return
		case "collection":
			//TODO: implement me
		default:
			log.Error("unrecognized command")
		}
	},
}

func sendRequest(url string, method string, payload interface{}, response interface{}) error {
	client := http.Client{}
	var b bytes.Buffer
	var err error
	var errResp responder.ErrorResponse

	if payload != nil {
		if err = json.NewEncoder(&b).Encode(payload); err != nil {
			return err
		}

	}
	req, err := http.NewRequest(method, url, &b)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		if err = json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return err
		}
		return errResp
	}

	if response != nil {
		if err = json.NewDecoder(resp.Body).Decode(response); err != nil {
			return err
		}
	}
	return nil
}

// AddBook calls the api to add a book
func AddBook(isbn, title, author, description string, publishedYear int, genres []string) error {
	book := models.Book{
		ISBN:        isbn,
		Title:       title,
		Author:      author,
		Description: description,
		Metadata:    models.Metadata{Genres: genres},
	}
	if publishedYear != 0 {
		book.PublishedAt = time.Date(publishedYear, 0, 0, 0, 0, 0, 0, time.UTC)
	}
	if err := sendRequest(fmt.Sprintf("http://%s/books", bookmanagerURL), http.MethodPost, &book, nil); err != nil {
		return err
	}
	return nil
}

// ViewBook calls the api to view a book
func ViewBook(isbn string, bookManagerURL string) error {
	return nil
}

// EditBook calls the api to edit a book
func EditBook(isbn, title, author, description string, publishedYear int, genres []string) error {
	return nil
}

// RemoveBook calls the api to remove a book
func RemoveBook(isbn string, bookManagerURL string) error {
	return nil
}

// ViewBooks calls the api to view book with filter criteria
func ViewBooks(title, author, description string, publishedYear int, genres []string) ([]models.Book, error) {
	var books []models.Book
	if err := sendRequest(fmt.Sprintf("http://%s/books", bookmanagerURL), http.MethodGet, nil, &books); err != nil {
		return nil, err
	}
	return books, nil
}

// AddCollection calls the api to add a collection
func AddCollection(name, description string) (models.Collection, error) {
	collection := models.Collection{
		Name:        name,
		Description: description,
	}
	if err := sendRequest(fmt.Sprintf("http://%s/collections", bookmanagerURL), http.MethodPost, &collection, &collection); err != nil {
		return collection, err
	}
	return collection, nil
}

// ViewCollection calls the api get the details of a collection
func ViewCollection(id int) (models.Collection, error) {
	var collection models.Collection
	if err := sendRequest(fmt.Sprintf("http://%s/collections/%d", bookmanagerURL, id), http.MethodGet, &collection, nil); err != nil {
		return models.Collection{}, err
	}
	return collection, nil
}

// ViewCollections calls the api get the details of all collections
func ViewCollections() ([]models.Collection, error) {
	var collections []models.Collection
	if err := sendRequest(fmt.Sprintf("http://%s/collections", bookmanagerURL), http.MethodGet, &collections, nil); err != nil {
		return nil, err
	}
	return collections, nil
}

var (
	bookmanagerURL string
	isbn           string
	title          string
	author         string
	description    string
	publishedYear  int
	genres         []string
)

var (
	collectionName        string
	collectionDescription string
)

func init() {
	bookmanagerURL = os.Getenv("BOOKMANAGER_URL")
	if bookmanagerURL == "" {
		log.Fatal("BOOKMANAGER_URL is required")
	}
	addCmd.Flags().StringVar(&isbn, "isbn", "", "isbn of the book")
	addCmd.Flags().StringVar(&title, "title", "", "title of the book")
	addCmd.Flags().StringVar(&author, "author", "", "author of the book")
	addCmd.Flags().StringVar(&description, "description", "", "description of the book")
	addCmd.Flags().IntVar(&publishedYear, "published", 0, "year the book was published")
	addCmd.Flags().StringSliceVar(&genres, "genres", []string{}, "genres of the book")

	addCmd.Flags().StringVar(&collectionName, "name", "", "name of the collection")
	addCmd.Flags().StringVar(&collectionDescription, "collection-description", "", "description of the collection")

	viewCmd.Flags().StringVar(&isbn, "isbn", "", "isbn of the book")
	viewCmd.Flags().StringVar(&title, "title", "", "title of the book")
	viewCmd.Flags().StringVar(&author, "author", "", "author of the book")
	viewCmd.Flags().StringVar(&description, "description", "", "description of the book")
	viewCmd.Flags().IntVar(&publishedYear, "published", 0, "year the book was published")
	viewCmd.Flags().StringSliceVar(&genres, "genres", []string{}, "genres of the book")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(viewCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
