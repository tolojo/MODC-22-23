var express = require('express');
var router = express.Router();
var mUsers = require("../models/usersModel");

router.post('/login', async function (req, res, next) {
  let nome = req.body.nome;
  let password = req.body.password;
  let result = await mUsers.login(nome, password);
  res.status(result.status).send(result.result);
});

router.post("/adduser", async function (req, res, next) {
  let user = req.body;
  let result = await mUsers.AddUser(user);
  res.status(result.status).send(result);
});

module.exports = router;
