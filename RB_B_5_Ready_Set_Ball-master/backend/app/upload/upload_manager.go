package upload

import (
	"bytes"
	"io"
	"mime/multipart"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/constraints"

	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/validation"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/backend/yoda"
	"git.linux.iastate.edu/309Fall2017/RB_B_5_Ready_Set_Ball/backend/app/model"
)

// NewFile initiates a new file upload to the server
// data contains the data to be stored and metadata
// Right now we only support image/png. TODO: server nits should add more
func NewFile(s model.Store, data map[string]interface{}) error {
	// we don't check for failed type assertions because we set the types in yoda
	name := data["filename"].(string)
	if !validation.IsValidImageName(name) {
		return yoda.ErrInvalidFileName
	}

	size := data["filesize"].(int64)
	if size > constraints.MaxProfilePicSize {
		return yoda.ErrFileTooLarge
	}

	f := data["file"].(multipart.File)
	fBuf, err := fileToBuffer(&f)
	if err != nil {
		return err
	}

	file := fBuf.Bytes()
	if !validation.IsValidFile(file, "profilepic") {
		return yoda.ErrInvalidFile
	}

	username := data["username"].(string)
	if username != "" {
		c := s.GetUsers()
		user, err := model.GetUserByUsername(c, username)
		if err != nil {
			return err
		}

		err = user.ChangeAvatar(c, file)
		if err != nil {
			return err
		}
	}

	return nil
}

// handleUpload returns the content of a multipart file as a bytes Buffer
func fileToBuffer(f *multipart.File) (*bytes.Buffer, error) {
	fBuf := new(bytes.Buffer)

	_, err := io.Copy(fBuf, *f)
	if err != nil {
		return nil, err
	}

	return fBuf, nil
}
