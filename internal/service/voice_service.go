package service

import (
	"context"
	"errors"
	"io"

	"github.com/matheusabido/open-diet-api/internal/dto"
	"github.com/matheusabido/open-diet-api/internal/repository"
	"github.com/matheusabido/open-diet-api/internal/voice"
	"github.com/matheusabido/open-diet-api/internal/whisper"
	"github.com/sirupsen/logrus"
)

const maxAliasLength = 100

var ErrVoiceNotUnderstood = errors.New("voice not understood")

type VoiceService struct {
	transcriber *whisper.Transcriber
	aliasRepo   *repository.FoodAliasRepository
	foodRepo    *repository.FoodRepository
}

func NewVoiceService(transcriber *whisper.Transcriber, aliasRepo *repository.FoodAliasRepository, foodRepo *repository.FoodRepository) *VoiceService {
	return &VoiceService{transcriber: transcriber, aliasRepo: aliasRepo, foodRepo: foodRepo}
}

func (s *VoiceService) FoodItem(ctx context.Context, userID uint64, audio io.Reader) (*dto.VoiceFoodItemDTO, error) {
	text, err := s.transcriber.Transcribe(ctx, audio)
	if err != nil {
		return nil, err
	}
	logrus.Infof("voice transcription: %q", text)

	amount, itemUnit, name, ok := voice.Parse(text)
	if !ok || len([]rune(name)) > maxAliasLength {
		return nil, ErrVoiceNotUnderstood
	}

	alias, err := s.aliasRepo.FindOrCreate(ctx, userID, name)
	if err != nil {
		return nil, err
	}

	item := &dto.VoiceFoodItemDTO{AliasID: alias.ID, Name: alias.Name, Amount: amount, Unit: itemUnit}
	if alias.FoodID != nil {
		if item.Food, err = s.foodRepo.FindByID(ctx, *alias.FoodID); err != nil {
			return nil, err
		}
	}

	return item, nil
}
