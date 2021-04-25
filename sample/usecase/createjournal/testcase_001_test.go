package createjournal

import (
	"accounting/domain/entity"
	"accounting/domain/vo"
	"accounting/usecase/createjournal/mocks"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCase001(t *testing.T) {

	ctx := context.Background()

	mockOutport := mocks.Outport{}

	mockOutport.
		On("BeginTransaction", ctx).
		Return(ctx)

	mockOutport.
		On("CommitTransaction", ctx).
		Return()

	mockOutport.
		On("RollbackTransaction", ctx).
		Return()

	mockOutport.
		On("GenerateUUID", ctx).
		Return("123")

	mockOutport.
		On("FindAllLastJournalBalance", ctx, "A33", []string{"100", "101"}).
		Return(map[string]float64{
			"100": 100,
			"101": 236,
		}, nil)

	mockOutport.
		On("FindAllAccountSideByCodes", ctx, "A33", []string{"100", "101"}).
		Return(map[string]vo.AccountSide{
			"100": vo.ActivaAccountSideEnum,
			"101": vo.PassivaAccountSideEnum,
		}, nil)

	now, _ := time.Parse("060102150405", "210424095637")

	obj := entity.Journal{
		ID:          "123",
		BusinessID:  "A33",
		Date:        now,
		Description: "",
		JournalType: "",
		UserID:      "",
		JournalBalances: []*entity.JournalBalance{
			{
				ID:          "A33-100-210424095637-1",
				AccountCode: "100",
				Amount:      200,
				Balance:     300,
				Direction:   entity.DirectionIncreased,
				Sequence:    1,
			},
			{
				ID:          "A33-101-210424095637-2",
				AccountCode: "101",
				Amount:      200,
				Balance:     436,
				Direction:   entity.DirectionIncreased,
				Sequence:    2,
			},
		},
	}

	mockOutport.
		On("SaveJournal", ctx, &obj).
		Return(nil)

	usecase := NewUsecase(&mockOutport)
	_, err := usecase.Execute(ctx, InportRequest{
		BusinessID:  "A33",
		Date:        now,
		Description: "",
		JournalType: "",
		UserID:      "",
		JournalBalances: []JournalBalanceRequest{
			{AccountCode: "100", Nominal: 200},
			{AccountCode: "101", Nominal: 200},
		},
	})

	assert.Nil(t, err)

}
