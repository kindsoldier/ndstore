/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package bsbcont

import (
    "io"
    "ndstore/bstore/bsapi"
    "ndstore/dsrpc"
    "ndstore/dserr"
)

func (contr *Contr) SaveBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewSaveBlockParams()
    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    binSize     := context.BinSize()
    blockReader := context.BinReader()

    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    blockType   := params.BlockType

    blockSize   := params.BlockSize
    dataSize    := params.DataSize
    hashAlg     := params.HashAlg
    hashInit    := params.HashInit
    hashSum     := params.HashSum

    err = contr.store.SaveBlock(fileId, batchId, blockId,  blockType, blockSize, dataSize,
                                                hashAlg, hashInit, hashSum, blockReader, binSize)
    if err != nil {
        context.SendError(dserr.Err(err))
        return dserr.Err(err)
    }
    result := bsapi.NewSaveBlockResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) LoadBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewLoadBlockParams()
    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }

    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    blockType   := params.BlockType

    blockWriter := context.BinWriter()

    err = context.ReadBin(io.Discard)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }

    _, dataSize, err := contr.store.BlockParams(fileId, batchId, blockId, blockType)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewLoadBlockResult()
    err = context.SendResult(result, dataSize)
    if err != nil {
        return dserr.Err(err)
    }

    err = contr.store.LoadBlock(fileId, batchId, blockId, blockType, blockWriter)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) BlockExistsHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewBlockExistsParams()

    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    blockType   := params.BlockType

    exists, _, err := contr.store.BlockParams(fileId, batchId, blockId, blockType)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewBlockExistsResult()
    result.Exists = exists
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) CheckBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewCheckBlockParams()

    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    blockType   := params.BlockType

    correct, err := contr.store.CheckBlock(fileId, batchId, blockId, blockType)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewCheckBlockResult()
    result.Correct = correct
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (contr *Contr) DeleteBlockHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewDeleteBlockParams()

    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }
    fileId      := params.FileId
    batchId     := params.BatchId
    blockId     := params.BlockId
    blockType   := params.BlockType

    err = contr.store.DeleteBlock(fileId, batchId, blockId, blockType)
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewDeleteBlockResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (contr *Contr) PurgeAllHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewPurgeAllParams()

    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }

    err = contr.store.PurgeAll()
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewPurgeAllResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}


func (contr *Contr) ListBlocksHandler(context *dsrpc.Context) error {
    var err error
    params := bsapi.NewListBlocksParams()
    err = context.BindParams(params)
    if err != nil {
        return dserr.Err(err)
    }

    blocks, err := contr.store.ListBlocks()
    if err != nil {
        context.SendError(err)
        return dserr.Err(err)
    }
    result := bsapi.NewListBlocksResult()
    result.Blocks = blocks
    err = context.SendResult(result, 0)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}