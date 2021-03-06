/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "io/fs"
    "ndstore/dscom"
    "ndstore/fstore/fssrv/fsreg"
)


type Store struct {
    dataRoot    string
    dirPerm     fs.FileMode
    filePerm    fs.FileMode
    reg         dscom.IFSReg

    fileWCChan  chan byte
    blockWCChan chan byte
    batchWCChan chan byte
}

func NewStore(dataRoot string, reg *fsreg.Reg) *Store {
    var store Store
    store.dataRoot  = dataRoot
    store.dirPerm   = 0755
    store.filePerm  = 0644
    store.reg       = reg
    store.fileWCChan    = make(chan byte, 1024)
    store.batchWCChan   = make(chan byte, 1024)
    store.blockWCChan   = make(chan byte, 1024)
    return &store
}

func (store *Store) SetDirPerm(dirPerm fs.FileMode) {
    store.dirPerm = dirPerm
}

func (store *Store) SetFilePerm(filePerm fs.FileMode) {
    store.filePerm = filePerm
}
