package delete_test

import (
	"github.com/stretchr/testify/mock"

	"github.com/financial-manager/api/internal/application/account/delete/mocks"
)

// buildMockRepoHasTx creates a mocks.Repository pre-configured for one HasTransactions call only.
func buildMockRepoHasTx(id string, hasTx bool, err error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("HasTransactions", mock.Anything, id).Return(hasTx, err).Once()
	return m
}

// buildMockRepoDelete creates a mocks.Repository pre-configured for one HasTransactions
// (returning false, nil) and one Delete call.
func buildMockRepoDelete(id string, deleteErr error) *mocks.Repository {
	m := &mocks.Repository{}
	m.On("HasTransactions", mock.Anything, id).Return(false, nil).Once()
	m.On("Delete", mock.Anything, id).Return(deleteErr).Once()
	return m
}
