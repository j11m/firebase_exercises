package main

import (
	"context"
	"log"
	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
)

// Firestore istemcisi global olarak tanımlanır
var client *firestore.Client

// Firestore istemcisini başlatır
func initializeFirestoreClient() (*firestore.Client, error) {
	opt := option.WithCredentialsFile("C:/Users/cem_a/Downloads/firebase-adminsdk.json") // JSON dosyasının doğru yolu
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Kullanıcı ekleme (Create)
func addUser(c *fiber.Ctx) error {
	var user map[string]interface{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz JSON formatı"})
	}

	_, _, err := client.Collection("users").Add(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı eklenemedi"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Kullanıcı başarıyla eklendi"})
}

// Kullanıcıları listeleme (Read)
func listUsers(c *fiber.Ctx) error {
	iter := client.Collection("users").Documents(context.Background())
	docs, err := iter.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcılar listelenemedi"})
	}

	var users []map[string]interface{}
	for _, doc := range docs {
		data := doc.Data()
		data["id"] = doc.Ref.ID
		users = append(users, data)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

// Kullanıcı silme (Delete)
func deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	_, err := client.Collection("users").Doc(id).Delete(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı silinemedi"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Kullanıcı başarıyla silindi"})
}

// Kullanıcı güncelleme (Update)
func updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz JSON formatı"})
	}

	_, err := client.Collection("users").Doc(id).Set(context.Background(), updates, firestore.MergeAll)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı güncellenemedi"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Kullanıcı başarıyla güncellendi"})
}

func main() {
	// Firestore istemcisini başlat
	var err error
	client, err = initializeFirestoreClient()
	if err != nil {
		log.Fatalf("Firestore başlatılamadı: %v", err)
	}
	defer client.Close()

	// Fiber uygulamasını başlat
	app := fiber.New()

	// CRUD rotaları
	app.Post("/users", addUser)          // Kullanıcı ekle
	app.Get("/users", listUsers)         // Kullanıcıları listele
	app.Delete("/users/:id", deleteUser) // Kullanıcı sil
	app.Put("/users/:id", updateUser)    // Kullanıcı güncelle

	// API sunucusunu çalıştır
	log.Println("API çalışıyor... http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
