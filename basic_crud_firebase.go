package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// Firebase App'i başlatma
func initializeFirebaseApp() (*firebase.App, error) {
	opt := option.WithCredentialsFile("C:/Users/cem_a/Downloads/firebase-adminsdk.json") // Firebase key JSON dosyası ekleme
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("Firebase başlatılamadı: %v", err)
	}
	return app, nil
}

// adding sample data
func addSampleData(client *firestore.Client) error {
	users := []map[string]interface{}{
		{"name": "Ali", "age": 30, "city": "Adana"},
		{"name": "Ayşe", "age": 25, "city": "Mersin"},
		{"name": "Ahmet", "age": 35, "city": "İstanbul"},
	}

	for _, user := range users {
		_, _, err := client.Collection("users").Add(context.Background(), user)
		if err != nil {
			return fmt.Errorf("Örnek veri eklenemedi: %v", err)
		}
	}
	fmt.Println("Örnek veriler başarıyla eklendi!")
	return nil
}

// doc listing
func listDocuments(client *firestore.Client) ([]*firestore.DocumentSnapshot, error) {
	fmt.Println("\nVerileri listeliyorum...")
	iter := client.Collection("users").Documents(context.Background())
	docs, err := iter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("Veriler okunamadı: %v", err)
	}

	for _, doc := range docs {
		fmt.Printf("Doküman ID: %s, Veri: %v\n", doc.Ref.ID, doc.Data())
	}
	return docs, nil
}

// doc delete
func deleteRandomDocument(client *firestore.Client, docID string) error {
	_, err := client.Collection("users").Doc(docID).Delete(context.Background())
	if err != nil {
		return fmt.Errorf("Doküman silinemedi: %v", err)
	}
	fmt.Printf("Doküman ID: %s başarıyla silindi!\n", docID)
	return nil
}

// doc update
func updateDocument(client *firestore.Client, docID string) error {
	_, err := client.Collection("users").Doc(docID).Set(context.Background(), map[string]interface{}{
		"city": "Gaziantep",
		"age":  40,
	}, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("Doküman güncellenemedi: %v", err)
	}
	fmt.Printf("Doküman ID: %s başarıyla güncellendi!\n", docID)
	return nil
}

func main() {
	// Firebase uygulamasını başlatma
	app, err := initializeFirebaseApp()
	if err != nil {
		log.Fatalf("Uygulama başlatılamadı: %v", err)
	}

	// Firestore istemcisi oluşturma
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Firestore başlatılamadı: %v", err)
	}
	defer client.Close()

	// 1. Örnek veriler eklenir (eğer zaten yoksa)
	err = addSampleData(client)
	if err != nil {
		log.Fatalf("Hata oluştu: %v", err)
	}

	// 2. Tüm verileri listeleme
	docs, err := listDocuments(client)
	if err != nil {
		log.Fatalf("Hata oluştu: %v", err)
	}

	// 3. Bir doküman seçtip0 siliyor ve tekrar listeleyip farka bakıyoruz
	if len(docs) > 0 {
		err = deleteRandomDocument(client, docs[0].Ref.ID) // İlk dokümanı siliyoruz
		if err != nil {
			log.Fatalf("Hata oluştu: %v", err)
		}

		// Sildikten sonra tekrar listele
		_, err = listDocuments(client)
		if err != nil {
			log.Fatalf("Hata oluştu: %v", err)
		}
	} else {
		fmt.Println("Silinecek doküman bulunamadı!")
	}

	// 4. Bir doküman güncelleyip sonra ise tekrar listeleyip farka bakıyoruz
	if len(docs) > 1 {
		err = updateDocument(client, docs[1].Ref.ID) // İkinci dokümanı güncelliyoruz
		if err != nil {
			log.Fatalf("Hata oluştu: %v", err)
		}

		// Güncellemeden sonra tekrar listele
		_, err = listDocuments(client)
		if err != nil {
			log.Fatalf("Hata oluştu: %v", err)
		}
	} else {
		fmt.Println("Güncellenecek doküman bulunamadı!")
	}
}
