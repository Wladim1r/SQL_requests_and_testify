package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	defer db.Close()

	parcel := getTestParcel()

	// add
	num, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, num)

	// get
	newParcel, err := store.Get(num)
	require.NoError(t, err)

	parcel.Number = newParcel.Number
	assert.Equal(t, parcel, newParcel)

	// delete
	err = store.Delete(num)
	require.NoError(t, err)

	_, err = store.Get(num)
	require.Error(t, err)
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	defer db.Close()

	parcel := getTestParcel()

	// add
	num, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, num)

	// set address
	newAddress := "new test address"
	err = store.SetAddress(num, newAddress)
	require.NoError(t, err)

	// check
	newParcel, err := store.Get(num)
	require.NoError(t, err)
	assert.NotEqual(t, parcel.Address, newParcel.Address)
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	parcel := getTestParcel()
	defer db.Close()

	// add
	num, err := store.Add(parcel)
	require.NoError(t, err)
	assert.NotEmpty(t, num)

	// set status
	newStatus := ParcelStatusSent
	err = store.SetStatus(num, newStatus)
	require.NoError(t, err)

	// check
	newParcel, err := store.Get(num)
	require.NoError(t, err)
	assert.NotEqual(t, parcel.Status, newParcel.Status)
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err)
	store := NewParcelStore(db)
	defer db.Close()

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// задаём всем посылкам один и тот же идентификатор клиента
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// add
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err)

		// обновляем идентификатор добавленной посылки
		parcels[i].Number = id

		// сохраняем добавленную посылку в структуру map, чтобы её можно было легко достать по идентификатору посылки
		parcelMap[id] = parcels[i]
	}

	// get by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	assert.Len(t, storedParcels, len(parcels))

	// check
	for _, parcel := range storedParcels {
		assert.Equal(t, parcelMap[parcel.Number], parcel)
	}
}
