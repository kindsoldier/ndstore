/*
 * Copyright 2022 Oleg Borodin  <borodin@unix7.org>
 */
package bsureg

import (
    "ndstore/bstore/bscom"
    "ndstore/dslog"
)

const usersSchema = `
    DROP TABLE IF EXISTS users;
    CREATE TABLE IF NOT EXISTS users (
        login       TEXT,
        pass        TEXT,
        state       TEXT,
        role        TEXT
    );
    DROP INDEX IF EXISTS user_idx;
    CREATE UNIQUE INDEX IF NOT EXISTS user_idx
        ON users (login);`


func (reg *Reg) AddUserDescr(login, pass, state, role string) error {
    var err error
    request := `
        INSERT INTO users(login, pass, state, role)
        VALUES ($1, $2, $3, $4);`
    _, err = reg.db.Exec(request, login, pass, state, role)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) UpdateUserDescr(login, pass, state, role string) error {
    var err error
    dslog.LogDebug(login, pass, state)
    request := `
        UPDATE users
        SET pass = $1, state = $2, role = $3
        WHERE login = $4;`
    _, err = reg.db.Exec(request, pass, state, role, login)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) RenewUserDescr(descr *bscom.UserDescr) error {
    var err error
    request := `
        UPDATE users
        SET pass = $1, state = $2, role = $3
        WHERE login = $4;`
    _, err = reg.db.Exec(request, descr.Pass, descr.State, descr.Role, descr.Login)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) GetUserDescr(login string) (*bscom.UserDescr, bool, error) {
    var err error
    var exists bool
    var user *bscom.UserDescr
    request := `
        SELECT login, pass, state, role
        FROM users
        WHERE login = $1
        LIMIT 1;`
    users := make([]*bscom.UserDescr, 0)
    err = reg.db.Select(&users, request, login)
    if err != nil {
        return user, exists, err

    }
    if len(users) > 0 {
        exists = true
        user = users[0]
    }
    return user, exists, err
}

func (reg *Reg) UserDescrExists(login string) (bool, error) {
    var err error
    var exists bool
    request := `
        SELECT login, pass, state, role
        FROM users
        WHERE login = $1
        LIMIT 1;`
    users := make([]*bscom.UserDescr, 0)
    err = reg.db.Select(&users, request, login)
    if err != nil {
        return exists, err
    }
    if len(users) > 0 {
        exists = true
    }
    return exists, err
}

func (reg *Reg) DeleteUserDescr(login string) error {
    var err error
    request := `
        DELETE FROM users
        WHERE login = $1;`
    _, err = reg.db.Exec(request, login)
    if err != nil {
        return err
    }
    return err
}

func (reg *Reg) ListUserDescrs() ([]*bscom.UserDescr, error) {
    var err error
    request := `
        SELECT login, pass, state, role
        FROM users;`
    users := make([]*bscom.UserDescr, 0)
    err = reg.db.Select(&users, request)
    if err != nil {
        return users, err
    }
    return users, err
}