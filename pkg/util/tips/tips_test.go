package tips

import (
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockFile struct {
	*os.File
	mockClose func() error
	clean     func(name string) error
}

func (m *mockFile) Close() (err error) {
	defer func() {
		if m.clean != nil {
			err = errors.Join(err, m.clean(m.Name()))
		}
	}()
	if m.mockClose != nil {
		return m.mockClose()
	}
	return m.File.Close()
}

func TestDoSomething(t *testing.T) {
	tests := []struct {
		name        string
		fileExists  bool
		file        *mockFile
		wantErr     error
		expectedErr error
	}{
		{
			name:       "file exists",
			fileExists: true,
			file: &mockFile{
				File:  os.NewFile(1, "test_file.txt"),
				clean: os.Remove,
			},
			wantErr: nil,
		},
		// {
		// 	name:        "file does not exist",
		// 	fileExists:  false,
		// 	wantErr:     os.ErrNotExist,
		// 	expectedErr: os.ErrNotExist,
		// },
		// {
		// 	name:        "error closing file",
		// 	fileExists:  true,
		// 	wantErr:     errors.New("failed to close file"),
		// 	expectedErr: errors.New("failed to close file"),
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file if needed
			// var file *mockFile
			// var err error
			// if tt.fileExists {
			// 	file.File, err = os.CreateTemp("", "test_file.txt")
			// 	assert.NoError(t, err)
			// 	defer os.Remove(file.Name())
			// }

			// Mock the file.Close function if needed
			// if tt.expectedErr != nil {
			// 	originalClose := file.Close
			// 	file.Close = func() error {
			// 		defer func() { file.Close = originalClose }()
			// 		err := originalClose()
			// 		if err != nil {
			// 			return errors.Join(err, tt.expectedErr)
			// 		}
			// 		return tt.expectedErr
			// 	}
			// }
			err := doSomething()
			println(err)

			// Check the error
			if tt.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func TestSingleFlight(t *testing.T) {
	go UsingSingleFlight("key")
	go UsingSingleFlight("key")
	go UsingSingleFlight("key")

	// the next calls will shared the result of the first executed call if the wait timeout is less than the actual time duration of the function call
	// if the time durtaion of the function call is 3s, the sleep time is 2s, the next calls shared the result
	// if the time durtaion of the function call is 3s, the sleep time is 5s, the shared result won't happen
	time.Sleep(5 * time.Second)

	go UsingSingleFlight("key")
	go UsingSingleFlight("key")
	go UsingSingleFlight("key")

	time.Sleep(3 * time.Second)
}

func TestLoadConfigOnce(t *testing.T) {
	GetConfig()
	GetConfig()
}

func TestGetConfigOnce(t *testing.T) {
	assert.Nil(t, instance)
	GetConfigOnce()
	GetConfigOnce()
	assert.NotNil(t, instance)
}

func TestErrorGroup(t *testing.T) {
	errorGroup()
}
func TestMaxProcs(t *testing.T) {
	maxprocs()
}

type IntLockable Lockable[int]

func TestLockable(t *testing.T) {
	// directly use
	var safeUser Lockable[user]
	safeUser.SetValue(user{name: "test"})
	log.Println("user: ", safeUser.GetValue())

	// or make a new type
	var safeInt Lockable[int]
	safeInt.SetValue(1)
	log.Println("int: ", safeInt.GetValue())
}
