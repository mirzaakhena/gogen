package gogencommand

import "github.com/mirzaakhena/gogen2/util"

type InitializeModel struct {
}

func NewInitializeModel() (Commander, error) {
	return &InitializeModel{}, nil
}

func (i InitializeModel) Run() error {

	err := util.CreateFolderIfNotExist("application/apperror")
	if err != nil {
		return err
	}

	err = util.CreateFolderIfNotExist("application/registry")
	if err != nil {
		return err
	}

	err = util.CreateFolderIfNotExist("domain/entity")
	if err != nil {
		return err
	}

	err = util.CreateFolderIfNotExist("domain/vo")
	if err != nil {
		return err
	}

	err = util.CreateFolderIfNotExist("domain/repository")
	if err != nil {
		return err
	}

	err = util.CreateFolderIfNotExist("domain/service")
	if err != nil {
		return err
	}

	err = util.CreateFolderIfNotExist("gateway")
	if err != nil {
		return err
	}

	err = util.CreateFolderIfNotExist("usecase")
	if err != nil {
		return err
	}

	err = util.CreateFolderIfNotExist("controller")
	if err != nil {
		return err
	}

	err = util.CreateFolderIfNotExist("infrastructure/util")
	if err != nil {
		return err
	}

	err = util.CreateFolderIfNotExist("infrastructure/log")
	if err != nil {
		return err
	}

	InitiateLog()

	InitiateError()

	InitiateHelper()

	return nil
}
