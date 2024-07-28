package folders

import "github.com/gofrs/uuid"

type FetchFolderRequest struct {
	OrgID uuid.UUID
}

type FetchFolderResponse struct {
	Folders []*Folder
}

type FetchFolderRequestWithPag struct {
	OrgID     uuid.UUID
	PageLimit int
	Token     string
}

type FetchFolderResponseWithPag struct {
	Folders   []*Folder
	NextToken string
	PrevToken string
}
