package sessions

import (
	"testing"
)

// FIXME: this causes import cycle
func TestCreateSession(t *testing.T) {
	// connector := *GetTestDB()
	// tests := []struct {
	// 	testName string // description of this test case
	// 	// Named input parameters for target function.
	// 	name             string
	// 	workingDirectory string
	// 	connector        common.Connector
	// 	want             int
	// 	wantErr          bool
	// }{
	// 	{
	// 		testName:         "Name and PWD",
	// 		name:             "Test",
	// 		workingDirectory: "~/debianWSL/",
	// 		connector:        connector,
	// 		want:             1,
	// 		wantErr:          false,
	// 	},
	// 	{
	// 		testName:         "Name and no PWD",
	// 		name:             "Test2",
	// 		workingDirectory: "",
	// 		connector:        connector,
	// 		want:             1,
	// 		wantErr:          false,
	// 	},
	// 	{
	// 		testName:         "No name",
	// 		name:             "",
	// 		workingDirectory: "",
	// 		connector:        connector,
	// 		want:             0,
	// 		wantErr:          true,
	// 	},
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.testName, func(t *testing.T) {
	// 		got, gotErr := CreateSession(tt.name, tt.workingDirectory, tt.connector)
	// 		if gotErr != nil {
	// 			if !tt.wantErr {
	// 				t.Errorf("CreateSession() failed: %v", gotErr)
	// 			}
	// 			return
	// 		}
	// 		if tt.wantErr {
	// 			t.Fatal("CreateSession() succeeded unexpectedly")
	// 		}
	// 		if got != tt.want {
	// 			t.Errorf("CreateSession() = %v, want %v", got, tt.want)
	// 		}
	// 	})
	// }
}
