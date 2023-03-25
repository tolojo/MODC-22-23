const res = require("express/lib/response");
var pool = require("./connection");

module.exports.login = async function (nome, password) {
    try {
        let sql = "Select * from utilizadores where use_nome = '"+nome+"' and use_password = '"+password+"'";
        console.log(sql)
        let result = await pool.query(sql);
        console.log(result.rows)
        if (result.rows.length > 0)
            return {
                status: 200,
                result: result.rows[0]
            };
        else return {
            status: 401,
            result: {
                msg: "Nome ou password incorreta"
            }
        };
    } catch (err) {
        console.log(err);
        return {
            status: 500,
            result: err
        };
    }
}

module.exports.AddUser = async function (user) {
 try {
    console.log(user);
    let sql = "Insert into utilizadores (use_nome,use_email,use_password) values($1,$2,$3)";
        result = await pool.query(sql, [user.nome,user.email,user.password]);
        let utilizador = result.rows;     
        console.log(result);
        return {
            status: 200,
            result: utilizador
        };
        
    } catch (err) {
        console.log(err);
        return {
            status: 500,
            result: err
        };
    }
}
