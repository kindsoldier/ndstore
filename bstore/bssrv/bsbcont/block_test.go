/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsbcont

import (
    "bytes"
    "math/rand"
    "testing"
    "path/filepath"

    "ndstore/bstore/bsapi"
    "ndstore/bstore/bssrv/bsblock"
    "ndstore/bstore/bssrv/bsbreg"
    "ndstore/dsrpc"

    "github.com/stretchr/testify/assert"
)




func Test_Block_SaveLoadDelete(t *testing.T) {
    var err error

    rootDir := t.TempDir()
    path := filepath.Join(rootDir, "blocks.db")
    reg := bsbreg.NewReg()
    err = reg.OpenDB(path)
    assert.NoError(t, err)
    err = reg.MigrateDB()
    assert.NoError(t, err)

    store := bsblock.NewStore(rootDir, reg)
    assert.NoError(t, err)

    contr := NewContr(store)

    data := make([]byte, 1024 * 1024)
    rand.Read(data)

    reader := bytes.NewReader(data)
    size := int64(len(data))

    params := bsapi.NewSaveBlockParams()
    params.FileId       = 2
    params.BatchId      = 3
    params.BlockId      = 4
    result := bsapi.NewSaveBlockResult()

    err = dsrpc.LocalPut(bsapi.SaveBlockMethod, reader, size, params, result, nil, contr.SaveBlockHandler)
    assert.NoError(t, err)

    writer := bytes.NewBuffer(make([]byte, 0))

    err = dsrpc.LocalGet(bsapi.LoadBlockMethod, writer, params, result, nil, contr.LoadBlockHandler)
    assert.NoError(t, err)
    assert.Equal(t, len(data), len(writer.Bytes()))
    assert.Equal(t, data, writer.Bytes())

    err = dsrpc.LocalExec(bsapi.DeleteBlockMethod, params, result, nil, contr.DeleteBlockHandler)
    assert.NoError(t, err)

    err = reg.CloseDB()
    assert.NoError(t, err)
}

func Test_Block_Hello(t *testing.T) {
    var err error

    rootDir := t.TempDir()
    path := filepath.Join(rootDir, "blocks.db")
    reg := bsbreg.NewReg()

    err = reg.OpenDB(path)
    assert.NoError(t, err)

    err = reg.MigrateDB()
    assert.NoError(t, err)

    store := bsblock.NewStore(rootDir, reg)
    assert.NoError(t, err)

    contr := NewContr(store)

    helloResp := GetHelloMsg
    params := bsapi.NewGetHelloParams()
    params.Message = GetHelloMsg
    result := bsapi.NewGetHelloResult()
    err = dsrpc.LocalExec(bsapi.GetHelloMethod, params, result, nil, contr.GetHelloHandler)

    assert.NoError(t, err)
    assert.Equal(t, helloResp, result.Message)

    err = reg.CloseDB()
    assert.NoError(t, err)

}