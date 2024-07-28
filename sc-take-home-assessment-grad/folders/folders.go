package folders

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gofrs/uuid"
)

/*
Input: pointer to the type FetchFolderRequest.
Output: tuple of the pointer to the type FetchFolderResponse,
with a slice of pointers to the org's folders and error msg to indicate fail.

What the function does:
  - Retrieves all folders for the specific organisation ID by calling FetchAllFoldersByOrgID.
  - Converts the retrieved folders to a slice of Folder structs.
  - Converts the slice of Folder structs back to a slice of pointers to Folder structs.
  - Constructs a FetchFolderResponse containing the slice of pointers to Folder structs and returns it.

Suggested Improvements:
  - The current function name is confusing in its meaning. Could change it to getFolderResponse
    Without reading its name, we could have intepreted it as get all folders from all organisations.
  - Handle errors returned from FetchAllFoldersByOrgID.
  - The current code converts the pointers to structs and back which is redundant.
  - Abstract away the function and only handle the response in this function.
*/
func GetAllFolders(req *FetchFolderRequest) (*FetchFolderResponse, error) {
	// validate if orgId is valid
	if !IsValidUUID(req.OrgID.String()) {
		return nil, errors.New("invalid OrgID: must be valid UUID")
	}

	// Fetch all folders by organisation ID
	folders, err := FetchAllFoldersByOrgID(req.OrgID)
	if err != nil {
		return nil, err
	}

	// If no folders found, return an error
	if len(folders) == 0 {
		return nil, fmt.Errorf("no folders found for OrgID: %s", req.OrgID.String())
	}

	// Create the folder response and return
	response := &FetchFolderResponse{Folders: folders}
	return response, nil
}

/*
input: organisation ID of type uuid.UUID.
output: slice of pointers to Folder structs and an error.
what the function does:
  - Retrieves the sample data of the folders.
  - Filters the folders by organisation id.
  - Returns the filtered list of folders by orgid.
  - Returns an error if there is an issue fetching the folders.

improvements suggested:
  - Adding error checking on folder id conforming to the uuidv4 format.
*/
func FetchAllFoldersByOrgID(orgID uuid.UUID) ([]*Folder, error) {
	folders := GetSampleData()

	// Filter folders based off OrgID and validate uuid
	resFolder := []*Folder{}
	for _, folder := range folders {
		if folder.OrgId == orgID {
			if !IsValidUUID(folder.Id.String()) {
				return nil, fmt.Errorf("folder with non-valid ID found for OrgID: %s on folder named %s", orgID.String(), folder.Name)
			}
			resFolder = append(resFolder, folder)
		}
	}
	return resFolder, nil
}

// checks if the uuidv4 is valid or not
func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-f0-9]{8}-[a-f0-9]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[a-f0-9]{12}$")
	return r.MatchString(uuid)
}
