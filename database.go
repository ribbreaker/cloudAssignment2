package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

//Initializes the database client
func (db *FirestoreDatabase) Init() error {
	db.Ctx = context.Background()
	var err error
	sa := option.WithCredentialsFile("./cloudassignment2-f6ffb-firebase-adminsdk-0s3q7-af72ed8b91.json")
	db.Client, err = firestore.NewClient(db.Ctx, db.ProjectID, sa)
	if err != nil {
		fmt.Printf("Error in FirebaseDatabase.Init() function: %v\n", err)
		log.Fatal(err, "Error in FirebaseDatabase.Save()")
	}
	return nil
}

//closes the database client
func (db *FirestoreDatabase) Close() {
	_ = db.Client.Close()
}

//Saves to the firebase client
func (db *FirestoreDatabase) Save(s *Webhook) error {
	ref := db.Client.Collection(db.CollectionName).NewDoc()
	s.Id = ref.ID
	_, err := ref.Set(db.Ctx, s)
	if err != nil {
		fmt.Println("ERROR saving student to Firestore DB: ", err)
		log.Fatal(err, "Error in FirebaseDatabase.Save()")
	}
	return nil
}

//Deletes from the database client
func (db *FirestoreDatabase) Delete(s *Webhook) error {
	docRef := db.Client.Collection(db.CollectionName).Doc(s.Id)
	_, err := docRef.Delete(db.Ctx)
	if err != nil {
		fmt.Printf("ERROR deleting student (%v) from Firestore DB: %v\n", s, err)
		log.Fatal(err, "Error in FirebaseDatabase.Save()")
	}
	return nil
}
