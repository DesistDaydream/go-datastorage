package main

import "testing"

func TestInsert(t *testing.T) {
	db := InitDB()
	Insert(db)
}

func TestQuery(t *testing.T) {
	db := InitDB()
	Query(db)
}

func TestUpdate(t *testing.T) {
	db := InitDB()
	Update(db)
}

func TestDelete(t *testing.T) {
	db := InitDB()
	Delete(db)
}
