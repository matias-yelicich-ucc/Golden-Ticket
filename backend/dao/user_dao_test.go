package dao

import (
	"testing"

	"golden-ticket/backend/domain"
)

func TestUserDAO(t *testing.T) {
	setupTestDB(t)
	dao := NewUserDAO()

	user := &domain.User{
		Nombre:   "Mati",
		Apellido: "Yelicich",
		Email:    "mati@ucc.edu.ar",
		Password: "hash",
		Rol:      "cliente",
		DNI:      "44555666",
	}

	if err := dao.Create(user); err != nil {
		t.Fatalf("expected create success, got %v", err)
	}

	byEmail, err := dao.GetByEmail("mati@ucc.edu.ar")
	if err != nil || byEmail.Email != "mati@ucc.edu.ar" {
		t.Fatalf("expected get by email success, got user=%+v err=%v", byEmail, err)
	}

	byDNI, err := dao.GetByDNI("44555666")
	if err != nil || byDNI.DNI != "44555666" {
		t.Fatalf("expected get by dni success, got user=%+v err=%v", byDNI, err)
	}

	if _, err := dao.GetByEmail("missing@ucc.edu.ar"); err == nil {
		t.Fatalf("expected error for missing email")
	}
	if _, err := dao.GetByDNI("00000000"); err == nil {
		t.Fatalf("expected error for missing dni")
	}
}
