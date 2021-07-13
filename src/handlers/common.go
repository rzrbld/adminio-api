package handlers

import (
	madmin "github.com/minio/madmin-go"
	clients "github.com/rzrbld/adminio-api/clients"
)

// clients
var madmClnt = clients.MadmClnt
var minioClnt = clients.MinioClnt

var err error

type Policy struct {
	policyName   string `json:"policyName"`
	policyString string `json:"policyString"`
}

type policySet struct {
	policyName string `json:"policyName"`
	entityName string `json:"entityName"`
	isGroup    string `json:"isGroup"`
}

type UserStatus struct {
	accessKey string               `json:"accessKey"`
	status    madmin.AccountStatus `json:"status"`
}

type User struct {
	accessKey string `json:"accessKey"`
	secretKey string `json:"secretKey"`
}
