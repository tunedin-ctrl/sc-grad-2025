package folders_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetAllFolders(t *testing.T) {
	tests := []struct {
		name    string
		orgID   uuid.UUID
		wantErr bool
		errMsg  string
	}{
		{
			name:    "DefaultOrgID test case",
			orgID:   uuid.FromStringOrNil(folders.DefaultOrgID),
			wantErr: false,
		},
		{
			name:    "NIL orgID test case",
			orgID:   uuid.FromStringOrNil(uuid.Nil.String()),
			wantErr: true,
			errMsg:  "invalid OrgID: must be valid UUID",
		},
		{
			name:    "orgID is incorrect format",
			orgID:   uuid.FromStringOrNil("00000-58cc-4372-a567-0e02b2c3d479"),
			wantErr: true,
			errMsg:  "invalid OrgID: must be valid UUID",
		},
		{
			name:    "the orgID contains no folders",
			orgID:   uuid.FromStringOrNil("f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			wantErr: true,
			errMsg:  "no folders found for OrgID",
		},
		{
			name:    "orgId that exists",
			orgID:   uuid.FromStringOrNil("d835788f-08e3-4951-ae49-54f8da1b7d42"),
			wantErr: true,
			errMsg:  "no folders found for OrgID",
		},
		{
			name:    "orgId that has folder id nil err",
			orgID:   uuid.FromStringOrNil("3b9a868b-8cd9-4b6b-ba23-fd1e08f3e2fa"),
			wantErr: true,
			errMsg:  "folder with non-valid ID found for OrgID: 3b9a868b-8cd9-4b6b-ba23-fd1e08f3e2fa on folder named nil uuid v4",
		},
		{
			name:    "orgId that has incorrect uuidv4 format on folder id",
			orgID:   uuid.FromStringOrNil("c1556e17-b7c0-45a3-a6ae-9546248fb17c"),
			wantErr: true,
			errMsg:  "folder with non-valid ID found for OrgID: c1556e17-b7c0-45a3-a6ae-9546248fb17c on folder named incorrect uuid v4 format",
		},
	}
	// execution loop, check for err
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := &folders.FetchFolderRequest{OrgID: tt.orgID}
			resp, err := folders.GetAllFolders(req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				for _, folder := range resp.Folders {
					assert.Equal(t, tt.orgID, folder.OrgId)
				}
			}
		})
	}
}
