package folders_test

import (
	"testing"

	"github.com/georgechieng-sc/interns-2022/folders"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetFoldersWithPag(t *testing.T) {
	tests := []struct {
		name      string
		orgID     uuid.UUID
		pageLimit int
		token     string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Valid request with default org ID",
			orgID:     uuid.FromStringOrNil(folders.DefaultOrgID),
			pageLimit: 10,
			token:     "",
			wantErr:   false,
		},
		{
			name:      "Invalid org ID",
			orgID:     uuid.FromStringOrNil(uuid.Nil.String()),
			pageLimit: 10,
			token:     "",
			wantErr:   true,
			errMsg:    "invalid OrgID: must be valid UUID",
		},
		{
			name:      "Page limit too small",
			orgID:     uuid.FromStringOrNil(folders.DefaultOrgID),
			pageLimit: 0,
			token:     "",
			wantErr:   true,
			errMsg:    "invalid page limit: must be non-negative and below 1000 items",
		},
		{
			name:      "Page limit too large",
			orgID:     uuid.FromStringOrNil(folders.DefaultOrgID),
			pageLimit: 1001,
			token:     "",
			wantErr:   true,
			errMsg:    "invalid page limit: must be non-negative and below 1000 items",
		},
		{
			name:      "Valid request with token",
			orgID:     uuid.FromStringOrNil(folders.DefaultOrgID),
			pageLimit: 10,
			token:     folders.EncodeToken(2),
			wantErr:   false,
		},
		{
			name:      "Invalid token format",
			orgID:     uuid.FromStringOrNil(folders.DefaultOrgID),
			pageLimit: 10,
			token:     "invalid_token",
			wantErr:   true,
			errMsg:    "invalid token",
		},
		{
			name:      "No folders for org ID",
			orgID:     uuid.FromStringOrNil("f47ac10b-58cc-4372-a567-0e02b2c3d479"),
			pageLimit: 10,
			token:     "",
			wantErr:   true,
			errMsg:    "no folders found for OrgID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &folders.FetchFolderRequestWithPag{
				OrgID:     tt.orgID,
				PageLimit: tt.pageLimit,
				Token:     tt.token,
			}
			resp, err := folders.GetFoldersWithPag(req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Folders)
				assert.LessOrEqual(t, len(resp.Folders), tt.pageLimit)
				for _, folder := range resp.Folders {
					assert.Equal(t, tt.orgID, folder.OrgId)
				}
				if len(resp.Folders) == tt.pageLimit {
					assert.NotEmpty(t, resp.NextToken)
				}
				if tt.token != "" {
					assert.NotEmpty(t, resp.PrevToken)
				}
			}
		})
	}
}
