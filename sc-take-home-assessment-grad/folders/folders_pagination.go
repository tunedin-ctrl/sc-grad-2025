package folders

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/gofrs/uuid"
)

/*
Input: pointer to the type FetchFolderRequest.
Output: tuple of the pointer to the type FetchFolderResponse,
with a slice of pointers to the org's folders and error msg to indicate fail.

What the function does:
  - Retrieves all folders for the specific organisation ID by calling FetchAllFoldersByOrgIDsInPag.
  - Converts the retrieved folders to a slice of Folder structs.
  - Converts the slice of Folder structs back to a slice of pointers to Folder structs.
  - Constructs a FetchFolderResponse containing the slice of pointers to Folder structs and returns it with a next token and prev token to help with pagination.
*/
func GetFoldersWithPag(req *FetchFolderRequestWithPag) (*FetchFolderResponseWithPag, error) {
	// validate if orgId is valid
	if !IsValidUUID(req.OrgID.String()) {
		return nil, errors.New("invalid OrgID: must be valid UUID")
	}

	//check if page limit is larger than the limit or negative
	if req.PageLimit <= 0 || req.PageLimit >= dataSetSize {
		return nil, errors.New("invalid page limit: must be non-negative and below 1000 items")
	}

	// Fetch all folders by organisation ID
	folders, nextToken, prevToken, err := FetchAllFoldersByOrgIDsInPag(req.OrgID, req.PageLimit, req.Token)
	if err != nil {
		return nil, err
	}

	// If no folders found, return an error
	if len(folders) == 0 {
		return nil, fmt.Errorf("no folders found for OrgID: %s", req.OrgID.String())
	}

	// Create the folder response and return
	response := &FetchFolderResponseWithPag{Folders: folders, NextToken: nextToken, PrevToken: prevToken}
	return response, nil
}

/*
Inputs:
  - orgID: The UUID of the organisation.
  - PageLimit: The maximum number of folders to return in a single page.
  - Token: The pagination token for the current request.

Outputs:
  - A slice of Folder pointers containing the paginated results.
  - A string containing the next page token (empty if there are no more results).
  - A string containing the previous page token (empty if this is the first page).
  - An error if any issues occur during the process.

What the function does:
  - retrieves a paginated list of folders for a specific organisation.
  - filters folders by organisation ID, applies pagination, and generates tokens for obscurity.
*/
func FetchAllFoldersByOrgIDsInPag(orgID uuid.UUID, PageLimit int, Token string) ([]*Folder, string, string, error) {
	folders := GetSampleData()

	// Filter folders based off OrgID and validate uuid
	resFolder := []*Folder{}
	for _, folder := range folders {
		if folder.OrgId == orgID {
			if !IsValidUUIDv4(folder.Id.String()) {
				return nil, "", "", fmt.Errorf("folder with non-valid ID found for OrgID: %s on folder named %s", orgID.String(), folder.Name)
			}
			resFolder = append(resFolder, folder)
		}
	}

	// Pagination logic
	startIndex := 0
	if Token != "" {
		index, err := DecodeToken(Token)
		if err != nil {
			return nil, "", "", fmt.Errorf("invalid token: %v", err)
		}
		startIndex = index
	}

	if startIndex < 0 {
		return nil, "", "", errors.New("invalid token: token must be non-negative")
	}

	endIndex := startIndex + PageLimit
	if endIndex > len(resFolder) {
		endIndex = len(resFolder)
	}

	// Slice the results for pagination
	paginatedFolders := resFolder[startIndex:endIndex]

	// Determine the next and previous tokens
	nextToken := ""
	if endIndex < len(resFolder) {
		nextToken = EncodeToken(endIndex)
	}

	prevToken := ""
	if startIndex > 0 {
		prevIndex := startIndex - PageLimit
		if prevIndex < 0 {
			prevIndex = 0
		}
		prevToken = EncodeToken(prevIndex)
	}

	return paginatedFolders, nextToken, prevToken, nil
}

// checks if the uuidv4 is valid or not
func IsValidUUIDv4(uuid string) bool {
	r := regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[a-f0-9]{12}$")
	return r.MatchString(uuid)
}

// encode the indexes
func EncodeToken(index int) string {
	tokenStr := strconv.Itoa(index)

	// Encode the string to base64
	encodedToken := base64.StdEncoding.EncodeToString([]byte(tokenStr))

	return encodedToken
}

// decode the indexes
func DecodeToken(token string) (int, error) {
	// Decode the base64 string
	decodedBytes, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return 0, fmt.Errorf("invalid token format: %v", err)
	}
	// Convert the decoded bytes to a string
	decodedStr := string(decodedBytes)
	index, _ := strconv.Atoi(decodedStr)

	return index, nil
}
