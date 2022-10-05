// admin
db.auth("mongoadm", "mongoadm");

// user
userdb = db.getSiblingDB("cadet_system");
userdb.createUser({
  user: "cadetuser",
  pwd: "123qweASD!A",
  roles: [{ role: "readWrite", db: "cadet_system" }],
  mechanisms: ["SCRAM-SHA-1"],
  passwordDigestor: "client",
});
userdb.auth("cadetuser", "123qweASD!A");
