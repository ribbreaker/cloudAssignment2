package main

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
)

//*********************************common structs**************************************
type IdStruct struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Languages struct {
	Languages []string `json:"languages"`
	Auth      bool     `json:"auth"`
}

type Status struct {
	StatusGit      int    `json:"gitlab"`
	StatusDatabase int    `json:"database"`
	Version        string `json:"version"`
	Uptime         int64  `json:"uptime"`
}

type Commit struct {
	Repos []Repository `json:"repos"`
	Auth  bool         `json:"auth"`
}

type Repository struct {
	Repository string `json:"path_with_namespace"`
	Commits    int    `json:"commits"`
}

//***********************************Webhooks structs******************************************

type Webhook struct {
	Id    string `json:"id"`
	Event string `json:"event"`
	Url   string `json:"url"`
	Time  int64  `json:"time"`
}

type Invocation struct {
	Event  string    `json:"event"`
	Params []string  `json:"params"`
	Time   time.Time `json:"time"`
}

type FirestoreDatabase struct {
	ProjectID      string
	CollectionName string
	Ctx            context.Context
	Client         *firestore.Client
}

type Database interface {
	Init() error
	Close()
	Save(*Webhook) error
	Delete(*Webhook) error
	ReadByID(string) (Webhook, error)
}
