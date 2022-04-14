let ovsdb = db.getSiblingDB('ovsdb');
if (db.getUser("ovsdbuser") == null){
    db.createUser({user :"ovsdbuser",pwd:"notsecure",roles : ["dbAdmin"]});
}

db.createCollection('auth');
db.createCollection('sequence');
db.createCollection('tasks');
db.createCollection('resources');
db.createCollection('journal');