const express = require("express");
const mysql = require("mysql");
const bodyParser = require("body-parser");
var cors = require("cors");

const app = express();
app.use(bodyParser.json());
app.use(cors());

const db = mysql.createConnection({
  host: process.env.MYSQL_HOST,
  user: process.env.MYSQL_USER,
  password: process.env.MYSQL_PASSWORD,
  database: "covidhub"
});
const port = 8080;

db.connect(err => {
  if (err) {
    throw err;
  }
  console.log("MySQL connected.");
});

app.get('/', (req, res) => {
  return res.status(200).json({response: "working"})
})

/*
A MOCK AUTHENTICAION AND AUTHORIZATION SERVICE
*/

/*
  JSON FORMAT
{
  "email": ""
  "password": ""
}
*/
app.post("/login", (req, res) => {
  let sql = "SELECT * FROM covidhub.users WHERE email = ?;";
  db.query(sql, req.body.email, (err, result) => {
    if (err)
      return res
        .status(500)
        .json({ response: { msgBody: "Database error", msgError: true } });

    if (result.length === 0) {
      return res
        .status(500)
        .json({ response: { msgBody: "User does not exist", msgError: true } });
    }

    if (result[0]["loggedIn"] == true) {
      return res.status(500).json({
        response: { msgBody: "User is already logged in", msgError: true }
      });
    }

    if (result[0]["password"] === req.body.password) {
      sql = "UPDATE covidhub.users SET loggedIn = true WHERE email = ?;";
      db.query(sql, req.body.email, err => {
        if (err)
          return res
            .status(500)
            .json({ response: { msgBody: "Database error", msgError: true } });
        return res.status(200).json({
          response: {
            msgBody: "Logged in user",
            msgError: false,
            user: result[0]
          }
        });
      });
    }
  });
});

/*
  JSON FORMAT
{
  "email": ""
  "password": ""
}
*/
app.post("/logout", (req, res) => {
  let sql = "SELECT * FROM covidhub.users WHERE email = ?;";
  db.query(sql, req.body.email, (err, result) => {
    if (err)
      return res
        .status(500)
        .json({ response: { msgBody: "Database error", msgError: true } });

    if (
      result.length !== 0 &&
      result[0]["password"] === req.body.password &&
      result[0]["loggedIn"] == true
    ) {
      sql = "UPDATE covidhub.users SET loggedIn = false WHERE email = ?;";
      db.query(sql, req.body.email, err => {
        if (err)
          return res
            .status(500)
            .json({ response: { msgBody: "Database error", msgError: true } });
        return res
          .status(200)
          .json({ response: { msgBody: "Logged out user", msgError: false } });
      });
    } else {
      return res
        .status(500)
        .json({ response: { msgBody: "Invalid user logout", msgError: true } });
    }
  });
});

/*
  JSON FORMAT
{
  email: "Michasffhi19@gmail.com"
  firstname: "ddssf"
  lastname: "fdff"
  password: "hhb%7&j5"
  phonenumber: "+12345667789"
}
*/

app.post("/register", (req, res) => {
  let sql = "SELECT * FROM covidhub.users WHERE email = ?;";
  db.query(sql, req.body.email, (err, result) => {
    if (err)
      return res
        .status(500)
        .json({ response: { msgBody: "Database error", msgError: true } });
    if (result.length !== 0) {
      return res
        .status(500)
        .json({ response: { msgBody: "User already exists", msgError: true } });
    } else {
      sql =
        "INSERT INTO covidhub.users (lastname, firstname, email, password, phoneNumber, loggedIn) VALUES (?, ?, ?, ?, ?, false);";
      db.query(
        sql,
        [
          req.body.lastname,
          req.body.firstname,
          req.body.email,
          req.body.password,
          req.body.phonenumber
        ],
        (err, result) => {
          if (err)
            return res.status(500).json({
              response: { msgBody: "Database error", msgError: true }
            });
          return res
            .status(200)
            .json({ response: { msgBody: "User created", msgError: false } });
        }
      );
    }
  });
});

/*
  JSON FORMAT
{
  "email": ""
  "password": ""
}
*/
app.get("/authorized", (req, res) => {
  let sql = "SELECT * FROM covidhub.users WHERE email = ?;";
  db.query(sql, req.body.email, (err, result) => {
    if (err)
      return res
        .status(500)
        .json({ response: { msgBody: "Database error", msgError: true } });
    if (
      result.length !== 0 &&
      result[0]["password"] === req.body.password &&
      result[0]["loggedIn"] == true
    ) {
      return res
        .status(200)
        .json({ response: { msgBody: "User authorized", msgError: false } });
    } else {
      return res
        .status(401)
        .json({ response: { msgBody: "User unauthorized", msgError: true } });
    }
  });
});

app.listen(port, () => {
  console.log(`Example app listening on port ${port}!`);
});
