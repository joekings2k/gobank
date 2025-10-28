package token

import (
	"testing"
	"time"

	"github.com/joekings2k/gobank/util"
	"github.com/stretchr/testify/require"
)

func TestPasteoMaker (t *testing.T){
	maker, err:= NewPasetoMaker(util.RandomString(32))
	require.NoError(t,err)
	username := util.RandomOwner()
	duration := time.Minute
	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token,payload, err := maker.CreateToken(username,duration)
	require.NoError(t,err)
	require.NotEmpty(t,token)
	require.NotEmpty(t,payload)

	payload ,err = maker.VerifyToken(token)
	require.NoError(t,err)
	require.NotEmpty(t,payload)

	require.NotZero(t,payload.ID)
	require.Equal(t,username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t,expiredAt, payload.ExpiredAt,time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker,err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t,err)

	token, payload, err := maker.CreateToken(util.RandomOwner(), -time.Minute)
	require.NoError(t,err)
	require.NotEmpty(t,token)
	require.NotEmpty(t,payload)

	payload,err = maker.VerifyToken(token)
	require.Error(t,err)
	require.EqualError(t,err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestPastetoSecretkeyLength(t *testing.T) {
	_, err := NewPasetoMaker(util.RandomString(12))
	require.Error(t,err)
}

func TestPasteoInvalidToken(t *testing.T) {
	maker, err:= NewPasetoMaker(util.RandomString(32))
	require.NoError(t,err)
	username := util.RandomOwner()
	duration := time.Minute
	
	

	token, payload, err := maker.CreateToken(username,duration)
	require.NoError(t,err)
	require.NotEmpty(t,token)
	require.NotEmpty(t,payload)

	payload ,err = maker.VerifyToken("")
	require.Error(t,err)
	require.Empty(t,payload)
}