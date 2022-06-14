/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */

package fsrec

import (
    "errors"
    "ndstore/dscom"
    "ndstore/dserr"

)

func (store *Store) SeedBStores() error {
    var err error

    bStores, err := store.reg.ListBStoreDescrs()
    if err != nil {
        return dserr.Err(err)
    }

    if len(bStores) > 0 {
        return dserr.Err(err)
    }

    const address   = "127.0.0.1"
    const login     = "admin"
    const pass      = "admin"
    ports := []string{ "5101", "5102", "5103" }
    for _, port := range ports {
        _, err = store.reg.AddBStoreDescr(address, port, login, pass, dscom.BStateNormal)
        if err != nil {
            return dserr.Err(err)
        }
    }
    return dserr.Err(err)
}

func (store *Store) AddBStore(userName, address, port, login, pass string) error {
    var err error

    role, err := store.reg.GetUserRole(userName)
    if role != URoleAdmin {
        err = errors.New("insufficient rights for adding bStore")
        return dserr.Err(err)
    }

    _, err = store.reg.AddBStoreDescr(address, port, login, pass, dscom.BStateNormal)
    if err != nil {
        return dserr.Err(err)
    }
    return dserr.Err(err)
}

func (store *Store) GetBStore(address, port string) (*dscom.BStoreDescr, error) {
    var err error
    bStore, err := store.reg.GetBStoreDescr(address, port)
    return bStore, dserr.Err(err)
}

func (store *Store) UpdateBStore(userName, address, port, login, pass string) error {
    var err error

    role, err := store.reg.GetUserRole(userName)
    if role != URoleAdmin {
        err = errors.New("insufficient rights for updating bStore")
        return dserr.Err(err)
    }

    ok, err := validatePass(pass)
    if !ok {
        return dserr.Err(err)
    }
    err = store.reg.UpdateBStoreDescr(address, port, login, pass, dscom.BStateNormal)
    return dserr.Err(err)
}

func (store *Store) ListBStores(userName string) ([]*dscom.BStoreDescr, error) {
    var err error
    bStores := make([]*dscom.BStoreDescr, 0)

    role, err := store.reg.GetUserRole(userName)
    if role != URoleAdmin {
        err = errors.New("insufficient rights for listing bStores")
        return bStores, dserr.Err(err)
    }

    bStores, err = store.reg.ListBStoreDescrs()
    //for i := range BStores {
    //    BStores[i].Pass = "xxxxx"
    //}
    if err != nil {
        return bStores, dserr.Err(err)
    }
    return bStores, dserr.Err(err)
}

func (store *Store) DeleteBStore(userName, address, port string) error {
    var err error
    role, err := store.reg.GetUserRole(userName)
    if role != URoleAdmin {
        err = errors.New("insufficient rights for delete bStore")
        return dserr.Err(err)
    }
    err = store.reg.DeleteBStoreDescr(address, port)
    return dserr.Err(err)
}
