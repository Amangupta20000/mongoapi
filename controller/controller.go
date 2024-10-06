package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Amangupta20000/mongoapi/model"
	"github.com/gorilla/mux"

	// "github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
)

// Struct to hold YAML config
type Config struct {
	Database struct {
		DBName         string `yaml:"dbName"`
		CollectionName string `yaml:"collectionName"`
	} `yaml:"database"`
}

// MOST IMPORTANT
var collection *mongo.Collection
var config Config

func loadConfig() error {
	// Use os.ReadFile instead of ioutil.ReadFile
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the YAML file into the config struct
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return fmt.Errorf("error unmarshalling config file: %v", err)
	}
	return nil
}

// Connect with MongoDB
func init() {

	// Load the .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	// Get MongoDB URL from environment variables
	connectionString := os.Getenv("MONGODB_URL")
	if connectionString == "" {
		log.Fatal("MONGODB_URL not set in environment variables")
	}

	// Load the YAML config
	// err = loadConfig()
	err := loadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// Client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")

	collection = client.Database(config.Database.DBName).Collection(config.Database.CollectionName)

	// If collection instance is ready
	fmt.Println("Collection instance/reference is ready")
}

// MongoDB Helpers -file

// insert 1 record
func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("inserted 1 movie with id : ", inserted.InsertedID)
}

// unodate a record
func updateOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)

	filter := bson.M{"_id": id}

	update := bson.M{"$set": bson.M{"watched": true}}
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("modified count : ", result.ModifiedCount)
}

// delete a record
func deleteOneMovie(movieId string) error {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		return fmt.Errorf("incorrect ID Type : %s ", err)
	}
	// matching our id with mongo:
	filter := bson.M{"_id": id}

	checkMovie := getOneMovie(movieId)
	fmt.Println("check movie", checkMovie)
	if len(checkMovie) == 0 {
		return fmt.Errorf("movie with id %s not found", movieId)
	}

	deleteCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if deleteCount.DeletedCount == 0 {
		return fmt.Errorf("no movie deleted with id %s", movieId)
	}
	fmt.Println("No of movies Deleted : ", deleteCount.DeletedCount)
	return nil
}

// delete all records from mongodb
func deleteAllMovie() int64 {
	deleteCount, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("No of movies Deleted : ", deleteCount.DeletedCount)
	return deleteCount.DeletedCount
}

// show all records
func getAllMovies() []primitive.M {
	cursor, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var movies []primitive.M

	for cursor.Next(context.Background()) {
		var movie bson.M
		err := cursor.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	defer cursor.Close(context.Background())
	return movies
}

// find one record
func getOneMovie(movieId string) primitive.M {
	// Convert the movieId from string to ObjectID
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Fatal(err)
	}

	// Create a filter to find the document by _id
	filter := bson.M{"_id": id}

	var movie bson.M

	// Find one document based on the filter
	err = collection.FindOne(context.Background(), filter).Decode(&movie)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No document found
			fmt.Println("No document was found with that ID")
			return nil
		}
		// log.Fatal(err) // handle other potential errors
	}

	return movie
}

// Actual Controler -file
func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allMovies := getAllMovies()

	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	fmt.Println("create movie started")
	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)

	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)

}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	params := mux.Vars(r)

	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	params := mux.Vars(r)

	err := deleteOneMovie(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"message": "movie deleted successfully",
		"id":      params["id"],
	})
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}

func FindOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	params := mux.Vars(r)
	movie := getOneMovie(params["id"])

	json.NewEncoder(w).Encode(movie)
}

var validCities = map[string]bool{
	"delhi":     true,
	"mumbai":    true,
	"bangalore": true,
	"chennai":   true,
	"kolkata":   true,
	// Add more cities as needed
}

func isValidCity(city string) bool {
	_, exists := validCities[city]
	return exists
}

func CheckWeather(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	params := mux.Vars(r)
	city, exists := params["city"] // Get the "city" parameter
	if !exists {
		http.Error(w, "City parameter is missing", http.StatusBadRequest)
		return
	}

	if !isValidCity(city) {
		http.Error(w, "Invalid city", http.StatusBadRequest)
		return
	}

	// Sample weather data for demonstration purposes
	weather := "25 deg" // Hardcoded for demonstration, replace this with real data fetching logic
	response := map[string]string{
		"city":    city,
		"weather": weather,
	}

	// Return the response as JSON
	json.NewEncoder(w).Encode(response)
}
