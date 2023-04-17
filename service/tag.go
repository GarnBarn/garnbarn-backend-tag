package service

import (
	"github.com/GarnBarn/gb-tag-service/model"
	"github.com/GarnBarn/gb-tag-service/repository"
	"github.com/pquerna/otp/totp"
	"github.com/sirupsen/logrus"
)

type tag struct {
	tagRepository repository.Tag
}

type Tag interface {
	GetAllTag(author string) ([]model.Tag, error)
	CreateTag(tag *model.Tag) error
	UpdateTag(tagId int, tagUpdateRequest *model.UpdateTagRequest) (*model.Tag, error)
	GetTagById(tagId int) (model.TagPublic, error)
	DeleteTag(tagId int) error
	IsTagExist(tagId int) bool
}

func NewTagService(tagRepository repository.Tag) Tag {
	return &tag{
		tagRepository: tagRepository,
	}
}

func (t *tag) GetAllTag(author string) ([]model.Tag, error) {
	return t.tagRepository.GetAllTag(author)
}

func (t *tag) CreateTag(tag *model.Tag) error {

	// Create the otp secret
	totpKeyResult, err := totp.Generate(totp.GenerateOpts{Issuer: "GarnBarn", AccountName: "GarnBarn"})
	if err != nil {
		logrus.Error(err)
		return err
	}
	totpPrivateKey := totpKeyResult.Secret()
	logrus.Info(totpPrivateKey)

	tag.SecretKeyTotp = totpPrivateKey

	return t.tagRepository.Create(tag)
}

func (t *tag) UpdateTag(tagId int, tagUpdateRequest *model.UpdateTagRequest) (*model.Tag, error) {
	// Get current tag
	tag, err := t.tagRepository.GetByID(tagId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Update the tag object
	tagUpdateRequest.UpdateTag(tag)

	// Update the data in db.
	err = t.tagRepository.Update(tag)
	return tag, err
}

func (t *tag) GetTagById(tagId int) (model.TagPublic, error) {
	tag, err := t.tagRepository.GetByID(tagId)
	if err != nil {
		logrus.Error(err)
		return model.TagPublic{}, err
	}

	return tag.ToTagPublic(true), nil
}

func (t *tag) DeleteTag(tagId int) error {
	logrus.Info("Check delete tag")
	defer logrus.Info("Complete delete tag")
	err := t.tagRepository.DeleteTag(tagId)
	return err
}

func (t *tag) IsTagExist(tagId int) bool {
	_, err := t.GetTagById(tagId)
	if err != nil {
		logrus.Warn(err)
	}
	return err == nil
}
