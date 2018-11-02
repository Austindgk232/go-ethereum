// Copyright 2018 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
)

type AuditLogger struct {
	log log.Logger
	api ExternalAPI
}

func (l *AuditLogger) List(ctx context.Context) ([]common.Address, error) {
	l.log.Info("List", "type", "request", "metadata", MetadataFromContext(ctx).String())
	res, e := l.api.List(ctx)
	l.log.Info("List", "type", "response", "data", res)

	return res, e
}

func (l *AuditLogger) New(ctx context.Context) (accounts.Account, error) {
	return l.api.New(ctx)
}

func (l *AuditLogger) SignTransaction(ctx context.Context, args SendTxArgs, methodSelector *string) (*ethapi.SignTransactionResult, error) {
	sel := "<nil>"
	if methodSelector != nil {
		sel = *methodSelector
	}
	l.log.Info("SignTransaction", "type", "request", "metadata", MetadataFromContext(ctx).String(),
		"tx", args.String(),
		"methodSelector", sel)

	res, e := l.api.SignTransaction(ctx, args, methodSelector)
	if res != nil {
		l.log.Info("SignTransaction", "type", "response", "data", common.Bytes2Hex(res.Raw), "error", e)
	} else {
		l.log.Info("SignTransaction", "type", "response", "data", res, "error", e)
	}
	return res, e
}

func (l *AuditLogger) SignData(ctx context.Context, contentType string, addr common.MixedcaseAddress, data interface{}) (hexutil.Bytes, error) {
	l.log.Info("SignData", "type", "request", "metadata", MetadataFromContext(ctx).String(),
		"addr", addr.String(), "data", data, "content-type", contentType)
	b, e := l.api.SignData(ctx, contentType, addr, data)
	l.log.Info("SignData", "type", "response", "data", common.Bytes2Hex(b), "error", e)
	return b, e
}

//func (l *AuditLogger) SignTypedData(ctx context.Context, addr common.MixedcaseAddress, data TypedData) (hexutil.Bytes, error) {
//	l.log.Info("SignTypedData", "type", "request", "metadata", MetadataFromContext(ctx).String(),
//		"addr", addr.String(), "data", data)
//	b, e := l.api.SignTypedData(ctx, addr, data)
//	l.log.Info("SignTypedData", "type", "response", "data", common.Bytes2Hex(b), "error", e)
//	return b, e
//}

func (l *AuditLogger) EcRecover(ctx context.Context, contentType string, data hexutil.Bytes, sig hexutil.Bytes) (common.Address, error) {
	l.log.Info("EcRecover", "type", "request", "metadata", MetadataFromContext(ctx).String(),
		"data", common.Bytes2Hex(data), "sig", common.Bytes2Hex(sig), "content-type", contentType)
	b, e := l.api.EcRecover(ctx, contentType, data, sig)
	l.log.Info("EcRecover", "type", "response", "address", b.String(), "error", e)
	return b, e
}

func (l *AuditLogger) Export(ctx context.Context, addr common.Address) (json.RawMessage, error) {
	l.log.Info("Export", "type", "request", "metadata", MetadataFromContext(ctx).String(),
		"addr", addr.Hex())
	j, e := l.api.Export(ctx, addr)
	// In this case, we don't actually log the json-response, which may be extra sensitive
	l.log.Info("Export", "type", "response", "json response size", len(j), "error", e)
	return j, e
}

//func (l *AuditLogger) Import(ctx context.Context, keyJSON json.RawMessage) (Account, error) {
//	// Don't actually log the json contents
//	l.log.Info("Import", "type", "request", "metadata", MetadataFromContext(ctx).String(),
//		"keyJSON size", len(keyJSON))
//	a, e := l.api.Import(ctx, keyJSON)
//	l.log.Info("Import", "type", "response", "addr", a.String(), "error", e)
//	return a, e
//}

func NewAuditLogger(path string, api ExternalAPI) (*AuditLogger, error) {
	l := log.New("api", "signer")
	handler, err := log.FileHandler(path, log.LogfmtFormat())
	if err != nil {
		return nil, err
	}
	l.SetHandler(handler)
	l.Info("Configured", "audit log", path)
	return &AuditLogger{l, api}, nil
}
