/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fdcont

import (
    "io"
    "ndstore/fstore/fsapi"
    "ndstore/dsrpc"
)


func (contr *Contr) SaveFileHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewSaveFileParams()

    err = context.BindParams(params)
    if err != nil {
        return err
    }

    fileSize   := context.BinSize()
    fileReader := context.BinReader()

    filePath := params.FilePath
    err = contr.store.SaveFile(filePath, fileReader, fileSize)
    if err != nil {
        context.SendError(err)
        return err
    }

    result := fsapi.NewSaveFileResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) LoadFileHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewLoadFileParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }

    filePath := params.FilePath
    fileWriter := context.BinWriter()

    err = context.ReadBin(io.Discard)
    if err != nil {
        context.SendError(err)
        return err
    }

    fileSize, err := contr.store.FileExists(filePath)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := fsapi.NewLoadFileResult()
    err = context.SendResult(result, fileSize)
    if err != nil {
        return err
    }

    err = contr.store.LoadFile(filePath, fileWriter)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) DeleteFileHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewDeleteFileParams()

    err = context.BindParams(params)
    if err != nil {
        return err
    }
    filePath   := params.FilePath
    err = contr.store.DeleteFile(filePath)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := fsapi.NewDeleteFileResult()
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}

func (contr *Contr) ListFilesHandler(context *dsrpc.Context) error {
    var err error
    params := fsapi.NewListFilesParams()
    err = context.BindParams(params)
    if err != nil {
        return err
    }
    dirPath   := params.DirPath

    entries, err := contr.store.ListFiles(dirPath)
    if err != nil {
        context.SendError(err)
        return err
    }
    result := fsapi.NewListFilesResult()
    result.Entries = entries
    err = context.SendResult(result, 0)
    if err != nil {
        return err
    }
    return err
}