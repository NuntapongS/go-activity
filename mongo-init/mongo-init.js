db.createUser({
    user: "root",
    pwd: "password",
    roles: [{
        role: "readWrite",
        db: "goecho"
    }]
})

db.staffs.createIndex({ username: 1 }, { unique: true })