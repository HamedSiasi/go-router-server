/* Description of user for the UTM server.
 *
 * Copyright (C) u-blox Melbourn Ltd
 * u-blox Melbourn Ltd, Melbourn, UK
 *
 * All rights reserved.
 *
 * This source file is the sole property of u-blox Melbourn Ltd.
 * Reproduction or utilization of this source in whole or part is
 * forbidden without the written consent of u-blox Melbourn Ltd.
 */

package models

import (
    "crypto/md5"
    "encoding/hex"
    "errors"
    "io"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

/// Definition of a user
type User struct {
    ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
    Company   string        `bson:"company" json:"company"`
    FirstName string        `bson:"firstName" json:"firstName"`
    LastName  string        `bson:"lastName" json:"lastName"`
    UserName  string        `bson:"userName" json:"userName"`
    Email     string        `bson:"email" json:"email"`
    Password  string        `bson:"password" json:"password"`
}

/// Add a user to the database
func (u *User) NewUser(db *mgo.Database, company, firstName, lastName, userName, email, password string) {
    u.Company = company
    u.FirstName = firstName
    u.LastName = lastName
    u.UserName = userName
    u.Email = email
    u.ID = bson.NewObjectId()
    h := md5.New()
    io.WriteString(h, password)
    u.Password = hex.EncodeToString(h.Sum(nil))
    c := db.C("users")

    c.Insert(&u)
}

/// Get a user from the database based on ID
func (u *User) Get(db *mgo.Database, id string) error {
    if bson.IsObjectIdHex(id) {
        return db.C("users").FindId(bson.ObjectIdHex(id)).One(&u)
    } else {
        return errors.New("It is not ID")
    }
}

/// Get all the users for a company
func (u *User) GetCompanyUsers(db *mgo.Database, company string) (users []User) {

    err := db.C("users").Find(bson.M{"company": company}).All(&users)

    if err != nil {
        return nil
    }
    return users
}

/// Get all the UUIDs for a user
func (u *User) GetUserUuids(db *mgo.Database, id string) (uuids []Uuid) {

    err := db.C("uuids").Find(bson.M{"user": id}).Limit(50).All(&uuids)

    if err != nil {
        return nil
    }
    return uuids
}

/// Check if a user exists in the database
func (u *User) Authenticate(db *mgo.Database, email string, password string) error {
    h := md5.New()
    io.WriteString(h, password)
    hex_password := hex.EncodeToString(h.Sum(nil))
    err := db.C("users").Find(map[string]string{
        "password": hex_password,
        "email":    email,
    }).One(&u)
    return err
}

/// TODO
// I think this gets all the response messages for devices owned by a given user, based on their e-mail
func (u *User) GetUEs(db *mgo.Database, email string) (ues []AmqpResponseMessage) {

    deviceName := email + "_SRT"
    err := db.C("responsemessages").Find(bson.M{"device_name": deviceName}).Limit(50).All(&ues)

    if err != nil {
        return nil
    }
    return ues
}

/* End Of File */
